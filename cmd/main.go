package main

import (
	"math/rand"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"log"
	"math"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
	"github.com/letnando/game-of-life/internal"
)

const (
	windowWidth  = 1024
	windowHeight = 768
)

const cellsPerRow = 50


func main() {
	log.Print("Started game")



	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		log.Fatal(err)
	}

	grid := make([][]gameoflife.Cell, cellsPerRow, cellsPerRow)
	for i, _ := range grid {
		grid[i] = make([]gameoflife.Cell, cellsPerRow, cellsPerRow)
		for j, _ := range grid[i] {
			if rand.Intn(100) > 90 {
				(&grid[i][j]).Revive()
			}
		}
	}

	cellSize := windowWidth / 50

	for !win.Closed() {
		win.Clear(colornames.Black)
		drawGrid(cellSize, win)

		for i, _ := range grid {
			for j, _ := range grid[i] {
				if rand.Intn(100) > 98 {
					(&grid[i][j]).Revive()
				} else {
					(&grid[i][j]).Kill()
				}
			}
		}

		for rowNumber, row := range grid {
			for columnNumber, cell := range row {
				if cell.IsAlive() {
					imd := imdraw.New(nil)
					imd.Clear()
					imd.Color = colornames.White
					x := float64(rowNumber * cellSize)
					y := float64(columnNumber * cellSize)

					imd.Push(pixel.V(x, y), pixel.V(x + float64(cellSize), y + float64(cellSize)))
					imd.Rectangle(0)
					imd.Draw(win)
				}
			}
		}

		win.Update()
	}
}

func drawGrid(cellSize int, win *pixelgl.Window) {
	fields := int(math.Abs(float64(windowWidth / cellSize)))

	for i := 0; i < fields; i++ {
		imd, _ := drawVerticalLine(i * cellSize)
		imd.Draw(win)

		imd, _ = drawHorizontalLine(i * cellSize)
		imd.Draw(win)
	}
}

func drawVerticalLine(x int) (*imdraw.IMDraw, error) {

	imd := imdraw.New(nil)
	imd.Clear()
	imd.Color = colornames.White
	imd.Push(pixel.V(float64(x), 0), pixel.V(float64(x), windowHeight))
	imd.Line(1)

	return imd, nil
}

func drawHorizontalLine(y int) (*imdraw.IMDraw, error) {

	imd := imdraw.New(nil)
	imd.Clear()
	imd.Color = colornames.White
	imd.Push(pixel.V(0, float64(y)), pixel.V(windowWidth, float64(y)))
	imd.Line(1)

	return imd, nil
}
