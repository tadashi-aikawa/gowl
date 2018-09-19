package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

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
	CloneRepository(word string, seq int) error
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
	res, _, err := h.client.Search.Repositories(context.Background(), word, nil)
	if err != nil {
		return errors.Wrap(err, "Fail to search repositories.")
	}

	for i, repo := range res.Repositories {
		fmt.Printf(
			"%-5v*%-8v%-15v%-15v%-60v%-15v\n",
			i+1, repo.GetStargazersCount(), repo.GetID(), repo.GetLanguage(), repo.GetFullName(), repo.GetUpdatedAt(),
		)
	}

	return nil
}

// CloneRepository clones repository which first matched word.
func (h *Handler) CloneRepository(word string, seq int) error {
	res, _, err := h.client.Search.Repositories(context.Background(), word, nil)
	if err != nil {
		return errors.Wrap(err, "Fail to search repositories.")
	}

	// default value
	if seq == 0 {
		seq = 1
	}
	fullName := res.Repositories[seq-1].GetFullName()
	cloneURL := res.Repositories[seq-1].GetCloneURL()
	dst := path.Join(os.Getenv("GOPATH"), "src", "github.com", fullName)

	fmt.Printf("Clone %v to %v\n", cloneURL, dst)
	err = exec.Command("git", "clone", cloneURL, dst).Run()
	if err != nil {
		return errors.Wrap(err, "Fail to clone "+cloneURL)
	}

	return nil
}
