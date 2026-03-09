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
		t.Fatalf("era esperado o titulo Nobody Care mas recebeu: \n%v", output)
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

	if !strings.Contains(output, "10-") {
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
