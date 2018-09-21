package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Repository struct {
	FullName string
	CloneURL string
	Language string
	License  string
	Star     int
}

func (r *Repository) fromGithub(gr *github.Repository) *Repository {
	return &Repository{
		FullName: gr.GetFullName(),
		CloneURL: gr.GetCloneURL(),
		Language: gr.GetLanguage(),
		License:  gr.GetLicense().GetName(),
		Star:     gr.GetStargazersCount(),
	}
}

// IHandler is IF of Handler.
type IHandler interface {
	Init(token string, editor string)
	SearchRepositories(word string) ([]Repository, error)
}

// GitHubHandler handles a command.
type GitHubHandler struct {
	client *github.Client
	editor string
}

// GitHubHandler handles a command.
type BitbucketServerHandler struct {
	client *github.Client
	editor string
}

func createGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func listRepositories(dir string) ([]string, error) {
	var gitDirs []string
	err := filepath.Walk(dir, func(cpath string, info os.FileInfo, err error) error {
		if _, err := os.Stat(filepath.Join(cpath, ".git")); err == nil {
			gitDirs = append(gitDirs, cpath)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Fail to search repositories.")
	}

	return gitDirs, nil
}

// Init initializes.
func (h *GitHubHandler) Init(token string, editor string) {
	h.client = createGithubClient(token)
	h.editor = editor
}

func (h *BitbucketServerHandler) Init(token string, editor string) {
	h.client = createGithubClient(token)
	h.editor = editor
}

func NewGithubHandler() IHandler {
	return &GitHubHandler{}
}

func NewBitbucketServerHandler() IHandler {
	return &BitbucketServerHandler{}
}

// SearchRepositories search repositories.
func (h *GitHubHandler) SearchRepositories(word string) ([]Repository, error) {
	res, _, err := h.client.Search.Repositories(context.Background(), word, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to search repositories.")
	}

	var repos []Repository
	for _, ghrepo := range res.Repositories {
		var r Repository
		repos = append(repos, *r.fromGithub(&ghrepo))
	}

	return repos, nil
}

// SearchRepositories search repositories.
func (h *BitbucketServerHandler) SearchRepositories(word string) ([]Repository, error) {
	res, _, err := h.client.Search.Repositories(context.Background(), word, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to search repositories.")
	}

	fmt.Println("Bitbucketだよーーー")

	var repos []Repository
	for _, ghrepo := range res.Repositories {
		var r Repository
		repos = append(repos, *r.fromGithub(&ghrepo))
	}

	return repos, nil
}
