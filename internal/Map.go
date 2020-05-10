package gameoflife

func Map(c []*Cell, f func(cell *Cell) [10]byte) []byte {
	m := make([]byte, len(c)*2)
	for _, v := range c {
		r := f(v)
		m = append(m, r[:]...)
	}
	return m
}

func Flatten(c [][]*Cell) []*Cell {
	f := make([]*Cell, 0)
	for _, v := range c {
		f = append(f, v[:]...)
	}
	return f
}

func FlatMap(c [][]*Cell, f func(cell *Cell) [10]byte) []byte {
	return Map(Flatten(c), f)
}
