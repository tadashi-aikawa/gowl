package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Handler handles a command.
type Handler struct {
	client *github.Client
}

// IHandler is IF of Handler.
type IHandler interface {
	Init(token string)
	SearchRepositories(word string) error
}

func createGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// Init initializes.
func (h *Handler) Init(token string) {
	h.client = createGithubClient(token)
}

// SearchRepositories search and output repositories.
func (h *Handler) SearchRepositories(word string) error {
	repos, _, err := h.client.Search.Repositories(context.Background(), word, nil)
	if err != nil {
		return errors.Wrap(err, "Fail to search repositories.")
	}

	for _, repo := range repos.Repositories {
		fmt.Printf(
			"*%-8v%-15v%-15v%-60v%-15v\n",
			repo.GetStargazersCount(), repo.GetID(), repo.GetLanguage(), repo.GetFullName(), repo.GetUpdatedAt(),
		)
	}

	return nil
}
