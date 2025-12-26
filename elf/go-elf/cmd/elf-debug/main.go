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

		elfReader := &elf.Reader{elfFile, fileData}
		return elf.Print(elfReader)
	case "write":
		return elf.Write(nil, "output.elf")
	default:
		return fmt.Errorf("unknown action '%s'", action)
	}

	// ph := ProgramHeader64{}
	// sh := SectionHeader64{}
	// fmt.Printf("program header size: %d\n", unsafe.Sizeof(ph))
	// fmt.Printf("section header size: %d\n", unsafe.Sizeof(sh))
}
