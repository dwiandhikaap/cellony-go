package ent

import (
	"cellony/game/assets"
	comp "cellony/game/gameplay/component"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type CreatePathNodeOptions struct {
	X      float64
	Y      float64
	Radius float64
}

func CreatePathNode(world donburi.World, opt *CreatePathNodeOptions) donburi.Entity {
	pathNode := world.Create(comp.PathNode, comp.Position, comp.Sprite)
	pathNodeEntry := world.Entry(pathNode)

	comp.Position.Get(pathNodeEntry).X = opt.X
	comp.Position.Get(pathNodeEntry).Y = opt.Y

	rawSpriteKey := "circle64"
	assetsKey := fmt.Sprintf("pathnode%d", int(opt.Radius))

	if assets.AssetsInstance.Sprites[assetsKey] == nil {
		nodeImage := ebiten.NewImage(64, 64)
		nodeImage.Clear()

		tintOp := &ebiten.DrawImageOptions{}
		tintOp.ColorScale.SetR(0.0)
		tintOp.ColorScale.SetG(0.8)
		tintOp.ColorScale.SetB(0.5)
		tintOp.ColorScale.SetA(0)

		nodeImage.DrawImage(assets.AssetsInstance.Sprites[rawSpriteKey], tintOp)

		assets.AssetsInstance.Sprites[assetsKey] = nodeImage
	}

	println(assets.AssetsInstance.Sprites[assetsKey])

	comp.Sprite.Get(pathNodeEntry).Sprite = assets.AssetsInstance.Sprites[assetsKey]
	comp.Sprite.Get(pathNodeEntry).Scale = opt.Radius / 32
	comp.Sprite.Get(pathNodeEntry).Z = 1
	comp.Sprite.Get(pathNodeEntry).Opacity = 0.7

	return pathNode
}
