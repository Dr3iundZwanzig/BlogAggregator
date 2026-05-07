package main

import (
	"fmt"
	"os"
)

func cli(st *state, commandsMap commands) error {
	userInput := os.Args
	if len(userInput) < 2 {
		return fmt.Errorf("not enough arguments provided")
	}
	newCommand := command{
		name:      userInput[1],
		arguments: userInput[2:],
	}
	err := commandsMap.run(st, newCommand)
	if err != nil {
		return err
	}
	return nil
}
