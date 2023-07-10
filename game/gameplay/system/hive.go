package system

import (
	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	ent "autocell/game/gameplay/entity"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/mroth/weightedrand/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

func HiveSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(comp.Hive),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		hive := comp.Hive.Get(entry)

		if hive.SpawnCountdown > 0 {
			hive.SpawnCountdown -= 1.0 / 60
			return
		}

		hive.SpawnCountdown = hive.SpawnCooldown

		cx := comp.Position.Get(entry).X
		cy := comp.Position.Get(entry).Y

		hiveColor := comp.Color.Get(entry)
		cellColor := color.RGBA{}
		cellColor.R = hiveColor.R
		cellColor.G = hiveColor.G
		cellColor.B = hiveColor.B

		class, _ := weightedrand.NewChooser[comp.CellClass](
			weightedrand.NewChoice(comp.Wanderer, hive.WandererOdd),
			weightedrand.NewChoice(comp.Worker, hive.WorkerOdd),
			weightedrand.NewChoice(comp.Soldier, hive.SoldierOdd),
		)

		// get globalstate
		speed := 30
		globalState, ok := donburi.NewQuery(
			filter.Contains(comp.GlobalState),
		).First(ecs.World)
		if ok {
			speed = comp.GlobalState.Get(globalState).CellSpeed
		}

		for i := 0; i < hive.SpawnCount; i++ {
			op := &ent.CreateCellOptions{
				X:               cx,
				Y:               cy,
				Speed:           float64(speed),
				Color:           cellColor,
				HiveID:          entry.Entity(),
				Health:          300,
				Class:           class.Pick(),
				PheromoneChance: 1.0 / (60 * 3),
			}

			ent.CreateCellEntity(ecs.World, op)
		}
	})
}

func HiveRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Vertices),
			filter.Contains(comp.Indices),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		vertices := comp.Vertices.Get(entry).Vertices
		indices := comp.Indices.Get(entry).Indices
		screen := cam.Surface

		op := &ebiten.DrawTrianglesOptions{}

		translatedVertices := make([]ebiten.Vertex, len(vertices))
		for i, v := range vertices {
			translatedVertices[i] = ebiten.Vertex{
				DstX:   float32(float64(v.DstX) - cam.X + float64(config.Video.Width/2)/cam.Scale),
				DstY:   float32(float64(v.DstY) - cam.Y + float64(config.Video.Height/2)/cam.Scale),
				SrcX:   v.SrcX,
				SrcY:   v.SrcY,
				ColorR: v.ColorR,
				ColorG: v.ColorG,
				ColorB: v.ColorB,
				ColorA: v.ColorA,
			}
		}
		screen.DrawTriangles(translatedVertices, indices, whiteSubImage, op)
	})
}
