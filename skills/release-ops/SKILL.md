---
name: release-ops
description: Standardisiert Release-Ablauf fuer dieses Projekt (test, build, commit, push, tag, release-check) mit Sicherheits- und Hygiene-Checks.
---

## Zweck

Nutze diesen Skill immer dann, wenn ein Release oder ein "fertig machen und veroeffentlichen" ansteht.

## Ablauf

1. Arbeitsstand pruefen
```bash
git status --short --ignored
git branch --show-current
```

2. Qualitaet pruefen
```bash
go test ./...
go build -o guard ./cmd/guard
```

3. Security/Hygiene kurz pruefen
```bash
rg -n "api_key|secret|token|password|bearer" -i --glob '!*.env' --glob '!*/.git/*' .
```
Hinweis: Demo-/Test-Treffer gesondert bewerten.

4. Commit erstellen
- Nur beabsichtigte Dateien stagen.
- Konventionelle Commit-Message verwenden.

5. Push auf Hauptbranch
```bash
git push
```

6. Release-Tag setzen und pushen
```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

7. Release-Workflow pruefen
```bash
gh run list --repo gestalten-sass/agent-credential-guard --limit 6
gh release view vX.Y.Z --repo gestalten-sass/agent-credential-guard --json url,assets
```

## Checkliste vor Abschluss

- Tests erfolgreich
- Build erfolgreich
- Keine ungewollten Dateien im Commit
- Tag gepusht
- Release-Assets vorhanden (amd64 + arm64 + sha256)
