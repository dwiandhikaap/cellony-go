package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"cellony/game/config"
	"cellony/game/gameplay"
	input "cellony/game/input"
	"cellony/game/menu"
	"cellony/game/scene"
)

type Game struct {
	sceneManager *scene.SceneManager
}

func CreateGame() *Game {
	ebiten.SetVsyncEnabled(false)

	sceneManager := scene.SceneManager{
		Scenes: []scene.Scene{},
	}

	menuScene := menu.CreateMenuScene(&sceneManager)
	gameplayScene := gameplay.CreateWorldScene()

	sceneManager.Scenes = append(sceneManager.Scenes, menuScene)
	sceneManager.Scenes = append(sceneManager.Scenes, gameplayScene)

	g := Game{
		sceneManager: &sceneManager,
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
