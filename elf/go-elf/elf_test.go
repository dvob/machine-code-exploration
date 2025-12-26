package elf

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
	var (
		virtualAddress uint64 = 0x401000
		code                  = []byte{
			0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, // mov $60, %rax
			0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00, // mov $33, %rdi
			0x0f, 0x05, // syscall
		}
	)
	elfBinary := Write(virtualAddress, virtualAddress, code)

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.elf")

	err := os.WriteFile(outputPath, elfBinary, 0755)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(outputPath)

	_, err = cmd.Output()
	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 33 {
			t.Fatalf("expected exit code 33 got %d: %s", exitErr.ExitCode(), err)
		}
	} else if err == nil {
		t.Fatal("expected error code 33")
	} else {
		t.Fatalf("other error returned: %s", err)
	}
}
