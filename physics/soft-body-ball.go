package physics

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type SoftBodyBall struct {
	CenterBody            *box2d.B2Body
	EdgeBodies            []*box2d.B2Body
	edgeBodyRenderOffsets []*box2d.B2Vec2

	color rl.Color
}

func NewSoftBodyBall(world *box2d.B2World, center box2d.B2Vec2, radius float64, components int) *SoftBodyBall {
	ballRadius := radius / 8
	ball := &SoftBodyBall{color: rl.Red}

	centerBallBodyDef := box2d.NewB2BodyDef()
	centerBallBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	centerBallBodyDef.Position = center

	ball.CenterBody = world.CreateBody(centerBallBodyDef)

	centerBallShape := box2d.NewB2CircleShape()
	centerBallShape.SetRadius(ballRadius)

	centerBallFixture := ball.CenterBody.CreateFixture(centerBallShape, 1)
	centerBallFixture.SetRestitution(0.75)
	centerBallFixture.SetUserData(ball)

	ball.EdgeBodies = make([]*box2d.B2Body, 0, components)
	ball.edgeBodyRenderOffsets = make([]*box2d.B2Vec2, 0, components)

	// Make ball bodies and join to center
	for i := 0; i < components; i++ {
		angleRad := (math.Pi * 2) * (float64(i) / float64(components))

		edgeBallBodyDef := box2d.NewB2BodyDef()
		edgeBallBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
		edgeBallBodyDef.Position = box2d.MakeB2Vec2(center.X+(math.Sin(angleRad)*radius), center.Y+(math.Cos(angleRad)*radius))

		edgeBallBody := world.CreateBody(edgeBallBodyDef)

		edgeBallShape := box2d.NewB2CircleShape()
		edgeBallShape.SetRadius(ballRadius)

		edgeBallFixture := edgeBallBody.CreateFixture(edgeBallShape, 1)
		edgeBallFixture.SetRestitution(0.75)
		edgeBallFixture.SetFriction(0.8)

		ball.EdgeBodies = append(ball.EdgeBodies, edgeBallBody)
		ball.edgeBodyRenderOffsets = append(ball.edgeBodyRenderOffsets, box2d.NewB2Vec2(-math.Sin(angleRad)*ballRadius, -math.Cos(angleRad)*ballRadius))

		edgeCenterJointDef := box2d.MakeB2DistanceJointDef()
		edgeCenterJointDef.BodyA = ball.CenterBody
		edgeCenterJointDef.BodyB = edgeBallBody
		edgeCenterJointDef.CollideConnected = true
		edgeCenterJointDef.FrequencyHz = 4
		edgeCenterJointDef.DampingRatio = 0.5
		edgeCenterJointDef.Length = radius
		world.CreateJoint(&edgeCenterJointDef)
	}

	for i := range ball.EdgeBodies {
		e1 := ball.EdgeBodies[i]
		e2 := ball.EdgeBodies[(i+1)%components]

		deltaPos := box2d.MakeB2Vec2(e1.GetPosition().X-e2.GetPosition().X, e1.GetPosition().Y-e2.GetPosition().Y)

		edgeToEdgeJoint := box2d.MakeB2DistanceJointDef()
		edgeToEdgeJoint.BodyA = e1
		edgeToEdgeJoint.BodyB = e2
		edgeToEdgeJoint.CollideConnected = true
		edgeToEdgeJoint.FrequencyHz = 4
		edgeToEdgeJoint.DampingRatio = 0.5
		edgeToEdgeJoint.Length = deltaPos.Length()
		world.CreateJoint(&edgeToEdgeJoint)
	}

	return ball
}

func (s *SoftBodyBall) Draw(dt float32, body *box2d.B2Body, fixture *box2d.B2Fixture) {
	centerPos := rl.NewVector2(float32(s.CenterBody.GetPosition().X), float32(s.CenterBody.GetPosition().Y))

	for i := range s.EdgeBodies {
		e1 := s.EdgeBodies[i]
		e1Offset := s.edgeBodyRenderOffsets[i]

		e2 := s.EdgeBodies[(i+1)%len(s.EdgeBodies)]
		e2Offset := s.edgeBodyRenderOffsets[(i+1)%len(s.EdgeBodies)]

		e1Pos := rl.NewVector2(float32(e1.GetPosition().X+e1Offset.X), float32(e1.GetPosition().Y+e1Offset.Y))
		e2Pos := rl.NewVector2(float32(e2.GetPosition().X+e2Offset.X), float32(e2.GetPosition().Y+e2Offset.Y))

		rl.DrawTriangle(e1Pos, e2Pos, centerPos, s.color)
	}
}
