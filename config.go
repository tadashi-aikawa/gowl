package main

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// Service is information of Github, Bitbucket, and so on.
type Service struct {
	Token        *string
	UserName     *string
	Password     *string
	MailAddress  *string
	BaseURL      *string
	Prefix       *string
	UseSSH       bool
	OverrideUser bool
}

// Config configuration
type Config struct {
	Editors         map[string]string
	Browser         string
	Root            string
	GitHub          Service
	BitbucketServer Service
}

// CreateConfig creates configurations from .gowlconfig(toml)
func CreateConfig() (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return Config{}, errors.Wrap(err, "Home directory is not found.")
	}

	configPath := filepath.Join(home, ".gowlconfig")

	var conf Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
