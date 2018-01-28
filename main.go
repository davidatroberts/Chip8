package main

import (
	"chip8/cpu"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	cycleSpeed = 60
)

var (
	pixels        []byte
	ticksPerFrame float64
)

type timer struct {
	startTicks uint32
}

func newTimer() *timer {
	return &timer{0}
}

func (t *timer) start() {
	t.startTicks = sdl.GetTicks()
}

func (t *timer) getTicks() uint32 {
	return sdl.GetTicks() - t.startTicks
}

func main() {
	// init SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Printf("Failed to init SDL: %v", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	// create the window
	window, err := sdl.CreateWindow(
		"Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 320, 160, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Failed to create window: %v", err)
		os.Exit(1)
	}
	defer window.Destroy()

	// create the renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %v", err)
		os.Exit(1)
	}
	defer renderer.Destroy()

	// create the texture to render to
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, 320, 160)
	if err != nil {
		fmt.Printf("Failed to create the texture: %v", err)
		os.Exit(1)
	}

	sdl.ShowCursor(sdl.DISABLE)

	pixels = make([]byte, 320*160*4)

	// frame capping setup
	ticksPerFrame = 1000.0 / float64(cycleSpeed)
	capTimer := newTimer()
	countedFrames := 0
	fpsTimer := newTimer()
	fpsTimer.start()

	// init the CPU
	cpu := cpu.CPU{}
	cpu.Init()

	var event sdl.Event
	running := true
	for running {
		// start frame rate cap timer
		capTimer.start()

		// process events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		}

		// execute a single CPU cycle
		cpu.ExecuteIteration()

		// update the buffer
		texture.Update(nil, pixels, 320*4)

		// update the screen
		renderer.Clear()
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		countedFrames++

		// delay next cycle to maintain ~60Hz
		frameTicks := capTimer.getTicks()
		if float64(frameTicks) < ticksPerFrame {
			sdl.Delay(uint32(ticksPerFrame - float64(frameTicks)))
		}

	}
}
