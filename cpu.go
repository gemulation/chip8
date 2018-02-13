package chip8

type CPU struct {
	stack [StackSize]uint16
	v     [RegSize]uint16
	sp    uint8
	pc    uint16
	i     uint16
}

func NewCPU() *CPU {
	return &CPU{
		pc: ProgramLocation,
	}
}

func (cpu *CPU) ReadInstruction(ram *RAM) Instruction {
	// read 2bytes integer in little endian format
	val := (uint16(ram.data[cpu.pc]) << 8) | uint16(ram.data[cpu.pc+1])
	if val == 0 {
		return nil
	}
	instruction := &BaseInstruction{cpu: cpu, ram: ram, val: val, pc: cpu.pc}
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
		return &Load{instruction}

	}
	return instruction
}
