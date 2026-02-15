package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registercommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.registercommands[cmd.Name]
	if !ok {
		return fmt.Errorf("comando %s nao esta retistrado", cmd.Name)
	}
	return function(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registercommands[name] = f
}
