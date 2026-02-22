package main

import (
	"context"
	"testing"
)

func TestFetchFeed_success(t *testing.T) {
	ctx := context.Background()
	response, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		t.Fatalf("erro ao bucar o feed %v", err)
	}

	if response.Channel.Title != "Lane's Blog" {
		t.Fatalf("channel title (%s) diferente do esperado: Lane's Blog", response.Channel.Title)
	}
}

func TestFetchFeed_fail(t *testing.T) {
	ctx := context.Background()
	_, err := fetchFeed(ctx, "https://www.wagslane/index.xml")
	if err == nil {
		t.Fatalf("esperava encontra um erro aqui uma vez que a url esta errada")
	}
}
