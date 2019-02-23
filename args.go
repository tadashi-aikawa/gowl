package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.7.0"
const usage = `Gowl.

Usage:
  gowl get [-f | --force] [-r | --recursive] [-s | --shallow] [-B | --bitbucket-server]
  gowl edit [-e=<editor> | --editor=<editor>]
  gowl web
  gowl list
  gowl purge
  gowl -h | --help
  gowl -V | --version

Options:
  -e --editor=<editor>        Use editor [default: default]
  -f --force                  Force remove and reclone if exists
  -r --recursive              Clone recursively
  -s --shallow                Use shallow clone
  -B --bitbucket-server       Use Bitbucket Server
  -h --help                   Show this screen.
  -V --version                Show version.
  `

// Args created by CLI args
type Args struct {
	CmdGet   bool `docopt:"get"`
	CmdEdit  bool `docopt:"edit"`
	CmdWeb   bool `docopt:"web"`
	CmdList  bool `docopt:"list"`
	CmdPurge bool `docopt:"purge"`

	Editor          string `docopt:"--editor"`
	Force           bool   `docopt:"--force"`
	Recursive       bool   `docopt:"--recursive"`
	Shallow         bool   `docopt:"--shallow"`
	BitbucketServer bool   `docopt:"--bitbucket-server"`
}

// CreateArgs creates Args
func CreateArgs(usage string, argv []string, version string) (Args, bool, error) {
	parser := &docopt.Parser{
		HelpHandler:  docopt.PrintHelpOnly,
		OptionsFirst: false,
	}

	opts, err := parser.ParseArgs(usage, argv, version)
	if err != nil {
		return Args{}, false, errors.Wrap(err, "Fail to parse arguments.")
	}

	if len(opts) == 0 {
		return Args{}, true, nil
	}

	var args Args
	opts.Bind(&args)

	return args, false, nil
}
