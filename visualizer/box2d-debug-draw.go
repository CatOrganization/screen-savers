package visualizer

import (
	"fmt"
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

var bodyTypeColorMapping = map[uint8]rl.Color{
	box2d.B2BodyType.B2_dynamicBody:   rl.Pink,
	box2d.B2BodyType.B2_staticBody:    rl.DarkGreen,
	box2d.B2BodyType.B2_kinematicBody: rl.DarkBlue,
}

func DebugDrawWorld(world *box2d.B2World) {
	for body := world.GetBodyList(); body != nil; body = body.GetNext() {
		DebugDrawBody(body)

		//center := body.GetWorldCenter()
		//rl.DrawCircle(int32(center.X), int32(center.Y), 2, rl.RayWhite)
	}
}

func DebugDrawBody(body *box2d.B2Body) {
	for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
		DebugDrawShape(body, fixture.GetShape())
	}
}

func DebugDrawShape(body *box2d.B2Body, shape box2d.B2ShapeInterface) {
	switch shape.(type) {
	case *box2d.B2CircleShape:
		DebugDrawCircleShape(body, shape.(*box2d.B2CircleShape))
	case *box2d.B2PolygonShape:
		DebugDrawPolygonShape(body, shape.(*box2d.B2PolygonShape))
	case *box2d.B2EdgeShape:
		DebugDrawEdgeShape(body, shape.(*box2d.B2EdgeShape))
	default:
		fmt.Printf("unknown shape: %T", shape)
	}
}

func DebugDrawCircleShape(body *box2d.B2Body, circle *box2d.B2CircleShape) {
	worldCenter := body.GetWorldPoint(circle.M_p)
	color := bodyTypeColorMapping[body.M_type]

	rl.DrawCircle(int32(worldCenter.X), int32(worldCenter.Y), float32(circle.GetRadius()), lightenColor(color))
	rl.DrawCircleLines(int32(worldCenter.X), int32(worldCenter.Y), float32(circle.GetRadius()), color)
}

func DebugDrawPolygonShape(body *box2d.B2Body, polygon *box2d.B2PolygonShape) {
	for i := 0; i < polygon.M_count; i++ {
		v1Index := (i - 1 + polygon.M_count) % polygon.M_count

		worldV1 := body.GetWorldPoint(polygon.M_vertices[v1Index])
		worldV2 := body.GetWorldPoint(polygon.M_vertices[i])
		worldCentroid := body.GetWorldPoint(polygon.M_centroid)

		v1 := rl.Vector2{
			X: float32(worldV1.X),
			Y: float32(worldV1.Y),
		}

		v2 := rl.Vector2{
			X: float32(worldV2.X),
			Y: float32(worldV2.Y),
		}

		centroid := rl.Vector2{
			X: float32(worldCentroid.X),
			Y: float32(worldCentroid.Y),
		}

		c := bodyTypeColorMapping[body.M_type]
		rl.DrawTriangle(centroid, v2, v1, lightenColor(c))
		rl.DrawLineV(v1, v2, c)

	}
}

func DebugDrawEdgeShape(body *box2d.B2Body, edge *box2d.B2EdgeShape) {
	// TODO: handle v0 and v3?
	worldV1 := body.GetWorldPoint(edge.M_vertex1)
	worldV2 := body.GetWorldPoint(edge.M_vertex2)

	rl.DrawLine(int32(worldV1.X), int32(worldV1.Y), int32(worldV2.X), int32(worldV2.Y), bodyTypeColorMapping[body.M_type])
}

func lightenColor(c color.RGBA) color.RGBA {
	return rl.NewColor(c.R, c.G, c.B, c.A/5)
}
