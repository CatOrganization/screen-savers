package visualizer

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
)

type Plinko struct {
	cfg Config

	world *box2d.B2World

	timeSinceLastBox float32
}

func (p *Plinko) BeginContact(contact box2d.B2ContactInterface) {
	if listener, ok := contact.GetFixtureA().GetUserData().(ContactListener); ok {
		listener.OnBeginContact(contact)
	}

	if listener, ok := contact.GetFixtureB().GetUserData().(ContactListener); ok {
		listener.OnBeginContact(contact)
	}
}

func (p *Plinko) EndContact(contact box2d.B2ContactInterface) {
	if listener, ok := contact.GetFixtureA().GetUserData().(ContactListener); ok {
		listener.OnEndContact(contact)
	}

	if listener, ok := contact.GetFixtureB().GetUserData().(ContactListener); ok {
		listener.OnEndContact(contact)
	}
}

func (p *Plinko) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold)     {}
func (p *Plinko) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {}

func NewPlinko(cfg Config) *Plinko {
	world := box2d.MakeB2World(box2d.MakeB2Vec2(0, 98))

	wallsDef := box2d.NewB2BodyDef()
	wallsDef.Type = box2d.B2BodyType.B2_staticBody
	wallsDef.Position.Set(0, 0)
	wallsBody := world.CreateBody(wallsDef)

	leftWallShape := box2d.NewB2EdgeShape()
	leftWallShape.Set(box2d.MakeB2Vec2(0, 0), box2d.MakeB2Vec2(0, float64(cfg.WindowHeight)))
	wallsBody.CreateFixture(leftWallShape, 1)

	rightWallShape := box2d.NewB2EdgeShape()
	rightWallShape.Set(box2d.MakeB2Vec2(float64(cfg.WindowWidth), 0), box2d.MakeB2Vec2(float64(cfg.WindowWidth), float64(cfg.WindowHeight)))
	wallsBody.CreateFixture(rightWallShape, 1)

	for y := 0; y < int(cfg.WindowHeight); y += 100 {
		for x := 0; x < int(cfg.WindowWidth); x += 100 {
			pegBodyDef := box2d.NewB2BodyDef()
			pegBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
			pegBodyDef.Position.Set(float64(x), float64(y)+25)
			if y%200 == 0 {
				pegBodyDef.Position.X += 50
			}

			pegBody := world.CreateBody(pegBodyDef)
			pegShape := box2d.NewB2CircleShape()
			pegShape.SetRadius(5)
			pegFixture := pegBody.CreateFixture(pegShape, 1)
			pegFixture.SetRestitution(0.2)
			pegFixture.SetUserData(NewPlinkoPegDrawableContactListener(pegShape))

			jointDef := box2d.MakeB2MouseJointDef()
			jointDef.SetBodyA(wallsBody)
			jointDef.SetBodyB(pegBody)

			jointDef.Target.Set(pegBody.GetPosition().X, pegBody.GetPosition().Y)
			jointDef.MaxForce = 2500 * pegBody.GetMass()
			world.CreateJoint(&jointDef)
		}
	}

	plinko := &Plinko{
		cfg:   cfg,
		world: &world,
	}

	world.SetContactListener(plinko)
	return plinko
}

func (p *Plinko) Update(dt float32) error {
	p.world.Step(float64(dt), 6, 2)

	// Remove off-screen bodies
	for body := p.world.GetBodyList(); body != nil; body = body.GetNext() {
		if body.GetWorldCenter().Y > float64(p.cfg.WindowHeight)+100 {
			p.world.DestroyBody(body)
		}
	}

	p.timeSinceLastBox += dt
	if p.timeSinceLastBox > 0.5 {
		// Add new Bodies
		boxDef := box2d.NewB2BodyDef()
		boxDef.Type = box2d.B2BodyType.B2_dynamicBody
		boxDef.Position.Set(rand.Float64()*float64(p.cfg.WindowWidth), -100)
		boxDef.AngularVelocity = math.Pi * rand.Float64()
		boxBody := p.world.CreateBody(boxDef)

		boxShape := box2d.NewB2PolygonShape()
		boxShape.SetAsBox(5, float64(rand.Intn(15)+10))
		boxFixture := boxBody.CreateFixture(boxShape, 1)
		boxFixture.SetUserData(NewPlinkoPelletDrawableContactListener(boxShape))

		p.timeSinceLastBox = 0
	}

	return nil
}

func (p *Plinko) Draw(debug bool) error {
	rl.ClearBackground(rl.Black)

	if debug {
		DebugDrawWorld(p.world)
	}

	for body := p.world.GetBodyList(); body != nil; body = body.GetNext() {
		for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
			if drawable, ok := fixture.GetUserData().(Drawable); ok {
				drawable.Draw(rl.GetFrameTime(), body, fixture)
			}
		}
	}

	return nil
}
