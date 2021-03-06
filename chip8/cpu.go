package chip8

type CPU struct {
	stack [StackSize]uint16
	v     [RegSize]uint8
	sp    byte
	pc    uint16
	i     uint16
	dt    uint16
	st    uint16
}

func NewCPU() *CPU {
	return &CPU{
		pc: ProgramLocation,
	}
}
func (cpu *CPU) UpdateTimers() {
	if cpu.dt > 0 {
		cpu.dt--
	}
	if cpu.st > 0 {
		if cpu.st == 1 {
			// TODO: play sound
		}
		cpu.st--
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
	case 0xE:
		switch val & 0xFF {
		case 0x9E:
			return &SkipKey{instruction}
		case 0xA1:
			return &SkipNotKey{instruction}
		}
	case 0xF:
		switch val & 0xFF {
		case 0x07:
			return &GetDelayTimer{instruction}
		case 0x0A:
			return &WaitKey{instruction}
		case 0x15:
			return &SetDelayTimer{instruction}
		case 0x18:
			return &SetSoundTimer{instruction}
		case 0x1E:
			return &AddI{instruction}
		case 0x29:
			return &LoadSprite{instruction}
		case 0x33:
			return &StoreBCD{instruction}
		case 0x55:
			return &WriteMemory{instruction}
		case 0x65:
			return &ReadMemory{instruction}
		}
	}
	return instruction
}
