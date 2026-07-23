package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var builtins map[string]func(args []string)

func init() {
	builtins = map[string]func(args []string){
		"echo": builtinEcho,
		"type": builtinType,
		"pwd":  builtinPwd,
		"cd":   builtinCd,
	}
}

func exitCode(args []string) int {
	if len(args) == 0 {
		return 0
	}
	code, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit: %s: numeric argument required\n", args[0])
		return 2
	}
	return code
}

func builtinEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func builtinPwd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		return
	}
	fmt.Println(dir)
}

func builtinCd(args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cd:", err)
		return
	}

	arg := "~"
	if len(args) > 0 {
		arg = args[0]
	}

	path := arg
	switch {
	case path == "~":
		path = home
	case strings.HasPrefix(path, "~/"):
		path = filepath.Join(home, path[2:])
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", arg)
	}
}

func builtinType(args []string) {
	for _, name := range args {
		if _, ok := builtins[name]; ok || name == "exit" {
			fmt.Println(name + " is a shell builtin")
			continue
		}
		if path, err := lookupPath(name); err == nil {
			fmt.Println(name + " is " + path)
			continue
		}
		fmt.Println(name + ": not found")
	}
}
