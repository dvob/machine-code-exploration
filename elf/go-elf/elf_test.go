package elf

import (
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
)

func Test_writeElf(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.elf")
	err := Write(nil, outputPath)
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
