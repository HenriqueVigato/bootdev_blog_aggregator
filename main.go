package main

import (
	"log"
	"os"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("err: ", err)
	}

	programState := &state{
		cfg: cfg,
	}

	cmds := &commands{
		registercommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: gator <command> [args...]")
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal("erro com a chamda da funcionalidado: ", err)
	}
}
