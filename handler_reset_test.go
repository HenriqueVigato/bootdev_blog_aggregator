package main

import (
	"context"
	"testing"
)

func TestHandlerReset(t *testing.T) {
	s, cmd, _ := setupTestState(t)

	cmd.Name = "register"
	cmd.Args = []string{"test_user_delete"}

	err := handlerRegister(s, cmd)
	if err != nil {
		t.Fatalf("nao era esperado nenhum erro aqui %v", err)
	}

	user, err := s.db.GetUser(context.Background(), "test_user_delete")
	if err != nil {
		t.Logf("user: %v", user)
		t.Fatalf("nao era esperado nenhum erro agoria na hora de ver se o usuario teste foi criado")
	}

	err = handlerReset(s, cmd)
	if err != nil {
		t.Logf("resetUser err: %v", err)
		t.Fatalf("nao era esperado nenhum erro aqui")
	}

	_, err = s.db.GetUser(context.Background(), "test_user_delete")
	if err == nil {
		t.Logf("err %v", err)
		t.Fatalf("era esperado um erro pois o usuario foi excluido")
	}
}
