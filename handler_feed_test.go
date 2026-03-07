package main

import (
	"context"
	"strings"
	"testing"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
)

func TestAddFeed_emptyUser(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	_, err := capturaOutput(func() error {
		user, err := getUser(t, s)
		if err != nil {
			t.Errorf("%v", err)
		}
		return addFeed(s, cmd, user)
	})
	if err == nil {
		t.Error("era esperado uma mensagem de erro falando que esta faltando argumentos")
	}
}

func TestAddFeed_success(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"Hacker News RSS", "https://hnrss.org/newst"}
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = capturaOutput(func() error {
		return addFeed(s, cmd, user)
	})
	if err != nil {
		t.Logf("Err: %v", err)
		t.Errorf("Nao era esperado uma mensagem de erro")
	}
}

func TestAddFeed_duplicateUser(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"Hacker News RSS", "https://hnrss.org/newst"}
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = capturaOutput(func() error {
		return addFeed(s, cmd, user)
	})
	if err != nil {
		t.Logf("Err: %v", err)
		t.Errorf("nao era esperado um erro nesta etapa")
	}

	_, err = capturaOutput(func() error {
		return addFeed(s, cmd, user)
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
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = capturaOutput(func() error {
		return follow(s, cmd, user)
	})
	if err == nil {
		t.Errorf("era esperado uma mensagem avisando que esta faltando argumentos na chamada")
	}
}

func TestFollow_success(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"https://hnrss.org/newest"}
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("%v", err)
	}

	output, err := capturaOutput(func() error {
		return follow(s, cmd, user)
	})
	if err != nil {
		t.Logf("Erro inesperado: %v", err)
		t.Errorf("nao era esperado nenhum erro na chamada do capturaOutput")
	}

	if !strings.Contains(output, "Hacker News") {
		t.Errorf("era esperado o nome e o nome do dono do feed mas veio o output: \n%v", output)
	}
}

func TestFollowing_withArgs(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("%v", err)
	}
	cmd.Args = []string{"argumento"}

	err = following(s, cmd, user)

	if err == nil {
		t.Errorf("era esperado uma mensagem de erro avisando que a funcao following nao recebe nenhum argumento")
	}
}

func TestFollowing_success(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("%v", err)
	}

	output, err := capturaOutput(func() error {
		return following(s, cmd, user)
	})
	if err != nil {
		t.Errorf("erro na chamada da funcionalidade following: %v", err)
	}

	if !strings.Contains(output, "Lane's Blog") {
		t.Errorf("o output nao continha as informacoes esperadas. \nOutput: \n%v", output)
	}
}

func TestUnfollowing_insuficientArgument(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	user := database.User{}

	err := unfollow(s, cmd, user)
	if err == nil {
		t.Errorf("deveria conter uma mensagem que informa que o comando unfollow receber argumentos\n mas recebeu: \n%v", err)
	}
}

func TestUnfollowing_success(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	user, err := getUser(t, s)
	if err != nil {
		t.Errorf("erro ao buscar o usuario", err)
	}
	ctx := context.Background()
	followingFeeds := s.db.GetFeedFollowsForUser(ctx, user.ID)
}
