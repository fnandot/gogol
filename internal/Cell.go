package gameoflife

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"math/rand"
)

type Cell struct {
	Alive bool
	Age   uint8
	Vec   Vec
	Pol   Polygon
}

type Vec struct {
	X uint
	Y uint
}

func NewCell(width float64, X uint, Y uint) *Cell {

	alive := false

	if rand.Intn(100) > 90 {
		alive = true
	}

	pol := NewHexagon(uint(width), X, Y)

	return &Cell{
		Alive: alive,
		Pol:   pol[:],
		Age:   0,
		Vec: Vec{
			Y: Y,
			X: X,
		},
	}
}

func (c *Cell) Kill() {
	c.Alive = false
}

func (c *Cell) Revive() {
	c.Age = 0
	c.Alive = true
}

func (c *Cell) IsAlive() bool {
	return c.Alive
}

func (c *Cell) Tick() {
	if c.IsAlive() {
		c.Age++
	} else if c.Age > 0 {
		c.Age--
	}
}

func (c *Cell) Draw(t pixel.Target) {

	if false == c.IsAlive() && c.Age <= 0 {
		return
	}

	draw := imdraw.New(nil)

	draw.Color = color.RGBA{R: 0xff, G: 0xff - c.Age*10, B: 0xff - c.Age*20, A: 0xff}

	draw.Push(c.Pol[:]...)

	if c.IsAlive() {
		draw.Polygon(0)
	} else {
		draw.Polygon(1)
	}

	draw.Draw(t)
}

func (c *Cell) DrawEmpty(t pixel.Target) {
	draw := imdraw.New(nil)
	draw.Color = color.RGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xff}
	draw.Push(c.Pol[:]...)
	draw.Polygon(1)
	draw.Draw(t)
}
