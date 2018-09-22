package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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

		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
		s.Color("fgHiGreen")
		s.Start()
		repos, err := handler.SearchRepositories(word)
		s.Stop()
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

func clone(url string, dst string, shallow bool) error {
	fmt.Printf("Clone %v to %v\n", url, dst)

	var args []string
	if shallow {
		args = []string{"clone", "--depth", "1", url, dst}
	} else {
		args = []string{"clone", url, dst}
	}

	if err := execCommand(nil, "git", args...); err != nil {
		return errors.Wrap(err, "Fail to clone "+url)
	}

	return nil
}

// CmdGet executes `get`
func CmdGet(handler IHandler, force bool, shallow bool) error {
	repo, err := doRepositorySelection(handler)
	if repo.FullName == "" {
		return nil
	}
	if err != nil {
		return err
	}

	dst := filepath.Join(os.Getenv("GOPATH"), "src", handler.GetPrefix(), repo.FullName)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		clone(repo.CloneURL, dst, shallow)
	} else {
		if force {
			fmt.Printf("Remove %v\n", dst)
			if err := os.RemoveAll(dst); err != nil {
				return errors.Wrap(err, "Fail to remove "+dst)
			}
			clone(repo.CloneURL, dst, shallow)
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
func CmdEdit(handler IHandler, editor string) error {
	selection, err := selectLocalRepository()
	if selection == "" {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select a repository.")
	}

	if err := execCommand(nil, editor, selection); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to edit repository %v", selection))
	}

	return nil
}

// CmdWeb executes `web`
func CmdWeb(handler IHandler, browser string) error {
	selection, err := selectLocalRepository()
	if selection == "" {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select a repository.")
	}

	remoteURL, err := getCommandStdout(&selection, "git", "config", "--get", "remote.origin.url")
	if err != nil {
		return errors.Wrap(err, "Fail to get remote origin URL")
	}

	if err := execCommand(nil, browser, remoteURL); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to open repository %v", selection))
	}

	return nil
}
