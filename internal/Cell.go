package gameoflife

type Cell struct {
	alive bool
}

func (c *Cell) Kill()  {
	c.alive = false
}

func (c *Cell) Revive() {
	c.alive = true
}

func (c *Cell) IsAlive() bool {
	return c.alive
}