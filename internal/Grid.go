package gogol

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"log"
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
	immutableRules  Rule
}

type CoordinateVector [2]int

func NewGrid(cfg GameConfig) *Grid {

	cols := int(CellsInWidth(cfg.Width, cfg.Margin, cfg.CellSize))
	rows := int(CellsInHeight(cfg.Height, cfg.Margin, cfg.CellSize))

	log.Printf("Created new grid with %d rows and %d columns, so %d cells", rows, cols, rows * cols)

	grid := make([][]*Cell, rows, rows)

	for row := range grid {
		grid[row] = make([]*Cell, cols, cols)
		for col := 0; col < cols; col++ {
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

	immutableRules := []RuleConstructor{
		//NewDetectDeadBorders,
		NewDoNotDeadShadowIfNotNeighbors,
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
		rows:            rows,
		cols:            cols,
		ruleChain:       NewRuleChain(rules),
		immutableRules:  NewRuleChain(immutableRules),
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
		immutableRules:  g.immutableRules,
	}

	for row := range ng.Data {
		ng.Data[row] = make([]*Cell, len(g.Data[row]), len(g.Data[row]))
		copy(ng.Data[row], g.Data[row])
	}

	return &ng
}

func (g *Grid) Update(clock int64) (p int) {

	log.Printf("-------------------------------------------")

	population := 0
	ng := g.copy()


	for i := range g.Data {
		for _, c := range g.Data[i] {
			c.Update(clock)

			nc := g.ruleChain(c, g, g.GetNeighbors(c))
			//fmt.Printf("%s == %s\n", c.DeadVisited, nc.DeadVisited)

			if nc.IsAlive() {
				population++
			}

			ng.Data[c.Pos.row][c.Pos.col] = nc
		}
	}

	for i := range g.Data {
		for _, c := range g.Data[i] {
			ng.Data[c.Pos.row][c.Pos.col] = g.immutableRules(c, g, g.GetNeighbors(c))
		}
	}

	totalVisited := 0
	totalWithBorders := 0

	for i := range ng.Data {
		for _, c2 := range ng.Data[i] {
			if c2.DeadVisited {
				totalVisited++
			}

			if c2.DeadBorder {
				totalWithBorders++
			}
		}
	}

	log.Printf("[GRID] Total visited %d and with borders %d\n", totalVisited, totalWithBorders)

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
