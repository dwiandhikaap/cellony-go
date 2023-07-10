package system

import (
	comp "autocell/game/gameplay/component"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func HUDSystem(ecs *ecs.ECS) {
	menuQuery := donburi.NewQuery(filter.Contains(comp.HUD))

	firstEntry, ok := menuQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}

	hud := comp.HUD.Get(firstEntry)

	hud.Menu.ElmTexts["cellCountText"].Label = fmt.Sprintf("Cell Count: %d", getCellCounts(ecs))
	hud.Menu.ElmTexts["hiveResourceText"].Label = fmt.Sprintf("Hive Resource: %d", getHiveResource(ecs))

	hud.Menu.UI.Update()
}

func HUDRenderer(ecs *ecs.ECS, screen *ebiten.Image) {
	menuQuery := donburi.NewQuery(filter.Contains(comp.HUD))

	firstEntry, ok := menuQuery.FirstEntity(ecs.World)
	if !ok {
		return
	}

	hud := comp.HUD.Get(firstEntry)
	hud.Menu.Draw(screen)
}

func getCellCounts(ecs *ecs.ECS) int {
	cellQuery := donburi.NewQuery(filter.Contains(comp.Cell))
	return cellQuery.Count(ecs.World)
}

func getHiveResource(ecs *ecs.ECS) int {
	hiveQuery := donburi.NewQuery(filter.Contains(comp.Hive))
	firstEntry, ok := hiveQuery.FirstEntity(ecs.World)
	if !ok {
		return 0
	}

	return int(comp.Hive.Get(firstEntry).Resource)
}
