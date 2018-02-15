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
	cpu  *CPU
	ram  *RAM
	val  uint16
	addr uint16
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
	r.cpu.pc = r.cpu.stack[r.cpu.sp] // retrieve the program counter from the call stack
	r.cpu.sp--                       // decrement the stack
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
	j.cpu.pc = nnn
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
	c.cpu.sp++                       // increment the stack
	c.cpu.stack[c.cpu.sp] = c.cpu.pc // store the program counter on the call stack
	c.cpu.pc = nnn                   // set the program counter to the call address
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
	if s.cpu.v[x] == kk {
		s.cpu.pc += InstructionSize // skip one instruction
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
	if s.cpu.v[x] != kk {
		s.cpu.pc += InstructionSize // skip one instruction
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
	if s.cpu.v[x] == s.cpu.v[y] {
		s.cpu.pc += InstructionSize // skip one instruction
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
	l.cpu.v[x] = kk // load register
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
	a.cpu.v[x] += kk // add value to register
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
	l.cpu.v[x] = l.cpu.v[y] // copy register
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
	o.cpu.v[x] |= o.cpu.v[y] // bitwise OR
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
	a.cpu.v[x] &= a.cpu.v[y] // bitwise AND
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
	r.cpu.v[x] ^= r.cpu.v[y] // bitwise XOR
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
	xy := a.cpu.v[x] + a.cpu.v[y]

	// set VF with the carry
	a.cpu.v[0xF] = 0
	if xy > 0xFF {
		a.cpu.v[0xF] = 1
	}

	// only the lowest 8 bits of the result are kept, and stored in Vx.
	a.cpu.v[x] = xy & 0xFF
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
	xy := s.cpu.v[x] - s.cpu.v[y]

	// set VF with NOT borrow
	s.cpu.v[0xF] = 0
	if s.cpu.v[x] > s.cpu.v[y] {
		s.cpu.v[0xF] = 1
	}

	s.cpu.v[x] = xy
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

	s.cpu.v[0xF] = 0
	if s.cpu.v[x]&1 == 1 {
		s.cpu.v[0xF] = 1
	}

	s.cpu.v[x] /= 2
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
	yx := s.cpu.v[y] - s.cpu.v[x]

	// set VF with NOT borrow
	s.cpu.v[0xF] = 0
	if s.cpu.v[y] > s.cpu.v[x] {
		s.cpu.v[0xF] = 1
	}

	s.cpu.v[x] = yx
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

	s.cpu.v[0xF] = 0
	if (s.cpu.v[x]>>3)&1 == 1 {
		s.cpu.v[0xF] = 1
	}

	s.cpu.v[x] *= 2
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
	if s.cpu.v[x] != s.cpu.v[y] {
		s.cpu.pc += InstructionSize // skip one instruction
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
	l.cpu.i = nnn
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
	nnn := (j.val & 0xFFF) + j.cpu.v[0]
	j.cpu.pc = nnn
}

func (j *JumpV0) String() string {
	nnn := (j.val & 0xFFF) + j.cpu.v[0]
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
	r.cpu.v[x] = uint16(rand.Intn(255)) & kk // bitwise AND
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
	// x := (d.val >> 8) & 0xF
	// y := (d.val >> 4) & 0xF
	// n := (d.val) & 0xF
	// TODO: implement
}

func (d *Draw) String() string {
	x := (d.val >> 8) & 0xF
	y := (d.val >> 4) & 0xF
	n := (d.val) & 0xF
	return fmt.Sprintf("%04X - %04X - DRW V%X, V%X, %04X", d.addr, d.val, x, y, n)
}
