package main

import (
	"log"

	"github.com/Dr3iundZwanzig/BlogAggregator/internal/config"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	st := &state{
		config: &configStruct,
	}
	commandsStruct := commands{
		commandMap: map[string]func(*state, command) error{},
	}
	err = commandsStruct.register("login", handlerLogin)
	if err != nil {
		log.Fatal(err)
	}
	err = cli(st, commandsStruct)
	if err != nil {
		log.Fatal(err)
	}
}
