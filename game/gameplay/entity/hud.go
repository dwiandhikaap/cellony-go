package ent

import (
	"cellony/game/assets"
	comp "cellony/game/gameplay/component"
	"cellony/game/menu"
	"cellony/game/scene"
	"image/color"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
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
	titleImage := assets.AssetsInstance.Textures["brand-wide.png"]

	titleImageResized := resize.Resize(0, 100, *titleImage, resize.Lanczos3)
	titleImageEbiten := ebiten.NewImageFromImage(titleImageResized)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(
					widget.Insets{
						Top: 150,
					},
				),
				widget.RowLayoutOpts.Spacing(50),
			),
		),
	)

	header := widget.NewGraphic(
		widget.GraphicOpts.Image(titleImageEbiten),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				},
			),
		),
	)

	menuButtonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(20, 20),
			),
		),

		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				},
			),
		),
	)

	startButtonImage := loadButtonImage()
	startButton := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(startButtonImage),

		widget.ButtonOpts.Text("  Start  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			sceneManager.SelectScene(1)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				},
			),
		),
	)

	exitButton := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(startButtonImage),

		widget.ButtonOpts.Text("  Exit  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
				},
			),
		),
	)

	menuButtonContainer.AddChild(startButton)
	menuButtonContainer.AddChild(exitButton)

	rootContainer.AddChild(header)
	rootContainer.AddChild(menuButtonContainer)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &menu.Menu{
		SceneManager: sceneManager,
		UI:           &ui,
	}
}

// idle rgba(0, 171, 52, 0.8)
// hover rgba(0, 171, 52, 1)
// pressed rgba(0, 171, 52, 1)
func loadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.RGBA{R: 0, G: 159, B: 48, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 0, G: 138, B: 41, A: 255})
	pressed := image.NewNineSliceColor(color.RGBA{R: 0, G: 159, B: 48, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}
