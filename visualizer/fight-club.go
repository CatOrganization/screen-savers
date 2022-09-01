package visualizer

import (
	"fmt"
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmoghaddam385/stock-market-visualizer/physics"
	"math"
	"math/rand"
	"time"
)

const PhysicsToRenderScale = 10

var colors = []rl.Color{rl.RayWhite, rl.Blue, rl.Red, rl.Gold, rl.SkyBlue, rl.Purple, rl.Pink, rl.Beige}

type FightClub struct {
	cfg Config

	physicsWidth, physicsHeight float64

	world *box2d.B2World
	floor *box2d.B2Body

	sceneTime float32

	paused bool
}

func NewFightClub(cfg Config) *FightClub {
	world := box2d.MakeB2World(box2d.MakeB2Vec2(0, 9.8))

	physicsWidth := float64(cfg.WindowWidth / PhysicsToRenderScale)
	physicsHeight := float64(cfg.WindowHeight / PhysicsToRenderScale)

	wallsBodyDef := box2d.NewB2BodyDef()
	wallsBodyDef.Type = box2d.B2BodyType.B2_staticBody
	wallsBody := world.CreateBody(wallsBodyDef)

	floorEdge := box2d.NewB2EdgeShape()
	floorEdge.Set(box2d.MakeB2Vec2(0, physicsHeight), box2d.MakeB2Vec2(physicsWidth, physicsHeight))
	floorFixture := wallsBody.CreateFixture(floorEdge, 1)
	floorFixture.SetRestitution(0.5)
	floorFixture.SetFriction(0.75)

	leftWallEdge := box2d.NewB2EdgeShape()
	leftWallEdge.Set(box2d.MakeB2Vec2(0, physicsHeight), box2d.MakeB2Vec2(0, -physicsHeight))
	leftWallFixture := wallsBody.CreateFixture(leftWallEdge, 1)
	leftWallFixture.SetRestitution(0.5)
	leftWallFixture.SetFriction(0.75)

	rightWallEdge := box2d.NewB2EdgeShape()
	rightWallEdge.Set(box2d.MakeB2Vec2(physicsWidth, physicsHeight), box2d.MakeB2Vec2(physicsWidth, -physicsHeight))
	rightWallFixture := wallsBody.CreateFixture(rightWallEdge, 1)
	rightWallFixture.SetRestitution(0.5)
	rightWallFixture.SetFriction(0.75)

	boxBodyDef := box2d.NewB2BodyDef()
	boxBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	boxBodyDef.Position.Set(300, 400)
	boxBody := world.CreateBody(boxBodyDef)
	boxShape := box2d.NewB2PolygonShape()
	boxShape.SetAsBox(50, 50)
	boxBody.CreateFixture(boxShape, 1)

	//physics.NewSoftBodyBall(&world, box2d.MakeB2Vec2(350, 200), 50, 15)
	//physics.NewSoftBodyBall2(&world, box2d.MakeB2Vec2(325, 100), 40, 10).Color = rl.DarkGreen

	viz := &FightClub{
		cfg:           cfg,
		physicsWidth:  physicsWidth,
		physicsHeight: physicsHeight,
		world:         &world,
	}

	rand.Seed(time.Now().Unix())

	viz.GenerateObjects(5)
	viz.GenerateFloor(int(cfg.WindowWidth / 25))

	return viz
}

func (f *FightClub) Update(dt float32) error {
	if rl.IsKeyPressed(80) {
		f.paused = !f.paused
	}

	if f.paused {
		return nil
	}

	f.world.Step(float64(dt), 8, 4)
	//fmt.Println("-----")

	// Remove off-screen bodies
	for body := f.world.GetBodyList(); body != nil; body = body.GetNext() {
		if body.GetWorldCenter().Y > float64(f.cfg.WindowHeight)+100 {
			f.world.DestroyBody(body)
		}
	}

	f.sceneTime += dt

	// 30 seconds or the right arrow is pressed
	if f.sceneTime > 30 || rl.IsKeyPressed(32) {
		f.sceneTime = 0

		f.GenerateObjects(5 + rand.Intn(15))
		f.GenerateFloor(int(f.cfg.WindowHeight) / 25)
	}

	return nil
}

func (f *FightClub) Draw(debug bool) error {
	rl.ClearBackground(rl.Black)

	for body := f.world.GetBodyList(); body != nil; body = body.GetNext() {
		if drawable, ok := body.GetUserData().(Drawable); ok {
			drawable.Draw(rl.GetFrameTime(), body, nil)
		}

		for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
			if drawable, ok := fixture.GetUserData().(Drawable); ok {
				drawable.Draw(rl.GetFrameTime(), body, fixture)
			}
		}
	}

	if debug {
		b := f.world.GetBodyList().GetNext()
		v := b.GetLinearVelocity()
		b.SetLinearDamping(-1)

		rl.DrawText(fmt.Sprintf("linear velocity: (%.2f, %.2f)", v.X, v.Y), 10, 50, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("linear damping: %.2f", b.GetLinearDamping()), 10, 90, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("angular velocity: %.2f", b.GetAngularVelocity()), 10, 130, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("angular damping: %.2f", b.GetAngularDamping()), 10, 170, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("awake?  %v", b.IsAwake()), 10, 210, 32, rl.RayWhite)

		massData := &box2d.B2MassData{}
		b.GetMassData(massData)

		rl.DrawText(fmt.Sprintf("mass:  %.2f", massData.Mass), 10, 250, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("I:  %.2f", massData.I), 10, 290, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("inertia:  %.2f", b.GetInertia()), 10, 330, 32, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("force:  (%.2f, %.2f)", b.M_force.X, b.M_force.Y), 10, 370, 32, rl.RayWhite)

		DebugDrawWorld(f.world, PhysicsToRenderScale)
	}

	return nil
}

