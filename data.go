package main

//Person holds information for where a person is and their health status
type Person struct {
	x, y         float64 //position of the person
	status       int     //susceptible = 0, exposed = 1, infected = 2, recovered = 3, removed = 4, quaren =5
	areaX, areaY int     //what square it is in
	dayInfect    int     //number of days infected
	quarantine   int     //are they in quarantine? no = 0, yes = 1, yes but can leave = 2
	closeCom     int     //are they in a bordered-off community? no = 0, yes = 1, yes but can leave = 2
	dayCluster   int     //number of days spent in a hotspot or cluster
}

//Community represents a community (ie. neighborhood, county, country, etc.)
type Community struct {
	size    float64    //size of the board, this is a square for drawing purposes
	squares [][]Square //this is a board of squares
}

type Square struct {
	peopleList []*Person //contains a list of people in a square
	brderE     float64   //% of people to enter this square
	brderL     float64   //% of people to leave this square
	infection  float64   //infection threshold
}
