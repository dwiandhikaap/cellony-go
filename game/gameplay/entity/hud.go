package ent

import (
	comp "cellony/game/gameplay/component"
	"cellony/game/menu"
	"cellony/game/scene"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/yohamta/donburi"
)

func CreateHUDEntity(world donburi.World, sceneManager *scene.SceneManager) donburi.Entity {
	hud := world.Create(comp.HUD)

	hudEntry := world.Entry(hud)
	hudData := comp.HUD.Get(hudEntry)
	hudData.Menu = createHUDMenu(sceneManager)

	return hud
}

func createHUDMenu(sceneManager *scene.SceneManager) *menu.Menu {
	rootContainer := widget.NewContainer()

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &menu.Menu{
		SceneManager: sceneManager,
		UI:           &ui,
	}
}
