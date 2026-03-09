package main

import (
	"context"
	"fmt"
	"time"
)

func agg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("se espera um intervalo de busca dos feeds")
	}
	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("tempo invalide: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", interval)

	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {
		if err = scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feedsToFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("GetNextFeedToFetch error: %v", err)
	}

	if err = s.db.MarkFeedFetched(ctx, feedsToFetch.ID); err != nil {
		return fmt.Errorf("marking feed fetched error: %v", err)
	}

	feed, err := fetchFeed(ctx, feedsToFetch.Url)
	if err != nil {
		return fmt.Errorf("error in getting the feed: %v", err)
	}

	for i, item := range feed.Channel.Item {
		fmt.Printf("%d- %s\n", i, item.Title)
	}

	return nil
}
