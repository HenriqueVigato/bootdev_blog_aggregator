package main

import (
	"context"
	"fmt"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func createAndReturnFeedFollow(s *state, user database.User, feed database.Feed) (database.CreateFeedFollowRow, error) {
	ctx := context.Background()
	feedFollow := &database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	followFeed, err := s.db.CreateFeedFollow(ctx, *feedFollow)
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("createFollow error: %v", err)
	}
	return followFeed, nil
}

func addFeed(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) < 2 {
		return fmt.Errorf("o handler addFeed espera o nome do feed e a URL")
	}

	newFeed := &database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	createdFeed, err := s.db.CreateFeed(ctx, *newFeed)
	if err != nil {
		return err
	}
	_, err = createAndReturnFeedFollow(s, user, createdFeed)
	if err != nil {
		return fmt.Errorf("erro o criar o follow feed %v", err)
	}

	return nil
}

func getFeeds(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) > 0 {
		return fmt.Errorf("a funcao getFeeds nao espera nenhum argumento")
	}
	feeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar os feeds: %v", err)
	}

	for _, feed := range feeds {
		currentUser, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("erro ao buscar o nome do usuario %v", err)
		}
		fmt.Printf("Feed name: %s;\n Feed URL: %s;\n UserName: %s;\n", feed.Name, feed.Url, currentUser.Name)
	}

	return nil
}

func follow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) < 1 {
		return fmt.Errorf("se espera a url a ser seguida")
	}

	feed, err := s.db.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("getFeed error: %v", err)
	}

	followFeed, err := createAndReturnFeedFollow(s, user, feed)
	if err != nil {
		return fmt.Errorf("erro o criar o follow feed %v", err)
	}
	fmt.Printf("Follow feed: %s, User: %s ", followFeed.FeedName, followFeed.UserName)
	return nil
}

func following(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("o comando following nao recebe argumentos")
	}
	ctx := context.Background()
	followedFeed, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("erro na busca do followedFeed %v", err)
	}
	for _, feed := range followedFeed {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}
