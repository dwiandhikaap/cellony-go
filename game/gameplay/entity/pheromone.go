package ent

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"cellony/game/assets"
	comp "cellony/game/gameplay/component"
)

type CreatePheromoneOptions struct {
	X         float64
	Y         float64
	HiveID    donburi.Entity
	Intensity float64
	Activity  comp.Activity
}

func CreatePheromoneEntity(world donburi.World, options *CreatePheromoneOptions) donburi.Entity {
	pheromone := world.Create(comp.Position, comp.Pheromone, comp.Sprite)
	pheromoneEntry := world.Entry(pheromone)

	comp.Position.Get(pheromoneEntry).X = options.X
	comp.Position.Get(pheromoneEntry).Y = options.Y
	comp.Pheromone.Get(pheromoneEntry).HiveID = options.HiveID
	comp.Pheromone.Get(pheromoneEntry).Activity = options.Activity
	comp.Pheromone.Get(pheromoneEntry).Intensity = options.Intensity
	comp.Pheromone.Get(pheromoneEntry).MaxIntensity = options.Intensity
	comp.Pheromone.Get(pheromoneEntry).Activity = options.Activity

	assetsKey := fmt.Sprintf("phero-%d-%d", options.HiveID, options.Activity)
	pheroImage := assets.AssetsInstance.Sprites[assetsKey]
	if assets.AssetsInstance.Sprites[assetsKey] == nil {
		pheroImage = ebiten.NewImage(16, 16)
		pheroImage.Clear()

		hiveQuery := donburi.NewQuery(
			filter.Contains(comp.Hive),
		)

		var colorData comp.ColorData

		hiveQuery.Each(world, func(entry *donburi.Entry) {
			if entry.Entity() != options.HiveID {
				return
			}
			colorData = *comp.Color.Get(entry)
		})

		tintOp := &ebiten.DrawImageOptions{}
		tintOp.ColorScale.SetR(float32(colorData.R) / 255)
		tintOp.ColorScale.SetG(float32(colorData.G) / 255)
		tintOp.ColorScale.SetB(float32(colorData.B) / 255)
		tintOp.ColorScale.SetA(1)

		pheroImage.DrawImage(assets.AssetsInstance.Sprites["phero"], tintOp)

		assets.AssetsInstance.Sprites[assetsKey] = pheroImage
	}

	comp.Sprite.Get(pheromoneEntry).Sprite = pheroImage
	comp.Sprite.Get(pheromoneEntry).Z = 0
	comp.Sprite.Get(pheromoneEntry).Scale = 0.5
	comp.Sprite.Get(pheromoneEntry).Opacity = 0.5

	return pheromone
}
