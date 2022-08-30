package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmoghaddam385/stock-market-visualizer/visualizer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	screenWidth  = 800 * 2
	screenHeight = 600 * 2
)

type Renderer struct {
	v visualizer.Visualizer
}

func main() {
	cmd := cobra.Command{
		Use:               "screen-saver",
		Short:             "Run a visualization that looks like a screen saver",
		Long:              "Don't actually use this as a screen saver, it won't save your screen but it'll look cool.",
		PersistentPreRunE: rootPersistentPreRun,
		RunE:              runE,
	}

	cmd.Flags().Bool("fullscreen", false, "whether or not to start the program in fullscreen mode")
	cmd.Flags().Float64("scale", 1, "The scaling factor to use for visualizations")

	if err := cmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("exiting with error")
	}
}

func runE(_ *cobra.Command, args []string) error {

	cfg := visualizer.Config{
		WindowWidth:  800,
		WindowHeight: 600,
	}

	//rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint)
	rl.InitWindow(cfg.WindowWidth, cfg.WindowHeight, "some kinda visualizer")

	if viper.GetBool("fullscreen") {
		display := rl.GetCurrentMonitor()
		cfg.WindowWidth = int32(rl.GetMonitorWidth(display))
		cfg.WindowHeight = int32(rl.GetMonitorHeight(display))

		rl.SetWindowSize(int(cfg.WindowWidth), int(cfg.WindowHeight))
		rl.ToggleFullscreen()
		rl.HideCursor()
	}

	rl.SetTargetFPS(60)

	scale := viper.GetFloat64("scale")
	scaleCamera := rl.NewCamera2D(rl.Vector2{}, rl.Vector2{}, 0, float32(scale))

	cfg.WindowWidth = int32(float64(cfg.WindowWidth) / scale)
	cfg.WindowHeight = int32(float64(cfg.WindowHeight) / scale)
	cfg.ScaleFactor = float32(scale)
	//v := visualizer.NewPlinko(cfg)
	v := visualizer.NewFightClub(cfg)
	debug := false

	for !rl.WindowShouldClose() {
		if rl.GetCharPressed() == 'd' {
			debug = !debug
		}

		v.Update(rl.GetFrameTime())

		rl.BeginDrawing()

		rl.BeginMode2D(scaleCamera)
		v.Draw(debug)
		rl.EndMode2D()

		if debug {
			rl.DrawText(fmt.Sprintf("FPS: %.2f; %v", rl.GetFPS(), rl.GetKeyPressed()), 10, 10, 12, rl.RayWhite)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
	return nil
}

// bindViperFlagsPreRun binds the flags for a command in PreRunE.
// This has to be done in pre-run because it can only run for the command+subcommands that are actually going to execute.
func bindViperFlagsPreRun(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	return nil
}

func rootPersistentPreRun(cmd *cobra.Command, args []string) error {
	if err := bindViperFlagsPreRun(cmd, args); err != nil {
		return err
	}

	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.TraceLevel)
	}

	logrus.SetOutput(cmd.ErrOrStderr())
	return nil
}
