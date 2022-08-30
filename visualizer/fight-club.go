package visualizer

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmoghaddam385/stock-market-visualizer/physics"
)

type FightClub struct {
	cfg Config

	world *box2d.B2World

	paused bool
}

func NewFightClub(cfg Config) *FightClub {
	world := box2d.MakeB2World(box2d.MakeB2Vec2(0, 98))

	floorBodyDef := box2d.NewB2BodyDef()
	floorBodyDef.Type = box2d.B2BodyType.B2_staticBody
	floorBody := world.CreateBody(floorBodyDef)
	floorEdge := box2d.NewB2EdgeShape()
	floorEdge.Set(box2d.MakeB2Vec2(0, float64(cfg.WindowHeight)), box2d.MakeB2Vec2(float64(cfg.WindowWidth), float64(cfg.WindowHeight)))
	floorFixture := floorBody.CreateFixture(floorEdge, 1)
	floorFixture.SetRestitution(0.5)
	floorFixture.SetFriction(0.5)

	boxBodyDef := box2d.NewB2BodyDef()
	boxBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	boxBodyDef.Position.Set(300, 200)
	boxBody := world.CreateBody(boxBodyDef)
	boxShape := box2d.NewB2PolygonShape()
	boxShape.SetAsBox(50, 50)
	boxBody.CreateFixture(boxShape, 1)

	physics.NewSoftBodyBall(&world, box2d.MakeB2Vec2(350, 100), 50, 15)

	return &FightClub{
		cfg:   cfg,
		world: &world,
	}
}

func (f *FightClub) Update(dt float32) error {
	if rl.IsKeyPressed(32) {
		f.paused = !f.paused
	}

	if f.paused {
		return nil
	}

	f.world.Step(float64(dt), 6, 2)

	// Remove off-screen bodies
	for body := f.world.GetBodyList(); body != nil; body = body.GetNext() {
		if body.GetWorldCenter().Y > float64(f.cfg.WindowHeight)+100 {
			f.world.DestroyBody(body)
		}
	}

	return nil
}

func (f *FightClub) Draw(debug bool) error {
	rl.ClearBackground(rl.Black)

	for body := f.world.GetBodyList(); body != nil; body = body.GetNext() {
		for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
			if drawable, ok := fixture.GetUserData().(Drawable); ok {
				drawable.Draw(rl.GetFrameTime(), body, fixture)
			}
		}
	}

	if debug {
		DebugDrawWorld(f.world)
	}

	return nil
}
