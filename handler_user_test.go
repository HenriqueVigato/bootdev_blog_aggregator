package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
)

func setupTestState() (*state, command) {
	return &state{
			cfg: &config.Config{
				CurrentUserName: "test_user",
			},
		}, command{
			Name: "Test",
			Args: []string{},
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

func TestHandlerLogin_emptyUserName(t *testing.T) {
	setupTestEnv(t)
	s, cmd := setupTestState()
	output, err := capturaOutput(func() error {
		return handlerLogin(s, cmd)
	})

	if err == nil {
		t.Logf("output: %s", output)
		t.Fatalf("esperava um erro avisando que nao pode string vazia")
	}
}

func TestHandlerLogin_ValidUserName(t *testing.T) {
	setupTestEnv(t)
	s, cmd := setupTestState()
	cmd.Args = []string{"HenriqueVigato"}

	output, err := capturaOutput(func() error {
		return handlerLogin(s, cmd)
	})
	if err != nil {
		t.Fatalf("erro %v:", err)
	}

	cfg, err := config.Read()
	if err != nil {
		t.Logf("s.cfg.CurrentUserName %s", cfg.CurrentUserName)
		t.Fatalf("CurrentUserName diferento do nome definido")
	}

	if !strings.Contains(output, "been set") {
		t.Logf("output: %v", output)
		t.Fatalf("Esperava a mensagem de sucesso")
	}
}
