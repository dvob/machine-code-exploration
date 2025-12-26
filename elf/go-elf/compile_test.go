package elf

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCompile(t *testing.T) {
	var virtualAddress uint64 = 0x401000
	entryPoint, code := Compile(virtualAddress)

	elfBinary := Write(virtualAddress, entryPoint, code)

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.elf")

	err := os.WriteFile(outputPath, elfBinary, 0755)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(outputPath)

	out, err := cmd.Output()
	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 34 {
			t.Fatalf("expected exit code 34 got %d: %s", exitErr.ExitCode(), err)
		}
	} else if err == nil {
		t.Fatal("expected error code 34")
	} else {
		t.Fatalf("other error returned: %s", err)
	}

	expectedOutput := "Hello World!\n"
	if !bytes.Equal(out, []byte(expectedOutput)) {
		t.Fatalf("expected output %q, got %q", expectedOutput, out)
	}
}
