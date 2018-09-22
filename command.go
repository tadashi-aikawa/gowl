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

func execCommand(workdir *string, name string, arg ...string) error {
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

// selectLocalRepository returns repository path.
func selectLocalRepository() (string, error) {
	repoRoot := filepath.Join(os.Getenv("GOPATH"), "src")
	repoDirs, err := listRepositories(repoRoot)
	if err != nil {
		return "", errors.Wrap(err, "Fail to search repositories")
	}

	var selection string
	selectedPrompt := &survey.Select{
		Message:  "Choose a repository:",
		Options:  repoDirs,
		PageSize: 15,
	}
	survey.AskOne(selectedPrompt, &selection, nil)

	if selection == "" {
		return "", nil
	}

	return selection, nil
}

// selectLocalRepositories returns repository paths.
func selectLocalRepositories() ([]string, error) {
	repoRoot := filepath.Join(os.Getenv("GOPATH"), "src")
	repoDirs, err := listRepositories(repoRoot)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to search repositories")
	}

	var selections []string
	selectedPrompt := &survey.MultiSelect{
		Message:  "Choose repositories:",
		Options:  repoDirs,
		PageSize: 15,
	}
	survey.AskOne(selectedPrompt, &selections, nil)
	if len(selections) == 0 {
		return nil, nil
	}

	return selections, nil
}

// CmdGet executes `get`
func CmdGet(handler IHandler, force bool) error {
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
		if err := execCommand(nil, "git", "clone", repo.CloneURL, dst); err != nil {
			return errors.Wrap(err, "Fail to clone "+repo.CloneURL)
		}
	} else {
		if force {
			fmt.Printf("Remove %v\n", dst)
			if err := os.RemoveAll(dst); err != nil {
				return errors.Wrap(err, "Fail to remove "+dst)
			}
			fmt.Printf("Clone %v to %v\n", repo.CloneURL, dst)
			if err := execCommand(nil, "git", "clone", repo.CloneURL, dst); err != nil {
				return errors.Wrap(err, "Fail to clone "+repo.CloneURL)
			}
		} else {
			fmt.Printf("Checkout master %v\n", dst)
			if err := execCommand(&dst, "git", "checkout", "master"); err != nil {
				return errors.Wrap(err, "Fail to checkout master in "+dst)
			}
			fmt.Printf("Pull %v\n", dst)
			if err := execCommand(&dst, "git", "pull"); err != nil {
				return errors.Wrap(err, "Fail to pull in "+dst)
			}
		}
	}

	return nil
}

// CmdList executes `open`
func CmdList(handler IHandler) error {
	selection, err := selectLocalRepository()
	if selection == "" {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select a repository.")
	}

	fmt.Println(selection)
	return nil
}

// CmdEdit executes `edit`
func CmdEdit(handler IHandler, tool string) error {
	selections, err := selectLocalRepositories()
	if selections == nil {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select repositories.")
	}

	if err := execCommand(nil, tool, selections...); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to edit repository %v", selections))
	}

	return nil
}

// CmdWeb executes `web`
func CmdWeb(handler IHandler, browser string) error {
	selections, err := selectLocalRepositories()
	if selections == nil {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select repositories.")
	}

	var urls []string
	for _, s := range selections {
		remoteURL, err := getCommandStdout(&s, "git", "config", "--get", "remote.origin.url")
		if err != nil {
			return errors.Wrap(err, "Fail to get remote origin URL")
		}

		urls = append(urls, remoteURL)
	}

	if err := execCommand(nil, browser, urls...); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to open repository %v", selections))
	}

	return nil
}
