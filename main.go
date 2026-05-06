package main

import (
	"os"

	"github.com/Dr3iundZwanzig/BlogAggregator/internal/config"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		println(err)
		println("error reading")
		os.Exit(0)
	}

	err = configStruct.SetUser("Sven")
	if err != nil {
		println(err)
		println("error setting user")
		os.Exit(0)
	}

	configStruct, err = config.Read()
	if err != nil {
		println(err)
		println("error reading")
		os.Exit(0)
	}
	println("ConfigData:")
	println(*configStruct.CurrentUserName)
	println(configStruct.DbURL)
}
