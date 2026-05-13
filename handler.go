package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error invalid time argument")
	}
	time_req, err := time.ParseDuration("5s")
	if err != nil {
		return fmt.Errorf("error invalid internal time argument")
	}
	if timeBetweenRequests < time_req {
		return fmt.Errorf("minimum time is 10s between requests")
	}

	fmt.Printf("Collectionf feeds every: %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
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
		LastFetchedAt: sql.NullTime{
			Valid: false,
		},
		Name:    cmd.arguments[0],
		Url:     cmd.arguments[1],
		UsersID: currentUser.ID,
	}
	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("error creating feed")
	}
	printFeed(newFeed)
	cmd.arguments[0] = newFeed.Url
	err = handlerFollow(s, cmd, currentUser)
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds")
	}
	for _, feed := range feeds {
		printFeed(feed)
		feedCreator, err := s.db.GetUserFromId(context.Background(), feed.UsersID)
		if err != nil {
			return fmt.Errorf("error getting user by id")
		}
		fmt.Printf("* Feed creator:	 %s\n", feedCreator.Name)
	}
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UsersID)
}

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting feed from url")
	}

	currentTime := time.Now().UTC()
	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UsersID:   currentUser.ID,
		FeedID:    feed.ID,
	}
	createdFeedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("error creating feed follow")
	}
	fmt.Println(createdFeedFollow.UserName)
	fmt.Println(createdFeedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, currentUser database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("error getting follows: %v", err)
	}
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting feed")
	}
	toDelete := database.DeleteFeedFollowParams{
		UsersID: currentUser.ID,
		FeedID:  feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), toDelete)
	if err != nil {
		return fmt.Errorf("error unfollowing feed")
	}
	fmt.Printf("Unfollowed: %v with user: %v\n", feed.Name, currentUser.Name)
	return nil
}

func handlerBrowse(s *state, cmd command, currentUser database.User) error {
	limit := 2
	if len(cmd.arguments) > 0 {
		input, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return fmt.Errorf("invalid limit input: %v", cmd.arguments[0])
		}
		limit = input
	}
	param := database.GetPostsForUserParams{
		UsersID: currentUser.ID,
		Limit:   int32(limit),
	}
	userPosts, err := s.db.GetPostsForUser(context.Background(), param)
	if err != nil {
		return fmt.Errorf("error getting user posts")
	}
	for _, p := range userPosts {
		fmt.Println(p.Title)
	}
	return nil
}
