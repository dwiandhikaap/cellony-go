package ent

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"cellony/game/assets"
	comp "cellony/game/gameplay/component"
)

type CreateCellOptions struct {
	X        float64
	Y        float64
	Speed    float64
	Color    color.Color
	HiveID   donburi.Entity
	Activity comp.Activity

	Health float64
	Class  comp.CellClass

	PheromoneCooldown float64
	PheromoneChance   float64
}

func CreateCellEntity(world donburi.World, options *CreateCellOptions) donburi.Entity {
	cell := world.Create(comp.Cell, comp.Position, comp.Velocity, comp.Speed, comp.Parent, comp.CellActivity, comp.Sprite)
	cellEntry := world.Entry(cell)

	comp.Position.Get(cellEntry).X = options.X
	comp.Position.Get(cellEntry).Y = options.Y

	comp.Speed.Get(cellEntry).Speed = options.Speed

	comp.CellActivity.Get(cellEntry).Activity = options.Activity

	comp.Parent.Get(cellEntry).Id = options.HiveID

	angle := rand.Float64() * 2 * 3.14159
	comp.Velocity.Get(cellEntry).X = math.Cos(angle) * comp.Speed.Get(cellEntry).Speed
	comp.Velocity.Get(cellEntry).Y = math.Sin(angle) * comp.Speed.Get(cellEntry).Speed

	// assets key = "circle{hiveID}"
	assetsKey := "circle"
	if options.Class == comp.Gatherer {
		assetsKey = "square"
	} else if options.Class == comp.Soldier {
		assetsKey = "triangle"
	}

	rawSpriteKey := assetsKey + "64"
	assetsKey = fmt.Sprintf("%s%d", assetsKey, options.HiveID)

	cellImage := assets.AssetsInstance.Sprites[assetsKey]
	if assets.AssetsInstance.Sprites[assetsKey] == nil {
		cellImage = ebiten.NewImage(8, 8)
		cellImage.Clear()

		r, g, b, _ := options.Color.RGBA()

		tintOp := &ebiten.DrawImageOptions{}
		tintOp.ColorScale.SetR(float32(r) / 65535)
		tintOp.ColorScale.SetG(float32(g) / 65535)
		tintOp.ColorScale.SetB(float32(b) / 65535)
		tintOp.ColorScale.SetA(1)

		cellImage.DrawImage(assets.AssetsInstance.Sprites[rawSpriteKey], tintOp)

		assets.AssetsInstance.Sprites[assetsKey] = cellImage
	}

	comp.Sprite.Get(cellEntry).Sprite = cellImage
	comp.Sprite.Get(cellEntry).Z = 1
	comp.Sprite.Get(cellEntry).Scale = 1
	comp.Sprite.Get(cellEntry).Opacity = 1

	comp.Cell.Get(cellEntry).Health = options.Health
	comp.Cell.Get(cellEntry).Class = options.Class
	comp.Cell.Get(cellEntry).PheromoneChance = options.PheromoneChance

	return cell
}
