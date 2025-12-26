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
	if er.Header.SectionHeaderCount < uint16(stringTableOffset) {
		return "", fmt.Errorf("invalid elf file: section header string index to high")
	}
	stringOffset := int(er.SectionHeaders[stringTableOffset].Offset) + stringIndex

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
