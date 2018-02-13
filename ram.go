package chip8

type RAM struct {
	data [RamSize]byte
}

func NewRAM() *RAM {
	return &RAM{}
}

func (r *RAM) Load(rom *ROM) {
	for i, b := range rom.Data {
		r.data[ProgramLocation+i] = b
	}
}
