package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const EXIT = "exit"

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
		if command == EXIT {
			break
		}

		if strings.HasPrefix(command, "echo") {
			fmt.Fprintln(os.Stdout, command[5:])
		} else {
			fmt.Println(command + ": command not found")

		}

	}

}
