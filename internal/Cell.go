package gogol

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"math"
	"math/rand"
)

type Cell struct {
	Alive       bool
	age         uint8
	deadCounter uint16
	Pos         Coordinate
	pol         Polygon
}

type Coordinate struct {
	col int
	row int
}

func NewCell(pol Polygon, col int, row int) *Cell {

	alive := false

	if rand.Intn(100) > 90 {
		alive = true
	}

	return &Cell{
		Alive:       alive,
		pol:         pol,
		age:         0,
		deadCounter: math.MaxUint16,
		Pos: Coordinate{
			row: row,
			col: col,
		},
	}
}

func (c *Cell) Kill() {
	c.Alive = false
	c.deadCounter = 0
}

func (c *Cell) Born() {
	c.age = 0
	c.Alive = true
}

func (c *Cell) IsAlive() bool {
	return c.Alive
}

func (c *Cell) Update(clock int64) {
	if c.IsAlive() && c.age < math.MaxUint8 {
		c.age++
	} else if false == c.IsAlive() && c.deadCounter < math.MaxUint8 {
		c.deadCounter++
	}
}

func (c *Cell) DrawDead(t pixel.Target) bool {
	if false == c.IsAlive() && c.deadCounter < math.MaxUint8 {
		draw := imdraw.New(nil)

		step := uint8(30 - c.deadCounter/9)
		fmt.Printf("Alpha is %d\n", 30 - step)

		draw.Color = color.RGBA{R: 75, G: 0, B: 130, A: 0}

		draw.Push(c.pol[:]...)
		draw.Polygon(0)
		draw.Draw(t)
		return true
	}

	return false
}

func (c *Cell) Draw(t pixel.Target) {
	if false == c.IsAlive() {
		return
	}

	draw := imdraw.New(nil)
	draw.Color = color.RGBA{R: 200, G: safeUint8Mul(c.age, 20), B: 255, A: 0xff}

	draw.Push(c.pol[:]...)
	draw.Polygon(0)
	draw.Draw(t)
}

func safeUint8Mul(a, b uint8) uint8 {
	mul := a * b
	if mul < a {
		return math.MaxInt8
	}
	return mul
}

func (c *Cell) DrawEmpty(t pixel.Target) {
	draw := imdraw.New(nil)
	draw.Color = color.RGBA{R: 0x0F, G: 0x0F, B: 0x0F}
	draw.Push(c.pol[:]...)
	draw.Polygon(1)
	draw.Draw(t)
}
