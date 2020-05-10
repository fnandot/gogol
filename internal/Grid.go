package gameoflife

import (
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

type Grid struct {
	Matrix [][]*Cell
}

func NewGrid(cellSize float64, windowWidth int, windowHeight int) *Grid {

	cellWidth := 4 * (cellSize / float64(2) / math.Sqrt(float64(3)))
	cellsPerRow := int(float64(windowWidth)/(cellWidth*0.75)) - 1
	cellsPerCol := int(float64(windowHeight)/cellSize) - 1

	grid := make([][]*Cell, cellsPerCol, cellsPerCol)

	for row := range grid {
		grid[row] = make([]*Cell, cellsPerRow, cellsPerRow)
		for col := 0; col < cellsPerRow; col++ {
			grid[row][col] = NewCell(cellSize, uint(col), uint(row))
		}
	}

	return &Grid{
		Matrix: grid,
	}
}

func (g *Grid) Tick() (pop int, born int, died int) {

	var willDie []*Cell
	var willBorn []*Cell

	population := 0

	for i := range g.Matrix {
		for _, c := range g.Matrix[i] {

			c.Tick()

			if c.IsAlive() {
				population++
			}

			aliveNeighbours := sumAlive(g.getNeighbours(c)...)

			if false == c.IsAlive() && aliveNeighbours >= 2 && rand.Intn(100) > 75 { //Born
				//fmt.Printf("Cell(%d,%d) will born!\n", c.Vec.Y, c.Vec.X)
				willBorn = append(willBorn, c)
				population++
			} else if c.IsAlive() && (aliveNeighbours < 2 || aliveNeighbours > 3) && rand.Intn(100) > 35 { // Dies because underpopulation or overpopulation
				//fmt.Printf("Cell(%d,%d) will die!\n", c.Vec.Y, c.Vec.X)
				willDie = append(willDie, c)
				population--
			} else if c.IsAlive() && c.Age >= 20 { // Dies because of old
				willDie = append(willDie, c)
				population--
			} else if false == c.IsAlive() && rand.Intn(100000) > 99988 { //Born
				//fmt.Printf("Cell(%d,%d) will born!\n", c.Vec.Y, c.Vec.X)
				willBorn = append(willBorn, c)
				population++
			}
		}
	}

	for _, c := range willDie {
		c.Kill()
	}

	for _, c := range willBorn {
		c.Revive()
	}

	return population, len(willBorn), len(willDie)
}

func (g *Grid) getNeighbours(c *Cell) []*Cell {
	var n []*Cell
	maxCol := uint(len(g.Matrix[0]) - 1)
	maxRow := uint(len(g.Matrix) - 1)

	if c.Vec.Y > 0 {
		n = append(n, g.Matrix[c.Vec.Y-1][c.Vec.X]) //
		if c.Vec.X < maxRow {
			n = append(n, g.Matrix[c.Vec.Y-1][c.Vec.X+1]) //
		} else {
			n = append(n, g.Matrix[c.Vec.Y-1][0]) //
		}
	} else {
		n = append(n, g.Matrix[maxRow][c.Vec.X]) //
	}

	if c.Vec.X < maxCol {
		n = append(n, g.Matrix[c.Vec.Y][c.Vec.X+1]) //
	} else {
		n = append(n, g.Matrix[c.Vec.Y][0]) //
	}

	if c.Vec.X > 0 {
		n = append(n, g.Matrix[c.Vec.Y][c.Vec.X-1]) //
		if c.Vec.Y < maxRow {
			//n = append(n, g.Matrix[c.Vec.Y+1][c.Vec.X-1]) // Only Square Grid
		}
	} else {
		n = append(n, g.Matrix[c.Vec.Y][maxCol]) //
	}

	if c.Vec.Y < maxRow {
		n = append(n, g.Matrix[c.Vec.Y+1][c.Vec.X]) //

		if c.Vec.X < maxCol {
			n = append(n, g.Matrix[c.Vec.Y+1][c.Vec.X+1]) //
		} else {
			n = append(n, g.Matrix[c.Vec.Y+1][0]) //
		}
	} else {
		n = append(n, g.Matrix[0][c.Vec.X]) //
	}

	if c.Vec.X > 0 && c.Vec.Y > 0 {
		//n = append(n, g.Matrix[c.Vec.Y-1][c.Vec.X-1]) // Only Square Grid
	} else {
		//n = append(n, g.Matrix[maxRow][maxCol]) // Only Square Grid
	}
	return n
}

func sumAlive(cs ...*Cell) int {
	sum := 0
	for _, c := range cs {
		if c.IsAlive() {
			sum += 1
		}
	}
	return sum
}

func (g *Grid) Draw(t pixel.Target) {
	for i := range g.Matrix {
		for _, c := range g.Matrix[i] {
			g.Matrix[i][c.Vec.X].Draw(t)
		}
	}
}

func (g *Grid) DrawEmpty(t pixel.Target) {
	for i := range g.Matrix {
		for _, c := range g.Matrix[i] {
			g.Matrix[i][c.Vec.X].DrawEmpty(t)
		}
	}
}
