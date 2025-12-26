package main

import (
	"bytes"
	"fmt"
)

type ELFReader struct {
	*ELFFile
	Data []byte
}

// readString reads a string from a string table
func (er *ELFReader) readString(stringTableOffset, stringIndex int) (string, error) {
	stringOffset := stringTableOffset + stringIndex
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

func (er *ELFReader) readStringTable(stringTableOffset, size int) ([][]byte, error) {
	lastByte := stringTableOffset + size
	if lastByte > (len(er.Data) - 1) {
		return nil, fmt.Errorf("string offset out of bounds")
	}

	data := er.Data[stringTableOffset : stringTableOffset+size]
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

	return er.readString(int(er.SectionHeaders[er.Header.SectionHeaderStringIndex].Offset), int(er.SectionHeaders[sectionIndex].Name))
}
