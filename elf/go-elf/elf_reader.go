package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ELFReader struct {
	*ELFFile
	Data []byte
}

// readString reads a string from a string table
func (er *ELFReader) readString(sectionHeaderIndex int, stringIndex int) (string, error) {
	if sectionHeaderIndex > (len(er.SectionHeaders) - 1) {
		return "", fmt.Errorf("section header index too high")
	}

	sectionHeader := er.SectionHeaders[sectionHeaderIndex]

	if sectionHeader.Type != SHT_STRTAB {
		return "", fmt.Errorf("section header index %d is not of type string table", sectionHeaderIndex)
	}

	stringOffset := int(sectionHeader.Offset) + stringIndex
	if stringOffset > (len(er.Data) - 1) {
		return "", fmt.Errorf("string offset out of bounds")
	}

	strStart := er.Data[stringOffset:]
	before, _, ok := bytes.Cut(strStart, []byte{0x0})
	if !ok {
		return "", fmt.Errorf("invalid string table")
	}
	return string(before), nil
}

func (er *ELFReader) readStringTable(sectionHeaderIndex int) ([][]byte, error) {
	if sectionHeaderIndex > (len(er.SectionHeaders) - 1) {
		return nil, fmt.Errorf("section header index too high")
	}

	sectionHeader := er.SectionHeaders[sectionHeaderIndex]
	if sectionHeader.Type != SHT_STRTAB {
		return nil, fmt.Errorf("section header index %d is not of type string table", sectionHeaderIndex)
	}

	lastByte := int(sectionHeader.Offset) + int(sectionHeader.Size)
	if lastByte > (len(er.Data) - 1) {
		return nil, fmt.Errorf("string offset out of bounds")
	}

	data := er.Data[sectionHeader.Offset : sectionHeader.Offset+sectionHeader.Size]
	// string tables must start with a 0 byte
	if data[0] == 0x0 {
		data = data[1:]
	} else {
		return nil, fmt.Errorf("invalid string table does not start with null byte")
	}
	// string tables must end with a 0 byte
	if data[len(data)-1] == 0x0 {
		data = data[:len(data)-1]
	} else {
		return nil, fmt.Errorf("invalid string table does not end with null byte")
	}

	return bytes.Split(data, []byte{0x0}), nil
}

func (er *ELFReader) readSectionName(sectionIndex int) (string, error) {
	if er.Header.SectionHeaderStringIndex == 0 {
		return "", nil
	}

	if er.Header.SectionHeaderCount < er.Header.SectionHeaderStringIndex {
		return "", fmt.Errorf("invalid elf file: section header string index to high")
	}

	return er.readString(int(er.Header.SectionHeaderStringIndex), int(er.SectionHeaders[sectionIndex].Name))
}

func (er *ELFReader) sectionIndexByName(sectionName string) (int, bool) {
	for i := range er.SectionHeaders {
		currentSectionName, _ := er.readSectionName(i)
		if currentSectionName != "" && currentSectionName == sectionName {
			return i, true
		}
	}
	return 0, false
}

func (er *ELFReader) readSymbols() ([]Symbol64, error) {
	sectionName := ".symtab"
	index, ok := er.sectionIndexByName(sectionName)
	if !ok {
		return nil, fmt.Errorf("section %s not found", sectionName)
	}
	return er.readSymbolTable(index)
}

func (er *ELFReader) readDynSymbols() ([]Symbol64, error) {
	sectionName := ".dyntab"
	index, ok := er.sectionIndexByName(sectionName)
	if !ok {
		return nil, fmt.Errorf("section %s not found", sectionName)
	}
	return er.readSymbolTable(index)
}

func (er *ELFReader) readSymbolTable(sectionHeaderIndex int) ([]Symbol64, error) {

	if sectionHeaderIndex > (len(er.SectionHeaders) - 1) {
		return nil, fmt.Errorf("section header index too high")
	}

	sectionHeader := er.SectionHeaders[sectionHeaderIndex]

	if sectionHeader.Type != SHT_SYMTAB && sectionHeader.Type != SHT_DYNSYM {
		return nil, fmt.Errorf("section header index %d is not a symbol table", sectionHeaderIndex)
	}

	symbols := make([]Symbol64, sectionHeader.Size/sectionHeader.EntSize)

	data := bytes.NewBuffer(er.Data[sectionHeader.Offset : sectionHeader.Offset+sectionHeader.Size])

	err := binary.Read(data, binary.LittleEndian, &symbols)
	if err != nil {
		return nil, err
	}
	return symbols, err
}
