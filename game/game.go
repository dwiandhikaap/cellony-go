package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"cellony/game/assets"
	"cellony/game/gameplay"
	input "cellony/game/input"
)

type Game struct {
	sceneManager SceneManager
}

func CreateGame() *Game {
	err := assets.InitializeAssets()
	if err != nil {
		panic(err)
	}

	ebiten.SetVsyncEnabled(false)

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
	ScreenWidth  = 1280
	ScreenHeight = 720
)

func (g *Game) Update() error {
	input.Update()
	g.sceneManager.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
	fps := fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, fps)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
