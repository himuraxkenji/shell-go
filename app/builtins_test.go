package main

import (
	"io"
	"os"
	"path/filepath"
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

func TestBuiltinPwd(t *testing.T) {
	t.Run("prints the current working directory", func(t *testing.T) {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		out := captureStdout(t, func() {
			builtinPwd(nil)
		})
		if got := strings.TrimSpace(out); got != wd {
			t.Errorf("builtinPwd output = %q, want %q", got, wd)
		}
	})

	t.Run("reflects the directory after changing it", func(t *testing.T) {
		tmpDir := t.TempDir()
		origWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("os.Chdir failed: %v", err)
		}
		defer os.Chdir(origWd)

		resolvedTmpDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}

		out := captureStdout(t, func() {
			builtinPwd(nil)
		})
		if got := strings.TrimSpace(out); got != resolvedTmpDir {
			t.Errorf("builtinPwd output = %q, want %q", got, resolvedTmpDir)
		}
	})
}

func TestBuiltinCd(t *testing.T) {
	t.Run("changes to a valid relative/absolute path", func(t *testing.T) {
		tmpDir := t.TempDir()
		origWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		defer os.Chdir(origWd)

		builtinCd([]string{tmpDir})

		got, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		want, err := filepath.EvalSymlinks(tmpDir)
		if err != nil {
			t.Fatalf("filepath.EvalSymlinks failed: %v", err)
		}
		if got != want {
			t.Errorf("cwd = %q, want %q", got, want)
		}
	})

	t.Run("nonexistent path prints an error and does not change cwd", func(t *testing.T) {
		origWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		defer os.Chdir(origWd)

		stderr := captureStderr(t, func() {
			builtinCd([]string{"/this/path/does/not/exist/1234567890"})
		})
		want := "cd: /this/path/does/not/exist/1234567890: No such file or directory"
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr = %q, want it to contain %q", stderr, want)
		}

		got, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		if got != origWd {
			t.Errorf("cwd = %q, want unchanged %q", got, origWd)
		}
	})

	t.Run("no args or ~ goes to home directory", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("os.UserHomeDir failed: %v", err)
		}
		origWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		defer os.Chdir(origWd)

		builtinCd(nil)

		got, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		want, err := filepath.EvalSymlinks(home)
		if err != nil {
			t.Fatalf("filepath.EvalSymlinks failed: %v", err)
		}
		if got != want {
			t.Errorf("cwd = %q, want %q", got, want)
		}
	})

	t.Run("~/subpath expands to a subdir under HOME", func(t *testing.T) {
		fakeHome := t.TempDir()
		subDir := filepath.Join(fakeHome, "sub")
		if err := os.Mkdir(subDir, 0o755); err != nil {
			t.Fatalf("os.Mkdir failed: %v", err)
		}
		t.Setenv("HOME", fakeHome)

		origWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		defer os.Chdir(origWd)

		builtinCd([]string{"~/sub"})

		got, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd failed: %v", err)
		}
		want, err := filepath.EvalSymlinks(subDir)
		if err != nil {
			t.Fatalf("filepath.EvalSymlinks failed: %v", err)
		}
		if got != want {
			t.Errorf("cwd = %q, want %q", got, want)
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
