package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("o login handler espera um argumento user name")
	}
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("o usuario deve ja estar criado para poder fazer login")
	}

	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("o register handler esepara um argumento")
	}
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, cmd.Args[0])
	if err == nil {
		return fmt.Errorf("o usuario '%s' ja existe", cmd.Args[0])
	}

	if err != sql.ErrNoRows {
		return err
	}

	newUser := &database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	user, err := s.db.CreateUser(ctx, *newUser)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("The user was creatted: ", user)
	return nil
}

func getAllUsers(s *state, cmd command) error {
	_ = cmd
	ctx := context.Background()
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar ao buscar todos os usuarios %v", err)
	}
	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil
}
