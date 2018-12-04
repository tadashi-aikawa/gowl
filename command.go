package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-pipeline"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const gowlSite = "gowl.site"

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
	cmd.Stdin = os.Stdin
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
func selectLocalRepository(root string) (string, error) {
	repoRoot := filepath.Join(root)
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

func clone(url string, dst string, shallow bool, recursive bool) error {
	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth", "1")
	}
	if recursive {
		args = append(args, "--recursive")
	}
	args = append(args, url, dst)
	fmt.Printf("Exec: %v\n", strings.Join(args, " "))

	if err := execCommand(nil, "git", args...); err != nil {
		return errors.Wrap(err, "Fail to clone "+url)
	}

	return nil
}

func configureUser(dir string, name, mailAddress *string) error {
	if name != nil {
		fmt.Printf("Exec: git config user.name %v\n", *name)
		if err := execCommand(&dir, "git", "config", "user.name", *name); err != nil {
			return errors.Wrap(err, "Fail to config user.name "+*name)
		}
	}

	if mailAddress != nil {
		fmt.Printf("Exec: git config user.email %v\n", *mailAddress)
		if err := execCommand(&dir, "git", "config", "user.email", *mailAddress); err != nil {
			return errors.Wrap(err, "Fail to config user.email "+*mailAddress)
		}
	}

	return nil
}

// CmdGet executes `get`
func CmdGet(handler IHandler, root string, force bool, shallow bool, recursive bool) error {
	repo, err := doRepositorySelection(handler)
	if repo.FullName == "" {
		return nil
	}
	if err != nil {
		return err
	}

	var cloneURL string
	if handler.GetUseSSH() {
		cloneURL = repo.SSHCloneURL
	} else {
		cloneURL = repo.HTTPCloneURL
	}
	dst := filepath.Join(root, handler.GetPrefix(), repo.FullName)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		clone(cloneURL, dst, shallow, recursive)
		if err := execCommand(&dst, "git", "config", gowlSite, repo.SiteURL); err != nil {
			return errors.Wrap(err, "Fail to config gowl.site "+repo.SiteURL)
		}
		if handler.GetOverrideUser() {
			if err := configureUser(dst, handler.GetUserName(), handler.GetMailAddress()); err != nil {
				return errors.Wrap(err, "Fail to configure "+dst)
			}
		}
	} else {
		if force {
			fmt.Printf("Remove %v\n", dst)
			if err := os.RemoveAll(dst); err != nil {
				return errors.Wrap(err, "Fail to remove "+dst)
			}
			clone(cloneURL, dst, shallow, recursive)
			if err := execCommand(&dst, "git", "config", gowlSite, repo.SiteURL); err != nil {
				return errors.Wrap(err, "Fail to config gowl.site "+repo.SiteURL)
			}
			if handler.GetOverrideUser() {
				if err := configureUser(dst, handler.GetUserName(), handler.GetMailAddress()); err != nil {
					return errors.Wrap(err, "Fail to configure "+dst)
				}
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
func CmdList(handler IHandler, root string) error {
	repositories, err := listRepositories(root)
	if err != nil {
		return errors.Wrap(err, "Fail to search repository.")
	}

	fmt.Println(strings.Join(repositories, "\n"))
	return nil
}

// CmdEdit executes `edit`
func CmdEdit(handler IHandler, root string, editor string) error {
	selection, err := selectLocalRepository(root)
	if selection == "" {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select a repository.")
	}

	if err := execCommand(&selection, editor); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to edit repository %v", selection))
	}

	return nil
}

// CmdPurge purges...
func CmdPurge(handler IHandler, root string) error {
	selection, err := selectLocalRepository(root)
	if selection == "" {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select a repository.")
	}

	os.Chdir(selection)

	branches, err := pipeline.Output(
		[]string{"git", "branch", "--merged", "master"},
		[]string{"grep", "-vE", "'^\\*|master$|develop$'"},
	)

	if string(branches) == "" {
		print("No branches to purge.")
		return nil
	}

	print("Purge branches...\n" + string(branches))

	if _, err := pipeline.Output(
		[]string{"git", "branch", "--merged", "master"},
		[]string{"grep", "-vE", "'^\\*|master$|develop$'"},
		[]string{"xargs", "git", "branch", "-d"},
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to purge repositories %v", selection))
	}

	return nil
}

// CmdWeb executes `web`
func CmdWeb(handler IHandler, root string, browser string) error {
	selection, err := selectLocalRepository(root)
	if selection == "" {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Fail to select a repository.")
	}

	siteURL, err := getCommandStdout(&selection, "git", "config", "--get", "gowl.site")
	if siteURL == "" {
		return errors.Wrap(err, "Site URL is not set. [gowl.site in .git/config]")
	}
	if err != nil {
		return errors.Wrap(err, "Fail to get site URL")
	}

	if err := execCommand(nil, browser, siteURL); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Fail to open repository %v", selection))
	}

	return nil
}
