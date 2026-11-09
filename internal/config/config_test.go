package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestYAMLConfigLoadAndSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "neuroflow_yaml_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	yamlPath := filepath.Join(tmpDir, "config.yaml")

	// 1. Sauvegarder la config par défaut en YAML
	cfg := DefaultConfig()
	cfg.WorkDuration = 20 * time.Minute
	cfg.AIAssistant.Enabled = true
	cfg.AIAssistant.Provider = "copilot"
	cfg.WorkLife.DailyWinsGoal = 3

	if err := SaveConfigToFile(yamlPath, cfg); err != nil {
		t.Fatalf("Erreur SaveConfigToFile YAML: %v", err)
	}

	// 2. Recharger le fichier YAML
	loaded, err := LoadConfigFromFile(yamlPath)
	if err != nil {
		t.Fatalf("Erreur LoadConfigFromFile YAML: %v", err)
	}

	if loaded.WorkDuration != 20*time.Minute {
		t.Errorf("WorkDuration attendu 20m, obtenu %v", loaded.WorkDuration)
	}
	if !loaded.AIAssistant.Enabled {
		t.Errorf("AIAssistant.Enabled devrait être true")
	}
	if loaded.WorkLife.DailyWinsGoal != 3 {
		t.Errorf("DailyWinsGoal attendu 3, obtenu %d", loaded.WorkLife.DailyWinsGoal)
	}
	if loaded.HabitsDir != "./habits_workspace" {
		t.Errorf("HabitsDir attendu ./habits_workspace, obtenu %s", loaded.HabitsDir)
	}
}
