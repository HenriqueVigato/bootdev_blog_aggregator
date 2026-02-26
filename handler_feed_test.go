package main

import "testing"

func TestAddFeed_emptyUser(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	_, err := capturaOutput(func() error {
		return addFeed(s, cmd)
	})
	if err == nil {
		t.Error("era esperado uma mensagem de erro falando que esta faltando argumentos")
	}
}

func TestAddFeed_success(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"Hacker News RSS", "https://hnrss.org/newst"}

	_, err := capturaOutput(func() error {
		return addFeed(s, cmd)
	})
	if err != nil {
		t.Logf("Err: %v", err)
		t.Errorf("Nao era esperado uma mensagem de erro")
	}
}

func TestAddFeed_duplicateUser(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"Hacker News RSS", "https://hnrss.org/newst"}

	_, err := capturaOutput(func() error {
		return addFeed(s, cmd)
	})
	if err != nil {
		t.Logf("Err: %v", err)
		t.Errorf("nao era esperado um erro nesta etapa")
	}

	_, err = capturaOutput(func() error {
		return addFeed(s, cmd)
	})
	if err == nil {
		t.Logf("Err: %v", err)
		t.Errorf("era esperado um erro pois o feed esta duplicado")
	}
}
