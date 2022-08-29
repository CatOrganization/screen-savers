package visualizer

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
)

type Plinko struct {
	cfg Config

	world  *box2d.B2World
	shader rl.Shader

	maskShader     rl.Shader
	maskTextureLoc int32

	timeSinceLastBox float32
	paused           bool

	lightingTexture rl.RenderTexture2D
}

var lightTexture rl.Texture2D

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
	shader := rl.LoadShader("", "resources/shaders/glowing-lights.fs")
	screenHeightLoc := rl.GetShaderLocation(shader, "screen_height")
	scaleFactorLoc := rl.GetShaderLocation(shader, "scale")

	rl.SetShaderValue(shader, screenHeightLoc, []float32{float32(cfg.WindowHeight) * cfg.ScaleFactor}, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, scaleFactorLoc, []float32{cfg.ScaleFactor}, rl.ShaderUniformFloat)

	lightingTexture := rl.LoadRenderTexture(cfg.WindowWidth*int32(cfg.ScaleFactor), cfg.WindowHeight*int32(cfg.ScaleFactor))
	rl.BeginTextureMode(lightingTexture)
	rl.DrawRectangle(0, 0, cfg.WindowWidth/2, cfg.WindowHeight/2, rl.White)
	rl.DrawCircle(100, 100, 100, rl.Green)
	rl.EndTextureMode()

	lightTexture = rl.LoadTexture("resources/img/better_light_cone.png")

	maskShader := rl.LoadShader("", "resources/shaders/alpha-mask.fs")
	maskTextureLoc := rl.GetShaderLocation(maskShader, "texture1")
	//rl.SetShaderValueTexture(maskShader, maskTextureLoc, lightingTexture.Texture)

	//tmpShader = rl.LoadShader("", "resources/shaders/idk.fs")

	world := box2d.MakeB2World(box2d.MakeB2Vec2(0, 98))

	wallsDef := box2d.NewB2BodyDef()
	wallsDef.Type = box2d.B2BodyType.B2_staticBody
	wallsDef.Position.Set(0, 0)
	wallsBody := world.CreateBody(wallsDef)

	//leftWallShape := box2d.NewB2EdgeShape()
	//leftWallShape.Set(box2d.MakeB2Vec2(0, 0), box2d.MakeB2Vec2(0, float64(cfg.WindowHeight)))
	//wallsBody.CreateFixture(leftWallShape, 1)
	//
	//rightWallShape := box2d.NewB2EdgeShape()
	//rightWallShape.Set(box2d.MakeB2Vec2(float64(cfg.WindowWidth), 0), box2d.MakeB2Vec2(float64(cfg.WindowWidth), float64(cfg.WindowHeight)))
	//wallsBody.CreateFixture(rightWallShape, 1)

	var lightPositions []box2d.B2Vec2
	lightPositionsBaseLoc := rl.GetShaderLocation(shader, "light_positions")
	lightIntensitiesBaseLoc := rl.GetShaderLocation(shader, "light_intensities")
	maxLightDistanceLoc := rl.GetShaderLocation(shader, "max_light_distance")
	rl.SetShaderValue(shader, maxLightDistanceLoc, []float32{100.0}, rl.ShaderUniformFloat)

	for y := 0; y < int(cfg.WindowHeight); y += 100 {
		for x := 0; x < int(cfg.WindowWidth); x += 100 {
			lightIntensityLoc := lightIntensitiesBaseLoc + int32(len(lightPositions))
			lightPositionLoc := lightPositionsBaseLoc + int32(len(lightPositions))

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
			pegFixture.SetUserData(NewPlinkoPegDrawableContactListener(pegShape, shader, lightIntensityLoc, lightPositionLoc))

			jointDef := box2d.MakeB2MouseJointDef()
			jointDef.SetBodyA(wallsBody)
			jointDef.SetBodyB(pegBody)

			jointDef.Target.Set(pegBody.GetPosition().X, pegBody.GetPosition().Y)
			jointDef.MaxForce = 2500 * pegBody.GetMass()
			world.CreateJoint(&jointDef)

			rl.SetShaderValue(shader, lightIntensitiesBaseLoc+int32(len(lightPositions)), []float32{0}, rl.ShaderUniformFloat)
			rl.SetShaderValue(shader, lightPositionsBaseLoc+int32(len(lightPositions)), []float32{float32(pegBody.GetPosition().X), float32(pegBody.GetPosition().Y)}, rl.ShaderUniformVec2)
			lightPositions = append(lightPositions, pegBody.GetPosition())
		}
	}

	numLightsLoc := rl.GetShaderLocation(shader, "num_lights")
	rl.SetShaderValue(shader, numLightsLoc, []float32{float32(len(lightPositions))}, rl.ShaderUniformFloat)

	plinko := &Plinko{
		cfg:             cfg,
		shader:          shader,
		world:           &world,
		lightingTexture: lightingTexture,
		maskShader:      maskShader,
		maskTextureLoc:  maskTextureLoc,
	}

	world.SetContactListener(plinko)
	return plinko
}

