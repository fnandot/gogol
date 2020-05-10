package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/letnando/gogol/internal"
	"golang.org/x/image/colornames"
	"log"
	"time"
)

const (
	windowWidth  = 1024
	windowHeight = 768
)



func main() {
	log.Print("Started game")
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		log.Fatal(err)
	}

	grid := gameoflife.NewGrid(float64(6), windowWidth, windowHeight)

	canvas := pixelgl.NewCanvas(pixel.R(-windowWidth, -windowHeight, windowWidth, windowHeight))

	var (
		frames = 0
		second = time.NewTicker(time.Second)
	)

	imd := imdraw.New(nil)
	imd.Clear()
	imd.Color = colornames.White

	//ticker := time.NewTicker(500 * time.Millisecond)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, nil)
	grid.DrawSlot(canvas)

	for !win.Closed() {

		win.Clear(colornames.Black)
		canvas.Draw(win, pixel.IM)
		population := grid.Tick()
		batch.Clear()
		grid.Draw(batch)
		batch.Draw(win)

		frames++
		select {
		case <-second.C:
			win.SetTitle(fmt.Sprintf("%s | Population: %d | FPS: %d", cfg.Title, population, frames))
			frames = 0
		default:
		}

		win.Update()
	}
}

