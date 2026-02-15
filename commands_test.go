package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
)

func setupTestState() (*State, *Command) {
	return &State{
			cfg: &config.Config{
				CurrentUserName: "test_user",
			},
		}, &Command{
			Name: "Test",
			Args: []string{},
		}
}

func captruaOutput(fn func() error) (string, error) {
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
	s, cmd := setupTestState()
	output, err := captruaOutput(func() error {
		return handlerLogin(s, cmd)
	})
	if err != nil {
		t.Fatalf("erro com o handlerLogin: %v", err)
	}

	if !strings.Contains(output, "username") {
		t.Logf("output: %s", output)
		t.Fatalf("esperava um erro avisando que nao pode string vazia")
	}
}

func TestHandlerLogin_ValidUserName(t *testing.T) {
	s, cmd := setupTestState()
	cmd.Args = []string{"HenriqueVigato"}

	output, err := captruaOutput(func() error {
		return handlerLogin(s, cmd)
	})
	if err != nil {
		t.Fatalf("erro %v:", err)
	}

	if s.cfg.CurrentUserName != "HenriqueVigato" {
		t.Logf("s.cfg.CurrentUserName %s", s.cfg.CurrentUserName)
		t.Fatalf("CurrentUserName diferento do nome definido")
	}

	if !strings.Contains(output, "user has been set") {
		t.Logf("output: %v", output)
		t.Fatalf("Esperava a mensagem de sucesso")
	}
}
