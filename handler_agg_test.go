package main

import (
	"context"
	"strings"
	"testing"
	"time"
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

	if !strings.Contains(output, "5-") {
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

	output, err := capturaOutput(func() error {
		return scrapeFeeds(s)
	})
	if err != nil {
		t.Fatalf("erro ao scrapeFeeds: %v", err)
	}

	if !strings.Contains(output, "5-") {
		t.Fatalf("esperava o titulo dos feeds itens mas recebeu: \n%v", output)
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

	posts, err := s.db.GetPostForUser(ctx, user.ID)

	if len(posts) > 0 {
		t.Fatalf("nao deveria ter nenhum post no banco ainda")
	}

	if err = savePost(s, fetchedFeed.Channel.Item[0], feed); err != nil {
		t.Fatalf("erro ao salvar os posts no banco: %v", err)
	}

	posts, err = s.db.GetPostForUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("erro ao buscar os posts no banco de dados")
	}

	if len(posts) < 1 {
		t.Fatalf("deveria conter o post que acabamos de salvar: %v", posts)
	}
}
