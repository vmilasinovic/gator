package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// NOTE: CLEAN ARG LENGTH
func StartRepl(s *State, c *Commands) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Gator > ")
		scanner.Scan()
		cliInput := scanner.Text()
		cleanedInput := cleanInput(cliInput)

		if len(cleanedInput) < 2 && cleanedInput[0] != "exit" && cleanedInput[0] != "help" {
			continue
		}

		commandName := cleanedInput[0]
		args := cleanedInput[1:]

		if commandName == "help" {
			fmt.Println("")
			fmt.Println("Welcome to the Gator help menu!")
			fmt.Println("Here are your available commands:")

			for k, v := range c.Descriptions {
				fmt.Println("")
				fmt.Printf("  - %s: %s\n", k, v)
				fmt.Println("")
			}
		}

		if commandName == "exit" {
			os.Exit(0)
		}

		availableCommands := c.Get()
		_, ok := availableCommands[commandName]
		if !ok {
			if commandName != "help" {
				fmt.Println("Please enter a valid command")
			}
			continue
		}

		newCommand := Command{
			Name: commandName,
			Args: args,
		}

		if err := c.AllCommands[commandName](s, newCommand); err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(str string) []string {
	lowered := strings.ToLower(str)
	words := strings.Fields(lowered)
	return words
}
