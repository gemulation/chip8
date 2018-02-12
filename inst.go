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
	pos uint16
}

// Execute the instruction.
func (b *BaseInstruction) Execute() error {
	return nil
}

func (b *BaseInstruction) String() string {
	return fmt.Sprintf("%04X - %04X", b.pos, b.val)
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
	return fmt.Sprintf("%04X RET", r.val)
}

// Jump to location nnn.
// 1nnn - JP addr
// The interpreter sets the program counter to nnn.
type Jump struct{ *BaseInstruction }

// Execute the instruction.
func (j *Jump) Execute() error {
	// addr := ((j.val & 0xFFF) - ProgramLocation) / InstructionSize
	// j.cpu.pc = addr
	return nil
}

func (j *Jump) String() string {
	addr := ((j.val & 0xFFF) - ProgramLocation) / InstructionSize
	return fmt.Sprintf("%04d - %04X - JP %04X", j.pos, j.val, addr)
}

// Call subroutine at nnn.
// 2nnn - CALL addr
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack.
// The PC is then set to nnn.
type Call struct{ *BaseInstruction }

// Execute the instruction.
func (c *Call) Execute() error {
	addr := ((c.val & 0xFFF) - ProgramLocation) / InstructionSize
	c.cpu.stack[c.cpu.sp] = c.cpu.pc // store the program counter on the call stack
	c.cpu.pc = addr                  // set the program counter to the call address
	c.cpu.sp++                       // increment the stack
	return nil
}

func (c *Call) String() string {
	addr := ((c.val & 0xFFF) - ProgramLocation) / InstructionSize
	return fmt.Sprintf("%04X - %04X - CALL %04X", c.pos, c.val, addr)
}

// override fun clear() = builder.line("clear")
// override fun ret() = builder.line("ret")
// override fun jmp(address: Int) = builder.line("jmp 0x${address.hex}")
// override fun call(address: Int) = builder.line("call 0x${address.hex}")
// override fun jeq(reg: Int, value: Int) = builder.line("jeq v${reg.hex}, 0x${value.hex}")
// override fun jneq(reg: Int, value: Int) = builder.line("jneq v${reg.hex}, 0x${value.hex}")
// override fun jeqr(reg1: Int, reg2: Int) = builder.line("jeqr v${reg1.hex}, v${reg2.hex}")
// override fun set(reg: Int, value: Int) = builder.line("set v${reg.hex}, 0x${value.hex}")
// override fun add(reg: Int, value: Int) = builder.line("add v${reg.hex}, 0x${value.hex}")
// override fun setr(reg1: Int, reg2: Int) = builder.line("setr v${reg1.hex}, v${reg2.hex}")
// override fun or(reg1: Int, reg2: Int) = builder.line("or v${reg1.hex}, v${reg2.hex}")
// override fun and(reg1: Int, reg2: Int) = builder.line("and v${reg1.hex}, v${reg2.hex}")
// override fun xor(reg1: Int, reg2: Int) = builder.line("xor v${reg1.hex}, v${reg2.hex}")
// override fun addr(reg1: Int, reg2: Int) = builder.line("addr v${reg1.hex}, v${reg2.hex}")
// override fun sub(reg1: Int, reg2: Int) = builder.line("sub v${reg1.hex}, v${reg2.hex}")
// override fun shr(reg1: Int) = builder.line("shr v${reg1.hex}")
// override fun subb(reg1: Int, reg2: Int) = builder.line("subb v${reg1.hex}, v${reg2.hex}")
// override fun shl(reg1: Int) = builder.line("shl v${reg1.hex}")
// override fun jneqr(reg1: Int, reg2: Int) = builder.line("jneqr v${reg1.hex}, v${reg2.hex}")
// override fun seti(value: Int) = builder.line("seti 0x${value.hex}")
// override fun jmpv0(address: Int) = builder.line("jmpv0 0x${address.hex}")
// override fun rand(reg: Int, value: Int) = builder.line("rand v${reg.hex}, 0x${value.hex}")
// override fun draw(reg1: Int, reg2: Int, value: Int) = builder.line("draw v${reg1.hex}, v${reg2.hex}, 0x${value.hex}")
// override fun jkey(reg: Int) = builder.line("jkey v${reg.hex}")
// override fun jnkey(reg: Int) = builder.line("jnkey v${reg.hex}")
// override fun getdelay(reg: Int) = builder.line("getdelay v${reg.hex}")
// override fun waitkey(reg: Int) = builder.line("waitkey v${reg.hex}")
// override fun setdelay(reg: Int) = builder.line("setdelay v${reg.hex}")
// override fun setsound(reg: Int) = builder.line("setsound v${reg.hex}")
// override fun addi(reg: Int) = builder.line("addi v${reg.hex}")
// override fun spritei(reg: Int) = builder.line("spritei v${reg.hex}")
// override fun bcd(reg: Int) = builder.line("bcd v${reg.hex}")
// override fun push(reg: Int) = builder.line("push v0-v${reg.hex}")
// override fun pop(reg: Int) = builder.line("pop v0-v${reg.hex}")