func (p *Plinko) Update(dt float32) error {
	if rl.IsKeyPressed(32) {
		p.paused = !p.paused
	}

	if p.paused {
		return nil
	}

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

	p.prepareLightingMask()

	return nil
}

// WUMBO ON
func (p *Plinko) Draw(debug bool) error {
	//var (
	//	img    = rl.GenImageColor(800, 600, rl.Red)
	//	imgtex = rl.LoadTextureFromImage(img)
	//)

	rl.ClearBackground(rl.DarkGray)

	for body := p.world.GetBodyList(); body != nil; body = body.GetNext() {
		for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
			if drawable, ok := fixture.GetUserData().(Drawable); ok {
				drawable.Draw(rl.GetFrameTime(), body, fixture)
			}
		}
	}

	// TODO: try blend mode subtractive instead of this shader because the shader is sloooow
	//rl.BeginShaderMode(p.shader)
	//rl.DrawRectangle(0, 0, p.cfg.WindowWidth, p.cfg.WindowHeight, rl.Black)
	//rl.EndShaderMode()
	//
	rl.BeginShaderMode(p.maskShader)
	//rl.SetShaderValueTexture(p.maskShader, p.maskTextureLoc, p.lightingTexture.Texture)
	//rl.DrawRectangle(0, 0, p.cfg.WindowWidth, p.cfg.WindowHeight, rl.Black)
	rl.DrawTexture(p.lightingTexture.Texture, 0, 0, rl.White)
	rl.EndShaderMode()

	c := rl.White
	c.A = 120
	//rl.DrawTexture(p.lightingTexture.Texture, 0, 0, rl.White)
	//rl.DrawTexture(imgTex, 0, 0, rl.White)

	//rl.BeginShaderMode(tmpShader)
	//rl.DrawRectangle(0, 0, p.cfg.WindowWidth, p.cfg.WindowHeight, rl.White)
	//rl.EndShaderMode()

	if debug {
		DebugDrawWorld(p.world)
	}

	return nil
}

func (p *Plinko) prepareLightingMask() {
	rl.BeginTextureMode(p.lightingTexture)
	c := rl.NewColor(0, 0, 0, 0)
	rl.ClearBackground(c)
	//rl.DrawRectangle(0, 0, p.cfg.WindowWidth, p.cfg.WindowHeight, rl.Black)
	//rl.BeginBlendMode(rl.BlendCustom)

	for body := p.world.GetBodyList(); body != nil; body = body.GetNext() {
		for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
			if drawable, ok := fixture.GetUserData().(*PlinkoPegDrawableContactListener); ok {
				v := drawable.fadeInOut.Value()
				if v == 0 {
					continue
				}

				worldPoint := body.GetWorldCenter()
				offset := float32(75 * v)
				//pos := rl.NewVector2(offset+float32(worldPoint.X), offset+float32(p.cfg.WindowHeight)-float32(worldPoint.Y))
				//r := rl.NewRectangle(0, 0, 100, -100)

				//rl.DrawTexturePro(lightTexture, r, r, pos, 0, rl.White)
				//rl.DrawTextureRec(lightTexture, rl.NewRectangle(0, 0, 100, -100), pos, rl.White)

				rl.DrawTextureEx(lightTexture, rl.NewVector2(offset+float32(worldPoint.X), offset+float32(p.cfg.WindowHeight)-float32(worldPoint.Y)), 180, float32(v)*1.5, rl.White)
				//rl.DrawCircle(int32(worldPoint.X), p.cfg.WindowHeight-int32(worldPoint.Y), float32(v)*100, rl.White)
			}
		}
	}

	//rl.EndBlendMode()
	rl.EndTextureMode()
}
