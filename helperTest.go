package main

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func setupTestState(t *testing.T) (*state, command, string) {
	db := setupTestDB(t)
	s := &state{
		db: db,
		cfg: &config.Config{
			DBURL:           "",
			CurrentUserName: "test_user",
		},
	}
	cmd := command{
		Name: "Test",
		Args: []string{},
	}
	setupTestUser(t, s)
	setupTestFeeds(t, s)
	setupTestFollowing(t, s)
	return s, cmd, setupTestEnv(t)
}

func setupTestUser(t *testing.T, s *state) {
	ctx := context.Background()
	_, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      s.cfg.CurrentUserName,
	})
	if err != nil {
		t.Fatalf("erro ao criar o usuario de test %v", err)
	}
}

func setupTestEnv(t *testing.T) string {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	home, _ := os.UserHomeDir()

	testPath := filepath.Join(home, ".gatorconfig.json")

	err := os.WriteFile(testPath, []byte(`{"db_url":"postgres://example"}`), 0o644)
	if err != nil {
		t.Fatalf("erro ao escrever no arquivo de teste %v", err)
	}

	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	return tmpDir
}

func setupTestDB(t *testing.T) *database.Queries {
	dbURL := "postgres://gator_user:gator_pass@localhost:5433/gator?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("erro ao conectar no banco de testes: %v", err)
	}

	dbQueries := database.New(db)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM feeds")
		db.Exec("DELETE FROM feed_follows")
		db.Close()
	})
	return dbQueries
}

func capturaOutput(fn func() error) (string, error) {
	old := os.Stdout
	defer func() {
		os.Stdout = old
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func setupTestFeeds(t *testing.T, s *state) {
	t.Helper()
	ctx := context.Background()
	mainUser, _ := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	extraUsers := []string{"user_one", "user_two", "user_three"}
	for _, name := range extraUsers {
		extraState := &state{
			db:  s.db,
			cfg: &config.Config{CurrentUserName: name},
		}
		setupTestUser(t, extraState)
	}

	extraUser := make(map[string]database.User)
	for _, name := range extraUsers {
		extraUser[name], _ = s.db.GetUser(ctx, name)
	}

	feeds := []database.CreateFeedParams{
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Hacker News",
			Url:       "https://hnrss.org/newest",
			UserID:    extraUser["user_one"].ID,
		},
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Lane's Blog",
			Url:       "https://www.wagslane.dev/index.xml",
			UserID:    extraUser["user_two"].ID,
		},
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Go Blog",
			Url:       "https://go.dev/blog/feed.atom",
			UserID:    extraUser["user_three"].ID,
		},
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Ohh the lion",
			Url:       "https://the_lion_on_the_savane.org.altm",
			UserID:    mainUser.ID,
		},
	}

	for _, f := range feeds {
		_, err := s.db.CreateFeed(ctx, f)
		if err != nil {
			t.Fatalf("erro ao criar feed '%s': %v", f.Name, err)
		}
	}
}

func setupTestFollowing(t *testing.T, s *state) {
	t.Helper()
	ctx := context.Background()
	mainUser, _ := s.db.GetUser(ctx, s.cfg.CurrentUserName)

	urls := []string{
		"https://the_lion_on_the_savane.org.altm",
		"https://www.wagslane.dev/index.xml",
	}

	for _, url := range urls {
		feed, _ := s.db.GetFeedByURL(ctx, url)
		s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    mainUser.ID,
			FeedID:    feed.ID,
		})
	}
}
