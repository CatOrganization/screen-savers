package visualizer

type Visualizer interface {
	Update(deltaT float32) error
	Draw(debug bool) error
}

type Config struct {
	WindowWidth, WindowHeight int32
}
