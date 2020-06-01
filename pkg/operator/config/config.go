package config

import (
	"encoding/base64"
	"strings"
)

type Credentials struct {
	GithubAPIKey           string `yaml:"apiKey"`
	SlackToken             string `yaml:"slackToken"`
	SlackVerificationToken string `yaml:"slackVerificationToken"`
}

type Group []string

type OperatorConfig struct {
	Credentials Credentials `yaml:"credentials"`

	Groups map[string]Group `yaml:"groups"`

	// SlackChannel is a channel where the operator will post reports/etc.
	SlackChannel      string `yaml:"slackChannel"`
	SlackAdminChannel string `yaml:"slackAdminChannel"`

	CachePath string `yaml:"cachePath"`
}

// Anonymize makes a shallow copy of the config, suitable for dumping in logs (no sensitive data)
func (c *OperatorConfig) Anonymize() OperatorConfig {
	a := *c
	if key := a.Credentials.GithubAPIKey; len(key) > 0 {
		a.Credentials.GithubAPIKey = strings.Repeat("x", len(a.Credentials.DecodedAPIKey()))
	}
	if key := a.Credentials.SlackToken; len(key) > 0 {
		a.Credentials.SlackToken = strings.Repeat("x", len(a.Credentials.DecodedSlackToken()))
	}
	if key := a.Credentials.SlackVerificationToken; len(key) > 0 {
		a.Credentials.SlackVerificationToken = strings.Repeat("x", len(a.Credentials.DecodedSlackVerificationToken()))
	}
	return a
}

func decode(s string) string {
	if strings.HasPrefix(s, "base64:") {
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, "base64:"))
		if err != nil {
			return s
		}
		return string(data)
	}
	return s
}

func (b Credentials) DecodedAPIKey() string {
	return decode(b.GithubAPIKey)
}

func (b Credentials) DecodedSlackToken() string {
	return decode(b.SlackToken)
}

func (b Credentials) DecodedSlackVerificationToken() string {
	return decode(b.SlackVerificationToken)
}
