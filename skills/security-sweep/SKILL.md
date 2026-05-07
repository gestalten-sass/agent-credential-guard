---
name: security-sweep
description: Fuehrt einen fokussierten Security-Check fuer dieses Repo aus und liefert Befunde als Kritisch/Warnung/OK.
---

## Zweck

Nutze diesen Skill nach groesseren Aenderungen, vor Releases und bei Security-Verdacht.

## Scope

Standard-Scope: aktuelles Repo (`agent-credential-guard`).

## Pruefablauf

1. Secrets in Dateien (ohne .env)
```bash
rg -n "api_key|secret|token|password|bearer" -i --glob '!*.env' --glob '!*/.git/*' .
```

2. .env und Ignore-Regeln
```bash
cat .gitignore
find . -name ".env*" -not -path "*/.git/*"
```

3. Git-History auf Secret-Hinweise
```bash
git log --all --full-history -p -- '*.env' | head -120
git log --all -S "api_key" --oneline
```

4. Logs/Temp auf Secret-Hinweise
```bash
find . \( -name "*.log" -o -name "*.tmp" \) -print0 | xargs -0 rg -n -i "key|token|secret" 2>/dev/null
```

5. Berechtigungen .env-Dateien
```bash
find . -name ".env*" -exec ls -la {} \;
```

## Ausgabeformat

- `Kritisch`: sofort handeln (z. B. Key rotieren/widerrufen)
- `Warnung`: zeitnah verbessern
- `OK`: kein Befund

Wichtig: Secrets nie im Klartext wiedergeben, nur Dateipfad und ggf. Zeile nennen.
