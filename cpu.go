package chip8

type CPU struct {
	stack [StackSize]uint16
	v     [RegSize]uint16
	dt    uint8
	st    uint8
	sp    uint8
	pc    uint16
	i     uint16
}

func NewCPU() *CPU {
	return &CPU{
		pc: ProgramLocation,
	}
}

func (cpu *CPU) ReadInstruction(emulator *Emulator) Instruction {
	// read 2 bytes integer in little endian format
	val := (uint16(emulator.ram.data[cpu.pc]) << 8) | uint16(emulator.ram.data[cpu.pc+1])
	if val == 0 {
		return nil
	}
	instruction := &BaseInstruction{emulator: emulator, val: val, addr: cpu.pc}
	cpu.pc += InstructionSize

	switch (val >> 12) & 0xF {
	case 0x0:
		switch val {
		case 0x00E0:
			return &Clear{instruction}
		case 0x00EE:
			return &Return{instruction}
		}
	case 0x1:
		return &Jump{instruction}
	case 0x2:
		return &Call{instruction}
	case 0x3:
		return &SkipX{instruction}
	case 0x4:
		return &SkipNotX{instruction}
	case 0x5:
		return &SkipXY{instruction}
	case 0x6:
		return &LoadX{instruction}
	case 0x7:
		return &AddX{instruction}
	case 0x8:
		switch val & 0xF {
		case 0x0:
			return &LoadXY{instruction}
		case 0x1:
			return &OR{instruction}
		case 0x2:
			return &AND{instruction}
		case 0x3:
			return &XOR{instruction}
		case 0x4:
			return &AddXY{instruction}
		case 0x5:
			return &SubXY{instruction}
		case 0x6:
			return &SHR{instruction}
		case 0x7:
			return &SubN{instruction}
		case 0xE:
			return &SHL{instruction}
		}
	case 0x9:
		return &SkipNotXY{instruction}
	case 0xA:
		return &LoadI{instruction}
	case 0xB:
		return &JumpV0{instruction}
	case 0xC:
		return &RND{instruction}
	case 0xD:
		return &Draw{instruction}
	}
	return instruction
}
