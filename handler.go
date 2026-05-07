package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argument for command found")
	}
	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error setting user name")
	}
	println("user has been set")
	return nil
}
