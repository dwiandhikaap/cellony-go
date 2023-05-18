package ent

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"cellony/game/assets"
	"cellony/game/config"
	comp "cellony/game/gameplay/component"
	"cellony/game/graphics"
	"cellony/game/util"

	noise "github.com/ojrac/opensimplex-go"
)

type CreateCellOptions struct {
	X        float64
	Y        float64
	Speed    float64
	Color    color.Color
	HiveID   donburi.Entity
	Activity comp.Activity

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
	assetsKey := fmt.Sprintf("circle%d", options.HiveID)
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

		cellImage.DrawImage(assets.AssetsInstance.Sprites["circle64"], tintOp)

		assets.AssetsInstance.Sprites[assetsKey] = cellImage
	}

	comp.Sprite.Get(cellEntry).Sprite = cellImage
	comp.Sprite.Get(cellEntry).Z = 1
	comp.Sprite.Get(cellEntry).Scale = 1
	comp.Sprite.Get(cellEntry).Opacity = 1

	comp.Cell.Get(cellEntry).PheromoneChance = options.PheromoneChance

	return cell
}

func CreateHiveEntity(world donburi.World) donburi.Entity {
	hive := world.Create(comp.Position, comp.Vertices, comp.Indices, comp.Color, comp.Hive)
	hiveEntry := world.Entry(hive)

	radius := 64.0

	comp.Hive.Get(hiveEntry).SpawnCooldown = 1
	comp.Hive.Get(hiveEntry).SpawnCountdown = 0
	comp.Hive.Get(hiveEntry).SpawnCount = 30

	x := rand.Float64() * float64(config.Game.Width)
	y := rand.Float64() * float64(config.Game.Height)

	// padding
	xPadding := 0.2 * config.Game.Width
	yPadding := 0.2 * config.Game.Height

	x = util.RangeInterpolate(x, 0.0, config.Game.Width, xPadding, float64(config.Game.Width)-xPadding)
	y = util.RangeInterpolate(y, 0.0, config.Game.Height, yPadding, float64(config.Game.Height)-yPadding)

	comp.Position.Get(hiveEntry).X = x
	comp.Position.Get(hiveEntry).Y = y

	color := graphics.GenerateHiveColor()
	vs, is := graphics.GeneratePolygonVertices(float32(x), float32(y), color, radius, 16, 0.0)

	comp.Vertices.Get(hiveEntry).Vertices = vs
	comp.Indices.Get(hiveEntry).Indices = is

	r, g, b, _ := color.RGBA()
	comp.Color.Get(hiveEntry).R = uint8(r >> 8)
	comp.Color.Get(hiveEntry).G = uint8(g >> 8)
	comp.Color.Get(hiveEntry).B = uint8(b >> 8)

	// adjust map near hive
	mapQuery := donburi.NewQuery(
		filter.Contains(comp.Grid),
	)

	mapQuery.Each(world, func(entry *donburi.Entry) {
		grid := comp.Grid.Get(entry).Grid
		dirtyMask := comp.Grid.Get(entry).DirtyMask

		// Outer circle, reduce by 0.1 each steps
		for i := 0; i < 15; i++ {
			r := radius * (1.5 + float64(i)*0.1)
			indices := util.GetCircleLatticeArea(x/config.Game.TileSize, y/config.Game.TileSize, r/config.Game.TileSize)
			delta := 1 / 15.0
			for _, index := range indices {
				xIndex := int(index[0])
				yIndex := int(index[1])

				grid[xIndex][yIndex] = float32(util.Clamp(float64(grid[xIndex][yIndex])-delta, 0.0, 1.0))
				dirtyMask[xIndex][yIndex] = true
			}
		}
	})

	return hive
}

func CreateMapEntity(world donburi.World) {
	mapEntity := world.Create(comp.Grid, comp.Image)
	mapEntry := world.Entry(mapEntity)

	tileSize := config.Game.TileSize

	mapWidth := int(config.Game.Width / tileSize)
	mapHeight := int(config.Game.Height / tileSize)

	grid := make([][]float32, mapWidth)
	dirtyMask := make([][]bool, mapWidth)
	typeMask := make([][]bool, mapWidth)

	comp.Grid.Get(mapEntry).Grid = grid
	comp.Grid.Get(mapEntry).DirtyMask = dirtyMask
	comp.Grid.Get(mapEntry).TypeMask = typeMask
	comp.Image.Get(mapEntry).Img = ebiten.NewImage(int(config.Game.Width), int(config.Game.Height))

	terrainNoise := noise.NewNormalized(1)
	foodNoise := noise.NewNormalized(2)

	for i := 0; i < mapWidth; i++ {
		grid[i] = make([]float32, mapHeight)
		dirtyMask[i] = make([]bool, mapHeight)
		typeMask[i] = make([]bool, mapHeight)
		for j := 0; j < mapHeight; j++ {
			val := float32(terrainNoise.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > 0.45 {
				grid[i][j] = float32(util.RangeInterpolate(float64(val), 0.45, 1.0, 0.0, 1.0))
			} else {
				grid[i][j] = 0.0
			}
			dirtyMask[i][j] = true
			typeMask[i][j] = true
		}
	}

	for i := 0; i < mapWidth; i++ {
		for j := 0; j < mapHeight; j++ {
			val := float32(foodNoise.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > 0.8 {
				gridValue := float32(util.RangeInterpolate(float64(val), 0.5, 1.0, 0.0, 1.0))
				grid[i][j] = (grid[i][j]*0.25 + gridValue*1.75) / 2.0

				dirtyMask[i][j] = true
				typeMask[i][j] = false
			}
		}
	}
}

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
