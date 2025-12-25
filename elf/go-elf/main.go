package main

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"unsafe"
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

		elf, err := readElf(fileData)
		if err != nil {
			return err
		}

		printElf(elf)
		return nil
	case "write":
		return writeElf(nil, "output.elf")
	default:
		return fmt.Errorf("unknown action '%s'", action)
	}

	// ph := ProgramHeader64{}
	// sh := SectionHeader64{}
	// fmt.Printf("program header size: %d\n", unsafe.Sizeof(ph))
	// fmt.Printf("section header size: %d\n", unsafe.Sizeof(sh))
}

func readElf(data []byte) (*ELFFile, error) {
	ident := ELFIdentifier{}

	err := binary.Read(bytes.NewBuffer(data), binary.NativeEndian, &ident)
	if err != nil {
		return nil, err
	}

	if ident.Class != ELFCLASS64 {
		return nil, fmt.Errorf("32bit binary not supported")
	}

	var byteOrder binary.ByteOrder
	switch ident.Data {
	case ELFDATA2LSB:
		byteOrder = binary.LittleEndian
	case ELFDATA2MSB:
		byteOrder = binary.BigEndian
	default:
		return nil, fmt.Errorf("unknown data field %x", ident.Data)
	}

	elfHeader := &Header64{}
	err = binary.Read(bytes.NewBuffer(data), byteOrder, elfHeader)
	if err != nil {
		return nil, err
	}

	file := &ELFFile{
		Header:         elfHeader,
		ProgramHeaders: []ProgramHeader64{},
		SectionHeaders: []SectionHeader64{},
	}

	for i := 0; i < int(elfHeader.ProgramHeaderCount); i++ {
		programHeader := ProgramHeader64{}
		offset := elfHeader.ProgramHeaderOffset + uint64(i)*uint64(elfHeader.ProgramHeaderSize)
		end := offset + uint64(elfHeader.ProgramHeaderSize)
		err := binary.Read(bytes.NewBuffer(data[offset:end]), byteOrder, &programHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to read program header %d: %w", i, err)
		}
		file.ProgramHeaders = append(file.ProgramHeaders, programHeader)
	}
	for i := 0; i < int(elfHeader.SectionHeaderCount); i++ {
		sectionHeader := SectionHeader64{}
		offset := elfHeader.SectionHeaderOffset + uint64(i)*uint64(elfHeader.SectionHeaderSize)
		end := offset + uint64(elfHeader.SectionHeaderSize)
		err := binary.Read(bytes.NewBuffer(data[offset:end]), byteOrder, &sectionHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to read section header %d: %w", i, err)
		}
		file.SectionHeaders = append(file.SectionHeaders, sectionHeader)
	}
	return file, nil
}

