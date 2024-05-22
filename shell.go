package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	// create a scanner to read input from user
	scanner := bufio.NewScanner(os.Stdin)

	// TODO: env variables

	for {
		fmt.Print("$ ")
		scanner.Scan()

		// grab input from the user
		input := scanner.Text()

		// convert input into tokens
		tokens := strings.Fields(input)

		// if no input, continue
		if len(tokens) == 0 {
			continue
		}

		if handleOperators(tokens) {
			break
		}
	}
}

func handleOperators(tokens []string) bool {
	var wg sync.WaitGroup

	for i, token := range tokens {
		switch token {
		case "&&&":
			wg.Add(1)
			go func(cmdTokens []string) {
				defer wg.Done()
				execute(cmdTokens)
			}(tokens[:i])
			if i < len(tokens)-1 {
				return handleOperators(tokens[i+1:])
			}
			wg.Wait()
			return false
		case "&":
			go execute(tokens[:i])
			if i < len(tokens)-1 {
				return handleOperators(tokens[i+1:])
			}
			return false
		case "&&":
			if !execute(tokens[:i]) {
				return handleOperators(tokens[i+1:])
			}
			return false
		case "||":
			if execute(tokens[:i]) {
				return handleOperators(tokens[i+1:])
			}
			return false
		}
	}
	return execute(tokens)
}

func execute(tokens []string) bool {
	// always exit if the first token is exit
	if tokens[0] == "exit" {
		return true
	}

	// handle changing directory since this is not in os/exec
	if tokens[0] == "cd" {
		if err := os.Chdir(tokens[1]); err != nil {
			fmt.Println("Error changing directory", err)
		}
		return false
	}

	// handle running built-in commands and executables
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing command", err)
	}

	return false
}
