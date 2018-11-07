package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Repository used by gowl
type Repository struct {
	FullName     string
	SSHCloneURL  string
	HTTPCloneURL string
	Language     string
	License      string
	Star         int
}

func (r *Repository) fromGithub(gr *github.Repository) *Repository {
	return &Repository{
		FullName:     gr.GetFullName(),
		SSHCloneURL:  gr.GetSSHURL(),
		HTTPCloneURL: gr.GetCloneURL(),
		Language:     gr.GetLanguage(),
		License:      gr.GetLicense().GetName(),
		Star:         gr.GetStargazersCount(),
	}
}

func (r *Repository) fromBitbucketServer(bsr *BitbucketRepository) *Repository {
	// TODO: ssh option
	var httpURL string
	var sshURL string
	for _, x := range bsr.Links.Clone {
		switch x.Name {
		case "ssh":
			sshURL = x.Href
		case "http":
			httpURL = x.Href
		}
	}
	return &Repository{
		// Lower case for Bitbucket Server
		FullName:     bsr.GetFullName(),
		SSHCloneURL:  sshURL,
		HTTPCloneURL: httpURL,
		Language:     "UNKNOWN",
		License:      "No License",
		Star:         0,
	}
}

// IHandler is IF of Handler.
type IHandler interface {
	SearchRepositories(word string) ([]Repository, error)
	GetPrefix() string
	GetUseSSH() bool
}

// GitHubHandler handles a github command.
type GitHubHandler struct {
	client *github.Client
	prefix string
	useSSH bool
}

// BitbucketServerHandler handles a bitbucket command.
type BitbucketServerHandler struct {
	client *BitbucketClient
	prefix string
	useSSH bool
}

func createBitbucketClient(userName string, password string, baseURL string) *BitbucketClient {
	return &BitbucketClient{
		BaseURL:  baseURL,
		UserName: userName,
		Password: password,
	}
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

// NewGithubHandler creates Githubhandler
func NewGithubHandler(config Config) IHandler {
	return &GitHubHandler{
		client: createGithubClient(*config.GitHub.Token),
		prefix: "github.com",
		useSSH: config.GitHub.UseSSH,
	}
}

// NewBitbucketServerHandler creates BitbucketServerHandler
func NewBitbucketServerHandler(config Config) IHandler {
	return &BitbucketServerHandler{
		client: createBitbucketClient(*config.BitbucketServer.UserName, *config.BitbucketServer.Password, *config.BitbucketServer.BaseURL),
		prefix: *config.BitbucketServer.Prefix,
		useSSH: config.BitbucketServer.UseSSH,
	}
}

// GetPrefix gets prefix for local repository locastions
func (h *GitHubHandler) GetPrefix() string {
	return h.prefix
}

// GetPrefix gets prefix for local repository locastions
func (h *BitbucketServerHandler) GetPrefix() string {
	return h.prefix
}

// GetUseSSH gets whether use useSSH or not
func (h *GitHubHandler) GetUseSSH() bool {
	return h.useSSH
}

// GetUseSSH gets whether use useSSH or not
func (h *BitbucketServerHandler) GetUseSSH() bool {
	return h.useSSH
}

// SearchRepositories search repositories.
func (h *GitHubHandler) SearchRepositories(word string) ([]Repository, error) {
	res, _, err := h.client.Search.Repositories(context.Background(), word, &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
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
	res, err := h.client.searchRepositories(word)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to search repositories.")
	}

	var repos []Repository
	for _, bsrepo := range res.Values {
		var r Repository
		repos = append(repos, *r.fromBitbucketServer(&bsrepo))
	}

	return repos, nil
}
