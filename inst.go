package chip8

import (
	"fmt"
	"math/rand"
)

type Instruction interface {
	Execute()
	fmt.Stringer
}

type BaseInstruction struct {
	emulator *Emulator
	addr     uint16
	val      uint16
}

// Execute the instruction.
func (b *BaseInstruction) Execute() {
}

func (b *BaseInstruction) String() string {
	return fmt.Sprintf("%04X - %04X", b.addr, b.val)
}

// Clear the display.
// 00E0 - CLS
type Clear struct{ *BaseInstruction }

// Execute the instruction.
func (c *Clear) Execute() {
	c.emulator.display.Clear()
}

func (c *Clear) String() string {
	return fmt.Sprintf("%04X CLS", c.val)
}

// Return from a subroutine.
// 00EE - RET
// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
type Return struct{ *BaseInstruction }

// Execute the instruction.
func (r *Return) Execute() {
	r.emulator.cpu.pc = r.emulator.cpu.stack[r.emulator.cpu.sp] // retrieve the program counter from the call stack
	r.emulator.cpu.sp--                                         // decrement the stack
}

func (r *Return) String() string {
	return fmt.Sprintf("%04X - %04X - RET", r.addr, r.val)
}

// Jump to location nnn.
// 1nnn - JP addr
// The interpreter sets the program counter to nnn.
type Jump struct{ *BaseInstruction }

// Execute the instruction.
func (j *Jump) Execute() {
	nnn := j.val & 0xFFF
	j.emulator.cpu.pc = nnn
}

func (j *Jump) String() string {
	nnn := j.val & 0xFFF
	return fmt.Sprintf("%04X - %04X - JP %04X", j.addr, j.val, nnn)
}

// Call subroutine at nnn.
// 2nnn - CALL addr
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack.
// The PC is then set to nnn.
type Call struct{ *BaseInstruction }

// Execute the instruction.
func (c *Call) Execute() {
	nnn := c.val & 0xFFF
	c.emulator.cpu.sp++                                         // increment the stack
	c.emulator.cpu.stack[c.emulator.cpu.sp] = c.emulator.cpu.pc // store the program counter on the call stack
	c.emulator.cpu.pc = nnn                                     // set the program counter to the call address
}

func (c *Call) String() string {
	nnn := c.val & 0xFFF
	return fmt.Sprintf("%04X - %04X - CALL %04X", c.addr, c.val, nnn)
}

// SkipX skips next instruction if Vx = kk.
// 3xkk - SE Vx, byte
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
type SkipX struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipX) Execute() {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	if s.emulator.cpu.v[x] == kk {
		s.emulator.cpu.pc += InstructionSize // skip one instruction
	}
}

func (s *SkipX) String() string {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	return fmt.Sprintf("%04X - %04X - SE V%X, %04X", s.addr, s.val, x, kk)
}

// SkipNotX skips next instruction if Vx != kk.
// 4xkk - SNE Vx, byte
// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
type SkipNotX struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipNotX) Execute() {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	if s.emulator.cpu.v[x] != kk {
		s.emulator.cpu.pc += InstructionSize // skip one instruction
	}
}

func (s *SkipNotX) String() string {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	return fmt.Sprintf("%04X - %04X - SNE V%X, %04X", s.addr, s.val, x, kk)
}

// SkipXY skips next instruction if Vx = Vy.
// 5xy0 - SE Vx, Vy
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
type SkipXY struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipXY) Execute() {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	if s.emulator.cpu.v[x] == s.emulator.cpu.v[y] {
		s.emulator.cpu.pc += InstructionSize // skip one instruction
	}
}

func (s *SkipXY) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SE V%X, Vy%d", s.addr, s.val, x, y)
}

// LoadX sets Vx = kk.
// 6xkk - LD Vx, byte
// The interpreter puts the value kk into register Vx.
type LoadX struct{ *BaseInstruction }

// Execute the instruction.
func (l *LoadX) Execute() {
	x := (l.val >> 8) & 0xF
	kk := l.val & 0xFF
	l.emulator.cpu.v[x] = kk // load register
}

