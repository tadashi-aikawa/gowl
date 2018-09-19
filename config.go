package main

import (
	"github.com/BurntSushi/toml"
)

const configPath = ".gowlconfig"

// Service is information of Github, Bitbucket, and so on.
type Service struct {
	Token string
}

// Config configuration
type Config struct {
	Editor string
	GitHub Service
}

// CreateConfig creates configurations from .gowlconfig(toml)
func CreateConfig() (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
