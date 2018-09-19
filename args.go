package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.1.0"
const usage = `Gowl.

Usage:
  gowl repo <word>
  gowl -h | --help
  gowl --version

Options:
  <word>        Search word for repository.
  -h --help     Show this screen.
  --version     Show version.
  `

// Args created by CLI args
type Args struct {
	Repo bool   `docopt:"repo"`
	Word string `docopt:"<word>"`
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
