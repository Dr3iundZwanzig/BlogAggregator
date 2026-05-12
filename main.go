package main

import (
	"database/sql"
	"log"

	"github.com/Dr3iundZwanzig/BlogAggregator/internal/config"
	"github.com/Dr3iundZwanzig/BlogAggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", configStruct.DbURL)
	dbQueries := database.New(db)
	st := &state{
		config: &configStruct,
		db:     dbQueries,
	}

	commandsStruct := commands{
		commandMap: map[string]func(*state, command) error{},
	}
	commandsStruct.register("login", handlerLogin)
	commandsStruct.register("register", handlerRegister)
	commandsStruct.register("reset", handlerReset)
	commandsStruct.register("users", handlerUsers)
	commandsStruct.register("agg", handlerAgg)
	commandsStruct.register("addfeed", handlerAddFeed)
	commandsStruct.register("feeds", handlerFeeds)
	commandsStruct.register("follow", handlerFollow)
	commandsStruct.register("following", handlerFollowing)

	err = cli(st, commandsStruct)
	if err != nil {
		log.Fatal(err)
	}
}
