package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"agent-credential-guard/internal/gitutil"
	"gopkg.in/yaml.v3"
)

const ExampleYAML = `enabled_rules:
  - aws_access_key
  - github_token
  - private_key
  - generic_api_key
ignore_paths:
  - .env.local
  - secrets/*.txt
strict_mode: false
`

type Config struct {
	EnabledRules []string `yaml:"enabled_rules"`
	IgnorePaths  []string `yaml:"ignore_paths"`
	StrictMode   bool     `yaml:"strict_mode"`
}

func Load() (Config, error) {
	cfg := Config{}
	repoRoot, err := gitutil.RepoRoot()
	if err != nil {
		return cfg, err
	}

	p := filepath.Join(repoRoot, ".guard.yaml")
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	defer f.Close()

	return Parse(f)
}

func InitFile() (string, error) {
	repoRoot, err := gitutil.RepoRoot()
	if err != nil {
		return "", err
	}

	p := filepath.Join(repoRoot, ".guard.yaml")
	if _, err := os.Stat(p); err == nil {
		return "", fmt.Errorf("datei existiert bereits: %s", p)
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if err := os.WriteFile(p, []byte(ExampleYAML), 0o644); err != nil {
		return "", err
	}
	return p, nil
}

func Parse(r io.Reader) (Config, error) {
	cfg := Config{}
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
