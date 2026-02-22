package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	_ = cmd
	ctx := context.Background()

	err := s.db.ResetUser(ctx)
	if err != nil {
		return fmt.Errorf("erro ao resetar usuarios: %v", err)
	}
	fmt.Println("Usuarios deletados com sucesso")
	return nil
}
