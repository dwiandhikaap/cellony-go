package system

import (
	comp "cellony/game/gameplay/component"

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
