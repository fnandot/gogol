package main

import (
	"bytes"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/letnando/gogol/internal"
	"golang.org/x/image/colornames"
	"io/ioutil"
	"log"
	"math/rand"
	"compress/gzip"
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
		VSync:  true,
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
	grid.DrawEmpty(canvas)

	generation := 0

	record := false
	filename := fmt.Sprintf("gogol-%d.sim.gz", time.Now().Unix())

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Name = filename

	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	//if err != nil {
	//	panic(err)
	//}

	save := func(c *gameoflife.Cell) [10]byte {
		var r [10]byte

		r[0] = c.Age

		if c.Alive {
			r[1] = 0x1
		} else {
			r[1] = 0x0
		}

		//var xBin []byte
		//binary.LittleEndian.PutUint32(xBin, uint32(c.Vec.X))
		//
		//copy(r[:], xBin[:4])
		//
		//var yBin []byte
		//binary.LittleEndian.PutUint32(yBin, uint32(c.Vec.Y))
		//
		//copy(r[:], yBin[:4])

		return r
	}

	for !win.Closed() {

		if true == record {
			//_, err := f.Write(gameoflife.FlatMap(grid.Matrix, save))
			_, err := zw.Write(gameoflife.FlatMap(grid.Matrix, save))

			if nil != err {
				fmt.Printf("failed to buffer due %s", err.Error())
			}

			if buf.Len() >= 8192 {
				err = ioutil.WriteFile(filename, buf.Bytes(), 0666)

				if nil != err {
					fmt.Printf("failed to save due %s", err.Error())
				}

				buf.Reset()
			}
		}

		win.Clear(colornames.Black)
		canvas.Draw(win, pixel.IM)

		rand.Seed(time.Now().UnixNano())

		population, born, _ := grid.Tick()

		if born > 0 {
			generation++
		}

		batch.Clear()
		grid.Draw(batch)
		batch.Draw(win)

		frames++
		select {
		case <-second.C:
			win.SetTitle(fmt.Sprintf("%s | Population: %d | Generation: %d | FPS: %d", cfg.Title, population, generation, frames))
			frames = 0




		default:
		}

		win.Update()
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Finished game")
}
