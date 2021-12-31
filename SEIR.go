package main

import (
	"bufio"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//InitializeCommunity creates a community with a starting number of people in the community
//Inputs: Number of people, number of squares on a board, infection threshold,
//			# of communities who closed their borders, # of communities who are in quarantine but don't close borders,
//			# of communites who are in quarantine but close borders
//Returns: a list of communities that closed their borders
func (c *Community) InitializeCommunity(numPeople, slices int, ithreshold float64, numCloseB, numQuarantine int) ([][]int, [][]int, *Community) {

	//Create board with slices
	c.squares = make([][]Square, slices)
	for i := range c.squares {
		c.squares[i] = make([]Square, slices)
	}

	//Randomly select communities to close
	closeCom := make([][]int, 0)
	for i := 0; i < numCloseB; i++ {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(slices)
		rand.Seed(time.Now().UnixNano())
		y := rand.Intn(slices)
		if checkCom(x, y, closeCom) {
			closeCom = append(closeCom, []int{x, y})
		}
	}

	//Randomly select communities to close but NOT CLOSE borders (aka quarantine)
	quaranCom := make([][]int, 0)
	for i := 0; i < numQuarantine; i++ {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(slices)
		rand.Seed(time.Now().UnixNano())
		y := rand.Intn(slices)
		if checkCom(x, y, quaranCom) && checkCom(x, y, closeCom) {
			quaranCom = append(quaranCom, []int{x, y})
		}
	}

	for i := range closeCom {
		var p Person
		x := closeCom[i][0]
		y := closeCom[i][1]
		p.status = 0
		p.quarantine = 1
		p.closeCom = 1
		c.assignP(x, y, slices, p, 0, 0)
	}

	for i := range quaranCom {
		var p Person
		x := quaranCom[i][0]
		y := quaranCom[i][1]
		p.status = 0
		p.quarantine = 1
		p.closeCom = 0
		c.assignP(x, y, slices, p, 1, 0)
	}

	for i := 0; i < numPeople-len(closeCom)-len(quaranCom); i++ {
		var p Person
		//determine initial status based on prob of infection
		rand.Seed(time.Now().UnixNano())
		infectPercent := rand.Float64()
		if infectPercent < ithreshold {
			p.status = 2
		} else {
			p.status = 0
		}

		//randomly place the people, cuz ain't nobody in qurantine
		p.x = rand.Float64() * c.size
		p.y = rand.Float64() * c.size
		p.quarantine = 0
		p.closeCom = 0

		p.setSquare(c.size, slices)

		//change status if a person is in a bordered-off community or in quarantine community

		if p.inCom(closeCom, int(c.size), slices) {
			p.status = 0
			p.quarantine = 1
			p.closeCom = 1
		}

		if p.inCom(quaranCom, int(c.size), slices) {
			p.status = 0
			p.quarantine = 1
		}

		//create a square with people and add the people to the square
		if c.squares[p.areaX][p.areaY].peopleList == nil {
			var sq Square
			sq.peopleList = []*Person{}
			sq.brderE = 1.0
			sq.brderL = 1.0
			c.squares[p.areaX][p.areaY].peopleList = sq.peopleList
		}

		c.squares[p.areaX][p.areaY].peopleList = append(c.squares[p.areaX][p.areaY].peopleList, &p)

	}

	retC := c.copyBoard()
	return closeCom, quaranCom, retC
}

//InitializeCommunityFile initializes a community, creating communities with given entry and leaving % from a txt file
//(see README for format)
//Input: # of people, n squares for (n x n board), infection threshold
func (c *Community) InitializeCommunityFile(numPeople int, ithreshold float64, name string) ([][]int, [][]int, [][]int, int, []int, *Community) {
	//Read and process file
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var com, cluster [][]int
	var entP, levP, infectP []float64
	var pplL, dayClus []int
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var slices int
	for i, ln := range text {

		if i == 0 {
			slices, _ = strconv.Atoi(text[i])
		}

		//put the squares to border off in a list
		if i == 1 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s)%2 != 0 {
				panic("Please put in correct amount x,y coordinates")
			}

			for i := 0; i < len(s); i++ {
				x, _ := strconv.Atoi(s[i])
				y, _ := strconv.Atoi(s[i+1])
				com = append(com, []int{x, y})
				i++
			}
		}

		//put entry percentages in a list
		if i == 2 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s) != len(com) {
				panic("Please put the correct # of entry %")
			}
			for i := 0; i < len(s); i++ {
				e, _ := strconv.ParseFloat(s[i], 64)
				entP = append(entP, e)
			}
		}

		//put leave percentages in a list
		if i == 3 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s) != len(com) {
				panic("Please put the correct # of leaving %")
			}
			for i := 0; i < len(s); i++ {
				l, _ := strconv.ParseFloat(s[i], 64)
				levP = append(levP, l)
			}
		}

		//put infection percentages in a list
		if i == 4 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s) != len(com) {
				panic("Please put the correct # of infection %")
			}
			for i := 0; i < len(s); i++ {
				i, _ := strconv.ParseFloat(s[i], 64)
				infectP = append(infectP, i)
			}
		}

		//put number of people in a list
		if i == 5 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s) != len(com) {
				panic("Mismatch: Number of people don't match the number of squares")
			}
			for i := 0; i < len(s); i++ {
				p, _ := strconv.Atoi(s[i])
				pplL = append(pplL, p)
			}
		}

		//put number of cluster coordinates in a list
		if i == 6 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s)%2 != 0 {
				panic("Mismatch: Please enter the right number of x and y coordinates")
			}

			for i := 0; i < len(s); i++ {
				x, _ := strconv.Atoi(s[i])
				y, _ := strconv.Atoi(s[i+1])
				cluster = append(cluster, []int{x, y})
				i++
			}
		}

		//put number of days spent in each cluster in a list
		if i == 7 && len(ln) != 0 {
			s := strings.Split(ln, " ")
			if len(s) != len(cluster) {
				panic("Mismatch: Number of days don't match the number of clusters")
			}
			for i := 0; i < len(s); i++ {
				d, _ := strconv.Atoi(s[i])
				dayClus = append(dayClus, d)
			}
		}

		if i > 7 {
			panic("Please go look at file format in README")
		}
	}

	//Make sure that the cluster squares do not have an entry % of 0
	for i := range com {
		for j := range cluster {
			if com[i][0] == cluster[j][0] && com[i][1] == cluster[j][1] && entP[i] == 0 {
				log.Fatalf("Cluster: %v,%v, cannot have an entry percent of 0", cluster[j][0], cluster[j][1])
			}
		}
	}

	//Create board with slices
	c.squares = make([][]Square, slices)
	for i := range c.squares {
		c.squares[i] = make([]Square, slices)
	}

	peopleCount := 0
	emptySquareCount := 0
	emptySquare := make([][]int, 0)
	//ensures that there is at least one person in the square
	for i := range com {
		//make sure it's not a square where no one can enter
		if entP[i] != 0 {
			if len(pplL) > 0 {
				for j := 0; j < pplL[i]; j++ {
					var p Person
					//determine initial status based on prob of infection
					infectPercent := rand.Float64()
					if infectPercent < infectP[i] {
						p.status = 2
					} else {
						p.status = 0
					}

					x := com[i][0]
					y := com[i][1]

					p.quarantine = 2
					p.closeCom = 2
					c.assignP(x, y, slices, p, entP[i], levP[i])
					peopleCount++
				}
			}
		} else {
			emptySquareCount++
			x := com[i][0]
			y := com[i][1]
			emptySquare = append(emptySquare, []int{x, y})
			if c.squares[x][x].peopleList == nil {
				var sq Square
				sq.peopleList = []*Person{}
				sq.brderE = entP[i]
				sq.brderL = levP[i]
				c.squares[x][y].peopleList = []*Person{}
			}
		}
	}

	startSq := make([][]int, 0)
	//Remove all squares that are empty
	for i := 0; i < slices; i++ {
		for j := 0; j < slices; j++ {
			startSq = append(startSq, []int{i, j})
		}
	}
	avalSq := removeEmptySq(startSq, emptySquare)
	if len(pplL) != 0 {
		avalSq = removeEmptySq(avalSq, com)
	}

	for i := 0; i < numPeople+emptySquareCount-peopleCount; i++ {
		var p Person
		//determine initial status based on prob of infection
		infectPercent := rand.Float64()
		if infectPercent < ithreshold {
			p.status = 2
		} else {
			p.status = 0
		}

		//randomly place the people, cuz ain't nobody in qurantine
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(avalSq))
		rand.Seed(time.Now().UnixNano())
		p.areaHandle(avalSq[i][0], avalSq[i][1], c.size, float64(slices))
		p.quarantine = 0
		p.closeCom = 0

		p.setSquare(c.size, slices)

		//change status if a person is in a bordered-off community or in quarantine community

		if p.inCom(com, int(c.size), slices) {
			p.quarantine = 2
			p.closeCom = 2
		}

		//create a square with people and add the people to the square
		if c.squares[p.areaX][p.areaY].peopleList == nil {
			var sq Square
			sq.peopleList = []*Person{}
			sq.brderE = 1
			sq.brderL = 1
			c.squares[p.areaX][p.areaY].peopleList = sq.peopleList
		}

		c.squares[p.areaX][p.areaY].peopleList = append(c.squares[p.areaX][p.areaY].peopleList, &p)
	}

	drawSq := removeEmptySq(com, emptySquare)
	retC := c.copyBoard()
	return drawSq, emptySquare, cluster, slices, dayClus, retC
}

