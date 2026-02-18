package main

import (
	"fmt"
	"strings"
	"testing"
)

func newTestCommands(t *testing.T) *commands {
	t.Helper()

	cmd := &commands{
		registercommands: make(map[string]func(*state, command) error),
	}

	t.Cleanup(func() {
		cmd.registercommands = nil
	})
	return cmd
}

func TestRun_fail(t *testing.T) {
	s, c, _ := setupTestState(t)
	cmds := newTestCommands(t)
	c.Name = "test"

	err := cmds.run(s, c)
	if err == nil {
		t.Logf("command: %v", err)
		t.Errorf("deveria retornar um erro informando que o comando nao existe")
	}
}

func TestRun_success(t *testing.T) {
	s, c, _ := setupTestState(t)
	cmds := newTestCommands(t)
	c.Name = "test"
	cmds.registercommands[c.Name] = func(s *state, c command) error {
		return fmt.Errorf("state: %s, command name: %s", s.cfg.CurrentUserName, c.Name)
	}

	err := cmds.run(s, c)

	if !strings.Contains(err.Error(), "command name: test") {
		t.Logf("erro: %v", err)
		t.Errorf("deveria ter retornado o resultado da funcao")
	}
}

func TestRegister(t *testing.T) {
	_, c, _ := setupTestState(t)
	cmds := newTestCommands(t)

	c.Name = "test register"

	if len(cmds.registercommands) > 0 {
		t.Errorf("Ja existe comandos registrados, sendo que deveria estar limpo")
	}

	cmds.register(c.Name, func(*state, command) error {
		return nil
	})

	_, ok := cmds.registercommands[c.Name]

	if !ok {
		t.Errorf("Deveria ter cadastrado a funcao")
	}
}
