# Changelog

## v0.1.0 - 2026-05-07

### Added
- Go-CLI-Grundgeruest mit `scan`, `scan --env`, `scan --strict`, `hook install/remove`, `init`, `version`.
- Scanning fuer staged Git-Diffs und `.env*`.
- Regelset fuer typische Secret-Muster (AWS-Key, GitHub-Token, Private Keys, generische API-Keys).
- `.guard.yaml`-Konfiguration (`enabled_rules`, `ignore_paths`, `strict_mode`).
- Pre-commit Hook Installation/Entfernung mit verwaltetem Marker.
- Path-sicherer Hook-Aufruf (absoluter Binary-Pfad).
- CI-Workflow fuer `go test` und `go build`.
- Erste Unit-Tests fuer Config-Parsing und Scanner-Helfer.

### Notes
- Default ist Warnmodus (kein Block), `--strict` oder `strict_mode: true` blockiert.