//areaHandle given i,j, places a person randomly in that square
func (p *Person) areaHandle(i, j int, width, slices float64) {
	sqSize := width / slices
	x := float64(i * int(sqSize))
	y := float64(j * int(sqSize))
	p.x = x + rand.Float64()*sqSize
	p.y = y + rand.Float64()*sqSize
	p.setSquare(width, int(slices))
}

//removeEmptySq removes squares that no one should be in, aka has an entry percent of 0
func removeEmptySq(aSq, emptySq [][]int) (avalSq [][]int) {
	for i := range aSq {
		check := true
		for _, a := range emptySq {
			if a[0] == aSq[i][0] && a[1] == aSq[i][1] {
				check = false
			}
		}
		if check {
			avalSq = append(avalSq, aSq[i])
		}
	}
	return
}

//assignP assigns a person to a square in the community if there is no one there
//Input: a square coord, a person, and entry, leaving threshold
func (c *Community) assignP(x, y, slice int, p Person, entry, leave float64) {
	sq := c.size / float64(slice)
	if c.squares[x][y].peopleList == nil {
		var s Square
		s.peopleList = []*Person{}
		s.brderE = entry
		s.brderL = leave
		c.squares[x][y] = s
		rand.Seed(time.Now().UnixNano())
		p.x = float64(x)*sq + (rand.Float64() * sq)
		p.y = float64(y)*sq + (rand.Float64() * sq)
		p.setSquare(c.size, slice)
		c.squares[p.areaX][p.areaY].peopleList = append(c.squares[p.areaX][p.areaY].peopleList, &p)
	} else {
		rand.Seed(time.Now().UnixNano())
		p.x = float64(x)*sq + (rand.Float64() * sq)
		p.y = float64(y)*sq + (rand.Float64() * sq)
		p.setSquare(c.size, slice)
		c.squares[p.areaX][p.areaY].peopleList = append(c.squares[p.areaX][p.areaY].peopleList, &p)
	}
}

