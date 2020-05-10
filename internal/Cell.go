package gameoflife

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
)

type Cell struct {
	alive bool
	age   int
	col   int
	row   int
	hex   Hexagon
}

type Hexagon [6]pixel.Vec

func NewCell(width float64, height float64, col int, row int) *Cell {

	alive := false

	if rand.Intn(100) > 70 {
		alive = true
	}

	return &Cell{
		alive: alive,
		hex:   NewHexagon(width, height, col, row),
		age:   -1,
		row:   row,
		col:   col,
	}
}

func (c *Cell) Kill() {
	c.alive = false
}

func (c *Cell) Revive() {
	c.age = 0
	c.alive = true
}

func (c *Cell) IsAlive() bool {
	return c.alive
}

func (c *Cell) Tick() {
	if c.IsAlive() {
		c.age++
	} else if c.age > 0 {
		c.age--
	}
}

func (c *Cell) Draw(t pixel.Target) {

	if false == c.IsAlive() && c.age <= 0 {
		return
	}

	draw := imdraw.New(nil)

	if c.age < 3 {
		draw.Color = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	} else if c.age < 5 {
		draw.Color = colornames.Orange
	} else if c.age < 7 {
		draw.Color = colornames.Orangered
	} else {
		draw.Color = colornames.Red
	}

	draw.Push(c.hex[:]...)

	if c.IsAlive() {
		draw.Polygon(0)
	} else {
		draw.Polygon(1)
	}

	draw.Draw(t)
}

func (c *Cell) DrawSlot(t pixel.Target) {
	draw := imdraw.New(nil)
	draw.Color = colornames.Gray
	draw.Push(c.hex[:]...)
	draw.Polygon(1)
	draw.Draw(t)
}

func NewHexagon(width float64, height float64, col int, row int) Hexagon {

	x := float64(col)*(width*0.75) + width*0.5
	y := float64(row)*height + height/2

	if col%2 == 1 {
		y += height / 2
	}

	return Hexagon{
		pixel.V(x, y),
		pixel.V(x+width*0.50, y),
		pixel.V(x+width*0.75, y+height/2),
		pixel.V(x+width*0.50, y+height),
		pixel.V(x, y+height),
		pixel.V(x-width*0.25, y+height/2),
	}
}
