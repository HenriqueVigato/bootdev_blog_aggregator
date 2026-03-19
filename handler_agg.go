package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
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

	for _, item := range feed.Channel.Item {
		if err = savePost(s, item, feedsToFetch); err != nil {
			return fmt.Errorf("erro ao salvar o post no banco de dados: %v", err)
		}
	}

	return nil
}

func savePost(s *state, feedItem RSSItem, feed database.Feed) error {
	ctx := context.Background()
	tPublished := sql.NullTime{}
	t, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
	if err != nil {
		tPublished = sql.NullTime{
			Time:  t,
			Valid: true,
		}
	}

	post := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       feedItem.Title,
		Url:         feedItem.Link,
		Description: feedItem.Description,
		PublishedAt: tPublished,
		FeedID:      feed.ID,
	}

	_, err = s.db.CreatePost(ctx, post)
	if err != nil {
		if strings.Contains(err.Error(), "posts_url_key") {
			return nil
		}
		return fmt.Errorf("erro ao cadastrar o post no banco de dados: %v", err)
	}

	return nil
}

func browsePosts(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) < 1 {
		cmd.Args = []string{"2"}
	}
	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("erro ao converter o limit para int %v", err)
	}
	getPostsParams := database.GetPostForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	}

	posts, err := s.db.GetPostForUser(ctx, getPostsParams)
	if err != nil {
		return fmt.Errorf("erro ao buscar os posts no banco de dados %v", err)
	}
	for _, post := range posts {
		printPosts(post)
	}
	return nil
}

func printPosts(post database.Post) {
	fmt.Println("----------------------------------------")
	fmt.Printf("Title:       %s\n", post.Title)
	fmt.Printf("URL:         %s\n", post.Url)
	fmt.Printf("Description: %s\n", post.Description)
	fmt.Println("----------------------------------------")
}
