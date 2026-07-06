package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builtins map[string]func(args []string)

func init() {
	builtins = map[string]func(args []string){
		"echo": builtinEcho,
		"type": builtinType,
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
