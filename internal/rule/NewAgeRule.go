package rule


func NewAgeRule(next Rule) Rule {
	return func(c Cell, n []*Cell) *Cell {
		if c.IsAlive() && c.Age >= 20 { // Dies because of old
			//fmt.Printf("Cell(%d,%d) dies because of old age\n", c.Pos.col, c.Pos.row)
			c.Kill()
			return &c
		}

		return next(c, n)
	}
}