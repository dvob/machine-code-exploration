package main

import (
	"syscall"
	"unsafe"
)

func main() {
	// we cant run this code. exit seems to hang in runtime
	code := []byte{
		0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00,
		0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00,
		0x0f, 0x05,
	}

	// code = []byte{
	// 	0x48, 0xc7, 0xc0, 0x01, 0x0, 0x0, 0x0, // mov %rax,$0x1
	// 	0x48, 0xc7, 0xc7, 0x01, 0x0, 0x0, 0x0, // mov %rdi,$0x1
	// 	0x48, 0xc7, 0xc2, 0x0c, 0x0, 0x0, 0x0, // mov 0x13, %rdx
	// 	0x48, 0x8d, 0x35, 0x04, 0x0, 0x0, 0x0, // lea 0x4(%rip), %rsi
	// 	0x0f, 0x05, // syscall
	// 	0xc3, 0xcc, // ret
	// 	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, // Hello_(whitespace)
	// 	0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x0, 0x0a, // World!
	// }

	mmap, err := syscall.Mmap(-1, 0, len(code), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		panic(err)
	}

	copy(mmap, code)

	fnPtr := (uintptr)(unsafe.Pointer(&mmap))
	fn := *(*func())(unsafe.Pointer(&fnPtr))
	fn()

}
