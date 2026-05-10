package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Dr3iundZwanzig/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("User does not exist")
	}
	err = s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error setting user name")
	}
	fmt.Printf("User: %v has been set\n", cmd.arguments[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}
	currentTime := time.Now()
	User := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      cmd.arguments[0],
	}
	newDatabaseUser, err := s.db.CreateUser(context.Background(), User)
	if err != nil {
		return fmt.Errorf("error Register user: %v", err)
	}
	s.config.SetUser(newDatabaseUser.Name)
	fmt.Printf("User %v created\n", cmd.arguments[0])
	fmt.Println(newDatabaseUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting table")
	}
	fmt.Println("database resetted")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting users")
	}
	for _, user := range users {
		userName := user.Name
		if userName == s.config.CurrentUserName {
			fmt.Printf("* %v (current)\n", userName)
			continue
		}
		fmt.Printf("* %v\n", userName)
	}
	return nil
}