//checkCom checks if a community has already been designed to be a closed community
//true if community has not been yet assigned
func checkCom(x, y int, c [][]int) bool {
	for i := range c {
		if x == c[i][0] && y == c[i][1] {
			return false
		}
	}
	return true
}

//inCloseCom checks if someone is in a community given a list of communities
//true if they are in a community in the list
func (p *Person) inCom(Com [][]int, size, slice int) (check bool) {
	check = false
	for i := range Com {
		if p.areaX == Com[i][0] && p.areaY == Com[i][1] {
			check = true
			return
		}
	}
	return
}

//setSquare designates a person to a square given the community and the number of overall squares on the board
//p.areaX, p.areaX is the top left coordinate in a square
func (p *Person) setSquare(width float64, slices int) {
	p.areaX = int(math.Floor((p.x / width) * float64(slices)))
	p.areaY = int(math.Floor((p.y / width) * float64(slices)))
}

//infect determines whether to infect another person based on if someone is infected give infection rate
//Input: a person, infection rate
func (p1 *Person) infect(p2 *Person, ethreshold float64) {
	if p1.status == 2 && p2.status == 0 && p2.closeCom != 1 {
		exposedPercent := rand.Float64()
		if exposedPercent < ethreshold {
			p2.status = 1
		}
	}
}

