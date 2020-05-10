package gameoflife

import (
	"github.com/faiface/pixel"
	"math"
)

type Polygon []pixel.Vec

type Hexagon [6]pixel.Vec
type Square [4]pixel.Vec

func NewHexagon(size uint, X uint, Y uint) Hexagon {

	height := float64(size)
	width := 4 * (float64(size) / float64(2) / math.Sqrt(float64(3)))

	x := float64(X) * (width * 0.75) // + width*0.5
	y := float64(Y) * height         //+ height/2

	if X%2 == 1 {
		y += height / 2
	}

	return Hexagon{
		pixel.V(x, y),
		pixel.V(x+width*0.50, y),
		pixel.V(x+width*0.75, y+height/2),
		pixel.V(x+width*0.50, y+height),
		pixel.V(x, y+height),
		pixel.V(x-width*0.25, y+height/2),
	}
}

func NewSquare(size uint, X uint, Y uint) Square {

	s := float64(size)
	x := float64(X * size)
	y := float64(Y * size)

	return Square{
		pixel.V(x, y),
		pixel.V(x+s, y),
		pixel.V(x+s, y+s),
		pixel.V(x, y+s),
	}
}
