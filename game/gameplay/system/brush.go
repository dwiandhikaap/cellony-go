package system

import (
	"autocell/game/assets"
	"autocell/game/config"
	camera "autocell/game/gameplay/camera"
	comp "autocell/game/gameplay/component"
	ent "autocell/game/gameplay/entity"
	"autocell/game/util"
	bitmask "autocell/lib/bit"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	camlib "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var lastNodeX = -1
var lastNodeY = -1

func BrushSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	brushType := globalStateComp.CurrentBrush

	mouseCamX, mouseCamY := ebiten.CursorPosition()

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		_, scrollAmount := ebiten.Wheel()
		if scrollAmount > 0 {
			globalStateComp.BrushRadius = int(math.Min(100, float64(globalStateComp.BrushRadius+1)))
		} else if scrollAmount < 0 {
			globalStateComp.BrushRadius = int(math.Max(1, float64(globalStateComp.BrushRadius-1)))
		}
	}

	brushRadius := float64(globalStateComp.BrushRadius)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := camera.CursorWorldPosition()

		cx := int(cursorX)
		cy := int(cursorY)

		isFirstNode := lastNodeX == -1 || lastNodeY == -1
		isFarEnough := util.DistanceSquared(lastNodeX, lastNodeY, cx, cy) > 100

		if !isFirstNode && !isFarEnough && brushType != comp.BrushPheromone {
			return
		}

		mapQuery := donburi.NewQuery(
			filter.Contains(comp.Grid),
		)

		if brushType != comp.BrushPheromone {
			mapQuery.Each(ecs.World, func(entry *donburi.Entry) {
				grid := comp.Grid.Get(entry).Grid
				mask := comp.Grid.Get(entry).Mask

				// Outer circle, reduce by 0.1 each steps
				for i := 0; i < 15; i++ {
					r := brushRadius * (1.5 + float64(i)*0.1)
					indices := util.GetCircleLatticeArea(cursorX/config.Game.TileSize, cursorY/config.Game.TileSize, r/config.Game.TileSize)
					delta := 1 / 90.0
					for _, index := range indices {
						xIndex := int(index[0])
						yIndex := int(index[1])

						if xIndex < 0 || xIndex >= len(grid) || yIndex < 0 || yIndex >= len(grid[0]) {
							continue
						}

						if bitmask.HasBit(mask[xIndex][yIndex], comp.MarkedMask) {
							continue
						}

						// if cursor is above hud, do not draw
						hudX1 := 280
						hudX2 := 960
						hudY1 := 640
						hudY2 := 800

						hudX3 := 0
						hudX4 := 350
						hudY3 := 50
						hudY4 := 220

						if (mouseCamX) > hudX1 && (mouseCamX) < hudX2 && (mouseCamY) > hudY1 && (mouseCamY) < hudY2 {
							continue
						}

						if (mouseCamX) > hudX3 && (mouseCamX) < hudX4 && (mouseCamY) > hudY3 && (mouseCamY) < hudY4 {
							continue
						}

						if brushType == comp.BrushWall {
							playBrushSfx()
							mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.WallMask)
							grid[xIndex][yIndex] = float32(util.Clamp(float64(grid[xIndex][yIndex])+delta, 0.0, 1.0))
							mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.DirtyMask)
						} else if brushType == comp.BrushFood {
							playBrushSfx()
							mask[xIndex][yIndex] = bitmask.ClearBit(mask[xIndex][yIndex], comp.WallMask)
							grid[xIndex][yIndex] = float32(util.Clamp(float64(grid[xIndex][yIndex])+delta, 0.0, 1.0))
							mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.DirtyMask)
						} else if brushType == comp.BrushEraser {
							playBrushSfx()
							grid[xIndex][yIndex] = float32(util.Clamp(float64(grid[xIndex][yIndex])-delta, 0.0, 1.0))
							mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.DirtyMask)
						}
					}
				}
			})
		} else {
			cellQuery := donburi.NewQuery(
				filter.Contains(comp.Cell),
			)

			cellEntry, ok := cellQuery.First(ecs.World)
			if !ok {
				print("no cell")
				return
			}

			ent.CreatePheromoneEntity(ecs.World, &ent.CreatePheromoneOptions{
				X:         float64(cx),
				Y:         float64(cy),
				Activity:  comp.Wandering,
				HiveID:    comp.Parent.Get(cellEntry).Id,
				Intensity: 1,
			})
		}
		lastNodeX, lastNodeY = cx, cy
		return
	}

	lastNodeX, lastNodeY = -1, -1
}

func BrushRenderer(ecs *ecs.ECS, cam *camlib.Camera) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateComp := comp.GlobalState.Get(globalState)
	brushType := globalStateComp.CurrentBrush

	if brushType == comp.BrushNone || brushType == comp.BrushPheromone {
		return
	}

	brushRadius := globalStateComp.BrushRadius * 2

	// draw circle ring following cursor
	cursorX, cursorY := ebiten.CursorPosition()
	imgSize := int(brushRadius * 2)

	baseImage := *assets.AssetsInstance.Textures["circle8"]

	ebitenImg := ebiten.NewImageFromImage(baseImage)
	resizedEbitenImg := util.ResizeImage(ebitenImg, imgSize, imgSize)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.SetR(0.3)
	op.ColorScale.SetG(0.3)
	op.ColorScale.SetB(0.3)
	op.ColorScale.SetA(0.1)
	op.GeoM.Translate(float64(cursorX)/cam.Scale-(float64(brushRadius)), float64(cursorY)/cam.Scale-(float64(brushRadius)))

	cam.Surface.DrawImage(resizedEbitenImg, op)
}

var lastSeekTime int64 = 0

func playBrushSfx() {
	now := time.Now().UnixMilli()

	if now-lastSeekTime < 200 {
		return
	}

	lastSeekTime = now

	err := assets.AssetsInstance.Audio["brush"].Seek(time.Millisecond * 300)

	if err != nil {
		log.Println(err.Error())
		return
	}

	assets.AssetsInstance.Audio["brush"].Play()
}
