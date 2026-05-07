package hook

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"agent-credential-guard/internal/gitutil"
)

const marker = "# agent-credential-guard managed hook"

func Install() error {
	hookPath, err := gitutil.HookPath()
	if err != nil {
		return err
	}

	existing, err := os.ReadFile(hookPath)
	if err == nil && len(existing) > 0 && !strings.Contains(string(existing), marker) {
		return fmt.Errorf("bestehender pre-commit hook wird nicht ueberschrieben: %s", hookPath)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	content := strings.Join([]string{
		"#!/usr/bin/env bash",
		marker,
		"set -euo pipefail",
		"guard scan --env",
		"",
	}, "\n")

	if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(hookPath, []byte(content), 0o755); err != nil {
		return err
	}

	return nil
}

func Remove() (bool, error) {
	hookPath, err := gitutil.HookPath()
	if err != nil {
		return false, err
	}

	content, err := os.ReadFile(hookPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if !strings.Contains(string(content), marker) {
		return false, nil
	}

	if err := os.Remove(hookPath); err != nil {
		return false, err
	}

	return true, nil
}
