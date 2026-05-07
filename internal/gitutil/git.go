package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RepoRoot() (string, error) {
	out, err := runGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("repo-root nicht ermittelbar: %w", err)
	}
	return strings.TrimSpace(out), nil
}

func StagedDiff() (string, error) {
	out, err := runGit("diff", "--staged", "--no-color")
	if err != nil {
		return "", fmt.Errorf("staged diff nicht lesbar: %w", err)
	}
	return out, nil
}

func HookPath() (string, error) {
	out, err := runGit("rev-parse", "--git-path", "hooks/pre-commit")
	if err != nil {
		return "", fmt.Errorf("hook-pfad nicht ermittelbar: %w", err)
	}
	return strings.TrimSpace(out), nil
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}

	return stdout.String(), nil
}
