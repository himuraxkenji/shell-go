package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("$ ")
	if scanner.Scan() {
		command := scanner.Text()
		fmt.Printf("%s: command not found", command)

	}

}
