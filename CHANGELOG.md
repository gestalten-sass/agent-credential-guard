# Changelog

## v0.1.3 - 2026-05-07

### Added
- Globaler Hook-Flow erweitert: `guard hook install --global`, `guard hook remove --global`, `guard hook status --global`.
- Installer-Option `GUARD_AUTO_HOOK=1` fuer One-Liner-Setup inkl. globaler Hook-Aktivierung.
- E2E-Tests fuer globale Hook-Flows und `scan --env` mit `strict_mode`/`ignore_paths`.

### Changed
- README auf klaren Standard-Install-Flow mit optionalem Global-Hook-Schritt umgestellt.
- Hook-CLI akzeptiert `--global` in der natuerlichen Position (`hook install --global`).

## v0.1.2 - 2026-05-07

### Added
- Robustes YAML-Parsing via `gopkg.in/yaml.v3` mit Validierung unbekannter Felder.
- Zusätzliche Tests fuer Config-Parsing und Scanner-Verhalten.

### Changed
- Release-/CI-Workflows fuer Node24 vorbereitet; Cache-Warnungen reduziert.

## v0.1.1 - 2026-05-07

### Added
- Erste GitHub-Release-Pipeline fuer Linux-Binaries (`amd64`, `arm64`) inkl. SHA256-Dateien.
- Installer fuer Binary-Download aus GitHub Releases.
- `guard init` zum Anlegen einer Startkonfiguration.

### Changed
- Hook verwendet absoluten Binary-Pfad fuer robusten Aufruf.

## v0.1.0 - 2026-05-07

### Added
- Go-CLI-Grundgeruest mit `scan`, `scan --env`, `scan --strict`, `hook install/remove`, `version`.
- Scanning fuer staged Git-Diffs und `.env*`.
- Regelset fuer typische Secret-Muster (AWS-Key, GitHub-Token, Private Keys, generische API-Keys).
- `.guard.yaml`-Konfiguration (`enabled_rules`, `ignore_paths`, `strict_mode`).
- Pre-commit Hook Installation/Entfernung mit verwaltetem Marker.
- CI-Workflow fuer `go test` und `go build`.

### Notes
- Default ist Warnmodus (kein Block), `--strict` oder `strict_mode: true` blockiert.
