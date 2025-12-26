package elf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

func compile(cCode []byte, outputFile string) error {
	compiler, ok := os.LookupEnv("CC")
	if !ok {
		compiler = "gcc"
	}

	cmd := exec.Command(compiler, "-no-pie", "-x", "c", "-", "-o", outputFile)
	cmd.Stdin = bytes.NewBuffer(cCode)

	_, err := cmd.Output()
	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			return fmt.Errorf("compiler failed with exit code %d and stderr output %s", exitErr.ExitCode(), exitErr.Stderr)
		}
		return err
	}

	return nil
}

func Test_readSymbols(t *testing.T) {
	cCode, err := os.ReadFile("testdata/main.c")
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()

	outputFile := filepath.Join(tmpDir, "a.out")

	err = compile(cCode, outputFile)
	if err != nil {
		t.Fatal(err)
	}

	rawElfFile, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	elfFile, err := Read(rawElfFile)
	if err != nil {
		t.Fatal(err)
	}

	reader := &Reader{elfFile, rawElfFile}

	symbols, err := reader.readSymbols()
	if err != nil {
		t.Fatal(err)
	}

	// Get the .strtab section index for symbol name lookup
	strtabIndex, ok := reader.sectionIndexByName(".strtab")
	if !ok {
		t.Fatal("no .strtab section found")
	}

	// Collect all symbol names
	var symbolNames []string
	for _, sym := range symbols {
		name, err := reader.readString(strtabIndex, int(sym.Name))
		if err != nil {
			t.Fatalf("failed to read symbol name: %v", err)
		}
		symbolNames = append(symbolNames, name)
	}

	// Check for expected symbols
	expectedSymbols := []string{"counter", "main"}
	for _, expected := range expectedSymbols {
		found := slices.Contains(symbolNames, expected)
		if !found {
			t.Errorf("expected symbol %q not found in symbol table\nFound symbols: %v", expected, symbolNames)
		}
	}

	// Verify we found at least some symbols
	if len(symbols) == 0 {
		t.Error("no symbols found in .symtab")
	}
}
