package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.2.0"
const usage = `Gowl.

Usage:
  gowl get [-f | --force] [-r | --recursive] [-s | --shallow] [-B | --bitbucket-server]
  gowl edit [-e=<editor> | --editor=<editor>]
  gowl web
  gowl list
  gowl -h | --help
  gowl --version

Options:
  -e --editor=<editor>        Use editor [default: default]
  -f --force                  Force remove and reclone if exists
  -r --recursive              Clone recursively
  -s --shallow                Use shallow clone
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
	Force           bool   `docopt:"--force"`
	Recursive       bool   `docopt:"--recursive"`
	Shallow         bool   `docopt:"--shallow"`
	BitbucketServer bool   `docopt:"--bitbucket-server"`
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
