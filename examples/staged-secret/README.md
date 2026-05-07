# Beispiel: Secret im staged Diff

```bash
cd /pfad/zu/deinem/repo
echo 'TOKEN=1234567890abcdef' > demo.txt
git add demo.txt
guard scan
```

Erwartung: Treffer im `staged-diff`.
