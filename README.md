# agent-credential-guard

Linux-first Guardrail gegen Secret-Leaks in AI-Agent- und Git-Workflows.

## Build

```bash
go build -o guard ./cmd/guard
```

## Installation

```bash
go install ./cmd/guard
# oder
./scripts/install.sh
```

Hinweis: Stelle sicher, dass `$(go env GOPATH)/bin` in deinem `PATH` liegt.

## Schnellstart

```bash
guard init
guard hook install
```

## CLI

- `guard init`
- `guard scan`
- `guard scan --env`
- `guard scan --strict`
- `guard hook install`
- `guard hook remove`
- `guard version`

## Konfiguration

Optional im Repo-Root: `.guard.yaml`

- `enabled_rules`
- `ignore_paths`
- `strict_mode`

Beispiel: `.guard.yaml.example`

## Beispiele

- `examples/staged-secret/README.md`
- `examples/env-secret/README.md`

## Changelog

- `CHANGELOG.md`
