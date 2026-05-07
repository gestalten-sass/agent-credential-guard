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
	return installAt(hookPath)
}

func Remove() (bool, error) {
	hookPath, err := gitutil.HookPath()
	if err != nil {
		return false, err
	}
	return removeAt(hookPath)
}

func InstallGlobal() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	managedHooksPath := filepath.Join(home, ".config", "guard", "hooks")

	current, err := gitutil.GetGlobalHooksPath()
	if err != nil {
		return "", err
	}
	if current != "" && current != managedHooksPath {
		return "", fmt.Errorf("globaler core.hooksPath bereits gesetzt: %s", current)
	}

	if err := gitutil.SetGlobalHooksPath(managedHooksPath); err != nil {
		return "", err
	}
	if err := installAt(filepath.Join(managedHooksPath, "pre-commit")); err != nil {
		return "", err
	}

	return managedHooksPath, nil
}

func RemoveGlobal() (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	managedHooksPath := filepath.Join(home, ".config", "guard", "hooks")
	preCommit := filepath.Join(managedHooksPath, "pre-commit")

	removed, err := removeAt(preCommit)
	if err != nil {
		return false, err
	}

	current, err := gitutil.GetGlobalHooksPath()
	if err != nil {
		return removed, err
	}
	if current == managedHooksPath {
		if err := gitutil.UnsetGlobalHooksPath(); err != nil {
			return removed, err
		}
	}

	return removed, nil
}

func GlobalStatus() (bool, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, "", err
	}
	managedHooksPath := filepath.Join(home, ".config", "guard", "hooks")
	preCommit := filepath.Join(managedHooksPath, "pre-commit")

	current, err := gitutil.GetGlobalHooksPath()
	if err != nil {
		return false, "", err
	}
	if current != managedHooksPath {
		return false, current, nil
	}

	b, err := os.ReadFile(preCommit)
	if errors.Is(err, os.ErrNotExist) {
		return false, current, nil
	}
	if err != nil {
		return false, current, err
	}

	return strings.Contains(string(b), marker), current, nil
}

func installAt(hookPath string) error {
	existing, err := os.ReadFile(hookPath)
	if err == nil && len(existing) > 0 && !strings.Contains(string(existing), marker) {
		return fmt.Errorf("bestehender pre-commit hook wird nicht ueberschrieben: %s", hookPath)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	content := strings.Join([]string{
		"#!/usr/bin/env bash",
		marker,
		"set -euo pipefail",
		fmt.Sprintf("\"%s\" scan --env", exe),
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

func removeAt(hookPath string) (bool, error) {
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
