package ent

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	noise "github.com/ojrac/opensimplex-go"
	"github.com/yohamta/donburi"

	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	"autocell/game/util"
	bitmask "autocell/lib/bit"
)

func CreateMapEntity(world donburi.World) {
	mapEntity := world.Create(comp.Grid, comp.Image)
	mapEntry := world.Entry(mapEntity)

	tileSize := config.Game.TileSize

	mapWidth := int(config.Game.Width / tileSize)
	mapHeight := int(config.Game.Height / tileSize)

	grid := make([][]float32, mapWidth)
	mask := make([][]uint8, mapWidth)

	comp.Grid.Get(mapEntry).Grid = grid
	comp.Grid.Get(mapEntry).Mask = mask
	comp.Image.Get(mapEntry).Img = ebiten.NewImage(int(config.Game.Width), int(config.Game.Height))

	seed1 := rand.Int63()
	seed2 := rand.Int63()

	terrainNoise := noise.NewNormalized(seed1)
	foodNoise := noise.NewNormalized(seed2)

	var thresh float32 = float32(config.Game.WallThreshold)     // was 0.45
	var foodThresh float32 = float32(config.Game.FoodThreshold) // was 0.8

	for i := 0; i < mapWidth; i++ {
		grid[i] = make([]float32, mapHeight)
		mask[i] = make([]uint8, mapHeight)
		for j := 0; j < mapHeight; j++ {
			val := float32(terrainNoise.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > thresh {
				grid[i][j] = float32(util.RangeInterpolate(float64(val), float64(thresh), 1.0, 0.0, 1.0))
			} else {
				grid[i][j] = 0.0
			}

			mask[i][j] = bitmask.SetBit(mask[i][j], (comp.DirtyMask))
			mask[i][j] = bitmask.SetBit(mask[i][j], (comp.WallMask))
		}
	}

	for i := 0; i < mapWidth; i++ {
		for j := 0; j < mapHeight; j++ {
			val := float32(foodNoise.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > foodThresh {
				gridValue := float32(util.RangeInterpolate(float64(val), 0.5, 1.0, 0.0, 1.0))
				grid[i][j] = (grid[i][j]*0.25 + gridValue*1.75) / 2.0

				mask[i][j] = bitmask.SetBit(mask[i][j], comp.DirtyMask)
				mask[i][j] = bitmask.ClearBit(mask[i][j], comp.WallMask)
			}
		}
	}
}

func IsTileDirty(bitmask uint8) bool {
	return (bitmask & uint8(comp.DirtyMask)) != 0
}

func IsTileWall(bitmask uint8) bool {
	return (bitmask & uint8(comp.WallMask)) != 0
}

func IsTileFood(bitmask uint8) bool {
	return (bitmask & uint8(comp.WallMask)) == 0
}

func IsTileMarked(bitmask uint8) bool {
	return (bitmask & uint8(comp.MarkedMask)) != 0
}
