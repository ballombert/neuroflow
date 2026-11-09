package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// AIAssistantConfig configure l'assistant IA Copilot / Coach TDAH.
type AIAssistantConfig struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Provider   string `json:"provider" yaml:"provider"`       // "copilot", "builtin"
	CopilotBin string `json:"copilot_bin" yaml:"copilot_bin"` // Chemin vers le binaire copilot
	RolePrompt string `json:"role_prompt" yaml:"role_prompt"` // Directives de coaching TDAH
}

// WorkLifeConfig configure les rituels d'organisation professionnelle.
type WorkLifeConfig struct {
	DailyWinsGoal int    `json:"daily_wins_goal" yaml:"daily_wins_goal"` // Règle des 3 victoires quotidiennes
	ShutdownTime  string `json:"shutdown_time" yaml:"shutdown_time"`     // Heure du rituel de fin de travail (ex: "18:00")
	InboxLimitMin int    `json:"inbox_limit_min" yaml:"inbox_limit_min"` // Durée max du tri boîte de réception
}

// Config regroupe tous les paramètres de l'application NeuroFlow.
type Config struct {
	KanbanDir           string            `json:"kanban_dir" yaml:"kanban_dir"`
	HabitsDir           string            `json:"habits_dir" yaml:"habits_dir"`
	ArchiveDir          string            `json:"archive_dir" yaml:"archive_dir"`
	WorkDuration        time.Duration     `json:"work_duration_minutes" yaml:"-"`
	SoftLaunchDuration  time.Duration     `json:"soft_launch_minutes" yaml:"-"`
	ShortBreakDuration  time.Duration     `json:"short_break_minutes" yaml:"-"`
	LongBreakDuration   time.Duration     `json:"long_break_minutes" yaml:"-"`
	WorkDurationMin     int               `json:"-" yaml:"work_duration_minutes"`
	SoftLaunchMin       int               `json:"-" yaml:"soft_launch_minutes"`
	ShortBreakMin       int               `json:"-" yaml:"short_break_minutes"`
	LongBreakMin        int               `json:"-" yaml:"long_break_minutes"`
	LongBreakInterval   int               `json:"long_break_interval" yaml:"long_break_interval"`
	YouTubeMusicURL     string            `json:"youtube_music_url" yaml:"youtube_music_url"`
	AutoLaunchMusic     bool              `json:"auto_launch_music" yaml:"auto_launch_music"`
	SoundBell           bool              `json:"sound_bell" yaml:"sound_bell"`
	Theme               string            `json:"theme" yaml:"theme"`
	AIAssistant         AIAssistantConfig `json:"ai_assistant" yaml:"ai_assistant"`
	WorkLife            WorkLifeConfig    `json:"work_life" yaml:"work_life"`
}

// DefaultConfig fournit une configuration scientifiquement adaptée au TDAH et à la vie pro.
func DefaultConfig() *Config {
	return &Config{
		KanbanDir:          "./kanban_workspace",
		HabitsDir:          "./habits_workspace",
		ArchiveDir:         "./kanban_workspace/.archive",
		WorkDuration:       25 * time.Minute,
		SoftLaunchDuration: 5 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		WorkDurationMin:    25,
		SoftLaunchMin:      5,
		ShortBreakMin:      5,
		LongBreakMin:       15,
		LongBreakInterval:  4,
		YouTubeMusicURL:    "https://music.youtube.com/watch?v=jfKfPfyJRdk",
		AutoLaunchMusic:    false,
		SoundBell:          true,
		Theme:              "tokyo-night",
		AIAssistant: AIAssistantConfig{
			Enabled:    true,
			Provider:   "copilot",
			CopilotBin: "copilot",
			RolePrompt: "Tu es un coach expert en TDAH et en productivité bienveillante. Tu aides à décomposer les tâches professionnelles intimidantes en micro-étapes atomiques ridicules, sans culpabilité ni jargon superflu.",
		},
		WorkLife: WorkLifeConfig{
			DailyWinsGoal: 3,
			ShutdownTime:  "18:00",
			InboxLimitMin: 15,
		},
	}
}

// LoadConfigFromFile charge la configuration depuis un fichier YAML ou JSON selon son extension.
func LoadConfigFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			_ = SaveConfigToFile(path, cfg)
			return cfg, nil
		}
		return nil, err
	}

	if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return cfg, err
		}
	} else {
		if err := json.Unmarshal(data, cfg); err != nil {
			return cfg, err
		}
	}

	// Synchronisation des Durations
	if cfg.WorkDurationMin > 0 {
		cfg.WorkDuration = time.Duration(cfg.WorkDurationMin) * time.Minute
	}
	if cfg.SoftLaunchMin > 0 {
		cfg.SoftLaunchDuration = time.Duration(cfg.SoftLaunchMin) * time.Minute
	}
	if cfg.ShortBreakMin > 0 {
		cfg.ShortBreakDuration = time.Duration(cfg.ShortBreakMin) * time.Minute
	}
	if cfg.LongBreakMin > 0 {
		cfg.LongBreakDuration = time.Duration(cfg.LongBreakMin) * time.Minute
	}

	if cfg.ArchiveDir == "" {
		cfg.ArchiveDir = filepath.Join(cfg.KanbanDir, ".archive")
	}

	return cfg, nil
}

// SaveConfigToFile écrit la configuration au format YAML (recommandé) ou JSON.
func SaveConfigToFile(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	cfg.WorkDurationMin = int(cfg.WorkDuration / time.Minute)
	cfg.SoftLaunchMin = int(cfg.SoftLaunchDuration / time.Minute)
	cfg.ShortBreakMin = int(cfg.ShortBreakDuration / time.Minute)
	cfg.LongBreakMin = int(cfg.LongBreakDuration / time.Minute)

	var data []byte
	var err error

	if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
		data, err = yaml.Marshal(cfg)
	} else {
		data, err = json.MarshalIndent(cfg, "", "  ")
	}

	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadConfig reste rétro-compatible.
func LoadConfig(path string) (*Config, error) {
	return LoadConfigFromFile(path)
}

// SaveConfig reste rétro-compatible.
func SaveConfig(path string, cfg *Config) error {
	return SaveConfigToFile(path, cfg)
}
