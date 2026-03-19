package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
)

func TestAgg(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"5s"}

	output, err := capturaOutput(func() error {
		go agg(s, cmd)
		time.Sleep(2 * time.Second)
		return nil
	})
	if err != nil {
		t.Fatalf("erro com a funcao agg: %v", err)
	}

	if strings.Contains(output, "5-") {
		t.Fatalf("era esperado alguns titulos: \n%v", output)
	}
}

func TestScrapeFeeds(t *testing.T) {
	s, _, _ := setupTestState(t)
	ctx := context.Background()

	feeds, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		t.Logf("erro ao buscar os feeds")
	}

	if feeds.LastFetchedAt.Valid {
		t.Fatalf("Nao era esperado nenhum valor alem de null")
	}

	_, err = capturaOutput(func() error {
		return scrapeFeeds(s)
	})
	if err != nil {
		t.Fatalf("erro ao scrapeFeeds: %v", err)
	}

	feeds, err = s.db.GetFeedByURL(ctx, feeds.Url)
	if err != nil {
		t.Logf("erro ao buscar os feeds")
	}

	if !feeds.LastFetchedAt.Valid {
		t.Fatalf("era esperado a data que o feed foi buscado: %v", feeds.LastFetchedAt)
	}
}

func TestSavePost(t *testing.T) {
	s, _, _ := setupTestState(t)
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		t.Fatalf("erro ao buscar o usuario no banco de dados: %v", err)
	}

	feed, err := s.db.GetFeedByURL(ctx, "https://techcrunch.com/feed/")
	if err != nil {
		t.Fatalf("Erro ao buscar os feeds para teste")
	}

	fetchedFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		t.Fatalf("erro ou buscar o feed: %v", err)
	}

	getPostsParams := database.GetPostForUserParams{
		ID:    user.ID,
		Limit: 2,
	}

	posts, err := s.db.GetPostForUser(ctx, getPostsParams)
	if err != nil {
		t.Fatalf("erro ao buscar posts por usuario %v", err)
	}

	if len(posts) > 0 {
		t.Fatalf("nao deveria ter nenhum post no banco ainda")
	}

	if err = savePost(s, fetchedFeed.Channel.Item[0], feed); err != nil {
		t.Fatalf("erro ao salvar os posts no banco: %v", err)
	}

	posts, err = s.db.GetPostForUser(ctx, getPostsParams)
	if err != nil {
		t.Fatalf("erro ao buscar os posts no banco de dados")
	}

	if len(posts) < 1 {
		t.Fatalf("deveria conter o post que acabamos de salvar: %v", posts)
	}

	if err = savePost(s, fetchedFeed.Channel.Item[0], feed); err != nil {
		t.Fatalf("nao deveria apresentar erro ao constar uma url duplicada: %v ", err)
	}
}

func TestBrowsePosts(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"3"}
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		t.Fatalf("erro o buscar o usuario no banco de dados ")
	}
	feed, err := s.db.GetFeedByURL(ctx, "https://techcrunch.com/feed/")
	if err != nil {
		t.Fatalf("Erro ao buscar os feeds para teste")
	}

	fetchedFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		t.Fatalf("erro ou buscar o feed: %v", err)
	}

	for _, itens := range fetchedFeed.Channel.Item {
		if err = savePost(s, itens, feed); err != nil {
			t.Fatalf("erro: %v", err)
		}
	}

	output, err := capturaOutput(func() error {
		return browsePosts(s, cmd, user)
	})
	if err != nil {
		t.Fatalf("erro: %v", err)
	}

	if !strings.Contains(output, "Title:") {
		t.Fatalf("deveria conter os posts: \n%v", output)
	}
}
