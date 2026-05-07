# agent-credential-guard

Linux-first Guardrail gegen Secret-Leaks in AI-Agent- und Git-Workflows.

## Installation (ohne Go)

```bash
curl -fsSL https://raw.githubusercontent.com/gestalten-sass/agent-credential-guard/master/scripts/install.sh | bash
```

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
