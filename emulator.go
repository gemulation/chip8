package chip8

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

type Emulator struct {
	display *Display
	ram     *RAM
	cpu     *CPU
	rom     *ROM
}

func NewEmulator(rom *ROM) *Emulator {
	return &Emulator{
		display: NewDisplay(),
		ram:     NewRAM(),
		cpu:     NewCPU(),
		rom:     rom,
	}
}

func (emulator *Emulator) Run() {
	pixelgl.Run(func() {
		emulator.display.Init()
		emulator.ram.Load(emulator.rom)
		for {
			instruction := emulator.cpu.ReadInstruction(emulator)
			if instruction == nil {
				break
			}
			fmt.Println(instruction)
			// instruction.Execute()
			emulator.display.Update()
		}
	})
}
