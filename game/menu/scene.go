package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScene struct {
	menu *Menu
}

func CreateMenuScene() *MenuScene {
	return &MenuScene{
		menu: NewMenu(),
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
