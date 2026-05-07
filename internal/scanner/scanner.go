package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"agent-credential-guard/internal/config"
	"agent-credential-guard/internal/gitutil"
)

type Options struct {
	IncludeEnv bool
	Config     config.Config
}

type Finding struct {
	Rule     string
	Location string
}

type Report struct {
	Findings []Finding
}

type rule struct {
	name string
	re   *regexp.Regexp
	hint string
}

var rules = []rule{
	{name: "aws_access_key", re: regexp.MustCompile(`AKIA[0-9A-Z]{16}`), hint: "AWS-Key sofort in AWS deaktivieren/rotieren und aus Git-Historie entfernen."},
	{name: "github_token", re: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9_]{20,}`), hint: "GitHub-Token sofort widerrufen und neu erzeugen."},
	{name: "private_key", re: regexp.MustCompile(`-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----`), hint: "Private Key sofort ersetzen und nie im Repo speichern."},
	{name: "generic_api_key", re: regexp.MustCompile(`(?i)(api[_-]?key|token|secret)\s*[:=]\s*["']?[A-Za-z0-9_\-]{16,}`), hint: "Secret in .env halten, Wert rotieren und Datei gitignoren."},
}

func Run(opts Options) (Report, error) {
	report := Report{Findings: []Finding{}}
	activeRules := selectedRules(opts.Config)

	files, err := gitutil.StagedFiles()
	if err != nil {
		return report, err
	}
	for _, f := range files {
		if ignored(f, opts.Config.IgnorePaths) {
			continue
		}
		d, err := gitutil.StagedDiffForFile(f)
		if err != nil {
			return report, err
		}
		report.Findings = append(report.Findings, scanText(activeRules, "staged-diff:"+f, d)...)
	}

	if opts.IncludeEnv {
		envFindings, err := scanEnvFiles(activeRules, opts.Config.IgnorePaths)
		if err != nil {
			return report, err
		}
		report.Findings = append(report.Findings, envFindings...)
	}

	return dedupe(report), nil
}

func HintForRule(ruleName string) string {
	for _, r := range rules {
		if r.name == ruleName {
			return r.hint
		}
	}
	return "Secret pruefen, rotieren und aus Versionierung entfernen."
}

func scanEnvFiles(activeRules []rule, ignorePaths []string) ([]Finding, error) {
	repoRoot, err := gitutil.RepoRoot()
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(repoRoot, ".env*"))
	if err != nil {
		return nil, err
	}

	findings := make([]Finding, 0)
	for _, p := range files {
		rel := filepath.Base(p)
		if ignored(rel, ignorePaths) {
			continue
		}
		info, err := os.Stat(p)
		if err != nil || info.IsDir() {
			continue
		}

		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}

		s := bufio.NewScanner(f)
		lineNo := 0
		for s.Scan() {
			lineNo++
			line := s.Text()
			findings = append(findings, scanLine(activeRules, rel, lineNo, line)...)
		}
		_ = f.Close()

		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	return findings, nil
}

func scanText(activeRules []rule, location, text string) []Finding {
	findings := make([]Finding, 0)
	for _, r := range activeRules {
		if r.re.FindStringIndex(text) != nil {
			findings = append(findings, Finding{Rule: r.name, Location: location})
		}
	}
	return findings
}

func scanLine(activeRules []rule, file string, lineNo int, line string) []Finding {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return nil
	}

	findings := make([]Finding, 0)
	for _, r := range activeRules {
		if r.re.FindStringIndex(line) != nil {
			findings = append(findings, Finding{Rule: r.name, Location: fmt.Sprintf("%s:%d", file, lineNo)})
		}
	}
	return findings
}

func selectedRules(cfg config.Config) []rule {
	if len(cfg.EnabledRules) == 0 {
		return rules
	}
	allowed := map[string]struct{}{}
	for _, name := range cfg.EnabledRules {
		allowed[name] = struct{}{}
	}
	res := make([]rule, 0)
	for _, r := range rules {
		if _, ok := allowed[r.name]; ok {
			res = append(res, r)
		}
	}
	return res
}

func ignored(path string, ignore []string) bool {
	for _, p := range ignore {
		if p == path {
			return true
		}
		if ok, _ := filepath.Match(p, path); ok {
			return true
		}
	}
	return false
}

func dedupe(report Report) Report {
	seen := make(map[string]struct{})
	res := Report{Findings: make([]Finding, 0, len(report.Findings))}
	for _, f := range report.Findings {
		k := f.Rule + "::" + f.Location
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		res.Findings = append(res.Findings, f)
	}
	return res
}
