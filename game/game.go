package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"cellony/game/config"
	"cellony/game/gameplay"
	input "cellony/game/input"
	"cellony/game/menu"
)

type Game struct {
	sceneManager SceneManager
}

func CreateGame() *Game {
	ebiten.SetVsyncEnabled(false)

	g := Game{
		sceneManager: SceneManager{
			scenes: []Scene{
				menu.CreateMenuScene(),
				gameplay.CreateWorldScene(),
			},
		},
	}

	return &g
}

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
	return int(config.Video.Width), int(config.Video.Height)
}
