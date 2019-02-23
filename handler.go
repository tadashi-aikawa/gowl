package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func contains(elms []string, elm string) bool {
	for _, e := range elms {
		if e == elm {
			return true
		}
	}

	return false
}

// Repository used by gowl
type Repository struct {
	FullName     string
	SSHCloneURL  string
	HTTPCloneURL string
	SiteURL      string
	Language     string
	License      string
	Star         int
}

func (r *Repository) fromGithub(gr *github.Repository) *Repository {
	return &Repository{
		FullName:     gr.GetFullName(),
		SSHCloneURL:  gr.GetSSHURL(),
		HTTPCloneURL: gr.GetCloneURL(),
		SiteURL:      gr.GetHTMLURL(),
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
		SiteURL:      bsr.Links.Self[0].Href,
		Language:     "UNKNOWN",
		License:      "No License",
		Star:         0,
	}
}

// IHandler is IF of Handler.
type IHandler interface {
	SearchRepositories(word string) ([]Repository, error)
	GetPrefix() string
	GetUserName() *string
	GetMailAddress() *string
	GetUseSSH() bool
	GetOverrideUser() bool
}

// GitHubHandler handles a github command.
type GitHubHandler struct {
	client       *github.Client
	prefix       string
	userName     *string
	mailAddress  *string
	useSSH       bool
	overrideUser bool
}

// BitbucketServerHandler handles a bitbucket command.
type BitbucketServerHandler struct {
	client       *BitbucketClient
	prefix       string
	userName     *string
	mailAddress  *string
	useSSH       bool
	overrideUser bool
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

func listRepositories(dirs []string) ([]string, error) {
	var gitDirs []string
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(cpath string, info os.FileInfo, err error) error {
			if contains([]string{"node_modules", "_build", "build", "dist"}, filepath.Base(cpath)) {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(cpath, ".git")); err == nil {
				gitDirs = append(gitDirs, cpath)
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrap(err, "Fail to search repositories.")
		}
	}

	return gitDirs, nil
}

// NewGithubHandler creates Githubhandler
func NewGithubHandler(config Config) IHandler {
	return &GitHubHandler{
		client:       createGithubClient(*config.GitHub.Token),
		prefix:       "github.com",
		userName:     config.GitHub.UserName,
		mailAddress:  config.GitHub.MailAddress,
		useSSH:       config.GitHub.UseSSH,
		overrideUser: config.GitHub.OverrideUser,
	}
}

// NewBitbucketServerHandler creates BitbucketServerHandler
func NewBitbucketServerHandler(config Config) IHandler {
	return &BitbucketServerHandler{
		client:       createBitbucketClient(*config.BitbucketServer.UserName, *config.BitbucketServer.Password, *config.BitbucketServer.BaseURL),
		prefix:       *config.BitbucketServer.Prefix,
		userName:     config.BitbucketServer.UserName,
		mailAddress:  config.BitbucketServer.MailAddress,
		useSSH:       config.BitbucketServer.UseSSH,
		overrideUser: config.BitbucketServer.OverrideUser,
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

// GetUserName gets user name
func (h *GitHubHandler) GetUserName() *string {
	return h.userName
}

// GetUserName gets user name
func (h *BitbucketServerHandler) GetUserName() *string {
	return h.userName
}

// GetMailAddress gets user mail address
func (h *GitHubHandler) GetMailAddress() *string {
	return h.mailAddress
}

// GetMailAddress gets user mail address
func (h *BitbucketServerHandler) GetMailAddress() *string {
	return h.mailAddress
}

// GetUseSSH gets whether use useSSH or not
func (h *GitHubHandler) GetUseSSH() bool {
	return h.useSSH
}

// GetUseSSH gets whether use useSSH or not
func (h *BitbucketServerHandler) GetUseSSH() bool {
	return h.useSSH
}

// GetOverrideUser returns OverrideUser
func (h *GitHubHandler) GetOverrideUser() bool {
	return h.overrideUser
}

// GetOverrideUser returns OverrideUser
func (h *BitbucketServerHandler) GetOverrideUser() bool {
	return h.overrideUser
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