func (f *FightClub) GenerateObjects(numObjects int) {
	// Wipe any dynamic bodies and joints from the world
	for body := f.world.GetBodyList(); body != nil; body = body.GetNext() {
		if body.GetType() == box2d.B2BodyType.B2_dynamicBody {
			f.world.DestroyBody(body)
		}
	}

	for joint := f.world.GetJointList(); joint != nil; joint = joint.GetNext() {
		f.world.DestroyJoint(joint)
	}

	for i := 0; i < numObjects; i++ {
		switch rand.Intn(2) {
		case 0:
			//f.GenerateBox()
			f.GenerateSoftBall()
		case 1:
			f.GenerateBox()
		}
	}
}

func (f *FightClub) GenerateSoftBall() {
	radius := 2.5 + rand.Float64()*5
	bodyDef := f.randomizedBodyDef(radius + 1)

	opts := physics.SoftBodyBallOptions{
		Density:           radius,
		SpringFrequencyHz: 4 + rand.Float64()*10,
		SpringDampingHz:   0.1 + rand.Float64()*0.75,
		Restitution:       0.25 + (rand.Float64() * 0.5),
		Friction:          0.25 + (rand.Float64() * 0.75),

		Color:                randomColor(),
		PhysicsToRenderScale: PhysicsToRenderScale,
	}

	physics.NewSoftBodyBall(f.world, bodyDef, radius, int(radius*4), opts)
}

func (f *FightClub) GenerateBox() {
	width := 0.5 + rand.Float64()*7.5
	height := 0.5 + rand.Float64()*7.5

	body := f.world.CreateBody(f.randomizedBodyDef(math.Max(width, height) + 10))

	shape := box2d.NewB2PolygonShape()
	shape.SetAsBox(height, width)

	fixture := body.CreateFixture(shape, 1)
	fixture.SetFriction(0.25 + (rand.Float64() * 0.75))
	fixture.SetRestitution(0.1 + (rand.Float64() * 0.4))
	fixture.SetUserData(&FightClubBoxDrawable{shape: shape, color: randomColor()})
}

func (f *FightClub) randomizedBodyDef(xPosBuffer float64) *box2d.B2BodyDef {
	bodyDef := box2d.NewB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody

	bodyDef.Position.X = xPosBuffer + rand.Float64()*(f.physicsWidth-(2*xPosBuffer))
	bodyDef.Position.Y = -rand.Float64() * f.physicsHeight
	bodyDef.LinearVelocity.X = -50 + (rand.Float64() * 100)
	bodyDef.LinearVelocity.Y = 5 + rand.Float64()*5

	bodyDef.Angle = rand.Float64() * math.Pi
	bodyDef.AngularVelocity = rand.Float64() * math.Pi

	return bodyDef
}

func randomColor() rl.Color {
	return colors[rand.Intn(len(colors))]
}

func (f *FightClub) GenerateFloor(segments int) {
	if f.floor != nil {
		f.world.DestroyBody(f.floor)
	}

	floorDef := box2d.NewB2BodyDef()
	floorDef.Type = box2d.B2BodyType.B2_staticBody

	f.floor = f.world.CreateBody(floorDef)

	segmentLength := f.physicsWidth / float64(segments)
	heightVariance := f.physicsHeight / 2

	point := box2d.MakeB2Vec2(0, (f.physicsHeight-heightVariance)+(rand.Float64()*heightVariance))
	points := []rl.Vector2{rl.NewVector2(float32(point.X*PhysicsToRenderScale), float32(point.Y*PhysicsToRenderScale))}
	prevAngle := float64(0)

	for i := 0; i < segments; i++ {
		nextAngle := prevAngle + (-math.Pi / 4) + rand.Float64()*(math.Pi/2)
		nextAngle = math.Max(-math.Pi/3, math.Min(math.Pi/3, nextAngle))

		h := segmentLength * math.Tan(nextAngle)
		nextY := point.Y + h
		nextY = math.Min(f.physicsHeight-1, nextY)
		nextY = math.Max(f.physicsHeight-heightVariance, nextY)

		nextPoint := box2d.MakeB2Vec2(point.X+segmentLength, nextY)

		edgeShape := box2d.NewB2EdgeShape()
		edgeShape.Set(point, nextPoint)

		edgeFixture := f.floor.CreateFixture(edgeShape, 1)
		edgeFixture.SetRestitution(0.5)
		edgeFixture.SetFriction(0.75)

		prevAngle = math.Atan2(box2d.B2Vec2Sub(nextPoint, point).Y, box2d.B2Vec2Sub(nextPoint, point).X)

		point = nextPoint
		points = append(points, rl.NewVector2(float32(point.X*PhysicsToRenderScale), float32(point.Y*PhysicsToRenderScale)))
	}

	f.floor.SetUserData(&FightClubFloorDrawable{points: points, floorHeight: float32(f.cfg.WindowHeight)})
}
