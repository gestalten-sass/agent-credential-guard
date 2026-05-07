# Beispiel: Secret in .env

```bash
cd /pfad/zu/deinem/repo
echo 'API_KEY=1234567890abcdef' > .env
guard scan --env
```

Erwartung: Treffer in `.env:1`.
