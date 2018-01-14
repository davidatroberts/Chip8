package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	cycleSpeed             = 60
	programStartPosition   = 0x200
	programEndPosition     = 0xFFF
	interpreterEndPosition = 0x1FF
)

var (
	pixels []byte
	keys   []bool
	memory [4096]byte
	v      [16]byte
	i      uint16
	st     byte
	dt     byte
	pc     uint16
	sp     byte
	stack  [16]uint16

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

func init() {
	pc = programStartPosition
	i = 0
	sp = 0

	// should load the font set here
}

func executeCPUIteration() {
	// get the opcode
	opcode := uint16(memory[pc])<<8 | uint16(memory[pc+1])

	switch opcode & 0xF000 {
	case 0x1000:
		// jump to addess NNN
		addr := opcode & 0x0FFF
		pc = addr
	case 0x2000:
		// call subroutine at NNN
		stack[sp] = pc
		sp++
		pc = opcode & 0x0FFF
	case 0x3000:
		// skip next instruction if v[x] == nn
		nn := byte(opcode & 0x00FF)
		x := (opcode & 0x0F00) >> 8
		if v[x] == nn {
			pc += 4
		} else {
			pc += 2
		}
	case 0x4000:
		// skips next instruction if v[x] != nn
		nn := byte(opcode & 0x00FF)
		x := (opcode & 0x0F00) >> 8
		if v[x] != nn {
			pc += 4
		} else {
			pc += 2
		}
	case 0x5000:
		// skips the next instruction if v[x] == v[y]
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if v[x] == v[y] {
			pc += 4
		} else {
			pc += 2
		}
	case 0x6000:
		// sets v[x] = nn
		x := (opcode & 0x0F00) >> 8
		nn := byte(opcode & 0x00FF)
		v[x] = nn
		pc += 2
	case 0x7000:
		// adds nn to v[x]
		x := (opcode & 0x0F00) >> 8
		nn := byte(opcode & 0x00FF)
		v[x] += nn
		pc += 2
	case 0x8000:
		// TODO Vx/Vy register stuff
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		switch opcode & 0x000F {
		case 0x0000:
			// set v[x] = v[y]
			v[x] = v[y]
		case 0x0001:
			// set v[x] = v[x] OR v[y]
			v[x] = v[x] | v[y]
		case 0x0002:
			// set v[x] = v[x] AND v[y]
			v[x] = v[x] & v[y]
		case 0x0003:
			// set v[x] = v[x] XOR v[y]
			v[x] = v[x] ^ v[y]
		case 0x0004:
			// set v[x] = v[x] + v[y], set v[F] = carry.
			if uint16(v[x])+uint16(v[y]) > 255 {
				v[0xF] = 1
			} else {
				v[0xF] = 0
			}
			v[x] += v[y]
		case 0x0005:
			// set v[x] = v[x] - v[y], set v[F] = NOT borrow
			if v[x] > v[y] {
				v[0xF] = 1
			} else {
				v[0xF] = 0
			}
			v[x] -= v[y]
		case 0x0006:
			// set v[x] = v[y] SHR 1.
			v[0xF] = v[y] & 0x01
			v[y] = v[y] >> 1
			v[x] = v[y]
		case 0x0007:
			// set v[x] = v[y] - v[x], set v[F] = NOT borrow.
			if v[x] > v[y] {
				v[0xF] = 1
			} else {
				v[0xF] = 0
			}
			v[x] = v[y] - v[x]
		case 0x000E:
			// set v[x] = v[y] SHL 1.
			v[0xF] = v[y] & 0x80
			v[y] = v[y] << 1
			v[x] = v[y]
		}
		pc += 2
	case 0x9000:
		// skip next instruction if v[x] != v[y]
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 8
		if v[x] != v[y] {
			pc += 4
		} else {
			pc += 2
		}
	case 0xA000:
		// set i to address nnn
		addr := opcode & 0x0FFF
		i = addr
		pc += 2
	case 0xB000:
		// jump to address v[0]+nnn
		addr := opcode & 0x0FFF
		pc = addr + uint16(v[0])
	case 0xC000:
		// set v[x] to a random number with a mask of nn
		x := (opcode & 0x0F00) >> 8
		nn := byte(opcode & 0x00FF)
		rvalue := byte(rand.Intn(256))
		v[x] = rvalue & nn
		pc += 2
	case 0xD000:
		// TODO draws a sprite at vx/vy thats width of 8 pixels and height of n
		// x := (opcode & 0x0F00) >> 8
		// y := (opcode & 0x00F0) >> 4
		// n := opcode & 0x000F
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			// skip next instruction if key with the value of v[x] is pressed.
			x := (opcode & 0x0F00) >> 8
			if keys[v[x]] {
				pc += 4
			} else {
				pc += 2
			}
		case 0x00A1:
			// skip next instruction if key with the value of v[x] is not pressed.
			x := (opcode & 0x0F00) >> 8
			if !keys[v[x]] {
				pc += 4
			} else {
				pc += 2
			}
		}
	case 0xF000:
		x := (opcode & 0x0F00) >> 8
		switch opcode & 0x00FF {
		case 0x0007:
			// sets v[x] to the value of the delay timer.
			v[x] = dt
			pc += 2
		case 0x000A:
			// TODO wait for a key press, store the value of the key in v[x]
		case 0x0015:
			// sets the delay timer to v[x]
			dt = v[x]
			pc += 2
		case 0x0018:
			// sets the sound timer to v[x]
			st = v[x]
			pc += 2
		case 0x001E:
			// adds v[x] to i
			i += uint16(v[x])
			pc += 2
		case 0x0029:
			// TODO sets i to the location of the sprite for the character in v[x]
		case 0x0033:
			// TODO store BCD representation of v[x] in memory locations i, i+1, and i+2
		case 0x0055:
			// stores v[0] to v[x] (including v[x]) in memory starting at address i
			var xx uint16
			for xx = 0; xx <= x; xx++ {
				memory[xx] = v[xx]
				i++
			}
			pc += 2
		case 0x0065:
			// fills registers v[0] through v[x] from memory starting at location i
			var xx uint16
			for xx = 0; xx <= x; xx++ {
				v[xx] = memory[i]
				i++
			}
			pc += 2
		}
	default:
		fmt.Printf("Illegal instruction: %x", opcode)
		os.Exit(1)
	}
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
	init()

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
		executeCPUIteration()

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