//countRates counts the number of susectible, exposed, infected, and removed people given a board
func (c *Community) countRates() (s, e, i, r int) {
	for a := range c.squares {
		for j := range c.squares[a] {
			for _, p := range c.squares[a][j].peopleList {
				if p.status == 0 {
					s++
				} else if p.status == 1 {
					e++
				} else if p.status == 2 {
					i++
				} else if p.status == 3 || p.status == 4 {
					r++
				}
			}
		}
	}
	return
}

//collectSquares collect the people within aoe of a square
//Input: x,y of a square, aoe, area to infect, and number of slices of the board
func (c *Community) collectSquares(x, y int, aoe, slices int) (listP []*Person) {
	for i := Max(0, x-aoe); i < Min(slices, x+aoe+1); i++ {
		for j := Max(0, y-aoe); j < Min(slices, y+aoe+1); j++ {
			listP = append(listP, c.squares[i][j].peopleList...)
		}
	}
	return
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

//updateInfected determines who gets infected next
//Input: aoe, number of slices on board, min days infectious, exposed, infected, recovered percentage, lethalness
func (c *Community) updateInfected(aoe, slices, minInfect int, ethreshold, ithreshold, rthreshold, lethal float64) {
	for i := range c.squares {
		for j := range c.squares[i] {
			susceptible := c.collectSquares(i, j, aoe, slices)
			for _, p1 := range c.squares[i][j].peopleList {
				p1.update(susceptible, ithreshold, rthreshold, ethreshold, lethal, minInfect)
			}
		}
	}
}

//update checks if a person is exposed, if so turns into infected
//if person is infected, turns into removed/recovered
//Input: infected threshold, removed/recovered threhold, lethalness, min days infectious
func (p *Person) update(susceptible []*Person, ithreshold, rthreshold, ethreshold, lethality float64, minInfect int) {
	if p.status == 1 {
		infectPercent := rand.Float64()
		if infectPercent < ithreshold {
			p.status = 2
		}
	} else if p.status == 2 {
		if p.dayInfect > minInfect {
			recoverPercent := rand.Float64()
			if recoverPercent < rthreshold {
				livedPercent := rand.Float64()
				if livedPercent < lethality {
					p.status = 4
				} else {
					p.status = 3
				}

			}
		} else {
			for _, p2 := range susceptible {
				p.infect(p2, ethreshold)
			}
		}
		p.dayInfect += 1
	} else if p.closeCom == 1 {
		p.status = 0
	}
}

//copyBoard creates a deep copy of a community
func (c Community) copyBoard() (retC *Community) {
	retC = &Community{}
	retC.size = c.size
	retC.squares = make([][]Square, len(c.squares))
	for i := range retC.squares {
		retC.squares[i] = make([]Square, len(c.squares[i]))
	}
	for i := range c.squares {
		for j := range c.squares[i] {
			for _, p := range c.squares[i][j].peopleList {
				retC.squares[i][j].peopleList = append(retC.squares[i][j].peopleList, p.copyPerson())
			}
		}
	}
	return
}

//copyPerson returns a shallow copy of a person
func (p *Person) copyPerson() (retp *Person) {
	retp = &Person{}
	retp.x, retp.y = p.x, p.y
	retp.status = p.status
	retp.areaX, retp.areaY = p.areaX, p.areaY
	retp.dayInfect = p.dayInfect
	retp.closeCom = p.closeCom
	return
}

//Move randomly moves a person
//Input: step length - larger value means more movement, canvas size for checking purposes
func (p *Person) randMove(stepL float64, canvasWidth, slice int) {
	step := rand.Float64() * stepL
	angle := rand.Float64() * 2 * math.Pi
	p.x += step * math.Cos(angle)
	p.y += step * math.Sin(angle)

	//check if person is leaving canvas
	for p.x < 0 || p.x > float64(canvasWidth) {
		step = rand.Float64() * stepL
		angle = rand.Float64() * 2 * math.Pi
		p.x += step * math.Cos(angle)
	}
	for p.y < 0 || p.y > float64(canvasWidth) {
		step = rand.Float64() * stepL
		angle = rand.Float64() * 2 * math.Pi
		p.y += step * math.Sin(angle)

	}
	p.setSquare(float64(canvasWidth), slice)
}

// stayInSquare forces a person to stay in their square or community
func (p *Person) stayInSquare(stepL float64, canvasWidth, slice int) {
	step := rand.Float64() * stepL
	sqSize := float64(canvasWidth / slice)
	r := sqSize / 4
	for !(p.x > float64(p.areaX)*sqSize+r && p.x < (float64(p.areaX)*sqSize)+sqSize-r) {
		step := rand.Float64() * step
		angle := rand.Float64() * math.Pi * 2
		p.x += step * math.Cos(angle)
	}

	for !(p.y > float64(p.areaY)*sqSize+r && p.y < (float64(p.areaY)*sqSize)+sqSize-r) {
		step := rand.Float64() * step
		angle := rand.Float64() * math.Pi * 2
		p.y += step * math.Sin(angle)
	}
}

//moveAwayBorder checks if a person in close to a closed community and keeps outsiders from entering a community
func (p *Person) moveAwayBorder(comClose [][]int, stepL float64, canvasWidth, slice int) {
	p.setSquare(float64(canvasWidth), slice)
	for p.inCom(comClose, canvasWidth, slice) {
		sqSize := float64(canvasWidth / slice)
		step := rand.Float64() * stepL
		r := sqSize / 4
		for i := range comClose {
			//top left coordinate
			x := float64(comClose[i][0]) * sqSize
			y := float64(comClose[i][1]) * sqSize

			//left
			if p.x > x-r && p.x < x {
				angle := math.Pi/2 + rand.Float64()*math.Pi
				var change float64
				if math.Cos(angle) > 0 {
					change = -math.Cos(angle)
					p.x += step * change
				} else {
					p.x += step * math.Cos(angle)
				}

			}
			//right
			if p.x < (x+sqSize)+r && p.x > x+sqSize {
				angle := math.Pi/2 + rand.Float64()*math.Pi
				p.x += step * math.Abs(math.Cos(angle))
			}

			//bottom
			if p.y < (y+sqSize)+r && p.y > (y+sqSize) {
				angle := rand.Float64() * math.Pi
				var change float64
				if math.Sin(angle) > 0 {
					change = -math.Sin(angle)
					p.y += step * change
				} else {
					p.y += step * math.Sin(angle)
				}

			}

			//top
			if p.y > y-r && p.y < y {
				angle := math.Pi + rand.Float64()*math.Pi
				p.y += step * math.Abs(math.Sin(angle))
			}

			//check if they got into the square
			for p.x > x-r && p.x < x+sqSize+r && p.y > y-r && p.y < y+sqSize+r {
				angle := rand.Float64() * 2 * math.Pi
				p.x += step * math.Cos(angle)
				p.y += step * math.Sin(angle)
			}
		}
		p.setSquare(float64(canvasWidth), slice)
	}
}

//randMoveBoard randomly moves everyone in a community
//Input: step length - larger value means more movement, slices, closed-off comunities, quaranting communities, a cluster boolean,  coordinate to cluster at, and number of days to spend at a cluster
func (c *Community) randMoveBoard(stepL float64, slice int, closeCom, qCom [][]int, clust bool, clustX, clustY, clusDay int) {
	var aX, aY int
	for i := range c.squares {
		for j := range c.squares[i] {
			for _, p := range c.squares[i][j].peopleList {
				if p.status != 4 {
					if p.quarantine == 1 {
						aX = p.areaX
						aY = p.areaY
					}
					p.randMove(stepL, int(c.size), slice)
					if p.quarantine == 1 {
						p.areaX = aX
						p.areaY = aY
						p.stayInSquare(stepL, int(c.size), slice)
					} else if p.closeCom == 0 {
						if clust == true && p.dayCluster < clusDay {
							p.moveTo(float64(clustX), float64(clustY), float64(slice), c.size, stepL)
							p.dayCluster++
						}
						if len(closeCom) > 0 {
							p.moveAwayBorder(closeCom, stepL, int(c.size), slice)
						}
					}
				}
			}
		}
	}
}

//updateBoard updates the whole board
//Input: aoe = area to infect, slice = number to slice up board, canvasWidth = how big the board is, minInfect = min days infectious, ithreshold = infection threshold, stepL = step length, lethal = lethalness, closeCom = list of closed communities, qCom = list of quarantining communites, clust = cluster boolean, clustX, clustY = square to cluster at, clusDay = number of days to spend at a cluster, cS, cE, cI, cR = counts of susceptible, exposed, infected, and removed people
func (c *Community) updateBoard(aoe, slices, minInfect int, ethreshold, ithreshold, rthreshold, stepL, lethal float64, closeCom, qCom [][]int, clust bool, clustX, clustY, clusDay int) (retC *Community, cS, cE, cI, cR []int) {
	c.updateInfected(aoe, slices, minInfect, ethreshold, ithreshold, rthreshold, lethal)
	c.randMoveBoard(stepL, slices, closeCom, qCom, clust, clustX, clustY, clusDay)
	retC = (*c).copyBoard()
	s, e, i, r := c.countRates()
	cS = append(cS, s)
	cE = append(cE, e)
	cI = append(cI, i)
	cR = append(cR, r)
	return
}

//updateBoardFile updates the whole board given the entry and leaving percents specificed from a file
//Input: aoe = area to infect, slice = number to slice up board, minInfect = min days infectious, ithreshold = infection threshold, stepL = step length, lethal = lethalness, closeCom = list of closed communities, eCom = list of communities to not enter, cluster = list of squares to cluster at, dayCluster = days to spend at a cluster or hotspot, cS, cE, cI, cR = counts of susceptible, exposed, infected, and removed people
func (c *Community) updateBoardFile(aoe, slices, minInfect int, ethreshold, ithreshold, rthreshold, stepL, lethal float64, closeCom, eCom, cluster [][]int, dayClus []int) (retC *Community, cS, cE, cI, cR []int) {
	c.updateInfected(aoe, slices, minInfect, ethreshold, ithreshold, rthreshold, lethal)
	c.moveBoard(stepL, slices, closeCom, eCom, cluster, dayClus)
	retC = (*c).copyBoard()
	s, e, i, r := c.countRates()
	cS = append(cS, s)
	cE = append(cE, e)
	cI = append(cI, i)
	cR = append(cR, r)
	return
}

//moveBoard randomly moves everyone in a community
//Input: step length - larger value means more movement
func (c *Community) moveBoard(stepL float64, slice int, closeCom, eCom, cluster [][]int, dayClust []int) {
	var aX, aY int
	for i := range c.squares {
		for j := range c.squares[i] {
			for _, p := range c.squares[i][j].peopleList {
				if p.status != 4 {
					if p.closeCom == 2 {
						aX = p.areaX
						aY = p.areaY
					}
					p.randMove(stepL, int(c.size), slice)
					if p.closeCom == 2 {
						p.leaveMove(c.squares[i][j].brderL, stepL, int(c.size), slice, aX, aY, cluster, dayClust)
					} else if p.closeToCom(stepL, c.size, float64(slice), closeCom) {
						p.enterMove(c.squares[i][j].brderL, stepL, int(c.size), slice, closeCom, cluster, dayClust)
					}
					if p.closeToCom(stepL, c.size, float64(slice), eCom) {
						p.moveAwayBorder(eCom, stepL, int(c.size), slice)
					}
				}
			}
		}
	}
}

//leaveMove determines if a person in a closed square should leave or not based on the leaving percent
func (p *Person) leaveMove(percent, step float64, width, slice, aX, aY int, cluster [][]int, cThres []int) {
	var num int
	if len(cluster) > 0 {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(len(cluster))
	}

	decide := rand.Float64()
	if decide < percent && percent != 1 {
		p.areaX = aX
		p.areaY = aY
		p.stayInSquare(step, width, slice)
	} else if len(cluster) > 0 && len(cThres) > 0 && p.dayCluster < cThres[num] {
		p.moveTo(float64(cluster[num][0]), float64(cluster[num][1]), float64(slice), float64(width), step)
	} else {
		p.randMove(step, width, slice)
	}

}

//enterMove determines if a person should move into a closed square given they are close to it based on the entering percent
func (p *Person) enterMove(percent, step float64, width, slice int, close, cluster [][]int, cThres []int) {
	var num int
	if len(cluster) > 0 {
		rand.Seed(time.Now().UnixNano())
		num = rand.Intn(len(cluster))
	}

	decide := rand.Float64()
	if decide < percent {
		p.moveAwayBorder(close, step, width, slice)
	} else if len(cluster) > 0 && len(cThres) > 0 && p.dayCluster < cThres[num] {
		p.moveTo(float64(cluster[num][0]), float64(cluster[num][1]), float64(slice), float64(width), step)
	} else {
		p.randMove(step, width, slice)
	}
}

//closeToCom checks if someone is close to a community
func (p *Person) closeToCom(step, width, slice float64, comClose [][]int) bool {
	sqSize := width / slice
	for i := range comClose {
		//top left coordinate
		x := float64(comClose[i][0]) * sqSize
		y := float64(comClose[i][1]) * sqSize

		//left
		if p.x > x-step && p.x < x {
			return true
		}
		//right
		if p.x < (x+sqSize)+step && p.x > x+sqSize {
			return true
		}
		//bottom
		if p.y < (y+sqSize)+step && p.y < y+step {
			return true
		}
		//top
		if p.y > y-step && p.y < y {
			return true
		}

		//check if they got into the square
		for p.x > x-step && p.x < x+sqSize+step && p.y > y-step && p.y < y+sqSize+step {
			return true
		}
	}
	return false
}

// moveTo moves a person towards a square
// Input: x,y of a square (top left coordinate)
func (p *Person) moveTo(sx, sy, slices, width, step float64) {
	sqSize := width / slices
	//top left coordinate
	x := sx * sqSize
	y := sy * sqSize

	difX := math.Abs(p.x - x)
	difY := math.Abs(p.y - y)

	newX := 0.0
	newY := 0.0

	for difX > newX && difY > newY {
		//left
		if p.x < x {
			angle := math.Pi/2 + rand.Float64()*math.Pi
			p.x += step * math.Abs(math.Cos(angle))
		}

		//right
		if p.x > x+sqSize {
			angle := math.Pi/2 + rand.Float64()*math.Pi
			var change float64
			if math.Cos(angle) > 0 {
				change = -math.Cos(angle)
				p.x += step * change
			} else {
				p.x += step * math.Cos(angle)
			}
		}

		//bottom
		if p.y > y+sqSize {
			angle := rand.Float64() * math.Pi
			var change float64
			if math.Sin(angle) > 0 {
				change = -math.Sin(angle)
				p.y += step * change
			} else {
				p.y += step * math.Sin(angle)
			}
		}

		//top
		if p.y < y {
			angle := math.Pi + rand.Float64()*math.Pi
			p.y += step * math.Abs(math.Sin(angle))
		}

		//if in square, we good
		if p.y < y+sqSize && p.y > y && p.x < x+sqSize && p.x > x {
			p.dayCluster++
			break
		}

		newX = math.Abs(p.x - x)
		newY = math.Abs(p.y - y)
	}
	p.setSquare(width, int(slices))
}
