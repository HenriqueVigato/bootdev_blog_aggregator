package main

import (
	"fmt"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("o login handler espera um argumento user name")
	}
	err := config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}
