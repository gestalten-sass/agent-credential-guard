package scanner

import (
	"strings"
	"testing"

	"agent-credential-guard/internal/config"
)

func TestSelectedRules(t *testing.T) {
	cfg := config.Config{EnabledRules: []string{"github_token"}}
	rs := selectedRules(cfg)
	if len(rs) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rs))
	}
	if rs[0].name != "github_token" {
		t.Fatalf("unexpected rule: %s", rs[0].name)
	}
}

func TestIgnored(t *testing.T) {
	if !ignored(".env.local", []string{".env.local"}) {
		t.Fatalf("expected direct ignore match")
	}
	if !ignored("secrets/dev.txt", []string{"secrets/*.txt"}) {
		t.Fatalf("expected wildcard ignore match")
	}
	if ignored("README.md", []string{"secrets/*.txt"}) {
		t.Fatalf("did not expect ignore match")
	}
}

func TestScanLine(t *testing.T) {
	rs := selectedRules(config.Config{})
	findings := scanLine(rs, ".env", 1, "API_KEY=1234567890abcdef")
	if len(findings) == 0 {
		t.Fatalf("expected finding")
	}
}

func TestHintForRule(t *testing.T) {
	h := HintForRule("github_token")
	if !strings.Contains(strings.ToLower(h), "token") {
		t.Fatalf("unexpected hint: %s", h)
	}
}
