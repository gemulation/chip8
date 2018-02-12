package chip8

type CPU struct {
	stack [StackSize]uint16
	v     [RegSize]uint64
	sp    uint8
	pc    uint16
	i     uint16
}

func NewCPU() *CPU {
	return &CPU{}
}

func (cpu *CPU) ReadInstruction(ram *RAM) Instruction {
	i := ((cpu.pc * InstructionSize) + ProgramLocation)
	msb := uint16(ram.data[i])
	lsb := uint16(ram.data[i+1])
	val := (msb << 8) | lsb
	if val == 0 {
		return nil
	}
	instruction := &BaseInstruction{cpu: cpu, ram: ram, val: val, pos: cpu.pc}
	cpu.pc++

	switch (val >> 12) & 0xF {
	case 0x1:
		return &Jump{instruction}
	case 0x2:
		return &Call{instruction}
		// case 0x3:
		// 	se_bx_byte(self, inst)
		// case 0x4:
		// 	sne_bx_byte(self, inst)
		// case 0x5:
		// 	se_vx_vy(self, inst)
		// case 0x6:
		// 	ld_vx_byte(self, inst)
		// case 0x7:
		// 	add_vx_byte(self, inst)
		// case 0xA:
		// 	ld_i_addr(self, inst)
		// case 0xC:
		// 	rnd_vx_byte(self, inst)
		// case 0xD:
		// 	drw_x_y_n(self, inst)
	}
	return instruction
}
