package cpu

import "testing"

func TestJumpToAddress(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x1500
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != 0x0500 {
		t.Errorf("Expected: %x got: %x", 0x0500, c.pc)
	}
}

func TestSkipNextEqualsSkip(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x3110
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[1] = 16

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != ProgramStartPosition+4 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+4, c.pc)
	}
}

func TestSkipNextEqualsNotSkip(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x3110
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[1] = 15

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != ProgramStartPosition+2 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+4, c.pc)
	}
}

func TestSkipNextNotEqualsSkip(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x4203
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[2] = 4

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != ProgramStartPosition+4 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+4, c.pc)
	}
}

func TestSkipNextNotEqualsNotSkip(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x4203
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[2] = 3

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// exected
	if c.pc != ProgramStartPosition+2 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+2, c.pc)
	}
}

func TestSkipVXEqualsVYSkip(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x5120
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[1] = 128
	c.v[2] = 128

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != ProgramStartPosition+4 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+4, c.pc)
	}
}

func TestSkipVXEqualsVYNotSkip(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x5340
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[3] = 64
	c.v[4] = 128

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != ProgramStartPosition+2 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+2, c.pc)
	}
}

func TestSetVX(t *testing.T) {
	c := CPU{}
	c.Init()

	opcode := 0x6511
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)

	c.v[5] = 5

	// run
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}

	// expected
	if c.pc != ProgramStartPosition+2 {
		t.Errorf("Expected pc: %x got pc: %x", ProgramStartPosition+2, c.pc)
	}
	if c.v[5] != 0x11 {
		t.Errorf("Expected v[5]: %x, got v[5]: %x", 0x11, c.v[5])
	}
}
