package physics

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"math"
)

type SoftBodyBall struct {
	CenterBody *box2d.B2Body
	EdgeBodies []*box2d.B2Body
	edgeLength float64

	Color rl.Color
}

type SoftBodyBallOptions struct {
	Density float64

	SpringFrequencyHz float64
	SpringDampingHz   float64

	Restitution float64
	Friction    float64
}

func NewSoftBodyBall(world *box2d.B2World, bodyDef *box2d.B2BodyDef, radius float64, components int, options SoftBodyBallOptions) *SoftBodyBall {
	ball := &SoftBodyBall{Color: rl.Red}
	ball.edgeLength = math.Sin(math.Pi/float64(components)) * radius

	ball.CenterBody = world.CreateBody(bodyDef)
	ball.CenterBody.SetUserData(ball)
	center := bodyDef.Position

	centerBallShape := box2d.NewB2CircleShape()
	centerBallShape.SetRadius(5)

	centerBallFixture := ball.CenterBody.CreateFixture(centerBallShape, options.Density)
	centerBallFixture.SetRestitution(0.75)
	centerBallFixture.SetUserData(ball)

	ball.EdgeBodies = make([]*box2d.B2Body, 0, components)

	// Make ball bodies and join to center
	for i := 0; i < components; i++ {
		angleRad := (math.Pi * 2) * (float64(i) / float64(components))

		edgeBodyDef := box2d.NewB2BodyDef()
		edgeBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
		edgeBodyDef.Position = box2d.MakeB2Vec2(center.X+(math.Sin(angleRad)*radius), center.Y+(math.Cos(angleRad)*radius))

		positionNormal := box2d.MakeB2Vec2(math.Sin(angleRad)*radius, math.Cos(angleRad)*radius).Skew()
		edgeBodyDef.Angle = math.Atan2(positionNormal.Y, positionNormal.X)

		edgeBody := world.CreateBody(edgeBodyDef)

		edgeShape := box2d.NewB2PolygonShape()
		edgeShape.SetAsBox(ball.edgeLength, 2)

		edgeFixture := edgeBody.CreateFixture(edgeShape, options.Density)
		edgeFixture.SetRestitution(options.Restitution)
		edgeFixture.SetFriction(options.Friction)

		ball.EdgeBodies = append(ball.EdgeBodies, edgeBody)

		edgeCenterJointDef := box2d.MakeB2DistanceJointDef()
		edgeCenterJointDef.BodyA = ball.CenterBody
		edgeCenterJointDef.BodyB = edgeBody
		edgeCenterJointDef.CollideConnected = true
		edgeCenterJointDef.FrequencyHz = options.SpringFrequencyHz
		edgeCenterJointDef.DampingRatio = options.SpringDampingHz
		edgeCenterJointDef.Length = radius
		world.CreateJoint(&edgeCenterJointDef)
	}

	for i := range ball.EdgeBodies {
		e1 := ball.EdgeBodies[i]
		e2 := ball.EdgeBodies[(i+1)%components]

		//deltaPos := box2d.MakeB2Vec2(e1.GetPosition().X-e2.GetPosition().X, e1.GetPosition().Y-e2.GetPosition().Y)
		//halfDeltaPos := box2d.MakeB2Vec2(deltaPos.X/2, deltaPos.Y/2)

		edgeToEdgeJoint := box2d.MakeB2RevoluteJointDef()
		edgeToEdgeJoint.BodyA = e1
		edgeToEdgeJoint.LocalAnchorA = box2d.MakeB2Vec2(-ball.edgeLength, 0)
		edgeToEdgeJoint.BodyB = e2
		edgeToEdgeJoint.LocalAnchorB = box2d.MakeB2Vec2(ball.edgeLength, 0)
		edgeToEdgeJoint.CollideConnected = false
		world.CreateJoint(&edgeToEdgeJoint)
	}

	return ball
}

func (s *SoftBodyBall) Draw(dt float32, body *box2d.B2Body, fixture *box2d.B2Fixture) {
	ballCenterPos := s.CenterBody.GetPosition()
	centerPos := rl.NewVector2(float32(s.CenterBody.GetPosition().X), float32(s.CenterBody.GetPosition().Y))

	for i := range s.EdgeBodies {
		e1 := s.EdgeBodies[i]
		e2 := s.EdgeBodies[(i+1)%len(s.EdgeBodies)]

		e1Center := e1.GetPosition()
		e2Center := e2.GetPosition()

		e1Theta := math.Acos(box2d.B2Vec2Dot(e1Center, ballCenterPos) / (e1Center.Length() * ballCenterPos.Length()))
		e2Theta := math.Acos(box2d.B2Vec2Dot(e2Center, ballCenterPos) / (e2Center.Length() * ballCenterPos.Length()))

		e1Offset := e1.GetWorldPoint(box2d.MakeB2Vec2(-math.Sin(e1Theta)*2, -math.Cos(e1Theta)*2))
		e2Offset := e2.GetWorldPoint(box2d.MakeB2Vec2(-math.Sin(e2Theta)*2, -math.Cos(e2Theta)*2))

		e1Pos := rl.NewVector2(float32(e1Offset.X), float32(e1Offset.Y))
		e2Pos := rl.NewVector2(float32(e2Offset.X), float32(e2Offset.Y))

		rl.DrawTriangle(e1Pos, e2Pos, centerPos, lightenColor(s.Color))
		rl.DrawLineV(e1Pos, e2Pos, s.Color)
	}
}

func lightenColor(c color.RGBA) color.RGBA {
	return rl.NewColor(c.R, c.G, c.B, c.A/5)
}
