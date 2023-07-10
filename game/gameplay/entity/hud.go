package ent

import (
	"autocell/game/assets"
	comp "autocell/game/gameplay/component"
	"autocell/game/menu"
	"autocell/game/scene"
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

//var defaultCellSpeed = 30
//var defaultBrushRadius = 12
//var defaultCurrentBrush = 10
//var defaultCellCount = 0
//var defaultCellCooldown = 5000 // 5 sec

func CreateHUDEntity(ecs *ecs.ECS, sceneManager *scene.SceneManager) donburi.Entity {
	world := ecs.World
	hud := world.Create(comp.HUD)

	hudEntry := world.Entry(hud)
	hudData := comp.HUD.Get(hudEntry)
	hudData.Menu = createHUDMenu(ecs, sceneManager)

	return hud
}

func createHUDMenu(ecs *ecs.ECS, sceneManager *scene.SceneManager) *menu.Menu {
	buttonBg := loadButtonImage()

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(
					widget.Insets{
						Top: 0,
					},
				),
				widget.RowLayoutOpts.Spacing(50),
			),
		),
	)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	elements := make(map[string]*widget.Text)

	// holds sidemenucontainerleft and sidemenucontainerright horizontally
	sideMenuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(5, 5),
				widget.GridLayoutOpts.Padding(
					widget.Insets{
						Top: 30,
					},
				),
			),
		),
	)

	sideMenuBg := image.NewNineSliceColor(color.RGBA{10, 10, 10, 200})

	sideMenuContainerLeft := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(sideMenuBg),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Spacing(5, 10),
				widget.GridLayoutOpts.Padding(
					widget.Insets{
						Top:    10,
						Bottom: 10,
						Right:  10,
					},
				),
			),
		),

		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				},
			),
		),
	)

	sideMenuContainer.AddChild(sideMenuContainerLeft)

	// Cell count text
	cellCountText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			StretchHorizontal:  true,
		})),
		widget.TextOpts.Text("Cell Count: 0", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 20), color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{
			Left: 10,
		}),
	)
	elements["cellCountText"] = cellCountText
	sideMenuContainerLeft.AddChild(cellCountText)

	// Hive resource text
	hiveResourceText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			StretchHorizontal:  true,
		})),
		widget.TextOpts.Text("Hive Resource: 0", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 20), color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{
			Left: 10,
		}),
	)
	elements["hiveResourceText"] = hiveResourceText
	sideMenuContainerLeft.AddChild(hiveResourceText)

	cellSpeedSlider := createFuckingSlider("Cell Speed", 1, 100, 30, 1, func(val float64) {
		setCellSpeed(ecs, val)
	})

	hiveSpawnCountSlider := createFuckingSlider("Spawn Count", 0, 30, 10, 1, func(val float64) {
		setHiveSpawnCount(ecs, val)
	})

	hiveSpawnCooldown := createFuckingSlider("Spawn Cooldown", 100, 10000, 5000, 1000, func(val float64) {
		setHiveCooldown(ecs, val)
	})

	sideMenuContainerLeft.AddChild(cellSpeedSlider)
	sideMenuContainerLeft.AddChild(hiveSpawnCountSlider)
	sideMenuContainerLeft.AddChild(hiveSpawnCooldown)

	menuButtonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(7),
				widget.GridLayoutOpts.Spacing(20, 20),
				widget.GridLayoutOpts.Padding(
					widget.Insets{
						Top: 400,
					},
				),
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

	camBtn := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonBg),

		widget.ButtonOpts.Text("  None  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setBrushType(ecs, comp.BrushNone)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
				},
			),
		),
	)

	drawPheromoneBtn := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonBg),

		widget.ButtonOpts.Text("  Draw Pheromone  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setBrushType(ecs, comp.BrushPheromone)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				},
			),
		),
	)

	drawWallBtn := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonBg),

		widget.ButtonOpts.Text("  Draw Wall  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setBrushType(ecs, comp.BrushWall)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				},
			),
		),
	)

	drawFoodBtn := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonBg),

		widget.ButtonOpts.Text("  Draw Food  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setBrushType(ecs, comp.BrushFood)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
				},
			),
		),
	)

	eraserBtn := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonBg),

		widget.ButtonOpts.Text("  Eraser  ", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 24), &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Disabled: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setBrushType(ecs, comp.BrushEraser)
		}),

		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
				},
			),
		),
	)

	widget.NewRadioGroup(widget.RadioGroupOpts.Elements(
		camBtn, drawPheromoneBtn, drawWallBtn, drawFoodBtn, eraserBtn,
	))

	menuButtonContainer.AddChild(camBtn)
	menuButtonContainer.AddChild(drawPheromoneBtn)
	menuButtonContainer.AddChild(drawWallBtn)
	menuButtonContainer.AddChild(drawFoodBtn)
	menuButtonContainer.AddChild(eraserBtn)

	rootContainer.AddChild(sideMenuContainer)
	rootContainer.AddChild(menuButtonContainer)

	/* lowerMenuBackground := image.NewNineSliceColor(color.RGBA{255, 255, 255, 255})

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

	ui.Container.AddChild(lowerMenuContainer) */

	return &menu.Menu{
		SceneManager: sceneManager,
		UI:           &ui,
		ElmTexts:     elements,
	}
}

