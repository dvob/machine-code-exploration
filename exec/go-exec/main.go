package main

/*
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

void call_machine_code(void *code_ptr) {
    void (*fn)(void) = (void (*)(void))code_ptr;
    fn();
}
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

// funcval is Go's internal representation of a function value
// Copied from runtime/runtime2.go
// https://github.com/golang/go/blob/a23d1a4ebe5ca1f4964ad51a92d99edf5a95d530/src/runtime/runtime2.go#L179
type funcval struct {
	fn uintptr
}

func main() {
	// Machine code that calls exit_group(33)
	// exit_group (syscall 231) kills ALL threads, which is required for Go's multi-threaded runtime
	code := []byte{
		0x48, 0xc7, 0xc0, 0xe7, 0x00, 0x00, 0x00, // mov $231, %rax (exit_group syscall)
		0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00, // mov $33, %rdi (exit code)
		0x0f, 0x05, // syscall
	}

	mmap, err := syscall.Mmap(-1, 0, len(code), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		panic(err)
	}

	copy(mmap, code)

	direct := true
	if direct {
		// METHOD 1: Direct funcval construction (pure Go, no CGO)
		// Key insight: funcval must be HEAP-allocated, not stack-allocated

		fv := &funcval{
			fn: uintptr(unsafe.Pointer(&mmap[0])),
		}

		// Cast pointer-to-pointer to func() and dereference
		// This works because Go function values are internally pointers to funcval

		fn := *(*func())(unsafe.Pointer(&fv))

		// you can also do that without using the funval struct. its basically the same.
		// codePtr := new(uintptr)
		// *codePtr = uintptr(unsafe.Pointer(&mmap[0]))
		// fn := *(*func())(unsafe.Pointer(&codePtr))

		fmt.Println("Calling via direct funcval...")
		fn()
	} else {
		C.call_machine_code(unsafe.Pointer(&mmap[0]))
	}
}
