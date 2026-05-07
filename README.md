# agent-credential-guard

Linux-first Guardrail gegen Secret-Leaks in AI-Agent- und Git-Workflows.

## Projektstatus

`v0.1` (fruehe, bewusst fokussierte Version).

## Was es bewusst ist (v0.1)

- schneller lokaler Guardrail fuer Commit-Workflows
- Scan auf staged Diff und `.env*`
- einfache, transparente Regelbasis

## Was es bewusst nicht ist (v0.1)

- kein Enterprise-DLP- oder Compliance-Scanner
- kein Full-History- oder Full-Repo-Scanner
- keine perfekte Secret-Erkennung ohne False Positives/Negatives

## Installation (ohne Go)

```bash
curl -fsSL https://raw.githubusercontent.com/gestalten-sass/agent-credential-guard/master/scripts/install.sh | bash
```

Der Installer nutzt bevorzugt `/usr/local/bin`. Falls das nicht moeglich ist, installiert er nach `~/.local/bin` und traegt den PATH dauerhaft in deine Shell-Config ein.

## Optional: global fuer alle Repos aktivieren

```bash
guard hook install --global
guard hook status --global
```

## Power-User (One-Liner inkl. Auto-Hook)

```bash
curl -fsSL https://raw.githubusercontent.com/gestalten-sass/agent-credential-guard/master/scripts/install.sh | GUARD_AUTO_HOOK=1 bash
```

## CLI

- `guard init`
- `guard scan`
- `guard scan --env`
- `guard scan --strict`
- `guard hook install`
- `guard hook remove`
- `guard hook install --global`
- `guard hook remove --global`
- `guard hook status --global`
- `guard version`

## Konfiguration

Optional im Repo-Root: `.guard.yaml`

- `enabled_rules`
- `ignore_paths`
- `strict_mode`

Beispiel: `.guard.yaml.example`

## Entwicklung (mit Go)

```bash
go test ./...
go build -o guard ./cmd/guard
```

## Beispiele

- `examples/staged-secret/README.md`
- `examples/env-secret/README.md`

## Changelog

- `CHANGELOG.md`
