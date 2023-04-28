package menu

import (
	"cellony/game/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScene struct {
	menu *Menu
}

func CreateMenuScene(sceneManager *scene.SceneManager) *MenuScene {
	return &MenuScene{
		menu: NewMenu(sceneManager),
	}
}

func (s *MenuScene) Update() error {
	s.menu.Update()
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	s.menu.Draw(screen)
}

func (s *MenuScene) Open() {
	s.menu.Open()
}

func (s *MenuScene) Close() {
	s.menu.Close()
}
