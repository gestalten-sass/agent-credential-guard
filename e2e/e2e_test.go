package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var guardBin string

func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func TestMain(m *testing.M) {
	r := repoRoot()
	bin := filepath.Join(os.TempDir(), "guard-e2e-bin")
	build := exec.Command("go", "build", "-o", bin, "./cmd/guard")
	build.Dir = r
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		os.Exit(1)
	}
	guardBin = bin
	code := m.Run()
	_ = os.Remove(bin)
	os.Exit(code)
}

func runGuard(t *testing.T, cwd string, env []string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(guardBin, args...)
	if cwd == "" {
		cmd.Dir = repoRoot()
	} else {
		cmd.Dir = cwd
	}
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return string(out), 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return string(out), ee.ExitCode()
	}
	t.Fatalf("run failed: %v\n%s", err, string(out))
	return "", 1
}

func TestGlobalHookInstallStatusRemove(t *testing.T) {
	home := t.TempDir()
	env := []string{"HOME=" + home}

	out, code := runGuard(t, "", env, "hook", "status", "--global")
	if code == 0 {
		t.Fatalf("expected inactive global status, got success: %s", out)
	}

	out, code = runGuard(t, "", env, "hook", "install", "--global")
	if code != 0 {
		t.Fatalf("global install failed (%d): %s", code, out)
	}

	out, code = runGuard(t, "", env, "hook", "status", "--global")
	if code != 0 || !strings.Contains(out, "global aktiv") {
		t.Fatalf("expected global active, got (%d): %s", code, out)
	}

	out, code = runGuard(t, "", env, "hook", "remove", "--global")
	if code != 0 {
		t.Fatalf("global remove failed (%d): %s", code, out)
	}
}

func TestScanEnvHonorsIgnoreAndStrict(t *testing.T) {
	repo := t.TempDir()
	mustRun(t, repo, "git", "init")
	mustWrite(t, filepath.Join(repo, ".guard.yaml"), "ignore_paths:\n  - .env.ignore\nstrict_mode: true\n")
	mustWrite(t, filepath.Join(repo, ".env"), "API_KEY=EXAMPLE_E2E_KEY_123456\n")
	mustWrite(t, filepath.Join(repo, ".env.ignore"), "API_KEY=EXAMPLE_E2E_KEY_123456\n")

	out, code := runGuard(t, repo, nil, "scan", "--env")
	if code != 1 {
		t.Fatalf("expected strict exit 1, got %d: %s", code, out)
	}
	if !strings.Contains(out, ".env:1") {
		t.Fatalf("expected .env finding: %s", out)
	}
	if strings.Contains(out, ".env.ignore") {
		t.Fatalf("did not expect ignored file finding: %s", out)
	}
}

func mustRun(t *testing.T, cwd string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = cwd
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, string(out))
	}
}

func mustWrite(t *testing.T, p, content string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s failed: %v", p, err)
	}
}