func (l *LoadX) String() string {
	x := (l.val >> 8) & 0xF
	kk := l.val & 0xFF
	return fmt.Sprintf("%04X - %04X - LD V%X, %04X", l.addr, l.val, x, kk)
}

// AddX sets Vx = Vx + kk.
// 7xkk - ADD Vx, byte
// Adds the value kk to the value of register Vx, then stores the result in Vx.
type AddX struct{ *BaseInstruction }

// Execute the instruction.
func (a *AddX) Execute() {
	x := (a.val >> 8) & 0xF
	kk := a.val & 0xFF
	a.emulator.cpu.v[x] += kk // add value to register
}

func (a *AddX) String() string {
	x := (a.val >> 8) & 0xF
	kk := a.val & 0xFF
	return fmt.Sprintf("%04X - %04X - ADD V%X, %04X", a.addr, a.val, x, kk)
}

// LoadXY sets Vx = Vy.
// 8xy0 - LD Vx, Vy
// Stores the value of register Vy in register Vx.
type LoadXY struct{ *BaseInstruction }

// Execute the instruction.
func (l *LoadXY) Execute() {
	x := (l.val >> 8) & 0xF
	y := (l.val >> 4) & 0xF
	l.emulator.cpu.v[x] = l.emulator.cpu.v[y] // copy register
}

