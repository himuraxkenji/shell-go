package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	EXIT = "exit"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input", err)
			os.Exit(1)
		}

		tokens, err := tokenize(strings.TrimSpace(line))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if len(tokens) == 0 {
			continue
		}

		name, args := tokens[0], tokens[1:]

		if name == EXIT {
			os.Exit(exitCode(args))
		}

		if fn, ok := builtins[name]; ok {
			fn(args)
			continue
		}

		if path, err := lookupPath(name); err == nil {
			if err := runExternal(path, name, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		fmt.Println(name + ": command not found")
	}

}
