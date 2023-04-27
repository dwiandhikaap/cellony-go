package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Menu struct {
}

func NewMenu() *Menu {
	menu := &Menu{}

	return menu
}

func (s *Menu) Update() error {
	return nil
}

func (s *Menu) Draw(screen *ebiten.Image) {

}

func (s *Menu) Open() {
}

func (s *Menu) Close() {
}
