package system

import (
	comp "cellony/game/gameplay/component"
	"math"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func PathNodeSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(comp.PathNode),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		// set opacity to sin of time
		sprite := comp.Sprite.Get(entry)
		sprite.Opacity = 0.3 + 0.3*math.Abs(math.Sin(float64(time.Now().UnixMilli())/500.0))
	})
}
