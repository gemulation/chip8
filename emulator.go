package chip8

import "fmt"

type Emulator struct {
	ram *RAM
	cpu *CPU
	rom *ROM
}

func NewEmulator(rom *ROM) *Emulator {
	return &Emulator{
		ram: NewRAM(),
		cpu: NewCPU(),
		rom: rom,
	}
}

func (emulator *Emulator) Run() error {
	emulator.ram.Load(emulator.rom)
	for {
		instruction := emulator.cpu.ReadInstruction(emulator.ram)
		if instruction == nil {
			break
		}
		fmt.Println(instruction)
		// if err := instruction.Execute(); err != nil {
		// 	return err
		// }
	}
	return nil
}
