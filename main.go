package main

import (
	"fmt"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	document, _ := config.Read()
	fmt.Printf("%+v\n", document)
	_ = config.SetUser("HenriqueVigato")

	document, _ = config.Read()
	fmt.Printf("%+v\n", document)
}
