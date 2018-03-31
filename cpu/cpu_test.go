package cpu

import "testing"

func setOpcode(opcode int, c *CPU) {
	c.memory[c.pc] = byte(opcode >> 8)
	c.memory[c.pc+1] = byte(opcode & 0xFF)
}

func execute(t *testing.T, c *CPU) {
	if err := c.ExecuteIteration(); err != nil {
		t.Errorf("CPU error: %v", err)
	}
}

func checkPC(t *testing.T, c *CPU, value uint16) {
	if c.pc != value {
		t.Errorf("Expected: %x got: %x", value, c.pc)
	}
}

func checkRegister(t *testing.T, c *CPU, reg int, value byte) {
	if c.v[reg] != value {
		t.Errorf("Expected v[%d]: %x, got v[%d]: %x", reg, value, reg, c.v[reg])
	}
}

func TestJumpToAddress(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x1500, &c)

	execute(t, &c)

	checkPC(t, &c, 0x0500)
}

func TestSkipNextEqualsSkip(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x3110, &c)
	c.v[1] = 16

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
}

func TestSkipNextEqualsNotSkip(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x3110, &c)
	c.v[1] = 15

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
}

func TestSkipNextNotEqualsSkip(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x4203, &c)
	c.v[2] = 4

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
}

func TestSkipNextNotEqualsNotSkip(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x4203, &c)
	c.v[2] = 3

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
}

func TestSkipVXEqualsVYSkip(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x5120, &c)
	c.v[1] = 128
	c.v[2] = 128

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
}

func TestSkipVXEqualsVYNotSkip(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x5340, &c)
	c.v[3] = 64
	c.v[4] = 128

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
}

func TestSetVX(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x6511, &c)
	c.v[5] = 5

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 5, 0x11)
}

func TestAddVXPlusKK(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x7721, &c)
	c.v[7] = 7

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 7, 40)
}

func TestSetVXEqualsVY(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x89A0, &c)
	c.v[9] = 13
	c.v[0xA] = 44

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 9, c.v[0xA])
	checkRegister(t, &c, 9, 44)
}

func TestSetVXEqualsVXOrVY(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8011, &c)
	c.v[0] = 1
	c.v[1] = 2

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 3)
}

func TestSetVXEqualsVXAndVY(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8102, &c)
	c.v[1] = 1
	c.v[0] = 0

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 1, 0)

	c.Init()
	c.v[1] = 1
	c.v[0] = 1

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 1, 1)
}

func TestVXEqualsVXXorVY(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8013, &c)
	c.v[0] = 1
	c.v[1] = 0

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 1)

	c.Init()
	c.v[0] = 1
	c.v[1] = 1

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 0)

	c.Init()
	c.v[0] = 0
	c.v[1] = 0

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 0)
}

func TestVXEqualsVXPlusVYCarry(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8014, &c)
	c.v[0] = 5
	c.v[1] = 20

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 25)
	checkRegister(t, &c, 0xF, 0)

	c.Init()
	c.v[0] = 1
	c.v[1] = 255

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 1)
	checkRegister(t, &c, 0xF, 1)

	c.Init()
	c.v[1] = 129

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 130)
	checkRegister(t, &c, 0xF, 0)
}

func TestVXEqualsVXMinusVYBorrow(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8015, &c)
	c.v[0] = 200
	c.v[1] = 100

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 100)
	checkRegister(t, &c, 0xF, 1)

	c.Init()
	c.v[0] = 2
	c.v[1] = 250

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 8)
	checkRegister(t, &c, 0xF, 0)
}

func TestVXEqualsVYMinusVXBorrow(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8017, &c)
	c.v[0] = 51
	c.v[1] = 60

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 9)
	checkRegister(t, &c, 0xF, 1)

	c.Init()
	c.v[0] = 250
	c.v[1] = 2

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 8)
	checkRegister(t, &c, 0xF, 0)

}

func TestVXEqualsVYShiftRight(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8016, &c)
	c.v[0] = 0
	c.v[1] = 2

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 1)
	checkRegister(t, &c, 0xF, 0)

	c.Init()
	c.v[0] = 0
	c.v[1] = 9

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 4)
	checkRegister(t, &c, 0xF, 1)

}

func TestVXEqualsVYShiftLeft(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x801E, &c)
	c.v[0] = 0
	c.v[1] = 1

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 2)
	checkRegister(t, &c, 1, 2)
	checkRegister(t, &c, 0xF, 0)

	c.Init()
	c.v[0] = 0
	c.v[1] = 0xFF

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 254)
	checkRegister(t, &c, 1, 254)
	checkRegister(t, &c, 0xF, 1)
}

func TestSkipNextVXNotEqualVY(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x9010, &c)
	c.v[0] = 0
	c.v[1] = 1

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)

	c.Init()
	c.v[0] = 5
	c.v[1] = 5

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)

}
