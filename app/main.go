package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
)

var builtinCommands = map[string]struct{}{
	EXIT: {},
	ECHO: {},
	TYPE: {},
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input", err)
			os.Exit(1)
		}
		command = strings.TrimSpace(command)
		if strings.HasPrefix(command, TYPE) {
			if _, ok := builtinCommands[command[5:]]; ok {
				fmt.Fprintln(os.Stdout, command[5:]+" is a shell builtin")
			} else {
				fmt.Fprintln(os.Stderr, command[5:]+": not found")
			}

		} else if command == EXIT {
			break
		} else if strings.HasPrefix(command, ECHO) {
			fmt.Fprintln(os.Stdout, command[5:])
		} else {
			fmt.Println(command + ": command not found")

		}

	}

}
