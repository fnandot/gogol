package gogol

import (
	"log"
	"math/rand"
	"sync"
)

type Rule = func(c *Cell, g *Grid, n []*Cell) *Cell
type RuleConstructor = func(next Rule) Rule

func NewBornRule(next Rule) Rule {
	return func(c *Cell, g *Grid, n []*Cell) *Cell {
		aliveNeighbors := SumAliveCells(n...)

		if false == c.IsAlive() && aliveNeighbors >= 2 && rand.Intn(100) > 75 {
			c.Born()
			return c
		}

		return next(c, g, n)
	}
}

func NewPopulationRule(next Rule) Rule {
	return func(c *Cell, g *Grid, n []*Cell) *Cell {
		aliveNeighbors := SumAliveCells(n...)

		if c.IsAlive() && (aliveNeighbors < 2 || aliveNeighbors > 3) && rand.Intn(100) > 28 {
			c.Kill()
			return c
		}

		return next(c, g, n)
	}
}

func NewAgeRule(next Rule) Rule {
	return func(c *Cell, g *Grid, n []*Cell) *Cell {
		if c.IsAlive() && c.age >= 20 { // Dies because of old
			//fmt.Printf("Cell(%d,%d) dies because of old age\n", c.Pos.col, c.Pos.row)
			c.Kill()
			return c
		}

		return next(c, g, n)
	}
}

func NewSpontaneousGenerationRule(next Rule) Rule {
	return func(c *Cell, g *Grid, n []*Cell) *Cell {
		if false == c.IsAlive() && rand.Intn(100000) > 99988 { //Born
			c.Born()
			return c
		}

		return next(c, g, n)
	}
}

func NewDoNotDeadShadowIfNotNeighbors(next Rule) Rule {
	return func(c *Cell, g *Grid, n []*Cell) *Cell {

		if c.IsAlive() {
			return c
		}

		if c.DeadCounter > 0  && SumDeadWithShadowCells(n...) > 0 {
			c.DeadCounter--
		}

		return next(c, g, n)
	}
}

func NewDetectDeadBorders(next Rule) Rule {

	return func(c *Cell, g *Grid, n []*Cell) *Cell {

		if c.Pos.row != 0 || c.Pos.col != 0 {
			//os.Exit(1)
			return c
		}

		var wg sync.WaitGroup
		counter := g.cols * g.rows
		log.Printf("[BORDER] Starting in Cell(%d,%d)\n", c.Pos.row, c.Pos.col)
		wg.Add(1)
		go BorderDetector(c, g, &counter, &wg)
		wg.Wait()
		log.Printf("[BORDER] Remaining to visit after in Cell(%d,%d) %d\n", c.Pos.row, c.Pos.col, counter)

		return next(c, g, n)
	}
}

func BorderDetector(c *Cell, g *Grid, wg *int, wg2 *sync.WaitGroup) {

	defer wg2.Done()
	//fmt.Printf("Remaining to visit in Cell(%d,%d) %d\n", c.Pos.row, c.Pos.col, *wg)
	c.DeadVisited = true

	defer func() {
		*wg--
	}()

	nc := g.GetNeighbors(c)
	ncv := 0

	if HasBorder(nc...) {
		c.DeadBorder = true
		//fmt.Printf("Border detected for Cell(%d,%d)\n", c.Pos.row, c.Pos.col)
	}

	//fmt.Printf("Cell(%d,%d) has %d neightbors\n", c.Pos.row, c.Pos.col, len(nc))
	for _, nc := range nc {
		//fmt.Printf("Cell(%d,%d) neightbor Cell(%d,%d) is %s\n", c.Pos.row, c.Pos.col, nc.Pos.row, nc.Pos.col, nc.DeadVisited)
		if !nc.DeadVisited {
			ncv++
			wg2.Add(1)
			go BorderDetector(nc, g, wg, wg2)
		}
	}
	//time.Sleep(time.Millisecond * 500)
	//fmt.Printf("Cell(%d,%d) has %d visited neightbors\n", c.Pos.row, c.Pos.col, ncv)
}

func NewRuleChain(rules []RuleConstructor) Rule {
	rule := func(c *Cell, g *Grid, n []*Cell) *Cell {

		if c.Pos.row != 0 || c.Pos.col != 0 {
			//os.Exit(1)
			return c
		}

		totalVisited := 0
		totalWithBorders := 0

		for i := range g.Data {
			for _, c2 := range g.Data[i] {
				if c2.DeadVisited {
					totalVisited++
				}

				if c2.DeadBorder {
					totalWithBorders++
				}
			}
		}

		log.Printf("[NOOP] Total visited %d and with borders %d\n", totalVisited, totalWithBorders)

		return c
	}

	for _, ruleConstructor := range rules {
		rule = ruleConstructor(rule)
	}

	return rule
}
