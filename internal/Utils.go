package gogol

import "math"

func SumAliveCells(cs ...*Cell) int {
	sum := 0
	for _, c := range cs {
		if c.IsAlive() {
			sum += 1
		}
	}
	return sum
}

func SumDeadWithShadowCells(cs ...*Cell) int {
	sum := 0
	for _, c := range cs {
		if c.HasDeadShadow() {
			sum += 1
		}
	}
	return sum
}

func SumDeadCells(cs ...*Cell) int {
	sum := 0
	for _, c := range cs {
		if false == c.IsAlive() && c.DeadCounter == math.MaxUint8 {
			sum += 1
		}
	}
	return sum
}

func HasBorder(cs ...*Cell) bool {

	if SumDeadWithShadowCells(cs...) > 1 && SumDeadCells(cs...) > 1 {
		return true
	}

	return false
}