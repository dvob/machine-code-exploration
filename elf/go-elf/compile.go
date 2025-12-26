package elf

import (
	"encoding/binary"
)

type Compiler struct {
	startAddr uint64
	buf       []byte
}

// Compile generates machine code that:
// - Writes "Hello\n" to stdout
// - Increments a counter variable
// - Exits with the counter value as exit code
func Compile(startAddr uint64) (entryPoint uint64, code []byte) {
	c := &Compiler{
		startAddr: startAddr,
		buf:       make([]byte, 0),
	}

	// Data section
	str := "Hello World!\n"
	helloAddr := c.addString(str)
	counterAddr := c.addInt32(33)

	// Mark where code starts
	entryPoint = startAddr + uint64(len(c.buf))

	// Code section
	c.emitWrite(1, helloAddr, len(str))
	c.emitIncrementCounter(counterAddr)
	c.emitExit(counterAddr)

	return entryPoint, c.buf
}

func (c *Compiler) addString(s string) uint64 {
	addr := c.startAddr + uint64(len(c.buf))
	c.buf = append(c.buf, []byte(s)...)
	c.buf = append(c.buf, 0) // null terminator
	return addr
}

func (c *Compiler) addInt32(value int32) uint64 {
	addr := c.startAddr + uint64(len(c.buf))
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(value))
	c.buf = append(c.buf, b...)
	return addr
}

// mov r64, imm32 (sign-extended to 64-bit)
func (c *Compiler) emitMovRegImm32(reg byte, value uint32) {
	c.buf = append(c.buf, 0x48, 0xc7, 0xc0+reg)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, value)
	c.buf = append(c.buf, b...)
}

// mov r64, imm64
func (c *Compiler) emitMovRegImm64(reg byte, value uint64) {
	c.buf = append(c.buf, 0x48, 0xb8+reg)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, value)
	c.buf = append(c.buf, b...)
}

// mov r32, [abs32]
func (c *Compiler) emitMovRegMem32(reg byte, addr uint64) {
	c.buf = append(c.buf, 0x8b, (reg<<3)|0x04, 0x25)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(addr))
	c.buf = append(c.buf, b...)
}

// mov [abs32], r32
func (c *Compiler) emitMovMemReg32(addr uint64, reg byte) {
	c.buf = append(c.buf, 0x89, (reg<<3)|0x04, 0x25)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(addr))
	c.buf = append(c.buf, b...)
}

// add r32, imm8
func (c *Compiler) emitAddRegImm8(reg byte, value uint8) {
	c.buf = append(c.buf, 0x83, 0xc0+reg, value)
}

// syscall
func (c *Compiler) emitSyscall() {
	c.buf = append(c.buf, 0x0f, 0x05)
}

func (c *Compiler) emitWrite(fd int, bufAddr uint64, count int) {
	// mov rax, 1
	c.emitMovRegImm32(0, 1) // rax = syscall 1 (write)
	// mov rdi, fd
	c.emitMovRegImm32(7, uint32(fd)) // rdi = fd
	// mov rsi, bufAddr
	c.emitMovRegImm64(6, bufAddr) // rsi = buffer address
	// mov rdx, count
	c.emitMovRegImm32(2, uint32(count)) // rdx = count
	// syscall
	c.emitSyscall()
}

func (c *Compiler) emitIncrementCounter(addr uint64) {
	// mov eax, [addr]
	c.emitMovRegMem32(0, addr)
	// add eax, 1
	c.emitAddRegImm8(0, 1)
	// mov [addr], eax
	c.emitMovMemReg32(addr, 0)
}

func (c *Compiler) emitExit(counterAddr uint64) {
	// mov rax, 60
	c.emitMovRegImm32(0, 60) // rax = syscall 60 (exit)
	// mov edi, [counterAddr]
	c.emitMovRegMem32(7, counterAddr) // edi = counter value
	// syscall
	c.emitSyscall()
}
