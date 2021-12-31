package main

import (
	"flag"
	"fmt"
	"gifhelper"
	"image/png"
	"os"

	"github.com/Arafatk/glot"
)

func main() {
	//Initalize variables
	var numGens int
	var canvasWidth int
	var scalingFactor float64
	var frequency int
	var numPeople int
	var startinfected float64
	var exposed, infected, recovered float64
	var minInfect int
	var lethality float64
	var slices int
	var aoe int
	var stepLength float64
	var outputFilename, animOutputFile string
	var comBool bool
	var numCloseB int //We assume if a community closes their border, everyone in the community is healthy
	var numQuarantine int
	var file string
	var cluster bool
	var clusX, clusY, clusDay int
	var closedCom, qCom, sCom, eCom, clusCom [][]int
	var dayClus []int
	var countS, countE, countI, countR []int

	flag.IntVar(&numGens, "numGen", 100, "Number of steps to run the community.")
	flag.IntVar(&slices, "s", 10, "Number of squares or slices to split board.")
	flag.IntVar(&aoe, "aoe", 1, "Distance of infection | unit: slices")
	flag.Float64Var(&stepLength, "move", 50, "How much should people move?")
	flag.IntVar(&canvasWidth, "width", 1000, "Width (and height) of the image to create.")
	flag.Float64Var(&scalingFactor, "sf", 10, "a scaling factor for size of people")
	flag.Float64Var(&startinfected, "si", 0.01, "a starting % of the population that is infected")
	flag.Float64Var(&exposed, "e", 0.235, "likelihood for getting exposed given that are come in contact with someone infected")
	flag.Float64Var(&infected, "i", 0.157, "likelihood of turning infectious at a time step ")
	flag.Float64Var(&recovered, "r", 0.97, "likelihood of recovering at a time step")
	flag.IntVar(&minInfect, "minDay", 14, "min days infectious")
	flag.Float64Var(&lethality, "l", 0.027, "likelihood of dying")
	flag.StringVar(&outputFilename, "o", "out.png", "Name of PNG to output.")
	flag.StringVar(&animOutputFile, "a", "anim", "Animated GIF to write.")
	flag.IntVar(&frequency, "freq", 1, "Frame writing interval")
	flag.IntVar(&numPeople, "numPeop", 100, "Number of people in the community")
	flag.BoolVar(&comBool, "cBool", false, "Create communities (ie. countries), # of communites: slice x slice")
	flag.IntVar(&numCloseB, "numClose", 0, "Number of communities to close borders and quarantine")
	flag.IntVar(&numQuarantine, "numQ", 5, "Number of communities to quarantine but aren't closing borders")
	flag.StringVar(&file, "f", "", "Name of file to create special communities")
	flag.BoolVar(&cluster, "clus", false, "Should there be a cluster on the board?")
	flag.IntVar(&clusX, "clusX", 1, "What x coordinate should the people cluster at?")
	flag.IntVar(&clusY, "clusY", 1, "What y coordinate should the people cluster at?")
	flag.IntVar(&clusDay, "clusDay", 5, "How many days should a person spend in a hot spot?")
	flag.Parse()

	if comBool == true {
		aoe = 1
	}

	//Initalize the board for a community
	var board Community
	var w int
	board.size = float64(canvasWidth)

	if file == "" {
		//if there is no input file, use flags to create fully bordered off or quarantining communities
		var comList []*Community
		var b1 *Community
		if comBool == false {
			numCloseB = 0
			numQuarantine = 0
		}

		closedCom, qCom, b1 = board.InitializeCommunity(numPeople, slices, startinfected, numCloseB, numQuarantine)
		comList = append(comList, b1)
		//find initial counts of SEIR
		s1, e1, i1, r1 := board.countRates()
		countS = append(countS, s1)
		countE = append(countE, e1)
		countI = append(countI, i1)
		countR = append(countR, r1)

		//Create image of initial board
		img1 := board.DrawToCanvas(canvasWidth+int(scalingFactor)*2, slices, scalingFactor, comBool, closedCom, qCom)
		f1, _ := os.Create("first.png")
		png.Encode(f1, img1)
		fmt.Println("Created first image!")

		//Run the simulation for numGens # of generations and save to a list of board for gif generation

		for i := 0; i < numGens; i++ {
			b, cS, cE, cI, cR := board.updateBoard(aoe, slices, minInfect, exposed, infected, recovered, stepLength, lethality, closedCom, qCom, cluster, clusX, clusY, clusDay)
			comList = append(comList, b)
			countS = append(countS, cS...)
			countE = append(countE, cE...)
			countI = append(countI, cI...)
			countR = append(countR, cR...)
		}

		img2 := comList[1].DrawToCanvas(canvasWidth+int(scalingFactor)*2, slices, scalingFactor, comBool, closedCom, qCom)
		f2, _ := os.Create("second.png")
		png.Encode(f2, img2)
		fmt.Println("Created second image!")

		//Create final image
		img := board.DrawToCanvas(canvasWidth+int(scalingFactor)*2, slices, scalingFactor, comBool, closedCom, qCom)
		f, _ := os.Create(outputFilename)
		png.Encode(f, img)
		fmt.Println("Created final image!")

		//Create gif
		imgList := AnimateSystem(comList, canvasWidth, frequency, slices, scalingFactor, comBool, closedCom, qCom)
		gifhelper.ImagesToGIF(imgList, animOutputFile)
		fmt.Println("Created gif!")

	} else {
		var comList []*Community
		var b1 *Community
		//use file to create communities with a certain threshold of people leaving and entering
		comBool = true
		aoe = 1
		sCom, eCom, clusCom, w, dayClus, b1 = board.InitializeCommunityFile(numPeople, startinfected, file)
		comList = append(comList, b1)
		slices = w

		//Create image of initial board
		img1 := board.DrawToCanvas(canvasWidth+int(scalingFactor)*2, slices, scalingFactor, comBool, sCom, qCom)
		f1, _ := os.Create("first.png")
		png.Encode(f1, img1)
		fmt.Println("Created first image!")

		//Run the simulation for numGens # of generations and save to a list of board for gif generation
		for i := 0; i < numGens; i++ {
			b, cS, cE, cI, cR := board.updateBoardFile(aoe, slices, minInfect, exposed, infected, recovered, stepLength, lethality, sCom, eCom, clusCom, dayClus)
			comList = append(comList, b)
			countS = append(countS, cS...)
			countE = append(countE, cE...)
			countI = append(countI, cI...)
			countR = append(countR, cR...)
		}

		//Create final image
		img := board.DrawToCanvas(canvasWidth+int(scalingFactor)*2, slices, scalingFactor, comBool, sCom, qCom)
		f, _ := os.Create(outputFilename)
		png.Encode(f, img)
		fmt.Println("Created final image!")

		//Create gif
		imgList := AnimateSystem(comList, canvasWidth, frequency, slices, scalingFactor, comBool, sCom, qCom)
		gifhelper.ImagesToGIF(imgList, animOutputFile)
		fmt.Println("Created gif!")
	}

	//Make plots of the epidemic
	dimensions := 2
	persist := false
	debug := false
	plot, _ := glot.NewPlot(dimensions, persist, debug)

	pointGroupName := "Susceptible"
	style := "points"
	points := countS
	plot.AddPointGroup(pointGroupName, style, points)

	pointGroupName = "Exposed"
	style = "points"
	points = countE
	plot.AddPointGroup(pointGroupName, style, points)

	pointGroupName = "Infected"
	style = "points"
	points = countI
	plot.AddPointGroup(pointGroupName, style, points)

	pointGroupName = "Removed"
	style = "points"
	points = countR
	plot.AddPointGroup(pointGroupName, style, points)

	// A plot type used to make points/ curves and customize and save them as an image.
	plot.SetTitle("SEIR")
	// Optional: Setting the title of the plot
	plot.SetXLabel("Number of Generations")
	plot.SetYLabel("Number of People")
	// Optional: Setting label for X and Y axis
	plot.SetXrange(0, 100)
	plot.SetYrange(0, 100)
	// Optional: Setting axis ranges
	plot.SavePlot("SEIR.png")
}
