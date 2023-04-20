package game

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"

	"cellony/game/resources"
)

type GameUI struct {
	UI ebitenui.UI
}

func CreateGameUI() *GameUI {
	rootContainer := widget.NewContainer()

	var res, err = resources.CreateUIResource()
	if err != nil {
		panic(err)
	}

	// Create your UI here
	var face = res.Fonts.Face
	var _ = face

	button := widget.NewButton()
	rootContainer.AddChild(button)

	return &GameUI{
		UI: ebitenui.UI{
			Container:           rootContainer,
			DisableDefaultFocus: true,
		},
	}
}
