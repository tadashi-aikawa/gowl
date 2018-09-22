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

	var handler IHandler
	if args.BitbucketServer {
		handler = NewBitbucketServerHandler(config)
	} else {
		handler = NewGithubHandler(config)
	}

	switch true {
	case args.CmdGet:
		if err := CmdGet(handler, args.Force, args.Shallow); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `get`"))
		}
	case args.CmdEdit:
		if err := CmdEdit(handler, config.Editors[args.Editor]); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `edit`"))
		}
	case args.CmdWeb:
		if err := CmdWeb(handler, config.Browser); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `web`"))
		}
	case args.CmdList:
		if err := CmdList(handler); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `list`"))
		}
	}

}