func writeSmallElf(f *ELFFile, output string) error {
	var (
		virtualAddress = 0x401000
		usePadding     = false
		code           = []byte{
			0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00,
			0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00,
			0x0f, 0x05,
		}
	)
	f = &ELFFile{
		Header: &Header64{
			ELFIdentifier: ELFIdentifier{
				Magic:      MagicBytes,
				Class:      ELFCLASS64,
				Data:       ELFDATA2LSB,
				Version:    1,
				OSABI:      0x03, // Linux
				ABIVersion: 0,
				Padding:    [7]byte{},
			},
			Type:                     ET_EXEC,
			Machine:                  0x3e, // AMD x86-64
			Version:                  1,
			Entry:                    uint64(virtualAddress) + uint64(unsafe.Sizeof(Header64{})+unsafe.Sizeof(ProgramHeader64{})),
			ProgramHeaderOffset:      uint64(unsafe.Sizeof(Header64{})),
			SectionHeaderOffset:      0,
			Flags:                    0,
			EhSize:                   uint16(unsafe.Sizeof(Header64{})),
			ProgramHeaderSize:        uint16(unsafe.Sizeof(ProgramHeader64{})),
			ProgramHeaderCount:       1,
			SectionHeaderSize:        0,
			SectionHeaderCount:       0,
			SectionHeaderStringIndex: 0,
		},
		ProgramHeaders: []ProgramHeader64{
			{
				Type:  PT_LOAD,
				Flags: PF_R | PF_X,

				Offset:         0,                      //uint64(unsafe.Sizeof(Header64{}) + unsafe.Sizeof(ProgramHeader64{})),
				VirtualAddress: uint64(virtualAddress), //VADDR + uint64(unsafe.Sizeof(Header64{})+unsafe.Sizeof(ProgramHeader64{})),
				FileSize:       uint64(len(code)),
				MemorySize:     uint64(len(code)),
				Align:          4096,

				// this is what worked and is done by normal compilers
				// Offset:         4096, //uint64(unsafe.Sizeof(Header64{}) + unsafe.Sizeof(ProgramHeader64{})),
				// VirtualAddress: 0x401000,
				// FileSize:       uint64(len(code)),
				// MemorySize:     uint64(len(code)),
				// Align:          0,
			},
		},
		SectionHeaders: []SectionHeader64{},
	}

	var byteOrder binary.ByteOrder
	switch f.Header.Data {
	case ELFDATA2LSB:
		byteOrder = binary.LittleEndian
	case ELFDATA2MSB:
		byteOrder = binary.BigEndian
	default:
		return fmt.Errorf("invalid data type 0x%x", f.Header.Data)
	}

	buf := &bytes.Buffer{}
	err := binary.Write(buf, byteOrder, f.Header)
	if err != nil {
		return err
	}

	fmt.Println("header size:", buf.Len())

	for _, programHeader := range f.ProgramHeaders {
		err = binary.Write(buf, byteOrder, programHeader)
		if err != nil {
			return err
		}
		fmt.Println("ph header size:", buf.Len())
	}

	if usePadding {
		padding := 4096 - buf.Len()
		for i := 0; i < padding; i++ {
			buf.WriteByte(0)
		}
	}

	_, err = buf.Write(code)
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	_, err = outputFile.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = outputFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func writeElf(_ *ELFFile, output string) error {
	var (
		virtualAddress = 0x401000
		usePadding     = false
		code           = []byte{

			0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, // mov $60, %rax
			0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00, // mov $33, %rdi
			0x0f, 0x05, // syscall
		}
	)
	f := &ELFFile{
		Header: &Header64{
			ELFIdentifier: ELFIdentifier{
				Magic:      MagicBytes,
				Class:      ELFCLASS64,
				Data:       ELFDATA2LSB,
				Version:    1,
				OSABI:      0x03, // Linux
				ABIVersion: 0,
				Padding:    [7]byte{},
			},
			Type:                     ET_EXEC,
			Machine:                  0x3e, // AMD x86-64
			Version:                  1,
			Entry:                    uint64(virtualAddress) + uint64(unsafe.Sizeof(Header64{})+unsafe.Sizeof(ProgramHeader64{})),
			ProgramHeaderOffset:      uint64(unsafe.Sizeof(Header64{})),
			SectionHeaderOffset:      0,
			Flags:                    0,
			EhSize:                   uint16(unsafe.Sizeof(Header64{})),
			ProgramHeaderSize:        uint16(unsafe.Sizeof(ProgramHeader64{})),
			ProgramHeaderCount:       1,
			SectionHeaderSize:        0,
			SectionHeaderCount:       0,
			SectionHeaderStringIndex: 0,
		},
		ProgramHeaders: []ProgramHeader64{
			{
				Type:  PT_LOAD,
				Flags: PF_R | PF_X,

				Offset:         0,                      //uint64(unsafe.Sizeof(Header64{}) + unsafe.Sizeof(ProgramHeader64{})),
				VirtualAddress: uint64(virtualAddress), //VADDR + uint64(unsafe.Sizeof(Header64{})+unsafe.Sizeof(ProgramHeader64{})),
				FileSize:       uint64(len(code)),
				MemorySize:     uint64(len(code)),
				Align:          4096,

				// this is what worked and is done by normal compilers
				// Offset:         4096, //uint64(unsafe.Sizeof(Header64{}) + unsafe.Sizeof(ProgramHeader64{})),
				// VirtualAddress: 0x401000,
				// FileSize:       uint64(len(code)),
				// MemorySize:     uint64(len(code)),
				// Align:          0,
			},
		},
		SectionHeaders: []SectionHeader64{},
	}

	var byteOrder binary.ByteOrder
	switch f.Header.Data {
	case ELFDATA2LSB:
		byteOrder = binary.LittleEndian
	case ELFDATA2MSB:
		byteOrder = binary.BigEndian
	default:
		return fmt.Errorf("invalid data type 0x%x", f.Header.Data)
	}

	buf := &bytes.Buffer{}
	err := binary.Write(buf, byteOrder, f.Header)
	if err != nil {
		return err
	}

	fmt.Println("header size:", buf.Len())

	for _, programHeader := range f.ProgramHeaders {
		err = binary.Write(buf, byteOrder, programHeader)
		if err != nil {
			return err
		}
		fmt.Println("ph header size:", buf.Len())
	}

	if usePadding {
		padding := 4096 - buf.Len()
		for i := 0; i < padding; i++ {
			buf.WriteByte(0)
		}
	}

	_, err = buf.Write(code)
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	_, err = outputFile.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = outputFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func printElf(f *ELFFile) {
	printHeader := func(h *Header64) {
		fmt.Printf("class: %s\n", h.Class)
		fmt.Printf("data: %s\n", h.Data)
		fmt.Printf("version: %d\n", h.ELFIdentifier.Version)
		fmt.Printf("os abi: %d\n", h.OSABI)
		fmt.Printf("abi version: %d\n", h.ABIVersion)

		fmt.Printf("type: %s\n", h.Type)
		fmt.Printf("machine: 0x%x\n", h.Machine)
		fmt.Printf("version: %d\n", h.Version)
		fmt.Printf("entry: %x\n", h.Entry)

		fmt.Printf("header size: %x\n", h.EhSize)

		fmt.Printf("program header offset: 0x%x\n", h.ProgramHeaderOffset)
		fmt.Printf("program header count: %d\n", h.ProgramHeaderCount)
		fmt.Printf("program header size: %d\n", h.ProgramHeaderSize)

		fmt.Printf("section header offset: 0x%x\n", h.SectionHeaderOffset)
		fmt.Printf("section header count: %d\n", h.SectionHeaderCount)
		fmt.Printf("section header size: %d\n", h.SectionHeaderSize)

		fmt.Printf("section header string index: 0x%x\n", h.SectionHeaderStringIndex)
		fmt.Println()
	}
	printProgram := func(p ProgramHeader64) {
		fmt.Printf("type: %s\n", p.Type)
		fmt.Printf("flags: 0x%x\n", p.Flags)
		fmt.Printf("offset: 0x%x\n", p.Offset)
		fmt.Printf("virtual addr: 0x%x\n", p.VirtualAddress)
		fmt.Printf("physical addr: 0x%x\n", p.PhysicalAddress)
		fmt.Printf("file size: %d\n", p.FileSize)
		fmt.Printf("memory size: %d\n", p.MemorySize)
		fmt.Printf("align: 0x%x\n", p.Align)
		fmt.Println()
	}
	printSection := func(s SectionHeader64) {
		fmt.Printf("name: %v\n", s.Name)
		fmt.Printf("type: %v\n", s.Type)
		fmt.Printf("flags: %v\n", s.Flags)
		fmt.Printf("addr: 0x%x\n", s.Address)
		fmt.Printf("offset: 0x%x\n", s.Offset)
		fmt.Printf("size: %d\n", s.Size)
		fmt.Printf("link: %v\n", s.Link)
		fmt.Printf("info: %v\n", s.Info)
		fmt.Printf("addr align: 0x%x\n", s.AddressAlign)
		fmt.Printf("ent size: %v\n", s.EntSize)
		fmt.Println()
	}

	printHeader(f.Header)
	fmt.Printf("programm headers: %d\n", len(f.ProgramHeaders))
	for _, ph := range f.ProgramHeaders {
		printProgram(ph)
	}
	fmt.Println()
	fmt.Printf("section headers: %d\n", len(f.SectionHeaders))
	for _, sh := range f.SectionHeaders {
		printSection(sh)
	}
	fmt.Println()
}

func run1() error {
	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("usage: %s ELF-FILE", os.Args[0])
	}

	fileName := flag.Arg(0)

	elfFile, err := elf.Open(fileName)
	if err != nil {
		return err
	}

	fmt.Println("file header")
	fmt.Println(elfFile.FileHeader)

	fmt.Println("sections")
	for _, section := range elfFile.Sections {
		fmt.Println(section)
	}

	fmt.Println("progs")
	for _, prog := range elfFile.Progs {
		fmt.Println(prog)
	}

	return nil
}
