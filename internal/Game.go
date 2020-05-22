package gogol

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type GameConfig struct {
	CellSize float64
	Margin   float64
	Width    float64
	Height   float64
}

type Game struct {
	Grid        *Grid
	Population  int
	Generation  int
	Width       float64
	Height      float64
	batch       *pixel.Batch
	gameFont    *basicfont.Face
	gameMessage *text.Text
}

func NewGame(cfg GameConfig) *Game {
	grid := NewGrid(cfg)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, nil)

	font := basicfont.Face7x13
	basicTextAtlas := text.NewAtlas(font, text.ASCII)
	gameMessagesText := text.New(pixel.V(0, cfg.Height/2), basicTextAtlas)

	return &Game{
		Grid:        grid,
		batch:       batch,
		gameFont:    font,
		gameMessage: gameMessagesText,
		Width:       cfg.Width,
		Height:      cfg.Height,
	}
}

func (g *Game) IsRunning() bool {
	return 0 < g.Population
}

func (g *Game) WriteMessage(text string) {
	g.gameMessage.Clear()

	_, _ = g.gameMessage.WriteString(text)
}

func NewCenteredTextMatrix(gameMessagesText *text.Text, windowWidth, windowHeight float64, text string, font *basicfont.Face) pixel.Matrix {
	mat := pixel.IM.Scaled(gameMessagesText.Orig, 4)
	mat = mat.Moved(pixel.Vec{X: (windowWidth/2 - float64(len(text)/2-font.Width*4)) - gameMessagesText.Orig.X, Y: 0})
	return mat
}

func (g *Game) Draw(t pixel.Target) {
	g.batch.Clear()
	g.Grid.Draw(t)
	g.batch.Draw(t)
	if g.Generation < 50 {
		mat := NewCenteredTextMatrix(g.gameMessage, g.Width, g.Height, "test", g.gameFont)
		g.gameMessage.Draw(t, mat)
	}
}

func (g *Game) Update(clock int64) {
	g.Generation++
	g.Population = g.Grid.Update(clock)
}
