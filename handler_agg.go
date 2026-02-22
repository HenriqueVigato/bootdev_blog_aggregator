package main

import (
	"context"
	"fmt"
)

func agg(s *state, cmd command) error {
	_ = s
	_ = cmd
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil
}
