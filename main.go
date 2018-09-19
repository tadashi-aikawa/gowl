package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func createGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
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

	client := createGithubClient(config.GitHub.Token)

	repos, _, err := client.Search.Repositories(context.Background(), args.Word, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to search repositories."))
	}

	for _, repo := range repos.Repositories {
		fmt.Printf(
			"*%-8v%-15v%-15v%-60v%-15v\n",
			repo.GetStargazersCount(), repo.GetID(), repo.GetLanguage(), repo.GetFullName(), repo.GetUpdatedAt(),
		)
	}
}
