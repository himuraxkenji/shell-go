package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestLookupPath(t *testing.T) {
	t.Run("finds an existing binary in PATH", func(t *testing.T) {
		path, err := lookupPath("sh")
		if err != nil {
			t.Fatalf("lookupPath(\"sh\") unexpected error: %v", err)
		}
		if !strings.HasSuffix(path, "/sh") {
			t.Errorf("lookupPath(\"sh\") = %q, want path ending in /sh", path)
		}
	})

	t.Run("returns error for a nonexistent command", func(t *testing.T) {
		_, err := lookupPath("this_command_does_not_exist_1234567890")
		if err == nil {
			t.Fatal("lookupPath(nonexistent) expected error, got nil")
		}
	})
}

func TestRunExternal(t *testing.T) {
	shPath, err := lookupPath("sh")
	if err != nil {
		t.Skip("sh not available on PATH, skipping")
	}

	t.Run("executes the command and writes to stdout", func(t *testing.T) {
		out := captureStdout(t, func() {
			if err := runExternal(shPath, "sh", []string{"-c", "echo hello"}); err != nil {
				t.Fatalf("runExternal unexpected error: %v", err)
			}
		})
		if got := strings.TrimSpace(out); got != "hello" {
			t.Errorf("stdout = %q, want %q", got, "hello")
		}
	})

	t.Run("sets argv[0] to the typed name, not the resolved path", func(t *testing.T) {
		out := captureStdout(t, func() {
			if err := runExternal(shPath, "custom-name", []string{"-c", "echo $0"}); err != nil {
				t.Fatalf("runExternal unexpected error: %v", err)
			}
		})
		if got := strings.TrimSpace(out); got != "custom-name" {
			t.Errorf("argv[0] seen by child = %q, want %q", got, "custom-name")
		}
	})

	t.Run("returns error for a path that cannot be executed", func(t *testing.T) {
		err := runExternal("/nonexistent/path/to/binary", "binary", nil)
		if err == nil {
			t.Fatal("runExternal expected error for nonexistent path, got nil")
		}
	})
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	old := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading captured stdout failed: %v", err)
	}
	return string(out)
}
