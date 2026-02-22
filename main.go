package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("err: ", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: cfg,
	}

	cmds := &commands{
		registercommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	if len(os.Args) < 2 {
		log.Fatal("Usage: gator <command> [args...]")
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal("erro com a chamada da funcionalidade: ", err)
	}
}