func (l *LoadXY) String() string {
	x := (l.val >> 8) & 0xF
	y := (l.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - LD V%X, V%X", l.addr, l.val, x, y)
}

// OR sets Vx = Vx OR Vy.
// 8xy1 - OR Vx, Vy
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
// A bitwise OR compares the corrseponding bits from two values, and if either bit is 1,
// then the same bit in the result is also 1. Otherwise, it is 0.
type OR struct{ *BaseInstruction }

// Execute the instruction.
func (o *OR) Execute() {
	x := (o.val >> 8) & 0xF
	y := (o.val >> 4) & 0xF
	o.emulator.cpu.v[x] |= o.emulator.cpu.v[y] // bitwise OR
}

func (o *OR) String() string {
	x := (o.val >> 8) & 0xF
	y := (o.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - OR V%X, V%X", o.addr, o.val, x, y)
}

// AND sets Vx = Vx AND Vy.
// 8xy2 - AND Vx, Vy
// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
// A bitwise AND compares the corrseponding bits from two values, and if both bits are 1,
// then the same bit in the result is also 1. Otherwise, it is 0.
type AND struct{ *BaseInstruction }

// Execute the instruction.
func (a *AND) Execute() {
	x := (a.val >> 8) & 0xF
	y := (a.val >> 4) & 0xF
	a.emulator.cpu.v[x] &= a.emulator.cpu.v[y] // bitwise AND
}

func (a *AND) String() string {
	x := (a.val >> 8) & 0xF
	y := (a.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - AND V%X, V%X", a.addr, a.val, x, y)
}

// XOR sets Vx = Vx XOR Vy.
// 8xy3 - XOR Vx, Vy
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
// An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same,
// then the corresponding bit in the result is set to 1. Otherwise, it is 0.
type XOR struct{ *BaseInstruction }

// Execute the instruction.
func (r *XOR) Execute() {
	x := (r.val >> 8) & 0xF
	y := (r.val >> 4) & 0xF
	r.emulator.cpu.v[x] ^= r.emulator.cpu.v[y] // bitwise XOR
}

func (r *XOR) String() string {
	x := (r.val >> 8) & 0xF
	y := (r.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - XOR V%X, V%X", r.addr, r.val, x, y)
}

// AddXY sets Vx = Vx + Vy, set VF = carry.
// 8xy4 - ADD Vx, Vy
// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 0xFF,)
// VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
type AddXY struct{ *BaseInstruction }

// Execute the instruction.
func (a *AddXY) Execute() {
	x := (a.val >> 8) & 0xF
	y := (a.val >> 4) & 0xF
	xy := a.emulator.cpu.v[x] + a.emulator.cpu.v[y]

	// set VF with the carry
	a.emulator.cpu.v[0xF] = 0
	if xy > 0xFF {
		a.emulator.cpu.v[0xF] = 1
	}

	// only the lowest 8 bits of the result are kept, and stored in Vx.
	a.emulator.cpu.v[x] = xy & 0xFF
}

func (a *AddXY) String() string {
	x := (a.val >> 8) & 0xF
	y := (a.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - ADD V%X, V%X", a.addr, a.val, x, y)
}

// SubXY sets Vx = Vx - Vy, set VF = NOT borrow.
// 8xy5 - SUB Vx, Vy
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
type SubXY struct{ *BaseInstruction }

// Execute the instruction.
func (s *SubXY) Execute() {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	xy := s.emulator.cpu.v[x] - s.emulator.cpu.v[y]

	// set VF with NOT borrow
	s.emulator.cpu.v[0xF] = 0
	if s.emulator.cpu.v[x] > s.emulator.cpu.v[y] {
		s.emulator.cpu.v[0xF] = 1
	}

	s.emulator.cpu.v[x] = xy
}

func (s *SubXY) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SUB V%X, V%X", s.addr, s.val, x, y)
}

// SHR sets Vx = Vx SHR 1.
// 8xy6 - SHR Vx {, Vy}
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
type SHR struct{ *BaseInstruction }

// Execute the instruction.
func (s *SHR) Execute() {
	x := (s.val >> 8) & 0xF

	s.emulator.cpu.v[0xF] = 0
	if s.emulator.cpu.v[x]&1 == 1 {
		s.emulator.cpu.v[0xF] = 1
	}

	s.emulator.cpu.v[x] /= 2
}

func (s *SHR) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SHR V%X {, V%X}", s.addr, s.val, x, y)
}

// SubN set Vx = Vy - Vx, set VF = NOT borrow.
// 8xy7 - SUBN Vx, Vy
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
type SubN struct{ *BaseInstruction }

// Execute the instruction.
func (s *SubN) Execute() {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	yx := s.emulator.cpu.v[y] - s.emulator.cpu.v[x]

	// set VF with NOT borrow
	s.emulator.cpu.v[0xF] = 0
	if s.emulator.cpu.v[y] > s.emulator.cpu.v[x] {
		s.emulator.cpu.v[0xF] = 1
	}

	s.emulator.cpu.v[x] = yx
}

func (s *SubN) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SUBN V%X, V%X", s.addr, s.val, x, y)
}

// SHL sets Vx = Vx SHL 1.
// 8xyE - SHL Vx {, Vy}
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
type SHL struct{ *BaseInstruction }

// Execute the instruction.
func (s *SHL) Execute() {
	x := (s.val >> 8) & 0xF

	s.emulator.cpu.v[0xF] = 0
	if (s.emulator.cpu.v[x]>>3)&1 == 1 {
		s.emulator.cpu.v[0xF] = 1
	}

	s.emulator.cpu.v[x] *= 2
}

func (s *SHL) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SHL V%X {, V%X}", s.addr, s.val, x, y)
}

// SkipNotXY skips next instruction if Vx != Vy.
// 9xy0 - SNE Vx, Vy
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
type SkipNotXY struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipNotXY) Execute() {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	if s.emulator.cpu.v[x] != s.emulator.cpu.v[y] {
		s.emulator.cpu.pc += InstructionSize // skip one instruction
	}
}

func (s *SkipNotXY) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SNE V%X, V%X", s.addr, s.val, x, y)
}

// LoadI sets I = nnn.
// Annn - LD I, addr
// The value of register I is set to nnn.
type LoadI struct{ *BaseInstruction }

// Execute the instruction.
func (l *LoadI) Execute() {
	nnn := l.val & 0xFFF
	l.emulator.cpu.i = nnn
}

func (l *LoadI) String() string {
	nnn := l.val & 0xFFF
	return fmt.Sprintf("%04X - %04X - LD I, %04X", l.addr, l.val, nnn)
}

// JumpV0 jumps to location nnn + V0.
// Bnnn - JP V0, addr
// The program counter is set to nnn plus the value of V0.
type JumpV0 struct{ *BaseInstruction }

// Execute the instruction.
func (j *JumpV0) Execute() {
	nnn := (j.val & 0xFFF) + j.emulator.cpu.v[0]
	j.emulator.cpu.pc = nnn
}

