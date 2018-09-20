package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/github"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func toSelection(r github.Repository) string {
	return fmt.Sprintf("*%-5v %-45v %-10v %v", r.GetStargazersCount(), r.GetFullName(), r.GetLanguage(), r.GetLicense().GetKey())
}

func loopClone(handler Handler) (github.Repository, error) {
	for {
		word := ""
		prompt := &survey.Input{
			Message: "Search word:",
		}
		survey.AskOne(prompt, &word, nil)
		if word == "" {
			return github.Repository{}, nil
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

// CmdClone
func CmdClone(handler Handler) error {
	repo, err := loopClone(handler)
	if repo.GetID() == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	dst := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", repo.GetFullName())

	fmt.Printf("Clone %v to %v\n", repo.GetCloneURL(), dst)
	cmd := exec.Command("git", "clone", repo.GetCloneURL(), dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Fail to clone "+repo.GetCloneURL())
	}

	return nil
}

// CmdEdit
func CmdEdit(handler Handler, editor string) error {
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

	err = exec.Command(editor, selections...).Run()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to edit repository %v", selections))
	}

	return nil
}
