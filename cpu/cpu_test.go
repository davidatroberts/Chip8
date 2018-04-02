package cpu

import (
	"testing"
)

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
		t.Errorf("Expected PC: %x, got PC: %x", value, c.pc)
	}
}

func checkDT(t *testing.T, c *CPU, value byte) {
	if c.dt != value {
		t.Errorf("Expected DT: %x, got DT: %x", value, c.dt)
	}
}

func checkRegister(t *testing.T, c *CPU, reg int, value byte) {
	if c.v[reg] != value {
		t.Errorf("Expected v[%d]: %x, got v[%d]: %x", reg, value, reg, c.v[reg])
	}
}

func checkI(t *testing.T, c *CPU, value uint16) {
	if c.i != value {
		t.Errorf("Expected i: %x, got: %x", value, c.i)
	}
}

func checkMemory(t *testing.T, c *CPU, pos int, value byte) {
	if c.memory[pos] != value {
		t.Errorf("Expected mem[%x]: %d, got mem[%x]: %d", pos, value, pos, c.memory[pos])
	}
}

func TestJumpToAddress(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x1500, &c)

	execute(t, &c)

	checkPC(t, &c, 0x0500)
}

func TestJumpToAddressPlusV0(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xB010, &c)
	c.v[0] = 7

	execute(t, &c)

	checkPC(t, &c, 0x0017)
}

func TestSkipNextVXEqualsKK(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x3110, &c)
	c.v[1] = 16

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)

	c.Init()
	c.v[1] = 15

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
}

func TestSkipVXNotEqualsKK(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x4203, &c)
	c.v[2] = 4

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)

	c.Init()
	setOpcode(0x4203, &c)
	c.v[2] = 3

	execute(t, &c)
	checkPC(t, &c, ProgramStartPosition+2)
}

func TestSkipNextVXEqualsVY(t *testing.T) {
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

func TestSkipNextVXEqualsKeyPressed(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xE09E, &c)

	c.v[0] = 0xF
	c.keys[0xF] = true

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)

	c.Init()
	c.v[0] = 0xF
	c.keys[0xF] = false

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)

}

func TestSkipNextVXNotEqualsKeyPressed(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xE0A1, &c)

	c.v[0] = 0xF
	c.keys[0xF] = false

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)

	c.Init()
	c.v[0] = 0xF
	c.keys[0xF] = true

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

func TestVXEqualsVY(t *testing.T) {
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

func TestVXEqualsVXOrVY(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0x8011, &c)
	c.v[0] = 1
	c.v[1] = 2

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkRegister(t, &c, 0, 3)
}

func TestVXEqualsVXAndVY(t *testing.T) {
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

func TestSetIToAddr(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xA00F, &c)

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkI(t, &c, 0x000F)
}

func TestAddVXToI(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xA00F, &c)
	c.v[0] = 4

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkI(t, &c, 0x000F)

	setOpcode(0xF01E, &c)

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
	checkI(t, &c, 0x0013)
}

func TestVXEqualsRandomNumberWithMask(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xC0FF, &c)

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	if c.v[0] < 0 || c.v[0] >= 255 {
		t.Errorf("Expected v[0] range [0,256), got: %d", c.v[0])
	}

	c.Init()
	setOpcode(0xC0FE, &c)

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	if c.v[0] < 0 || c.v[0] >= 255 || c.v[0]%2 == 1 {
		t.Errorf("Expected v[0] range [0,256) and even, got: %d", c.v[0])
	}
}

func TestSetDTToVX(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xF015, &c)

	c.v[0] = 8

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkDT(t, &c, 8)

	// just run any command to get timer going
	setOpcode(0x6100, &c)
	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
	checkDT(t, &c, 7)
}

func TestSetVXToDT(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xF015, &c)

	c.v[0] = 8

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkDT(t, &c, 8)

	setOpcode(0xF107, &c)
	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
	checkRegister(t, &c, 1, 7)
}

func TestStoreV0ToVXInMemoryFromI(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xFF55, &c)

	for i := 0; i <= 0xF; i++ {
		c.v[i] = byte(i)
	}
	c.i = 0

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	for i := 0; i <= 0xF; i++ {
		checkMemory(t, &c, i, c.v[i])
	}
	checkI(t, &c, 16)
}

func TestStoreMemoryFromIInV0ToVX(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xFF65, &c)

	for i := 0; i <= 0xF; i++ {
		c.memory[200+i] = byte(i)
	}

	c.i = 200

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	for i := 0; i <= 0xF; i++ {
		checkRegister(t, &c, i, byte(i))
	}
	checkI(t, &c, 200+16)
}

func TestStoreBCDFromVXInI(t *testing.T) {
	c := CPU{}
	c.Init()
	setOpcode(0xF033, &c)

	c.v[0] = 128
	c.i = 10

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+2)
	checkMemory(t, &c, 10, 1)
	checkMemory(t, &c, 11, 2)
	checkMemory(t, &c, 12, 8)

	setOpcode(0xF033, &c)
	c.v[0] = 47

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+4)
	checkMemory(t, &c, 10, 0)
	checkMemory(t, &c, 11, 4)
	checkMemory(t, &c, 12, 7)

	setOpcode(0xF033, &c)
	c.v[0] = 4

	execute(t, &c)

	checkPC(t, &c, ProgramStartPosition+6)
	checkMemory(t, &c, 10, 0)
	checkMemory(t, &c, 11, 0)
	checkMemory(t, &c, 12, 4)
}
