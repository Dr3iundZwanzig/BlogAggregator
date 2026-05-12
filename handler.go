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
	currentTime := time.Now().UTC()
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

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed")
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("need more arguments: name and url")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user")
	}
	currentTime := time.Now().UTC()
	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UsersID:   currentUser.ID,
	}
	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("error creating feed")
	}
	fmt.Println(newFeed)
	return nil
}
