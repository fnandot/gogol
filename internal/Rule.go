package gogol

import "math/rand"

type Rule = func(c Cell, n []*Cell) *Cell
type RuleConstructor = func(next Rule) Rule

func NewBornRule(next Rule) Rule {
	return func(c Cell, n []*Cell) *Cell {
		aliveNeighbors := sumAlive(n...)

		if false == c.IsAlive() && aliveNeighbors >= 2 && rand.Intn(100) > 75 {
			c.Born()
			return &c
		}

		return next(c, n)
	}
}

func NewPopulationRule(next Rule) Rule {
	return func(c Cell, n []*Cell) *Cell {
		aliveNeighbors := sumAlive(n...)

		if c.IsAlive() && (aliveNeighbors < 2 || aliveNeighbors > 3) && rand.Intn(100) > 28 {
			c.Kill()
			return &c
		}

		return next(c, n)
	}
}

func NewAgeRule(next Rule) Rule {
	return func(c Cell, n []*Cell) *Cell {
		if c.IsAlive() && c.age >= 20 { // Dies because of old
			//fmt.Printf("Cell(%d,%d) dies because of old age\n", c.Pos.col, c.Pos.row)
			c.Kill()
			return &c
		}

		return next(c, n)
	}
}

func NewSpontaneousGenerationRule(next Rule) Rule {
	return func(c Cell, n []*Cell) *Cell {
		if false == c.IsAlive() && rand.Intn(100000) > 99988 { //Born
			c.Born()
			return &c
		}

		return next(c, n)
	}
}

func NewRuleChain(rules []RuleConstructor) Rule {
	rule := func(c Cell, n []*Cell) *Cell {
		return &c
	}

	for _, ruleConstructor := range rules {
		rule = ruleConstructor(rule)
	}

	return rule
}
