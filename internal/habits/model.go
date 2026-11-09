package habits

import (
	"fmt"
	"time"
)

// Habit représente une habitude avec micro-action et boucle de renforcement dopaminergique.
type Habit struct {
	ID          string    `json:"id"`
	FileName    string    `json:"file_name"`
	Path        string    `json:"path"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`     // "Pro", "Focus", "Energy"
	Trigger     string    `json:"trigger"`      // Déclencheur ("Si...")
	MicroAction string    `json:"micro_action"` // Micro-action ("Alors...")
	XPValue     int       `json:"xp_value"`     // Points d'XP par validation (10 à 25)
	Streak      int       `json:"streak"`       // Série en cours
	TotalXP     int       `json:"total_xp"`     // XP total accumulé sur cette habitude
	LastDoneDate string   `json:"last_done_date"` // Format "YYYY-MM-DD"
	History     []string  `json:"history"`      // Entrées du journal ("- [x] 2026-11-09 : Validé (+20 XP)")
	BodyLines   []string  `json:"body_lines"`   // Lignes brutes du Markdown
}

// IsDoneToday vérifie si l'habitude a déjà été validée aujourd'hui.
func (h *Habit) IsDoneToday() bool {
	today := time.Now().Format("2006-01-02")
	return h.LastDoneDate == today
}

// CategoryEmoji renvoie un emoji évocateur pour chaque catégorie.
func (h *Habit) CategoryEmoji() string {
	switch h.Category {
	case "Pro":
		return "💼"
	case "Focus":
		return "🧠"
	case "Energy":
		return "⚡"
	default:
		return "🎯"
	}
}

// PlayerProfile calcule le rang et la progression globale de l'utilisateur.
type PlayerProfile struct {
	TotalXP int
	Level   int
	Title   string
	NextXP  int
}

// CalculateProfile déduit le profil joueur à partir de la somme des XP de toutes les habitudes.
func CalculateProfile(habits []*Habit) PlayerProfile {
	totalXP := 0
	for _, h := range habits {
		totalXP += h.TotalXP
	}

	level := (totalXP / 100) + 1
	nextLevelXP := level * 100

	var title string
	switch {
	case level >= 10:
		title = "👑 Légende Dopaminergique"
	case level >= 7:
		title = "💎 Alchimiste de la Clarté"
	case level >= 5:
		title = "🌊 Surfeur d'Hyperfocus"
	case level >= 3:
		title = "🔥 Maître du Focus"
	case level >= 2:
		title = "⚡ Dompteur d'Inertie"
	default:
		title = "🌱 Apprenti de l'Élan"
	}

	return PlayerProfile{
		TotalXP: totalXP,
		Level:   level,
		Title:   title,
		NextXP:  nextLevelXP,
	}
}

// XPBar génère une barre de progression de niveau pour le profil.
func (p PlayerProfile) XPBar(width int) string {
	if width < 10 {
		width = 10
	}
	currentInLevel := p.TotalXP % 100
	ratio := float64(currentInLevel) / 100.0
	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	f := ""
	for i := 0; i < filled; i++ {
		f += "█"
	}
	e := ""
	for i := 0; i < empty; i++ {
		e += "░"
	}
	return fmt.Sprintf("%s%s (%d/100 XP)", f, e, currentInLevel)
}
