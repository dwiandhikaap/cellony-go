package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"cellony/game/gameplay"
)

type Game struct {
	sceneManager SceneManager
}

func CreateGame() *Game {
	g := Game{
		sceneManager: SceneManager{
			scenes: []Scene{
				gameplay.CreateWorldScene(),
			},
		},
	}

	return &g
}

var (
	ScreenWidth  = 640
	ScreenHeight = 480
)

func (g *Game) Update() error {
	g.sceneManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello World!")

	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
