package main

import (
	"fmt"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	document, err := config.Read()
	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	cmds := &commands{
		registercommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
}
