package chip8

import "fmt"

type Instruction interface {
	Execute() error
	fmt.Stringer
}

type BaseInstruction struct {
	cpu *CPU
	ram *RAM
	val uint16
	pc  uint16
}

// Execute the instruction.
func (b *BaseInstruction) Execute() error {
	return nil
}

func (b *BaseInstruction) String() string {
	return fmt.Sprintf("%04X - %04X", b.pc, b.val)
}

// Clear the display.
// 00E0 - CLS
type Clear struct{ *BaseInstruction }

// Execute the instruction.
func (c *Clear) Execute() error {
	return nil
}

func (c *Clear) String() string {
	return fmt.Sprintf("%04X CLS", c.val)
}

// Return from a subroutine.
// 00EE - RET
// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
type Return struct{ *BaseInstruction }

// Execute the instruction.
func (r *Return) Execute() error {
	r.cpu.sp--                       // decrement the stack
	r.cpu.pc = r.cpu.stack[r.cpu.sp] // retrieve the program counter from the call stack
	return nil
}

func (r *Return) String() string {
	return fmt.Sprintf("%04d - %04X - RET", r.pc, r.val)
}

// Jump to location nnn.
// 1nnn - JP addr
// The interpreter sets the program counter to nnn.
type Jump struct{ *BaseInstruction }

// Execute the instruction.
func (j *Jump) Execute() error {
	nnn := j.val & 0xFFF
	j.cpu.pc = nnn
	return nil
}

func (j *Jump) String() string {
	nnn := j.val & 0xFFF
	return fmt.Sprintf("%04d - %04X - JP %04X", j.pc, j.val, nnn)
}

// Call subroutine at nnn.
// 2nnn - CALL addr
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack.
// The PC is then set to nnn.
type Call struct{ *BaseInstruction }

// Execute the instruction.
func (c *Call) Execute() error {
	nnn := c.val & 0xFFF
	c.cpu.stack[c.cpu.sp] = c.cpu.pc // store the program counter on the call stack
	c.cpu.pc = nnn                   // set the program counter to the call address
	c.cpu.sp++                       // increment the stack
	return nil
}

func (c *Call) String() string {
	nnn := c.val & 0xFFF
	return fmt.Sprintf("%04X - %04X - CALL %04X", c.pc, c.val, nnn)
}

// SkipX skips next instruction if Vx = kk.
// 3xkk - SE Vx, byte
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
type SkipX struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipX) Execute() error {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	if s.cpu.v[x] == kk {
		s.cpu.pc += InstructionSize * 2 // skip one instruction
	}
	return nil
}

func (s *SkipX) String() string {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	return fmt.Sprintf("%04X - %04X - SE V%d, %04X", s.pc, s.val, x, kk)
}

// SkipNotX skips next instruction if Vx != kk.
// 4xkk - SNE Vx, byte
// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
type SkipNotX struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipNotX) Execute() error {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	if s.cpu.v[x] != kk {
		s.cpu.pc += InstructionSize * 2 // skip one instruction
	}
	return nil
}

func (s *SkipNotX) String() string {
	x := (s.val >> 8) & 0xF
	kk := s.val & 0xFF
	return fmt.Sprintf("%04X - %04X - SNE V%d, %04X", s.pc, s.val, x, kk)
}

// SkipXY skips next instruction if Vx = Vy.
// 5xy0 - SE Vx, Vy
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
type SkipXY struct{ *BaseInstruction }

// Execute the instruction.
func (s *SkipXY) Execute() error {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	if s.cpu.v[x] == s.cpu.v[y] {
		s.cpu.pc += InstructionSize * 2 // skip one instruction
	}
	return nil
}

func (s *SkipXY) String() string {
	x := (s.val >> 8) & 0xF
	y := (s.val >> 4) & 0xF
	return fmt.Sprintf("%04X - %04X - SE V%d, Vy%d", s.pc, s.val, x, y)
}

// Load sets Vx = kk.
// 6xkk - LD Vx, byte
// The interpreter puts the value kk into register Vx.
type Load struct{ *BaseInstruction }

// Execute the instruction.
func (l *Load) Execute() error {
	x := (l.val >> 8) & 0xF
	kk := l.val & 0xFF
	l.cpu.v[x] = kk // load register
	return nil
}

func (l *Load) String() string {
	x := (l.val >> 8) & 0xF
	kk := l.val & 0xFF
	return fmt.Sprintf("%04X - %04X - LD V%d, %04X", l.pc, l.val, x, kk)
}
