package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"agent-credential-guard/internal/gitutil"
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
	EnabledRules []string
	IgnorePaths  []string
	StrictMode   bool
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
	s := bufio.NewScanner(r)
	section := ""
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			continue
		}
		if strings.HasPrefix(line, "strict_mode:") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "strict_mode:"))
			cfg.StrictMode = (v == "true")
			continue
		}
		if strings.HasPrefix(line, "- ") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			v = strings.Trim(v, "\"'")
			switch section {
			case "enabled_rules":
				cfg.EnabledRules = append(cfg.EnabledRules, v)
			case "ignore_paths":
				cfg.IgnorePaths = append(cfg.IgnorePaths, v)
			}
		}
	}

	if err := s.Err(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
