# Projektplan v1: Linux Agent Credential Guard

Stand: 2026-05-07

## Zusammenfassung

Ziel ist ein Linux-first Open-Source-Tool, das Secret-Leaks in AI-Agent- und CLI-Workflows lokal vor dem Commit abfängt. Das Projekt wird bewusst nicht als allgemeiner Enterprise-Scanner gestartet, sondern als extrem einfaches Schutzwerkzeug für den Alltag in lokalen Git-Repositories.

## Positionierung

- Fokus: Agent-Workflow-Fokus (MCP, `.env`, staged Git-Diffs)
- Lizenz: MIT
- Startmodell: OSS first, kein Bezahlprodukt in den ersten 6 Monaten
- Plattformstart: Linux zuerst, macOS später als v1.1 oder v1.2

## v1 Scope

- Go-CLI als Single Binary
- Kein Cloud-Zwang
- Kein Daemon
- Keine GUI
- Kein Full-Repo-Scan in v1

### Öffentliche CLI-Kommandos

- `guard scan`
- `guard scan --env`
- `guard hook install`
- `guard hook remove`
- `guard version`

### Verhaltensdefault

- Standardmodus: Warnen, nicht blocken
- Optional: `--strict` blockiert Commit bei Treffern

### Konfiguration

- Datei: `.guard.yaml` im Repo-Root
- Erste Felder:
  - `enabled_rules`
  - `ignore_paths`
  - `strict_mode`

## Erkennungslogik v1

- Eigenes leichtes Regelset
- Erkennung typischer Secrets:
  - API-Keys
  - Access-Tokens
  - Private Keys
  - riskante `.env`-Einträge
- Fokus auf staged Diff + `.env*`

## Testplan v1

### Funktion

- Secret im staged Diff wird erkannt
- Secret in `.env` wird erkannt
- Unkritische ähnliche Strings sollen nicht unnötig alarmieren

### Hook

- `hook install` funktioniert in frischem Repo
- Commit mit Treffer zeigt Warnung
- Strict-Modus blockiert Commit

### UX

- Setup auf frischem Linux-System in unter 3 Minuten
- Ausgabe verständlich ohne lange Doku

## Roadmap nach v1

- v1.1: Context-Redaction-Modul für Agent/CLI-Ausgaben
- v1.2: macOS-Unterstützung
- v1.x: optional Windows, nur bei klarer Nachfrage

## Relevante Markthypothese

- Secret-Leaks sind ein stark wachsendes Problem
- AI-Nutzung steigt, Vertrauen in korrekte/sichere Ausgabe bleibt begrenzt
- Ein lokales, simples Guardrail vor Commits hat hohe praktische Relevanz