// idle rgba(0, 171, 52, 0.8)
// hover rgba(0, 171, 52, 1)
// pressed rgba(0, 171, 52, 1)
func loadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.RGBA{R: 0, G: 110, B: 27, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 0, G: 130, B: 20, A: 255})
	pressed := image.NewNineSliceColor(color.RGBA{R: 10, G: 169, B: 50, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func loadTrackImage() *widget.SliderTrackImage {
	idle := image.NewNineSliceColor(color.RGBA{R: 27, G: 27, B: 27, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})

	return &widget.SliderTrackImage{
		Idle:  idle,
		Hover: hover,
	}
}

func setBrushType(ecs *ecs.ECS, brushType comp.BrushType) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	globalStateComp.CurrentBrush = brushType
}

func setCellSpeed(ecs *ecs.ECS, speed float64) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	globalStateComp.CellSpeed = int(speed)
	globalStateComp.IsValueDirty = true
}

func setHiveSpawnCount(ecs *ecs.ECS, val float64) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	globalStateComp.SpawnCount = int(val)
	globalStateComp.IsValueDirty = true
}

func setHiveCooldown(ecs *ecs.ECS, val float64) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	globalStateComp.SpawnCooldown = val
	globalStateComp.IsValueDirty = true
}

/* func setBrushRadius(ecs *ecs.ECS, radius int) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	globalStateComp.BrushRadius = radius
} */

func createFuckingSlider(label string, min int, max int, currentVal int, divisor float64, callback func(value float64)) *widget.Container {
	buttonBg := loadButtonImage()

	sliderContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Spacing(10, 5),
			),
		),
	)

	sliderLabel := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			StretchHorizontal:  true,
		})),
		widget.TextOpts.Text(label, *assets.AssetsInstance.GetFont("FallingSkyCondensed", 20), color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{
			Left: 10,
		}),
	)

	sliderValue := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			StretchHorizontal:  true,
		})),
		widget.TextOpts.Text("1", *assets.AssetsInstance.GetFont("FallingSkyCondensed", 20), color.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
	)

	trackImage := loadTrackImage()
	slider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		}), widget.WidgetOpts.MinSize(150, 4)),
		widget.SliderOpts.MinMax(min, max),
		widget.SliderOpts.Images(trackImage, buttonBg),
		widget.SliderOpts.TrackOffset(0),
		widget.SliderOpts.FixedHandleSize(10),
		widget.SliderOpts.PageSizeFunc(func() int {
			return 1
		}),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			if divisor == 1 {
				sliderValue.Label = fmt.Sprintf("%d", args.Slider.Current)
				callback(float64(args.Slider.Current))
			} else {
				sliderValue.Label = fmt.Sprintf("%.2f", float64(args.Slider.Current)/divisor)
				callback(float64(args.Slider.Current) / divisor)
			}
		}),
	)
	slider.Current = currentVal
	sliderContainer.AddChild(sliderLabel)
	sliderContainer.AddChild(slider)
	sliderContainer.AddChild(sliderValue)

	return sliderContainer
}
