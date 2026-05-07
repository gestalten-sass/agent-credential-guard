package main

import (
	"flag"
	"fmt"
	"os"

	"agent-credential-guard/internal/config"
	"agent-credential-guard/internal/hook"
	"agent-credential-guard/internal/scanner"
)

const version = "v0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "scan":
		os.Exit(runScan(os.Args[2:]))
	case "hook":
		os.Exit(runHook(os.Args[2:]))
	case "init":
		os.Exit(runInit())
	case "version":
		fmt.Println(version)
		return
	case "-h", "--help", "help":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unbekanntes kommando: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func runInit() int {
	p, err := config.InitFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "init fehlgeschlagen: %v\n", err)
		return 1
	}
	fmt.Printf("guard: config angelegt: %s\n", p)
	return 0
}

func runScan(args []string) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	includeEnv := fs.Bool("env", false, "scannt zusaetzlich .env* dateien")
	strictFlag := fs.Bool("strict", false, "liefert exit code 1 bei treffern")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config konnte nicht geladen werden: %v\n", err)
		return 1
	}

	report, err := scanner.Run(scanner.Options{IncludeEnv: *includeEnv, Config: cfg})
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan fehlgeschlagen: %v\n", err)
		return 1
	}

	strict := *strictFlag || cfg.StrictMode
	if len(report.Findings) == 0 {
		fmt.Println("guard: keine secrets gefunden")
		return 0
	}

	fmt.Printf("guard: %d potentiell kritische treffer gefunden\n", len(report.Findings))
	for _, f := range report.Findings {
		fmt.Printf("- [%s] %s\n", f.Rule, f.Location)
		fmt.Printf("  -> fix: %s\n", scanner.HintForRule(f.Rule))
	}

	if strict {
		fmt.Println("guard: strict aktiv, commit wird blockiert")
		return 1
	}

	fmt.Println("guard: warnmodus aktiv, commit wird nicht blockiert")
	return 0
}

func runHook(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "hook erwartet ein subkommando: install|remove")
		return 1
	}

	switch args[0] {
	case "install":
		if err := hook.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "hook install fehlgeschlagen: %v\n", err)
			return 1
		}
		fmt.Println("guard: pre-commit hook installiert")
		return 0
	case "remove":
		removed, err := hook.Remove()
		if err != nil {
			fmt.Fprintf(os.Stderr, "hook remove fehlgeschlagen: %v\n", err)
			return 1
		}
		if removed {
			fmt.Println("guard: pre-commit hook entfernt")
		} else {
			fmt.Println("guard: kein guard-hook gefunden, nichts entfernt")
		}
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unbekanntes hook-subkommando: %s\n", args[0])
		return 1
	}
}

func printUsage() {
	fmt.Println("guard <kommando>")
	fmt.Println("")
	fmt.Println("Kommandos:")
	fmt.Println("  init")
	fmt.Println("  scan [--env] [--strict]")
	fmt.Println("  hook install")
	fmt.Println("  hook remove")
	fmt.Println("  version")
}
