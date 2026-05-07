package config

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	in := `enabled_rules:
  - "github_token"
  - private_key
ignore_paths:
  - .env.local
  - secrets/*.txt
strict_mode: true
`
	cfg, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if !cfg.StrictMode {
		t.Fatalf("expected strict_mode=true")
	}
	if len(cfg.EnabledRules) != 2 {
		t.Fatalf("expected 2 enabled rules, got %d", len(cfg.EnabledRules))
	}
	if cfg.EnabledRules[0] != "github_token" {
		t.Fatalf("unexpected first rule: %s", cfg.EnabledRules[0])
	}
	if len(cfg.IgnorePaths) != 2 {
		t.Fatalf("expected 2 ignore paths, got %d", len(cfg.IgnorePaths))
	}
}

func TestParseUnknownFieldFails(t *testing.T) {
	in := `strict_mode: false
unknown: value
`
	_, err := Parse(strings.NewReader(in))
	if err == nil {
		t.Fatalf("expected error for unknown field")
	}
}
