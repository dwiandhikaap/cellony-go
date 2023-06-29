package ent

import (
	comp "cellony/game/gameplay/component"
	"cellony/game/menu"
	"cellony/game/scene"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
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

	lowerMenuBackground := image.NewNineSliceColor(color.RGBA{255, 255, 255, 255})

	lowerMenuContainer := widget.NewGraphic(
		widget.GraphicOpts.ImageNineSlice(lowerMenuBackground),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				},
			),
		),
	)

	ui.Container.AddChild(lowerMenuContainer)

	return &menu.Menu{
		SceneManager: sceneManager,
		UI:           &ui,
	}
}
