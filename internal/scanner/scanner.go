package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"agent-credential-guard/internal/gitutil"
)

type Options struct {
	IncludeEnv bool
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
}

var rules = []rule{
	{name: "aws_access_key", re: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{name: "github_token", re: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9_]{20,}`)},
	{name: "private_key", re: regexp.MustCompile(`-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----`)},
	{name: "generic_api_key", re: regexp.MustCompile(`(?i)(api[_-]?key|token|secret)\s*[:=]\s*["']?[A-Za-z0-9_\-]{16,}`)},
}

func Run(opts Options) (Report, error) {
	report := Report{Findings: []Finding{}}

	diff, err := gitutil.StagedDiff()
	if err != nil {
		return report, err
	}

	report.Findings = append(report.Findings, scanText("staged-diff", diff)...)

	if opts.IncludeEnv {
		envFindings, err := scanEnvFiles()
		if err != nil {
			return report, err
		}
		report.Findings = append(report.Findings, envFindings...)
	}

	return dedupe(report), nil
}

func scanEnvFiles() ([]Finding, error) {
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
		info, err := os.Stat(p)
		if err != nil || info.IsDir() {
			continue
		}

		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(f)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			findings = append(findings, scanLine(filepath.Base(p), lineNo, line)...)
		}
		_ = f.Close()

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return findings, nil
}

func scanText(location, text string) []Finding {
	findings := make([]Finding, 0)
	for _, r := range rules {
		if r.re.FindStringIndex(text) != nil {
			findings = append(findings, Finding{Rule: r.name, Location: location})
		}
	}
	return findings
}

func scanLine(file string, lineNo int, line string) []Finding {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return nil
	}

	findings := make([]Finding, 0)
	for _, r := range rules {
		if r.re.FindStringIndex(line) != nil {
			findings = append(findings, Finding{
				Rule:     r.name,
				Location: fmt.Sprintf("%s:%d", file, lineNo),
			})
		}
	}
	return findings
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
