package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	UI *GameUI
}

var (
	ScreenWidth  = 640
	ScreenHeight = 480
)

func (g *Game) Update() error {
	g.UI.UI.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//ebitenutil.DebugPrint(screen, "Hello World!")
	g.UI.UI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
