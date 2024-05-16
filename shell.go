package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// create a scanner to read input from user
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ ")
		scanner.Scan()

		// grab input from the user
		input := scanner.Text()

		tokens := strings.Fields(input)
		if execute(tokens, 0, 0, 0, []int{0}) == 0 {
			break
		}
	}
}

func execute(tokens []string, status int, background int, parallel int, pids []int) int {
	// always exit if the first token is exit
	if tokens[0] == "exit" {
		return 0
	}

	// handle changing directory since this is not in os/exec
	if tokens[0] == "cd" {
		if err := os.Chdir(tokens[1]); err != nil {
			fmt.Println("Error changing directory", err)
		}
		return 1
	}

	// handle running built-in commands and executables
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing command", err)
	}

	// for _, token := range tokens {
	// 	fmt.Println("Token: ", token)
	// }

	return 1
}
