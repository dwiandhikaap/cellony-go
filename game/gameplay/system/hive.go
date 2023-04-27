package system

import (
	"cellony/game/config"
	"cellony/game/gameplay/comp"
	ent "cellony/game/gameplay/entity"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
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
		cx := comp.Position.Get(entry).X
		cy := comp.Position.Get(entry).Y

		hiveColor := comp.Color.Get(entry)
		cellColor := color.RGBA{}
		cellColor.R = hiveColor.R
		cellColor.G = hiveColor.G
		cellColor.B = hiveColor.B

		op := &ent.CreateCellOptions{
			X:      cx,
			Y:      cy,
			Speed:  50,
			Color:  cellColor,
			HiveID: entry.Entity(),
		}
		ent.CreateCellEntity(ecs.World, op)
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
