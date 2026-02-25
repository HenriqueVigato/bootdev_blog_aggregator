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
