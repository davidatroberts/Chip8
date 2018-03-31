package cpu

import (
	"fmt"
	"math/rand"
)

const (
	// ProgramStartPosition initial value of pc
	ProgramStartPosition   = 0x200
	programEndPosition     = 0xFFF
	interpreterEndPosition = 0x1FF
)

// CPU basic Chip CPU
type CPU struct {
	keys   []bool
	memory [4096]byte
	v      [16]byte
	i      uint16
	st     byte
	dt     byte
	pc     uint16
	sp     byte
	stack  [16]uint16
}

// Init sets the initial values of the registers
func (c *CPU) Init() {
	c.pc = ProgramStartPosition
	c.i = 0
	c.sp = 0
}

// ExecuteIteration executes a single iteration of the CPU
func (c *CPU) ExecuteIteration() error {
	// get the opcode
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	switch opcode & 0xF000 {
	case 0x1000:
		// jump to addess NNN
		addr := opcode & 0x0FFF
		c.pc = addr
	case 0x2000:
		// call subroutine at NNN
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = opcode & 0x0FFF
	case 0x3000:
		// skip next instruction if v[x] == nn
		nn := byte(opcode & 0x00FF)
		x := (opcode & 0x0F00) >> 8
		if c.v[x] == nn {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x4000:
		// skips next instruction if v[x] != nn
		nn := byte(opcode & 0x00FF)
		x := (opcode & 0x0F00) >> 8
		if c.v[x] != nn {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x5000:
		// skips the next instruction if v[x] == v[y]
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if c.v[x] == c.v[y] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x6000:
		// sets v[x] = nn
		x := (opcode & 0x0F00) >> 8
		nn := byte(opcode & 0x00FF)
		c.v[x] = nn
		c.pc += 2
	case 0x7000:
		// adds nn to v[x]
		x := (opcode & 0x0F00) >> 8
		nn := byte(opcode & 0x00FF)
		c.v[x] += nn
		c.pc += 2
	case 0x8000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		switch opcode & 0x000F {
		case 0x0000:
			// set v[x] = v[y]
			c.v[x] = c.v[y]
		case 0x0001:
			// set v[x] = v[x] OR v[y]
			c.v[x] = c.v[x] | c.v[y]
		case 0x0002:
			// set v[x] = v[x] AND v[y]
			c.v[x] = c.v[x] & c.v[y]
		case 0x0003:
			// set v[x] = v[x] XOR v[y]
			c.v[x] = c.v[x] ^ c.v[y]
		case 0x0004:
			// set v[x] = v[x] + v[y], set v[F] = carry.
			if uint16(c.v[x])+uint16(c.v[y]) > 255 {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] += c.v[y] % 255
		case 0x0005:
			// set v[x] = v[x] - v[y], set v[F] = NOT borrow
			if c.v[x] > c.v[y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] -= c.v[y]
		case 0x0006:
			// set v[x] = v[y] SHR 1.
			c.v[0xF] = c.v[y] & 0x01
			c.v[y] = c.v[y] >> 1
			c.v[x] = c.v[y]
		case 0x0007:
			// set v[x] = v[y] - v[x], set v[F] = NOT borrow.
			if c.v[y] > c.v[x] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] = c.v[y] - c.v[x]
		case 0x000E:
			// set v[x] = v[y] SHL 1.
			c.v[0xF] = (c.v[y] & 0x80) >> 7
			c.v[y] = c.v[y] << 1
			c.v[x] = c.v[y]
		}
		c.pc += 2
	case 0x9000:
		// skip next instruction if v[x] != v[y]
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if c.v[x] != c.v[y] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0xA000:
		// set i to address nnn
		addr := opcode & 0x0FFF
		c.i = addr
		c.pc += 2
	case 0xB000:
		// jump to address v[0]+nnn
		addr := opcode & 0x0FFF
		c.pc = addr + uint16(c.v[0])
	case 0xC000:
		// set v[x] to a random number with a mask of nn
		x := (opcode & 0x0F00) >> 8
		nn := byte(opcode & 0x00FF)
		rvalue := byte(rand.Intn(256))
		c.v[x] = rvalue & nn
		c.pc += 2
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
			if c.keys[c.v[x]] {
				c.pc += 4
			} else {
				c.pc += 2
			}
		case 0x00A1:
			// skip next instruction if key with the value of v[x] is not pressed.
			x := (opcode & 0x0F00) >> 8
			if !c.keys[c.v[x]] {
				c.pc += 4
			} else {
				c.pc += 2
			}
		}
	case 0xF000:
		x := (opcode & 0x0F00) >> 8
		switch opcode & 0x00FF {
		case 0x0007:
			// sets v[x] to the value of the delay timer.
			c.v[x] = c.dt
			c.pc += 2
		case 0x000A:
			// TODO wait for a key press, store the value of the key in v[x]
		case 0x0015:
			// sets the delay timer to v[x]
			c.dt = c.v[x]
			c.pc += 2
		case 0x0018:
			// sets the sound timer to v[x]
			c.st = c.v[x]
			c.pc += 2
		case 0x001E:
			// adds v[x] to i
			c.i += uint16(c.v[x])
			c.pc += 2
		case 0x0029:
			// TODO sets i to the location of the sprite for the character in v[x]
		case 0x0033:
			// TODO store BCD representation of v[x] in memory locations i, i+1, and i+2
		case 0x0055:
			// stores v[0] to v[x] (including v[x]) in memory starting at address i
			var xx uint16
			for xx = 0; xx <= x; xx++ {
				c.memory[xx] = c.v[xx]
				c.i++
			}
			c.pc += 2
		case 0x0065:
			// fills registers v[0] through v[x] from memory starting at location i
			var xx uint16
			for xx = 0; xx <= x; xx++ {
				c.v[xx] = c.memory[c.i]
				c.i++
			}
			c.pc += 2
		}
	default:
		return fmt.Errorf("Unimplemented instruction: %x", opcode)
	}

	return nil
}
