package main

import (
	"os"
	"os/exec"
)

func lookupPath(name string) (string, error) {
	return exec.LookPath(name)
}

func runExternal(fullPath, name string, args []string) error {
	cmd := exec.Command(fullPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
