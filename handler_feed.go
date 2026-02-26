package main

import (
	"context"
	"fmt"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func addFeed(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) < 2 {
		return fmt.Errorf("o handler addFeed espera o nome do feed e a URL")
	}
	getCurrentUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("erro ao buscar o usuario registrado: %v", err)
	}

	newFeed := &database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    getCurrentUser.ID,
	}

	_, err = s.db.CreateFeed(ctx, *newFeed)
	if err != nil {
		return err
	}
	return nil
}
