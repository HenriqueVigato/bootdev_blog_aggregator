package main

import (
	"context"
	"fmt"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("getUser error: %v", err)
		}
		return handler(s, cmd, user)
	}
}
