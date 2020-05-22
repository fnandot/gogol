package gogol

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type Debug struct {
	enabled        bool
	font           *basicfont.Face
	stats          DebugStats
	fpsText        *text.Text
	generationText *text.Text
	populationText *text.Text
}

type DebugStats struct {
	Fps        int
	Generation int
	Population int
}

func NewDebug(enabled bool, cfg pixelgl.WindowConfig) *Debug {

	windowHeight := cfg.Bounds.Max.Y

	font := basicfont.Face7x13
	basicTextAtlas := text.NewAtlas(font, text.ASCII)
	fps := text.New(pixel.V(10, windowHeight-20), basicTextAtlas)
	generation := text.New(pixel.V(10, windowHeight-float64(font.Height)-20), basicTextAtlas)
	population := text.New(pixel.V(10, windowHeight-float64(font.Height)*2-20), basicTextAtlas)

	return &Debug{
		enabled:        enabled,
		font:           font,
		fpsText:        fps,
		generationText: generation,
		populationText: population,
	}
}

func (d *Debug) Draw(t pixel.Target) {

	if false == d.enabled {
		return
	}

	d.fpsText.Draw(t, pixel.IM)
	d.populationText.Draw(t, pixel.IM)
	d.generationText.Draw(t, pixel.IM)

	d.fpsText.Clear()
	_, _ = fmt.Fprintf(d.fpsText, "FPS: %d", d.stats.Fps)

	d.populationText.Clear()
	_, _ = fmt.Fprintf(d.populationText, "Population: %d", d.stats.Population)

	d.generationText.Clear()
	_, _ = fmt.Fprintf(d.generationText, "Generation: %d", d.stats.Generation)
}

func (d *Debug) Update(stats DebugStats) {

	if false == d.enabled {
		return
	}

	d.stats = stats
}
