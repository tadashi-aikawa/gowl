package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/docopt/docopt-go"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const version = "0.1.0"
const usage = `Gowl.

Usage:
  gowl <word>
  gowl -h | --help
  gowl --version

Options:
  <word>        Search word for repository.
  -h --help     Show this screen.
  --version     Show version.
  `

const configPath = ".gowlconfig"

type args struct {
	Word string `docopt:"<word>"`
}
type service struct {
	Token string
}
type config struct {
	GitHub service
}

func createGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func createConfig() (config, error) {
	var conf config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return config{}, err
	}

	return conf, nil
}

func createArgs(usage string, argv []string, version string) (args, error) {
	parser := &docopt.Parser{
		HelpHandler:  docopt.PrintHelpOnly,
		OptionsFirst: false,
	}

	opts, err := parser.ParseArgs(usage, argv, version)
	if err != nil {
		return args{}, errors.Wrap(err, "Fail to parse arguments.")
	}

	var args args
	opts.Bind(&args)

	return args, nil
}

func main() {
	args, err := createArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to create arguments."))
	}

	config, err := createConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to load `.config`."))
	}

	client := createGithubClient(config.GitHub.Token)

	repos, _, err := client.Search.Repositories(context.Background(), args.Word, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to search repositories."))
	}

	for _, repo := range repos.Repositories {
		fmt.Printf(
			"★%-8v%-15v%-15v%-60v%-15v\n",
			repo.GetStargazersCount(), repo.GetID(), repo.GetLanguage(), repo.GetFullName(), repo.GetUpdatedAt(),
		)
	}
}
