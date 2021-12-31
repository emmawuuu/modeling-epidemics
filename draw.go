package main

import (
	"canvas"
	"fmt"
	"image"
)

//AnimateSystem takes a community and generates an list of images for animation to a gif
func AnimateSystem(c []*Community, canvasWidth, frequency, slices int, scalingFactor float64, comBool bool, closeCom, quaranCom [][]int) []image.Image {
	images := make([]image.Image, 0)

	// for every instance of a board, draw to canvas and grab the image
	for i := range c {
		if i%frequency == 0 {
			fmt.Println("frame", i)
			images = append(images, c[i].DrawToCanvas(canvasWidth, slices, scalingFactor, comBool, closeCom, quaranCom))
		}
	}

	return images
}

//DrawToCanvas generates the image corresponding to a canvas after drawing a community
func (com *Community) DrawToCanvas(canvasWidth, slices int, scalingFactor float64, comBool bool, closeCom, quaranCom [][]int) image.Image {

	// set a new square canvas
	c := canvas.CreateNewPalettedCanvas(canvasWidth, canvasWidth, nil)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	// declare colors
	darkGray := MakeColor(100, 100, 100)
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)
	green := MakeColor(0, 255, 0)
	//yellow := MakeColor(255, 255, 0)
	magenta := MakeColor(255, 0, 255)
	white := MakeColor(255, 255, 255)
	cyan := MakeColor(0, 200, 255)

	//Draws the communites
	if comBool == true {
		c.SetStrokeColor(white)
		DrawGridLines(c, slices)
	}

	//Draws closed communities
	if len(closeCom) > 0 {
		c.SetStrokeColor(cyan)
		c.SetLineWidth(5)
		if len(closeCom) > 50 {
			c.SetLineWidth(0.1)
		}
		drawClosedCom(c, slices, closeCom)
	}

	//Draws quarantining communities
	if len(quaranCom) > 0 {
		c.SetStrokeColor(magenta)
		c.SetLineWidth(5)
		drawClosedCom(c, slices, quaranCom)
	}

	// range over all the people and draw them.
	for _, row := range com.squares {
		for _, col := range row {
			for _, p := range col.peopleList {
				switch p.status {
				case 0:
					c.SetFillColor(green)
				case 1:
					c.SetFillColor(white)
				case 2:
					c.SetFillColor(red)
				case 3:
					c.SetFillColor(blue)
				case 4:
					c.SetFillColor(darkGray)
				case 5:
					c.SetFillColor(magenta)
				}

				c.Circle(p.x, p.y, scalingFactor)
				c.Fill()

				if p.quarantine == 1 {
					c.SetFillColor(magenta)
					c.Circle(p.x, p.y, 4*scalingFactor/5)
					c.Fill()
				}

			}
		}
	}

	return canvas.GetImage(c)
}

func DrawGridLines(pic canvas.Canvas, slices int) {
	w, h := pic.Width(), pic.Height()
	// first, draw vertical lines
	for i := 1; i < slices; i++ {
		y := i * w / slices
		pic.MoveTo(0.0, float64(y))
		pic.LineTo(float64(w), float64(y))
	}
	// next, draw horizontal lines
	for j := 1; j < slices; j++ {
		x := j * h / slices
		pic.MoveTo(float64(x), 0.0)
		pic.LineTo(float64(x), float64(h))
	}
	pic.Stroke()
}

//drawClosedCom draws the closed communities
func drawClosedCom(pic canvas.Canvas, slices int, closeCom [][]int) {
	w, h := pic.Width(), pic.Height()
	for i := range closeCom {
		//Assume top left
		//draw left line
		x0 := closeCom[i][0] * w / slices
		y0 := closeCom[i][1] * h / slices
		yL := y0 + h/slices
		pic.MoveTo(float64(x0), float64(y0))
		pic.LineTo(float64(x0), float64(yL))

		//draw top line
		xB := x0 + w/slices
		pic.MoveTo(float64(x0), float64(y0))
		pic.LineTo(float64(xB), float64(y0))

		//draw bottom line
		pic.MoveTo(float64(x0), float64(yL))
		pic.LineTo(float64(xB), float64(yL))

		//draw right line
		pic.MoveTo(float64(xB), float64(y0))
		pic.LineTo(float64(xB), float64(yL))
	}
	pic.Stroke()
}
