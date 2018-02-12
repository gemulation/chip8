package chip8_test

import (
	"testing"

	"github.com/gemu/chip8"
	"github.com/stretchr/testify/require"
)

func TestNewRom(t *testing.T) {
	rom, err := chip8.NewROM("../roms/pong.rom")

	require.Nil(t, err)
	require.NotNil(t, rom)
	require.Equal(t, "pong.rom", rom.Name)
	require.True(t, len(rom.Data) > 0)
}
