package visualizer

import (
	"fmt"
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

func DebugDrawWorld(world *box2d.B2World) {
	for body := world.GetBodyList(); body != nil; body = body.GetNext() {
		DebugDrawBody(body)
	}

	for joint := world.GetJointList(); joint != nil; joint = joint.GetNext() {
		DebugDrawJoint(joint)
	}
}

func DebugDrawJoint(joint box2d.B2JointInterface) {
	var anchorA, anchorB box2d.B2Vec2

	switch joint.(type) {
	case *box2d.B2DistanceJoint:
		distanceJoint := joint.(*box2d.B2DistanceJoint)

		anchorA = distanceJoint.M_bodyA.GetWorldPoint(distanceJoint.M_localAnchorA)
		anchorB = distanceJoint.M_bodyB.GetWorldPoint(distanceJoint.M_localAnchorB)
	case *box2d.B2MouseJoint:
		mouseJoint := joint.(*box2d.B2MouseJoint)

		anchorA = mouseJoint.M_targetA
		anchorB = mouseJoint.M_bodyB.GetWorldPoint(mouseJoint.M_localAnchorB)
	case *box2d.B2RevoluteJoint:
		revoluteJoint := joint.(*box2d.B2RevoluteJoint)

		anchorA = revoluteJoint.M_bodyA.GetWorldPoint(revoluteJoint.M_localAnchorA)
		anchorB = revoluteJoint.M_bodyB.GetWorldPoint(revoluteJoint.M_localAnchorB)
	// TODO: more joints
	default:
		anchorA = joint.GetBodyA().GetPosition()
		anchorB = joint.GetBodyB().GetPosition()
	}

	a := rl.NewVector2(float32(anchorA.X), float32(anchorA.Y))
	b := rl.NewVector2(float32(anchorB.X), float32(anchorB.Y))

	rl.DrawLineV(a, b, rl.SkyBlue)
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
		DebugDrawPolygonShape(body, shape.(*box2d.B2PolygonShape), colorForBody(body))
	case *box2d.B2EdgeShape:
		DebugDrawEdgeShape(body, shape.(*box2d.B2EdgeShape))
	default:
		fmt.Printf("unknown shape: %T", shape)
	}
}

func DebugDrawCircleShape(body *box2d.B2Body, circle *box2d.B2CircleShape) {
	worldCenter := body.GetWorldPoint(circle.M_p)
	color := colorForBody(body)

	rl.DrawCircle(int32(worldCenter.X), int32(worldCenter.Y), float32(circle.GetRadius()), lightenColor(color))
	rl.DrawCircleLines(int32(worldCenter.X), int32(worldCenter.Y), float32(circle.GetRadius()), color)
}

func DebugDrawPolygonShape(body *box2d.B2Body, polygon *box2d.B2PolygonShape, color rl.Color) {
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

		rl.DrawTriangle(centroid, v2, v1, lightenColor(color))
		rl.DrawLineV(v1, v2, color)
	}
}

func DebugDrawEdgeShape(body *box2d.B2Body, edge *box2d.B2EdgeShape) {
	// TODO: handle v0 and v3?
	worldV1 := body.GetWorldPoint(edge.M_vertex1)
	worldV2 := body.GetWorldPoint(edge.M_vertex2)

	rl.DrawLine(int32(worldV1.X), int32(worldV1.Y), int32(worldV2.X), int32(worldV2.Y), colorForBody(body))
}

func colorForBody(body *box2d.B2Body) color.RGBA {
	switch body.GetType() {
	case box2d.B2BodyType.B2_staticBody:
		return rl.DarkGreen
	case box2d.B2BodyType.B2_kinematicBody:
		return rl.DarkBlue
	case box2d.B2BodyType.B2_dynamicBody:
		if body.M_flags&box2d.B2Body_Flags.E_awakeFlag == box2d.B2Body_Flags.E_awakeFlag {
			return rl.Pink
		} else {
			return rl.DarkGray
		}
	}

	return rl.Red
}

func lightenColor(c color.RGBA) color.RGBA {
	return rl.NewColor(c.R, c.G, c.B, c.A/5)
}
