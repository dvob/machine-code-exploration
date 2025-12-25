package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"testing"
)

func Test_writeElf(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.elf")
	err := writeElf(nil, outputPath)
	if err != nil {
		t.Fatal(err)
	}

	// stdOut := &bytes.Buffer{}
	// stdErr := &bytes.Buffer{}
	cmd := exec.Command(outputPath)
	// cmd.Stdout = stdOut
	// cmd.Stderr = stdErr

	_, err = cmd.Output()
	exitErr := &exec.ExitError{}
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 33 {
			t.Fatalf("expected exit code 33 got %d", exitErr.ExitCode())
		}
	} else if err == nil {
		t.Fatal("expected error code 33")
	} else {
		t.Fatalf("other error returned: %s", err)
	}
}
