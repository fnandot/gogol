package gogol

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Grid struct {
	Data            [][]*Cell
	cellBatch       *pixel.Batch
	deadCellBatch   *pixel.Batch
	emptyCellCanvas *pixelgl.Canvas
	rows            int
	cols            int
	cellSize        float64
	ruleChain       Rule
}

type CoordinateVector [2]int

func NewGrid(cfg GameConfig) *Grid {

	cellsPerRow := int(CellsInWidth(cfg.Width, cfg.Margin, cfg.CellSize))
	cellsPerCol := int(CellsInHeight(cfg.Height, cfg.Margin, cfg.CellSize))

	grid := make([][]*Cell, cellsPerCol, cellsPerCol)

	for row := range grid {
		grid[row] = make([]*Cell, cellsPerRow, cellsPerRow)
		for col := 0; col < cellsPerRow; col++ {
			pol := NewHexagon(cfg.CellSize, cfg.Margin, uint(col), uint(row))
			grid[row][col] = NewCell(pol[:], col, row)
		}
	}

	rules := []RuleConstructor{
		//NewSpontaneousGenerationRule,
		NewPopulationRule,
		NewBornRule,
		NewAgeRule,
	}

	canvas := pixelgl.NewCanvas(pixel.R(-cfg.Width, -cfg.Height, cfg.Width, cfg.Height))
	for i := range grid {
		for _, c := range grid[i] {
			c.DrawEmpty(canvas)
		}
	}

	return &Grid{
		Data:            grid,
		cellSize:        cfg.CellSize,
		rows:            cellsPerCol,
		cols:            cellsPerRow,
		ruleChain:       NewRuleChain(rules),
		cellBatch:       pixel.NewBatch(&pixel.TrianglesData{}, nil),
		deadCellBatch:   pixel.NewBatch(&pixel.TrianglesData{}, nil),
		emptyCellCanvas: canvas,
	}
}

func (g *Grid) copy() *Grid {
	ng := Grid{
		Data:            make([][]*Cell, len(g.Data), len(g.Data)),
		cellSize:        g.cellSize,
		cellBatch:       g.cellBatch,
		rows:            g.rows,
		cols:            g.cols,
		ruleChain:       g.ruleChain,
		emptyCellCanvas: g.emptyCellCanvas,
		deadCellBatch:   g.deadCellBatch,
	}

	for row := range ng.Data {
		ng.Data[row] = make([]*Cell, len(g.Data[row]), len(g.Data[row]))
		copy(ng.Data[row], g.Data[row])
	}

	return &ng
}

func (g *Grid) Update(clock int64) (p int) {

	population := 0
	ng := g.copy()

	for i := range g.Data {
		for _, c := range g.Data[i] {
			c.Update(clock)
			ng.Data[c.Pos.row][c.Pos.col] = g.ruleChain(*c, g.GetNeighbors(c))
			if ng.Data[c.Pos.row][c.Pos.col].IsAlive() {
				population++
			}
		}
	}

	*g = *ng

	return population
}

func (g *Grid) GetNeighbors(c *Cell) []*Cell {
	var n []*Cell

	neighborsCoordinates := []CoordinateVector{{+1, +1}, {+1, 0}, {0, -1}, {-1, 0}, {-1, +1}, {0, +1}}

	if c.Pos.col%2 == 0 {
		neighborsCoordinates = []CoordinateVector{{+1, 0}, {+1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {0, +1}}
	}

	for _, coordinateVector := range neighborsCoordinates {
		newRow := c.Pos.row + coordinateVector[1]
		newCol := c.Pos.col + coordinateVector[0]
		if g.existsColumn(newCol) && g.existsRow(newRow) {
			n = append(n, g.Data[newRow][newCol])
		}
	}

	return n
}

func (g *Grid) existsRow(row int) bool {
	return row >= 0 && row <= g.rows-1
}

func (g *Grid) existsColumn(col int) bool {
	return col >= 0 && col <= g.cols-1
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

	g.emptyCellCanvas.Draw(t, pixel.IM)
	g.cellBatch.Clear()
	g.deadCellBatch.Clear()

	howManyDead := 0

	for _, i := range g.Data {
		for _, c := range i {
			if c.DrawDead(g.deadCellBatch) {
				howManyDead++
			}
			c.Draw(g.deadCellBatch)
		}
	}

	//fmt.Printf("%d are drawn\n", howManyDead)

	g.cellBatch.Draw(t)
	g.deadCellBatch.Draw(t)
}

func (g *Grid) PixelToMatrixCoordinate(v pixel.Vec) *Cell {
	//c := int(v.X / g.cellSize)
	//r := int(v.Y / g.cellSize)

	//c = size * sqrt(3) * (hex.col + 0.5 * (hex.row&1))
	//c = x * sqrt(3) * (z + 0.5 * (y))

	r := int((2 * g.cellSize) / (3 * v.X))
	c := int((0.57735 * g.cellSize) / (0.5*v.X + v.Y))

	//var y = size * (                         3./2 * hex.r)
	//var q = ( 2./3 * point.x                        ) / size
	fmt.Printf("Detected mouse at %d, %d\n", r, c)

	return g.Data[r][c]
}

func (g *Grid) DrawEmpty(t pixel.Target) {
	for i := range g.Data {
		for _, c := range g.Data[i] {
			c.DrawEmpty(t)
		}
	}
}
