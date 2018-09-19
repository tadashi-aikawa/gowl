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
	handler.Init(config.GitHub.Token, config.Editor)

	if args.Repo {
		if args.Clone {
			err := handler.CloneRepository(args.Word, args.Seq)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Fail to command `clone`"))
			}
		} else {
			err := handler.SearchRepositories(args.Word)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Fail to command `repo`"))
			}
		}
	} else if args.Edit {
		err := handler.EditRepository(args.Word, args.Seq)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `edit`"))
		}
	}
}
