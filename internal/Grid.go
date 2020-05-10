package gameoflife

import (
	"github.com/faiface/pixel"
	"math"
)

type Grid struct {
	Matrix [][]*Cell
}

func NewGrid(cellHeight float64, windowWidth int, windowHeight int) *Grid {

	cellWidth := 4 * (cellHeight / float64(2) / math.Sqrt(float64(3)))
	cellsPerRow := int(float64(windowWidth)/(cellWidth*0.75)) - 1
	cellsPerCol := int(float64(windowHeight)/cellHeight) - 1

	grid := make([][]*Cell, cellsPerCol, cellsPerCol)

	for row := range grid {
		grid[row] = make([]*Cell, cellsPerRow, cellsPerRow)
		for col := 0; col < cellsPerRow; col++ {
			grid[row][col] = NewCell(cellWidth, cellHeight, col, row)
		}
	}

	return &Grid{
		Matrix: grid,
	}
}

func (g *Grid) Tick() int {

	var toKill []*Cell
	var toBorn []*Cell

	aliveCells := 0

	for i := range g.Matrix {
		for _, c := range g.Matrix[i] {

			c.Tick()

			if c.IsAlive() {
				aliveCells++
			}

			var n []*Cell
			maxCol := len(g.Matrix[0]) - 1
			maxRow := len(g.Matrix) - 1

			if c.row > 0 {
				n = append(n, g.Matrix[c.row-1][c.col])
				if c.col < maxRow {
					n = append(n, g.Matrix[c.row-1][c.col+1])
				}
			}

			if c.col < maxCol {
				n = append(n, g.Matrix[c.row][c.col+1])
			}

			if c.col > 0 {
				n = append(n, g.Matrix[c.row][c.col-1])
				if c.row < maxRow {
					n = append(n, g.Matrix[c.row+1][c.col-1])
				}

			}

			if c.row < maxRow {
				n = append(n, g.Matrix[c.row+1][c.col])

				if c.col < maxCol {
					n = append(n, g.Matrix[c.row+1][c.col+1])
				}
			}

			if c.col > 0 && c.row > 0 {
				n = append(n, g.Matrix[c.row-1][c.col-1])
			}

			alive := sum(n...)

			if false == c.IsAlive() && alive == 3 { // Born
				toBorn = append(toBorn, c)
				aliveCells++
			} else if c.IsAlive() && (alive < 2 || alive > 3) { // Dies because underpopulation or overpopulation
				toKill = append(toKill, c)
				aliveCells--
			} else if c.age >= 15 { // Dies because of old
				toKill = append(toKill, c)
				aliveCells--
			}
		}
	}

	for _, c := range toKill {
		c.Kill()
	}

	for _, c := range toBorn {
		c.Revive()
	}

	return aliveCells
}

func sum(cs ...*Cell) int {
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
			g.Matrix[i][c.col].Draw(t)
		}
	}
}

func (g *Grid) DrawSlot(t pixel.Target) {
	for i := range g.Matrix {
		for _, c := range g.Matrix[i] {
			g.Matrix[i][c.col].DrawSlot(t)
		}
	}
}
