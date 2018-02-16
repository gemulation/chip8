package chip8

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type Emulator struct {
	keys    [KeyboardSize]bool
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
		// display
		emulator.display.Init()
		emulator.display.Clear()
		emulator.display.window.SetTitle(emulator.rom.Name)

		// memory
		emulator.ram.LoadRom(emulator.rom)
		emulator.ram.LoadFont(Font)

		for {
			instruction := emulator.cpu.ReadInstruction(emulator)
			if instruction == nil {
				break
			}
			fmt.Println(instruction)
			instruction.Execute()
			emulator.cpu.UpdateTimers()

			// slow down processor
			time.Sleep(1200 * time.Microsecond)
		}
	})
}
