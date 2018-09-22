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

	if args.CmdGet {
		if err := CmdGet(handler, args.Force); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `get`"))
		}
	} else if args.CmdEdit {
		if err := CmdEdit(handler, config.Tools[args.Tool]); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `edit`"))
		}
	} else if args.CmdWeb {
		if err := CmdWeb(handler, config.Browser); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `web`"))
		}
	} else if args.CmdList {
		if err := CmdList(handler); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `cd`"))
		}
	}
}
