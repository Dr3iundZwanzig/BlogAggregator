package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Dr3iundZwanzig/BlogAggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting user: %v", user)
		}
		return handler(s, cmd, user)
	}
}

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
