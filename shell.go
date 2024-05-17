package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	term "github.com/nsf/termbox-go"
)

func main() {
	// initialize keyboard
	// err := term.Init()
	// if err != nil {
	// 	panic(err)
	// }
	// defer term.Close()

	// initialize history
	// history := []string{}
	// historyIndex := 0

	// create a scanner to read input from user
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ ")
		scanner.Scan()

		// grab input from the user
		input := scanner.Text()
		// input := readInputWithHistory(&history, &historyIndex, scanner)

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

func readInputWithHistory(history *[]string, historyIndex *int, scanner *bufio.Scanner) string {
	var input strings.Builder

	for {
		fmt.Println("1")
		ev := term.PollEvent()

		if ev.Type == term.EventError {
			panic(ev.Err)
		}

		if ev.Key == term.KeyEnter {
			fmt.Println()
			fmt.Println("2")
			break
		} else if ev.Key == term.KeyArrowUp {
			if *historyIndex > 0 {
				*historyIndex--
				input.Reset()
				input.WriteString((*history)[*historyIndex])
				fmt.Print("\r$ " + input.String())
			}
		} else if ev.Key == term.KeyArrowDown {
			if *historyIndex < len(*history)-1 {
				*historyIndex++
				input.Reset()
				input.WriteString((*history)[*historyIndex])
				fmt.Print("\r$ " + input.String())
			} else {
				*historyIndex = len(*history)
				input.Reset()
				fmt.Print("\r$ ")
			}
		} else if ev.Key == term.KeyBackspace || ev.Key == term.KeyBackspace2 {
			if input.Len() > 0 {
				curr := input.String()[:input.Len()-1]
				input.Reset()
				input.WriteString(curr)
				fmt.Print("\r$ " + input.String() + " ")
				fmt.Print("\r$ " + input.String())
			}
		} else {
			input.WriteRune(ev.Ch)
			fmt.Print(string(ev.Ch))
		}
	}
	fmt.Println("3")
	command := input.String()
	if command != "" {
		*history = append(*history, command)
		*historyIndex = len(*history)
	}
	return command
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

	// for _, token := range tokens {
	// 	fmt.Println("Token: ", token)
	// }

	return false
}
