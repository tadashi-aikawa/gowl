package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func toSelection(r Repository) string {
	return fmt.Sprintf("*%-5v %-45v %-10v %v", r.Star, r.FullName, r.Language, r.License)
}

func doRepositorySelection(handler IHandler) (Repository, error) {
	for {
		word := ""
		prompt := &survey.Input{
			Message: "Search word:",
		}
		survey.AskOne(prompt, &word, nil)
		if word == "" {
			return Repository{}, nil
		}

		repos, err := handler.SearchRepositories(word)
		if err != nil {
			panic(err)
		}

		var selections []string
		for _, r := range repos {
			selections = append(selections, toSelection(r))
		}

		var selection string
		selectedPrompt := &survey.Select{
			Message:  "Choose a repository you want to clone:",
			Options:  selections,
			PageSize: 15,
		}
		survey.AskOne(selectedPrompt, &selection, nil)

		if selection != "" {
			for _, r := range repos {
				if toSelection(r) == selection {
					return r, nil
				}
			}
		}
	}
}

func execCommandIO(workdir *string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if workdir != nil {
		cmd.Dir = *workdir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCommandStdout(workdir *string, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	if workdir != nil {
		cmd.Dir = *workdir
	}

	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "Fail to command: %v %v in %v", name, strings.Join(arg, " "), *workdir)
	}

	return string(out), nil
}

// CmdGet executes `get`
func CmdGet(handler IHandler) error {
	repo, err := doRepositorySelection(handler)
	if repo.FullName == "" {
		return nil
	}
	if err != nil {
		return err
	}

	dst := filepath.Join(os.Getenv("GOPATH"), "src", handler.GetPrefix(), repo.FullName)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		fmt.Printf("Clone %v to %v\n", repo.CloneURL, dst)
		if err := execCommandIO(nil, "git", "clone", repo.CloneURL, dst); err != nil {
			return errors.Wrap(err, "Fail to clone "+repo.CloneURL)
		}
	} else {
		fmt.Printf("Pull %v\n", dst)
		if err := execCommandIO(&dst, "git", "checkout", "master"); err != nil {
			return errors.Wrap(err, "Fail to checkout master in "+dst)
		}
		if err := execCommandIO(&dst, "git", "pull"); err != nil {
			return errors.Wrap(err, "Fail to pull in "+dst)
		}
	}

	return nil
}

// CmdEdit executes `edit`
func CmdEdit(handler IHandler, editor string) error {
	githubDir := filepath.Join(os.Getenv("GOPATH"), "src")
	repoDirs, err := listRepositories(githubDir)
	if err != nil {
		return errors.Wrap(err, "Fail to search repositories")
	}

	var selections []string
	selectedPrompt := &survey.MultiSelect{
		Message:  "Choose repositories you want to edit:",
		Options:  repoDirs,
		PageSize: 15,
	}
	survey.AskOne(selectedPrompt, &selections, nil)
	if len(selections) == 0 {
		return nil
	}

	if err := execCommandIO(nil, editor, selections...); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to edit repository %v", selections))
	}

	return nil
}

// CmdWeb executes `open`
func CmdWeb(handler IHandler, browser string) error {
	githubDir := filepath.Join(os.Getenv("GOPATH"), "src")
	repoDirs, err := listRepositories(githubDir)
	if err != nil {
		return errors.Wrap(err, "Fail to search repositories")
	}

	var selections []string
	selectedPrompt := &survey.MultiSelect{
		Message:  "Choose repositories you want to open:",
		Options:  repoDirs,
		PageSize: 15,
	}
	survey.AskOne(selectedPrompt, &selections, nil)
	if len(selections) == 0 {
		return nil
	}

	var urls []string
	for _, s := range selections {
		remoteURL, err := getCommandStdout(&s, "git", "config", "--get", "remote.origin.url")
		if err != nil {
			return errors.Wrap(err, "Fail to get remote origin URL")
		}

		urls = append(urls, remoteURL)
	}

	if err := execCommandIO(nil, browser, urls...); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to open repository %v", selections))
	}

	return nil
}
