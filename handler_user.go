package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("O login espera apenas 1 argumento o 'User Name'")
		return fmt.Errorf("o login handler espera um argumento 'username'")
	}

	s.cfg.CurrentUserName = cmd.Args[0]
	fmt.Println("User has been set.")

	return nil
}
