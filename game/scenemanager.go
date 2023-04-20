package game

import "github.com/hajimehoshi/ebiten/v2"

var (
	transitionFrom = ebiten.NewImage(ScreenWidth, ScreenHeight)
	transitionTo   = ebiten.NewImage(ScreenWidth, ScreenHeight)
)

type GameState struct {
	SceneManager *SceneManager
}

type Scene interface {
	Update(state *GameState)
	Draw(screen *ebiten.Image)
}

type SceneManager struct {
	current         Scene
	next            Scene
	transitionCount int
}
