package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.1.0"
const usage = `Gowl.

Usage:
  gowl get [--bitbucket-server] [-f]
  gowl edit <editor>
  gowl web
  gowl -h | --help
  gowl --version

Options:
  <editor>                    Use editor
  -B --bitbucket-server       Use Bitbucket Server
  -f --force                  Force remove and reclone if exists
  -h --help                   Show this screen.
  --version                   Show version.
  `

// Args created by CLI args
type Args struct {
	CmdGet  bool `docopt:"get"`
	CmdEdit bool `docopt:"edit"`
	CmdWeb  bool `docopt:"web"`

	Editor          string `docopt:"<editor>"`
	BitbucketServer bool   `docopt:"--bitbucket-server"`
	Force           bool   `docopt:"--force"`
}

// CreateArgs creates Args
func CreateArgs(usage string, argv []string, version string) (Args, error) {
	parser := &docopt.Parser{
		HelpHandler:  docopt.PrintHelpOnly,
		OptionsFirst: false,
	}

	opts, err := parser.ParseArgs(usage, argv, version)
	if err != nil {
		return Args{}, errors.Wrap(err, "Fail to parse arguments.")
	}

	var args Args
	opts.Bind(&args)

	return args, nil
}
