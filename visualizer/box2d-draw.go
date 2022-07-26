package visualizer

import (
	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmoghaddam385/stock-market-visualizer/effect"
	"image/color"
	"time"
)

type Drawable interface {
	Draw(dt float32, body *box2d.B2Body, fixture *box2d.B2Fixture)
}

type ContactListener interface {
	OnBeginContact(contact box2d.B2ContactInterface)
	OnEndContact(contact box2d.B2ContactInterface)
}

type PlinkoPegDrawableContactListener struct {
	shape *box2d.B2CircleShape

	fadeInOut     *effect.FadeInOut
	totalContacts int
	temperature   int

	shader                  rl.Shader
	lightIntensityShaderLoc int32
	lightPositionShaderLoc  int32
}

func NewPlinkoPegDrawableContactListener(shape *box2d.B2CircleShape, shader rl.Shader, lightIntensityShaderLoc, lightPositionShaderLoc int32) *PlinkoPegDrawableContactListener {
	return &PlinkoPegDrawableContactListener{
		shape:                   shape,
		fadeInOut:               effect.NewFadeInOut(time.Second/5, time.Second*2),
		shader:                  shader,
		lightIntensityShaderLoc: lightIntensityShaderLoc,
		lightPositionShaderLoc:  lightPositionShaderLoc,
	}
}

func (p *PlinkoPegDrawableContactListener) OnBeginContact(contact box2d.B2ContactInterface) {
	p.totalContacts++
	p.fadeInOut.OnFadeInEvent()

	p.temperature += 10
	if p.temperature > 255 {
		p.temperature = 255
	}
}

func (p *PlinkoPegDrawableContactListener) OnEndContact(contact box2d.B2ContactInterface) {
	p.totalContacts--
	if p.totalContacts == 0 {
		p.fadeInOut.OnFadeOutEvent()
	}
}

// TODO: use shaders to simulate light glow around pellets that are in contact
func (p *PlinkoPegDrawableContactListener) Draw(dt float32, body *box2d.B2Body, fixture *box2d.B2Fixture) {
	p.fadeInOut.Step(dt)

	worldCenter := body.GetWorldPoint(p.shape.M_p)
	intensity := maxUint8(0, uint8(245-p.temperature))
	c := color.RGBA{
		R: 255,
		G: intensity,
		B: intensity,
		A: 255, // maxUint8(25, uint8(255*p.fadeInOut.Value())),
	}

	rl.SetShaderValue(p.shader, p.lightIntensityShaderLoc, []float32{float32(p.fadeInOut.Value())}, rl.ShaderUniformFloat)
	rl.SetShaderValue(p.shader, p.lightPositionShaderLoc, []float32{float32(body.GetPosition().X), float32(body.GetPosition().Y)}, rl.ShaderUniformVec2)

	rl.DrawCircle(int32(worldCenter.X), int32(worldCenter.Y), float32(p.shape.GetRadius()), lightenColor(c))
	rl.DrawCircleLines(int32(worldCenter.X), int32(worldCenter.Y), float32(p.shape.GetRadius()), c)

	p.temperature--
	if p.temperature < 0 {
		p.temperature = 0
	}
}

type PlinkoPelletDrawableContactListener struct {
	shape *box2d.B2PolygonShape

	fadeInOut     *effect.FadeInOut
	totalContacts int
	temperature   int
}

func NewPlinkoPelletDrawableContactListener(shape *box2d.B2PolygonShape) *PlinkoPelletDrawableContactListener {
	return &PlinkoPelletDrawableContactListener{
		shape:     shape,
		fadeInOut: effect.NewFadeInOut(time.Second/5, time.Second/2),
	}
}

func (p *PlinkoPelletDrawableContactListener) OnBeginContact(contact box2d.B2ContactInterface) {
	p.totalContacts++
	p.fadeInOut.OnFadeInEvent()

	p.temperature += 10
	if p.temperature > 255 {
		p.temperature = 255
	}
}

func (p *PlinkoPelletDrawableContactListener) OnEndContact(contact box2d.B2ContactInterface) {
	p.totalContacts--
	if p.totalContacts == 0 {
		p.fadeInOut.OnFadeOutEvent()
	}
}

func (p *PlinkoPelletDrawableContactListener) Draw(dt float32, body *box2d.B2Body, fixture *box2d.B2Fixture) {
	p.fadeInOut.Step(dt)

	intensity := maxUint8(0, uint8(245-p.temperature))
	c := color.RGBA{
		R: 255,
		G: intensity,
		B: intensity,
		A: 255, //maxUint8(25, uint8(255*p.fadeInOut.Value())),
	}

	for i := 0; i < p.shape.M_count; i++ {
		v1Index := (i - 1 + p.shape.M_count) % p.shape.M_count

		worldV1 := body.GetWorldPoint(p.shape.M_vertices[v1Index])
		worldV2 := body.GetWorldPoint(p.shape.M_vertices[i])
		worldCentroid := body.GetWorldPoint(p.shape.M_centroid)

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

		rl.DrawTriangle(centroid, v2, v1, lightenColor(c))
		rl.DrawLineV(v1, v2, c)
	}

	p.temperature--
	if p.temperature < 0 {
		p.temperature = 0
	}
}

type PlinkoPelletDrawable struct {
	shape *box2d.B2PolygonShape
}

func (p PlinkoPelletDrawable) Draw(body *box2d.B2Body, fixture *box2d.B2Fixture) {
	DebugDrawPolygonShape(body, p.shape)
}

func minUint8(a, b uint8) uint8 {
	if a < b {
		return a
	}

	return b
}

func maxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	}

	return b
}
