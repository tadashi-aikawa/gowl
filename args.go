package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.2.0-alpha"
const usage = `Gowl.

Usage:
  gowl get [-s | --shallow] [-f | --force] [-B | --bitbucket-server]
  gowl edit [-e=<editor> | --editor=<editor>]
  gowl web
  gowl list
  gowl -h | --help
  gowl --version

Options:
  -e --editor=<editor>        Use editor [default: default]
  -s --shallow                Use shallow clone
  -f --force                  Force remove and reclone if exists
  -B --bitbucket-server       Use Bitbucket Server
  -h --help                   Show this screen.
  --version                   Show version.
  `

// Args created by CLI args
type Args struct {
	CmdGet  bool `docopt:"get"`
	CmdEdit bool `docopt:"edit"`
	CmdWeb  bool `docopt:"web"`
	CmdList bool `docopt:"list"`

	Editor          string `docopt:"--editor"`
	BitbucketServer bool   `docopt:"--bitbucket-server"`
	Force           bool   `docopt:"--force"`
	Shallow         bool   `docopt:"--shallow"`
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
