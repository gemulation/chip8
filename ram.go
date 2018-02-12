package chip8

type RAM struct {
	data []byte
}

func NewRAM() *RAM {
	return &RAM{data: make([]byte, RamSize, RamSize)}
}

func (r *RAM) Load(rom *ROM) {
	for i, b := range rom.Data {
		r.data[ProgramLocation+i] = b
	}
}