func (j *JumpV0) String() string {
	nnn := (j.val & 0xFFF) + j.emulator.cpu.v[0]
	return fmt.Sprintf("%04d - %04X - JP V0, %04X", j.addr, j.val, nnn)
}

// RND sets Vx = random byte AND kk.
// Cxkk - RND Vx, byte
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk.
// The results are stored in Vx. See instruction 8xy2 for more information on AND.
type RND struct{ *BaseInstruction }

// Execute the instruction.
func (r *RND) Execute() {
	x := (r.val >> 8) & 0xF
	kk := r.val & 0xFF
	r.emulator.cpu.v[x] = uint16(rand.Intn(255)) & kk // bitwise AND
}

func (r *RND) String() string {
	x := (r.val >> 8) & 0xF
	kk := r.val & 0xFF
	return fmt.Sprintf("%04X - %04X - RND V%X, %04X", r.addr, r.val, x, kk)
}

// Draw displays n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
// Dxyn - DRW Vx, Vy, nibble
// The interpreter reads n bytes from memory, starting at the address stored in I.
// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
// Sprites are XORed onto the existing screen.
// If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.
// If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen.
// See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more information on the Chip-8 screen and sprites.
type Draw struct{ *BaseInstruction }

// Execute the instruction.
func (d *Draw) Execute() {
	x := d.emulator.cpu.v[(d.val>>8)&0xF]
	y := d.emulator.cpu.v[(d.val>>4)&0xF]
	height := d.val & 0xF

	d.emulator.cpu.v[0xF] = 0
	for yline := uint16(0); yline < height; yline++ {
		pixel := d.emulator.ram.data[d.emulator.cpu.i+yline]
		for xline := uint16(0); xline < 8; xline++ {
			if (pixel & (0x80 >> xline)) != 0 {
				// handle wrapping of screen
				x := (x + xline) % DisplayWidth
				y := (y + yline) % DisplayHeight
				index := x + (y * 64)

				// check collision
				if d.emulator.display.memory[index] == 1 {
					d.emulator.cpu.v[0xF] = 1 // collision
				}
				// set pixel
				d.emulator.display.memory[index] ^= 1
			}
		}
	}
}

func (d *Draw) String() string {
	x := (d.val >> 8) & 0xF
	y := (d.val >> 4) & 0xF
	n := (d.val) & 0xF
	return fmt.Sprintf("%04X - %04X - DRW V%X, V%X, %04X", d.addr, d.val, x, y, n)
}

// SkipKey skips next instruction if key with the value of Vx is pressed.
// Ex9E - SKP Vx
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
type SkipKey struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipKey) Execute() {
	x := (s.val >> 8) & 0xF
	if s.emulator.keys[x] {
		s.emulator.cpu.pc += InstructionSize // skip one instruction
	}
}

func (s *SkipKey) String() string {
	x := (s.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - SKP V%X", s.addr, s.val, x)
}

// SkipNotKey skips next instruction if key with the value of Vx is not pressed.
// ExA1 - SKNP Vx
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
type SkipNotKey struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipNotKey) Execute() {
	x := (s.val >> 8) & 0xF
	if !s.emulator.keys[x] {
		s.emulator.cpu.pc += InstructionSize // skip one instruction
	}
}

