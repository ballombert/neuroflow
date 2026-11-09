package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"neuroflow/internal/config"
	"neuroflow/internal/habits"
	"neuroflow/internal/kanban"
	"neuroflow/internal/tui"
)

const version = "1.1.0"

func main() {
	// Configuration UTF-8 sur Windows
	initWindowsConsole()

	configPath := flag.String("config", "config.yaml", "Chemin vers le fichier de configuration YAML ou JSON")
	kanbanPath := flag.String("kanban", "", "Chemin vers le répertoire Kanban (surcharge la config)")
	habitsPath := flag.String("habits", "", "Chemin vers le répertoire des Habitudes (surcharge la config)")
	workMin := flag.Int("work", 0, "Durée de la session de focus en minutes")
	softMin := flag.Int("soft", 0, "Durée de la rampe dopaminergique en minutes")
	breakMin := flag.Int("break", 0, "Durée de la pause courte en minutes")
	musicURL := flag.String("music", "", "URL de la playlist YouTube Music de concentration")
	showVersion := flag.Bool("version", false, "Afficher la version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("NeuroFlow v%s - Pomodoro Adaptatif, Kanban Fichiers, Habitudes & Assistant Copilot pour le TDAH\n", version)
		os.Exit(0)
	}

	// 1. Détection du fichier de configuration (config.yaml prioritaire, sinon config.json)
	activeConfigPath := *configPath
	if _, err := os.Stat(activeConfigPath); os.IsNotExist(err) && activeConfigPath == "config.yaml" {
		if _, errJSON := os.Stat("config.json"); errJSON == nil {
			activeConfigPath = "config.json"
		}
	}

	cfg, err := config.LoadConfigFromFile(activeConfigPath)
	if err != nil {
		fmt.Printf("Avertissement : Erreur lors du chargement de la config (%v), valeurs par défaut appliquées.\n", err)
		cfg = config.DefaultConfig()
	}

	// Surcharges éventuelles via les drapeaux CLI
	if *kanbanPath != "" {
		cfg.KanbanDir = *kanbanPath
	}
	if *habitsPath != "" {
		cfg.HabitsDir = *habitsPath
	}
	if *workMin > 0 {
		cfg.WorkDuration = time.Duration(*workMin) * time.Minute
	}
	if *softMin > 0 {
		cfg.SoftLaunchDuration = time.Duration(*softMin) * time.Minute
	}
	if *breakMin > 0 {
		cfg.ShortBreakDuration = time.Duration(*breakMin) * time.Minute
	}
	if *musicURL != "" {
		cfg.YouTubeMusicURL = *musicURL
	}

	// 2. Chargement ou initialisation du Board Kanban
	board, err := kanban.LoadBoard(cfg.KanbanDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur fatale lors du chargement du dossier Kanban : %v\n", err)
		os.Exit(1)
	}

	// 3. Chargement ou initialisation des Habitudes hors Kanban
	habitsList, err := habits.LoadHabits(cfg.HabitsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Avertissement lors du chargement des habitudes : %v\n", err)
		habitsList = make([]*habits.Habit, 0)
	}

	// 4. Lancement de l'interface TUI avec Bubble Tea
	app := tui.NewApp(cfg, board, habitsList)
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de l'application : %v\n", err)
		os.Exit(1)
	}
}
