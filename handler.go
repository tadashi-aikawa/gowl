package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Handler handles a command.
type Handler struct {
	client *github.Client
	editor string
}

// IHandler is IF of Handler.
type IHandler interface {
	Init(token string, editor string)
	SearchRepositories(word string) error
	CloneRepository(word string, seq int) error
	ChangeDirectory(word string, seq int) error
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
func (h *Handler) Init(token string, editor string) {
	h.client = createGithubClient(token)
	h.editor = editor
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
	dst := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", fullName)

	fmt.Printf("Clone %v to %v\n", cloneURL, dst)
	err = exec.Command("git", "clone", cloneURL, dst).Run()
	if err != nil {
		return errors.Wrap(err, "Fail to clone "+cloneURL)
	}

	return nil
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

// EditRepository edit repository which first matched word by specified editor.
func (h *Handler) EditRepository(word string, seq int) error {
	githubDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com")
	repoDirs, err := listRepositories(githubDir)
	if err != nil {
		return errors.Wrap(err, "Fail to search repositories")
	}

	var filteredRepoDirs []string
	for _, d := range repoDirs {
		if strings.Contains(d, word) {
			filteredRepoDirs = append(filteredRepoDirs, d)
		}
	}

	// default value
	if seq == 0 {
		seq = 1
	}

	target := filteredRepoDirs[seq-1]
	fmt.Printf("Edit %v.", target)

	err = exec.Command(h.editor, target).Run()
	if err != nil {
		return errors.Wrap(err, "Fail to edit repository "+target)
	}

	return nil
}