func (s *SkipNotKey) String() string {
	x := (s.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - SKNP V%X", s.addr, s.val, x)
}

// GetDelayTimer sets Vx = delay timer value.
// Fx07 - LD Vx, DT
// The value of DT is placed into Vx.
type GetDelayTimer struct{ *BaseInstruction }

// Execute the instruction.
func (g *GetDelayTimer) Execute() {
	x := (g.val >> 8) & 0xF
	g.emulator.cpu.v[x] = g.emulator.cpu.dt
}

func (g *GetDelayTimer) String() string {
	x := (g.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD V%X, DT", g.addr, g.val, x)
}

// WaitKey wait for a key press, store the value of the key in Vx.
// Fx0A - LD Vx, K
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
type WaitKey struct{ *BaseInstruction }

// Execute the instruction.
func (w *WaitKey) Execute() {
	// x := (w.val >> 8) & 0xF
}

func (w *WaitKey) String() string {
	x := (w.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD V%X, K", w.addr, w.val, x)
}

// SetDelayTimer sets delay timer = Vx.
// Fx15 - LD DT, Vx
// DT is set equal to the value of Vx.
type SetDelayTimer struct{ *BaseInstruction }

// Execute the instruction.
func (s *SetDelayTimer) Execute() {
	x := (s.val >> 8) & 0xF
	s.emulator.cpu.dt = s.emulator.cpu.v[x]
}

func (s *SetDelayTimer) String() string {
	x := (s.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD DT, V%X", s.addr, s.val, x)
}

// SetSoundTimer sets sound timer = Vx.
// Fx18 - LD ST, Vx
// ST is set equal to the value of Vx.
type SetSoundTimer struct{ *BaseInstruction }

// Execute the instruction.
func (s *SetSoundTimer) Execute() {
	x := (s.val >> 8) & 0xF
	s.emulator.cpu.st = s.emulator.cpu.v[x]
}

func (s *SetSoundTimer) String() string {
	x := (s.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD ST, V%X", s.addr, s.val, x)
}

// AddI sets I = I + Vx.
// Fx1E - ADD I, Vx
// The values of I and Vx are added, and the results are stored in I.
type AddI struct{ *BaseInstruction }

// Execute the instruction.
func (a *AddI) Execute() {
	x := (a.val >> 8) & 0xF
	a.emulator.cpu.i += a.emulator.cpu.v[x]
}

func (a *AddI) String() string {
	x := (a.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - ADD I, V%X", a.addr, a.val, x)
}

// LoadSprite sets I = location of sprite for digit Vx.
// Fx29 - LD F, Vx
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
// See section 2.4, Display, for more information on the Chip-8 hexadecimal font.
type LoadSprite struct{ *BaseInstruction }

// Execute the instruction.
func (l *LoadSprite) Execute() {
	x := (l.val >> 8) & 0xF
	l.emulator.cpu.i = l.emulator.cpu.v[x] * SpriteSize
}

func (l *LoadSprite) String() string {
	x := (l.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD F, V%X", l.addr, l.val, x)
}

// StoreBCD stores the BCD representation of Vx in memory locations I, I+1, and I+2.
// Fx33 - LD B, Vx
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I,
// the tens digit at location I+1, and the ones digit at location I+2.
type StoreBCD struct{ *BaseInstruction }

// Execute the instruction.
func (s *StoreBCD) Execute() {
	x := (s.val >> 8) & 0xF
	vx := s.emulator.cpu.v[x]
	i := s.emulator.cpu.i
	s.emulator.ram.data[i] = byte(vx / 100)
	s.emulator.ram.data[i+1] = byte((vx / 10) % 10)
	s.emulator.ram.data[i+2] = byte((vx % 100) % 10)
}

func (s *StoreBCD) String() string {
	x := (s.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD B, V%X", s.addr, s.val, x)
}

// WriteMemory stores registers V0 through Vx in memory starting at location I.
// Fx55 - LD [I], Vx
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
type WriteMemory struct{ *BaseInstruction }

// Execute the instruction.
func (w *WriteMemory) Execute() {
	x := (w.val >> 8) & 0xF
	for i := uint16(0); i <= x; i++ {
		w.emulator.ram.data[w.emulator.cpu.i+i] = byte(w.emulator.cpu.v[i])
	}
}

func (w *WriteMemory) String() string {
	x := (w.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD [I], V%X", w.addr, w.val, x)
}

// ReadMemory reads registers V0 through Vx from memory starting at location I.
// Fx65 - LD Vx, [I]
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
type ReadMemory struct{ *BaseInstruction }

// Execute the instruction.
func (l *ReadMemory) Execute() {
	x := (l.val >> 8) & 0xF
	for i := uint16(0); i <= x; i++ {
		l.emulator.cpu.v[i] = uint16(l.emulator.ram.data[l.emulator.cpu.i+i])
	}
}

func (l *ReadMemory) String() string {
	x := (l.val >> 8) & 0xF
	return fmt.Sprintf("%04X - %04X - LD V%X, [I]", l.addr, l.val, x)
}
