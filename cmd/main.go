package main

import (
	"flag"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	gogol "github.com/letnando/gogol/internal"
	"golang.org/x/image/colornames"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Print("Started game")

	res := flag.String("resolution", "1600x900", "resolution of the window in pixels [default=1600x900]")
	cs := flag.Float64("cellsize", float64(10), "cell size in pixels [default=10]")
	mr := flag.Float64("margin", float64(10), "margin of the grid in pixels [default=10]")
	fs := flag.Bool("fullscreen", false, "show on fullscreen")
	vs := flag.Bool("vsync", false, "activate vsync")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	spl := strings.Split(*res, "x")

	windowWidth, _ := strconv.ParseFloat(spl[0], 64)
	windowHeight, _ := strconv.ParseFloat(spl[1], 64)

	pixelgl.Run(func() {
		run(*fs, *debug, *vs, *cs, *mr, windowWidth, windowHeight)
	})
}

func run(fullscreen, debug, vsync bool, cellSize, margin, windowWidth, windowHeight float64) {

	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  vsync,
	}

	if true == fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		log.Fatal(err)
	}

	var (
		frames = 0
		second = time.NewTicker(time.Second)
	)

	rand.Seed(time.Now().UnixNano())

	debugMode := gogol.NewDebug(debug, cfg)
	game := gogol.NewGame(gogol.GameConfig{
		CellSize: cellSize,
		Margin:   margin,
		Width:    windowWidth,
		Height:   windowHeight,
	})

	for !win.Closed() {

		win.Clear(colornames.Black)

		game.WriteMessage("Game Started")

		//if win.Pressed(pixelgl.MouseButtonLeft) {
		//	mp := win.MousePosition()
		//	for _, c := range grid.GetNeighbors(grid.PixelToMatrixCoordinate(mp)) {
		//		c.Born()
		//	}
		//}

		game.Update(time.Now().UnixNano())

		if false == game.IsRunning() {
			game.WriteMessage("Game finished")
		}

		game.Draw(win)
		debugMode.Draw(win)

		frames++
		select {
		case <-second.C:
			win.SetTitle(fmt.Sprintf("%s | Population: %d | Generation: %d | FPS: %d", cfg.Title, game.Population, game.Generation, frames))
			debugMode.Update(gogol.DebugStats{Fps: frames, Generation: game.Generation, Population: game.Population})
			frames = 0
		default:
		}

		win.Update()
	}

	log.Printf("Finished game")
}
