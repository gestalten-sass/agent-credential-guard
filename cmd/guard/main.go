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
	global := false
	rest := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--global" {
			global = true
			continue
		}
		rest = append(rest, a)
	}
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "hook erwartet ein subkommando: install|remove|status")
		return 1
	}

	switch rest[0] {
	case "install":
		if global {
			p, err := hook.InstallGlobal()
			if err != nil {
				fmt.Fprintf(os.Stderr, "hook install --global fehlgeschlagen: %v\n", err)
				return 1
			}
			fmt.Printf("guard: globaler pre-commit hook installiert (%s)\n", p)
			return 0
		}
		if err := hook.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "hook install fehlgeschlagen: %v\n", err)
			return 1
		}
		fmt.Println("guard: pre-commit hook installiert")
		return 0
	case "remove":
		if global {
			removed, err := hook.RemoveGlobal()
			if err != nil {
				fmt.Fprintf(os.Stderr, "hook remove --global fehlgeschlagen: %v\n", err)
				return 1
			}
			if removed {
				fmt.Println("guard: globaler pre-commit hook entfernt")
			} else {
				fmt.Println("guard: kein globaler guard-hook gefunden")
			}
			return 0
		}
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
	case "status":
		if global {
			enabled, path, err := hook.GlobalStatus()
			if err != nil {
				fmt.Fprintf(os.Stderr, "hook status --global fehlgeschlagen: %v\n", err)
				return 1
			}
			if enabled {
				fmt.Printf("guard: global aktiv (%s)\n", path)
				return 0
			}
			if path == "" {
				fmt.Println("guard: global inaktiv (core.hooksPath nicht gesetzt)")
			} else {
				fmt.Printf("guard: global inaktiv (anderer core.hooksPath: %s)\n", path)
			}
			return 1
		}
		fmt.Println("guard: fuer lokalen status pruefe .git/hooks/pre-commit")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unbekanntes hook-subkommando: %s\n", rest[0])
		return 1
	}
}

func printUsage() {
	fmt.Println("guard <kommando>")
	fmt.Println("")
	fmt.Println("Kommandos:")
	fmt.Println("  init")
	fmt.Println("  scan [--env] [--strict]")
	fmt.Println("  hook install [--global]")
	fmt.Println("  hook remove [--global]")
	fmt.Println("  hook status [--global]")
	fmt.Println("  version")
}
