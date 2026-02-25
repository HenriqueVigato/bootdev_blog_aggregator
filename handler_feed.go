package main

import "fmt"

func addFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("o handler addFeed espera o nome do feed e a URL")
	}
	return nil
}
