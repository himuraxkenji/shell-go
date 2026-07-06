package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestExitCode(t *testing.T) {
	t.Run("no args defaults to 0", func(t *testing.T) {
		if got := exitCode(nil); got != 0 {
			t.Errorf("exitCode(nil) = %d, want 0", got)
		}
	})

	t.Run("parses a valid numeric arg", func(t *testing.T) {
		if got := exitCode([]string{"3"}); got != 3 {
			t.Errorf("exitCode([\"3\"]) = %d, want 3", got)
		}
	})

	t.Run("non-numeric arg returns 2 and prints an error", func(t *testing.T) {
		stderr := captureStderr(t, func() {
			if got := exitCode([]string{"abc"}); got != 2 {
				t.Errorf("exitCode([\"abc\"]) = %d, want 2", got)
			}
		})
		want := "exit: abc: numeric argument required"
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr = %q, want it to contain %q", stderr, want)
		}
	})
}

func TestBuiltinEcho(t *testing.T) {
	t.Run("joins args with spaces", func(t *testing.T) {
		out := captureStdout(t, func() {
			builtinEcho([]string{"hello", "world"})
		})
		if got := strings.TrimSpace(out); got != "hello world" {
			t.Errorf("builtinEcho output = %q, want %q", got, "hello world")
		}
	})

	t.Run("no args prints an empty line", func(t *testing.T) {
		out := captureStdout(t, func() {
			builtinEcho(nil)
		})
		if got := strings.TrimSpace(out); got != "" {
			t.Errorf("builtinEcho(nil) output = %q, want empty", got)
		}
	})
}

func TestBuiltinType(t *testing.T) {
	t.Run("reports a registered builtin", func(t *testing.T) {
		out := captureStdout(t, func() {
			builtinType([]string{"echo"})
		})
		want := "echo is a shell builtin"
		if got := strings.TrimSpace(out); got != want {
			t.Errorf("builtinType([\"echo\"]) = %q, want %q", got, want)
		}
	})

	t.Run("reports exit as a builtin even though it's not in the map", func(t *testing.T) {
		out := captureStdout(t, func() {
			builtinType([]string{"exit"})
		})
		want := "exit is a shell builtin"
		if got := strings.TrimSpace(out); got != want {
			t.Errorf("builtinType([\"exit\"]) = %q, want %q", got, want)
		}
	})

	t.Run("resolves an external command via PATH", func(t *testing.T) {
		path, err := lookupPath("sh")
		if err != nil {
			t.Skip("sh not available on PATH, skipping")
		}
		out := captureStdout(t, func() {
			builtinType([]string{"sh"})
		})
		want := "sh is " + path
		if got := strings.TrimSpace(out); got != want {
			t.Errorf("builtinType([\"sh\"]) = %q, want %q", got, want)
		}
	})

	t.Run("reports not found for an unknown command", func(t *testing.T) {
		out := captureStdout(t, func() {
			builtinType([]string{"this_command_does_not_exist_1234567890"})
		})
		want := "this_command_does_not_exist_1234567890: not found"
		if got := strings.TrimSpace(out); got != want {
			t.Errorf("builtinType(unknown) = %q, want %q", got, want)
		}
	})
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	old := os.Stderr
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading captured stderr failed: %v", err)
	}
	return string(out)
}
