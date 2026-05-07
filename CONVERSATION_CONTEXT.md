# Gesprächskontext und Entscheidungen

Stand: 2026-05-07

Dieses Dokument fasst die bisherige Unterhaltung zusammen, damit der nächste Kontext direkt anschließen kann.

## Ausgangspunkt

Der Wunsch war ein sinnvolles Linux-Projekt mit realem Nutzen und Marktchance, statt nur ein experimentelles Tool zu bauen.

## Diskutierte Projektoptionen

Es wurden mehrere Ideen verglichen. Favorisiert wurde:

- Linux Agent Credential Guard

Begründung:

- direkter Sicherheitsnutzen
- klarer Schmerzpunkt am Markt
- überschaubarer MVP
- hohe Relevanz für AI- und Git-Workflows

## Zentrale Entscheidungen

- v1 Plattform: Linux zuerst
- Umsetzung: Go CLI + Git Hooks
- Betriebsart: On-demand CLI
- Scope v1: staged Diff + `.env*`
- Default: Warnen statt hart blocken
- Lizenz: MIT
- Distribution: Open Source auf GitHub
- Positionierung: Agent-Workflow-Fokus
- Integrationsstrategie: eigenes leichtes Regelset statt Wrapper auf große bestehende Scanner

## Wichtiger Zusatz aus der Unterhaltung

Es gab einen konkreten Praxisfall:

- Agenten haben trotz Anweisung Secrets aus `.env` im Klartext ins Kontextfenster geschrieben.

Bewertung:

- gehört in dieselbe Produktfamilie
- aber nicht in v1-Kern
- geplante Erweiterung als eigenes Modul nach v1

Roadmap dazu:

- v1.1: `guard redact` oder äquivalentes Modul für Kontext-Redaction in Agent/CLI-Ausgaben

## Strategischer Grundansatz

- klein starten
- schnell reale Nutzung bekommen
- Fehlalarme niedrig halten
- danach modular erweitern

## Erwartungsrahmen Adoption

Realistische Ersteinschätzung aus der Unterhaltung:

- ohne aktiven Launch: niedrige bis moderate Downloads
- mit aktivem Community-Launch und Demos: signifikant höhere Reichweite möglich

## Nächster praktischer Schritt

- Projektstruktur initialisieren
- Go CLI Grundgerüst bauen
- erstes v1 Featurepaket umsetzen
  - `scan`
  - `scan --env`
  - `hook install/remove`

