package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const configPath = ".gowlconfig"

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

func main() {
	config, err := createConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failure to load `.config`."))
	}

	client := createGithubClient(config.GitHub.Token)

	repos, _, err := client.Search.Repositories(context.Background(), "miroir", nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failure to search repositories."))
	}

	for _, repo := range repos.Repositories {
		fmt.Printf(
			"★%-8v%-15v%-15v%-45v%-15v\n",
			repo.GetStargazersCount(), repo.GetID(), repo.GetLanguage(), repo.GetFullName(), repo.GetUpdatedAt(),
		)
	}
}
