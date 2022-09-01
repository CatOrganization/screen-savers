package visualizer

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type FightClubBoxDrawable struct {
	shape *box2d.B2PolygonShape
	color rl.Color
}

func (f *FightClubBoxDrawable) Draw(_ float32, body *box2d.B2Body, fixture *box2d.B2Fixture) {
	DebugDrawPolygonShape(body, f.shape, f.color, PhysicsToRenderScale)
}

type FightClubFloorDrawable struct {
	points []rl.Vector2

	floorHeight float32
}

func (f *FightClubFloorDrawable) Draw(dt float32, body *box2d.B2Body, fixture *box2d.B2Fixture) {
	for i := 1; i < len(f.points); i++ {
		p1 := f.points[i-1]
		p1Floor := rl.NewVector2(p1.X, f.floorHeight)

		p2 := f.points[i]
		p2Floor := rl.NewVector2(p2.X, f.floorHeight)

		rl.DrawLineV(p1, p2, rl.DarkGray)

		rl.DrawTriangle(p2Floor, p2, p1, lightenColor(rl.DarkGray))
		rl.DrawTriangle(p1Floor, p2Floor, p1, lightenColor(rl.DarkGray))
	}
}
