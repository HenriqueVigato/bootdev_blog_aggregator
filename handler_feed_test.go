package main

import (
	"strings"
	"testing"
)

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

func TestGetFeeds_failure(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"falha", "informacao inutil"}

	_, err := capturaOutput(func() error {
		return getFeeds(s, cmd)
	})
	if err == nil {
		t.Errorf("era esperado uma mensagem de erro uma vez que a funcao foi chamada com argumentos")
	}
}

func TestGetFeeds_success(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	setupTestFeeds(t, s)

	output, err := capturaOutput(func() error {
		return getFeeds(s, cmd)
	})
	if err != nil {
		t.Logf("Erro inesperado: %v", err)
		t.Errorf("nao era esperado nenhum erro na chamada do capturaOutput")
	}

	if !strings.Contains(output, "Hacker News") {
		t.Logf("output: \n%v", output)
		t.Errorf("era esperado receber o feeds no output")
	}
}

func TestFollow_withoutArgs(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	_, err := capturaOutput(func() error {
		return follow(s, cmd)
	})
	if err == nil {
		t.Errorf("era esperado uma mensagem avisando que esta faltando argumentos na chamada")
	}
}
