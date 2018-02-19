package chip8

type RAM struct {
	data [RamSize]byte
}

func NewRAM() *RAM {
	return &RAM{}
}

func (r *RAM) LoadRom(rom *ROM) {
	for i, b := range rom.Data {
		r.data[ProgramLocation+i] = b
	}
}

func (r *RAM) LoadFont(font [80]byte) {
	for i, f := range font {
		r.data[i] = f
	}
}
