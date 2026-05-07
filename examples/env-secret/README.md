# Beispiel: Secret in .env

```bash
cd /pfad/zu/deinem/repo
echo 'API_KEY=EXAMPLE_DEMO_KEY_123456' > .env
guard scan --env
```

Erwartung: Treffer in `.env:1`.
