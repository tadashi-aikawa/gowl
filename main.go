package main

import (
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func getRoot(config Config) string {
	if root := config.Root; root != "" {
		return root
	}

	if path := os.Getenv("GOPATH"); path != "" {
		return filepath.Join(path, "src")
	}

	home, err := homedir.Dir()
	if err == nil {
		return filepath.Join(home, ".gowlroot")
	}

	panic("Unexpected ERROR!!")
}

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
		if err := CmdGet(handler, getRoot(config), args.Force, args.Shallow); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `get`"))
		}
	case args.CmdEdit:
		if err := CmdEdit(handler, getRoot(config), config.Editors[args.Editor]); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `edit`"))
		}
	case args.CmdWeb:
		if err := CmdWeb(handler, getRoot(config), config.Browser); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `web`"))
		}
	case args.CmdList:
		if err := CmdList(handler, getRoot(config)); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `list`"))
		}
	}

}
