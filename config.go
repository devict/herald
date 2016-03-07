package main

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the herald's configuration
type Config struct {
	SlackToken string `envconfig:"SLACK_TOKEN"`
	MongoURI   string `envconfig:"MONGOLAB_URI"`
	DestChan   string `envconfig:"DEST_CHAN"`
}

// NewConfig parses a Config from the environment.
func NewConfig() (Config, error) {
	var c Config

	if err := envconfig.Process("", &c); err != nil {
		return Config{}, err
	}

	if c.SlackToken == "" {
		return Config{}, errors.New("Missing env var SLACK_TOKEN")
	}
	if c.MongoURI == "" {
		return Config{}, errors.New("Missing env var MONGOLAB_URI")
	}
	if c.DestChan == "" {
		return Config{}, errors.New("Missing env var DEST_CHAN")
	}

	return c, nil
}
