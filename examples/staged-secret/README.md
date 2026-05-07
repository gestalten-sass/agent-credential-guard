# Beispiel: Secret im staged Diff

```bash
cd /pfad/zu/deinem/repo
echo 'TOKEN=EXAMPLE_DEMO_TOKEN_123456' > demo.txt
git add demo.txt
guard scan
```

Erwartung: Treffer im `staged-diff`.
