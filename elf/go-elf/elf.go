package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

func Read(data []byte) (*File, error) {
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

	file := &File{
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

func Write(virtualAddress uint64, entryPoint uint64, code []byte) []byte {
	programHeaders := []ProgramHeader64{
		{
			Type:           PT_LOAD,
			Flags:          PF_R | PF_X | PF_W,
			Offset:         0, // will be set below
			VirtualAddress: uint64(virtualAddress),
			FileSize:       uint64(len(code)),
			MemorySize:     uint64(len(code)),
			Align:          0x1000,
		},
	}

	headerSize := uint64(unsafe.Sizeof(Header64{}) + uintptr(len(programHeaders))*unsafe.Sizeof(ProgramHeader64{}))
	padding := (4096 - headerSize%4096) % 4096
	codeStart := headerSize + padding

	programHeaders[0].Offset = codeStart

	f := &File{
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
			Entry:                    uint64(entryPoint),
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
		ProgramHeaders: programHeaders,
		SectionHeaders: []SectionHeader64{},
	}

	byteOrder := binary.LittleEndian

	buf := &bytes.Buffer{}
	err := binary.Write(buf, byteOrder, f.Header)
	if err != nil {
		panic(err)
	}

	for _, programHeader := range f.ProgramHeaders {
		err = binary.Write(buf, byteOrder, programHeader)
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < int(padding); i++ {
		buf.WriteByte(0)
	}

	_, err = buf.Write(code)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func Print(f *Reader) error {
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
	printSection := func(index int, s SectionHeader64) error {
		// print name if shstrtab is available
		fmt.Printf("type: %v\n", s.Type)
		fmt.Printf("flags: %v\n", s.Flags)
		fmt.Printf("addr: 0x%x\n", s.Address)
		fmt.Printf("offset: 0x%x\n", s.Offset)
		fmt.Printf("size: %d\n", s.Size)
		fmt.Printf("link: %v\n", s.Link)
		fmt.Printf("info: %v\n", s.Info)
		fmt.Printf("addr align: 0x%x\n", s.AddressAlign)
		fmt.Printf("ent size: %v\n", s.EntSize)
		switch s.Type {
		case SHT_STRTAB:
			if s.Size == 0 {
				fmt.Printf("strings: no strings in table\n")
			}
			strings, err := f.readStringTable(index)
			fmt.Printf("strings: %v\n", s.EntSize)
			if err != nil {
				return err
			}
			for _, str := range strings {
				fmt.Printf("  - '%s'\n", str)
			}
		case SHT_DYNSYM, SHT_SYMTAB:
			symbols, err := f.readSymbolTable(index)
			if err != nil {
				return err
			}

			if len(symbols) == 0 {
				fmt.Println("symbols: no symbols")

			}
			fmt.Println("symbols:")
			for i, symbol := range symbols {
				symbolName, err := f.readString(int(s.Link), int(symbol.Name))
				if err != nil {
					return err
				}
				fmt.Printf("  - index: %d\n", i)
				fmt.Printf("    name: %s\n", symbolName)
				fmt.Printf("    type: %s\n", symbol.SymbolType())
				fmt.Printf("    value: %d\n", symbol.Value)
				fmt.Printf("    size: %d\n", symbol.Size)
				fmt.Printf("    visibility: %s\n", symbol.SymbolVisibility())
				fmt.Printf("    binding: %s\n", symbol.SymbolBinding())
				fmt.Printf("    section header index: %d\n", symbol.SectionHeaderIndex)
			}
		}
		fmt.Println()
		return nil
	}

	printHeader(f.Header)
	fmt.Printf("programm headers: %d\n", len(f.ProgramHeaders))
	for i, ph := range f.ProgramHeaders {
		fmt.Printf("index: %d\n", i)
		printProgram(ph)
	}
	fmt.Println()
	fmt.Printf("section headers: %d\n", len(f.SectionHeaders))
	for i, sh := range f.SectionHeaders {
		fmt.Printf("index: %d\n", i)
		sectionName, err := f.readSectionName(i)
		if err != nil {
			return err
		}
		if sectionName != "" {
			fmt.Printf("name: %v\n", sectionName)
		}
		err = printSection(i, sh)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}
