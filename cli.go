package main

import (
	"fmt"
	"os"
)

func cli(st *state, commandsMap commands) error {
	userInput := os.Args
	userArguments := []string{}
	if len(userInput) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}
	if len(userInput) > 2 {
		userArguments = userInput[2:]
	}
	newCommand := command{
		name:      userInput[1],
		arguments: userArguments,
	}
	err := commandsMap.run(st, newCommand)
	if err != nil {
		return err
	}
	return nil
}
