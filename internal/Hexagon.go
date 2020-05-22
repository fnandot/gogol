package gogol

import (
	"github.com/faiface/pixel"
	"math"
)

type Polygon []pixel.Vec

type Hexagon [6]pixel.Vec

// odd-q vertical layout
func NewHexagon(size, margin float64, col, row uint) Hexagon {

	height := size
	width := HexagonWidth(size)

	x := float64(col)*(width*0.75) + width/4 + (margin) // + width*0.5
	y := float64(row)*height + margin

	if col%2 == 1 {
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

func HexagonWidth(height float64) float64 {
	return 4 * (height / float64(2) / math.Sqrt(float64(3)))
}

func HexagonHeight(width float64) float64 {
	return width * 0.8660254
}

func CellsInWidth(width, margin, cellheight float64) uint {
	return uint((width - HexagonWidth(cellheight)/4 - margin * 2) / (HexagonWidth(cellheight) * 0.75))
}

func CellsInHeight(height, margin, cellheight float64) uint {
	return uint((height - cellheight/2 - margin * 2) / cellheight)
}
