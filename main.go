package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func main() {
	args, err := CreateArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to create arguments."))
	}

	config, err := CreateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to load `.gowlconfig`."))
	}

	var handler Handler
	handler.Init(config.GitHub.Token)

	if args.Repo {
		if handler.SearchRepositories(args.Word) != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `repo`"))
		}
	} else {
		log.Fatal("Unexpected error.")
	}

}
