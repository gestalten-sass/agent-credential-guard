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

func StagedFiles() ([]string, error) {
	out, err := runGit("diff", "--staged", "--name-only")
	if err != nil {
		return nil, fmt.Errorf("staged dateien nicht lesbar: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	res := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		res = append(res, l)
	}
	return res, nil
}

func StagedDiffForFile(path string) (string, error) {
	out, err := runGit("diff", "--staged", "--no-color", "--", path)
	if err != nil {
		return "", fmt.Errorf("staged diff fuer %s nicht lesbar: %w", path, err)
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

func GetGlobalHooksPath() (string, error) {
	out, err := runGit("config", "--global", "--get", "core.hooksPath")
	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func SetGlobalHooksPath(path string) error {
	_, err := runGit("config", "--global", "core.hooksPath", path)
	if err != nil {
		return err
	}
	return nil
}

func UnsetGlobalHooksPath() error {
	_, err := runGit("config", "--global", "--unset", "core.hooksPath")
	if err != nil {
		if strings.Contains(err.Error(), "exit status 5") || strings.Contains(err.Error(), "exit status 1") {
			return nil
		}
		return err
	}
	return nil
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
