package main

import (
	"os"

	"github.com/gemu/chip8"
)

func main() {
	rom, err := chip8.NewROM(os.Args[1])
	if err != nil {
		panic(err)
	}

	emulator := chip8.NewEmulator(rom)
	if err := emulator.Run(); err != nil {
		panic(err)
	}
}
