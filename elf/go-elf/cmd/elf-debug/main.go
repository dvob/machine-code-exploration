package main

import (
	"flag"
	"fmt"
	"go-elf"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("missing action")
	}

	action := flag.Arg(0)
	switch action {
	case "read":
		if flag.NArg() < 2 {
			return fmt.Errorf("usage: %s read ELF-FILE", os.Args[0])
		}

		fileName := flag.Arg(1)

		fileData, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}

		elfFile, err := elf.Read(fileData)
		if err != nil {
			return err
		}

		elfReader := &elf.Reader{
			File: elfFile,
			Data: fileData,
		}
		return elf.Print(elfReader)

	case "write":
		var (
			virtualAddress uint64 = 0x401000
			code                  = []byte{
				0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, // mov $60, %rax
				0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00, // mov $33, %rdi
				0x0f, 0x05, // syscall
			}
		)
		elfBinary := elf.Write(virtualAddress, virtualAddress, code)
		return os.WriteFile("output.elf", elfBinary, 0755)

	case "compile":
		var (
			virtualAddress uint64 = 0x401000
		)
		entryPoint, code := elf.Compile(virtualAddress)

		elfBinary := elf.Write(virtualAddress, entryPoint, code)
		return os.WriteFile("output.elf", elfBinary, 0755)

	default:
		return fmt.Errorf("unknown action '%s'", action)
	}
}
