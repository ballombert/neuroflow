package habits

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	titleRegex    = regexp.MustCompile(`^#\s+(.+)$`)
	metaRegex     = regexp.MustCompile(`^>\s+\*\*(.+?)\s*:\*\*\s*(.+)$`)
	historyRegex  = regexp.MustCompile(`^-\s*\[([ xX])\]\s*(\d{4}-\d{2}-\d{2})\s*:\s*(.*)$`)
)

// InitDefaultHabits initialise le dossier et les habitudes conçues pour le TDAH et la vie pro.
func InitDefaultHabits(root string) error {
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	starterHabits := []struct {
		fileName string
		content  string
	}{
		{
			fileName: "kickoff-matinal-pro.md",
			content: `# 💼 Kickoff Pro : Choisir mes 3 victoires du jour

> **Catégorie :** Pro
> **Déclencheur (Si...) :** Si j'allume mon ordinateur le matin
> **Micro-Action (Alors...) :** Alors j'ouvre NeuroFlow et choisis mes 3 priorités d'impact
> **XP :** 20
> **Streak :** 0
> **Total XP :** 0
> **Dernière Validation :** 

### Principes Neuro-Inclusifs :
La 'Règle des 3' élimine la surcharge cognitive matinale. Le cerveau TDAH ne peut pas prioriser une liste de 40 éléments ; 3 victoires donnent un cap clair et gratifiant.

### Journal des Victoires :
`,
		},
		{
			fileName: "tri-inbox-chrono-pro.md",
			content: `# 💼 Tri Express Inbox Pro (15 min chrono)

> **Catégorie :** Pro
> **Déclencheur (Si...) :** Si j'ouvre mes emails ou Slack
> **Micro-Action (Alors...) :** Alors je lance un compte à rebours de 15 min chrono pour éviter le 'trou noir'
> **XP :** 15
> **Streak :** 0
> **Total XP :** 0
> **Dernière Validation :** 

### Principes Neuro-Inclusifs :
Les boîtes de réception sont des aspirateurs à dopamine passive. Borner temporellement le tri préserve l'énergie pour les tâches à forte valeur.

### Journal des Victoires :
`,
		},
		{
			fileName: "rituel-shutdown-pro.md",
			content: `# 💼 Rituel de Clôture Pro (Shutdown à 18h)

> **Catégorie :** Pro
> **Déclencheur (Si...) :** Si la fin de journée approche (18h00)
> **Micro-Action (Alors...) :** Alors je ferme mes onglets pros, range mon bureau et note la 1ère tâche de demain
> **XP :** 25
> **Streak :** 0
> **Total XP :** 0
> **Dernière Validation :** 

### Principes Neuro-Inclusifs :
Sans frontière nette, le cerveau TDAH continue de ruminer le travail toute la soirée. Définir formellement la fin du travail protège contre le burnout.

### Journal des Victoires :
`,
		},
		{
			fileName: "amorce-soft-launch.md",
			content: `# 🧠 Défi des 5 Minutes (Soft-Launch)

> **Catégorie :** Focus
> **Déclencheur (Si...) :** Si je ressens une paralysie devant une tâche intimidante
> **Micro-Action (Alors...) :** Alors j'appuie sur 's' pour faire seulement 5 minutes sans engagement
> **XP :** 15
> **Streak :** 0
> **Total XP :** 0
> **Dernière Validation :** 

### Principes Neuro-Inclusifs :
La friction d'amorce est neurobiologique. Réduire l'enjeu à 5 minutes désamorce la résistance de l'amygdale.

### Journal des Victoires :
`,
		},
		{
			fileName: "pause-eau-mouvement.md",
			content: `# ⚡ Reset Énergie & Hydratation

> **Catégorie :** Energy
> **Déclencheur (Si...) :** Si une pause Pomodoro commence
> **Micro-Action (Alors...) :** Alors je bois un grand verre d'eau et étire mon dos
> **XP :** 10
> **Streak :** 0
> **Total XP :** 0
> **Dernière Validation :** 

### Principes Neuro-Inclusifs :
L'hyperfocus fait souvent oublier la soif et la posture. Un reset physiologique réamorce la dopamine naturelle.

### Journal des Victoires :
`,
		},
	}

	for _, sh := range starterHabits {
		target := filepath.Join(root, sh.fileName)
		if _, err := os.Stat(target); os.IsNotExist(err) {
			_ = os.WriteFile(target, []byte(sh.content), 0644)
		}
	}

	return nil
}

// LoadHabits parcourt le dossier des habitudes et charge l'ensemble des fichiers Markdown.
func LoadHabits(root string) ([]*Habit, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := InitDefaultHabits(root); err != nil {
			return nil, err
		}
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var habitsList []*Habit
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			fullPath := filepath.Join(root, e.Name())
			data, err := os.ReadFile(fullPath)
			if err == nil {
				h := ParseHabitFile(string(data), e.Name(), fullPath)
				habitsList = append(habitsList, h)
			}
		}
	}

	if len(habitsList) == 0 {
		_ = InitDefaultHabits(root)
		return LoadHabits(root)
	}

	return habitsList, nil
}

// ParseHabitFile lit et extrait les métadonnées et l'historique d'un fichier Markdown d'habitude.
func ParseHabitFile(content, fileName, fullPath string) *Habit {
	h := &Habit{
		ID:        strings.TrimSuffix(fileName, ".md"),
		FileName:  fileName,
		Path:      fullPath,
		XPValue:   15,
		Category:  "Pro",
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		h.BodyLines = append(h.BodyLines, line)

		if h.Title == "" {
			if m := titleRegex.FindStringSubmatch(line); len(m) > 1 {
				h.Title = strings.TrimSpace(m[1])
			}
		}

		if m := metaRegex.FindStringSubmatch(line); len(m) == 3 {
			key := strings.TrimSpace(strings.ToLower(m[1]))
			val := strings.TrimSpace(m[2])

			switch {
			case strings.Contains(key, "catégorie") || strings.Contains(key, "categorie"):
				h.Category = val
			case strings.Contains(key, "déclencheur") || strings.Contains(key, "declencheur"):
				h.Trigger = val
			case strings.Contains(key, "action"):
				h.MicroAction = val
			case strings.Contains(key, "total xp"):
				if tx, err := strconv.Atoi(val); err == nil {
					h.TotalXP = tx
				}
			case key == "xp" || strings.HasPrefix(key, "xp"):
				if num, err := strconv.Atoi(strings.TrimPrefix(val, "+")); err == nil {
					h.XPValue = num
				}
			case strings.Contains(key, "streak"):
				if s, err := strconv.Atoi(val); err == nil {
					h.Streak = s
				}
			case strings.Contains(key, "validation"):
				h.LastDoneDate = val
			}
		}

		if m := historyRegex.FindStringSubmatch(line); len(m) == 4 {
			h.History = append(h.History, line)
		}
	}

	if h.Title == "" {
		h.Title = strings.Title(strings.ReplaceAll(h.ID, "-", " "))
	}

	return h
}

// SaveHabitFile régénère et sauvegarde le fichier Markdown de l'habitude avec ses métadonnées à jour.
func SaveHabitFile(h *Habit) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", h.Title))
	sb.WriteString(fmt.Sprintf("> **Catégorie :** %s\n", h.Category))
	sb.WriteString(fmt.Sprintf("> **Déclencheur (Si...) :** %s\n", h.Trigger))
	sb.WriteString(fmt.Sprintf("> **Micro-Action (Alors...) :** %s\n", h.MicroAction))
	sb.WriteString(fmt.Sprintf("> **XP :** %d\n", h.XPValue))
	sb.WriteString(fmt.Sprintf("> **Streak :** %d\n", h.Streak))
	sb.WriteString(fmt.Sprintf("> **Total XP :** %d\n", h.TotalXP))
	sb.WriteString(fmt.Sprintf("> **Dernière Validation :** %s\n\n", h.LastDoneDate))

	sb.WriteString("### Principes Neuro-Inclusifs :\n")
	sb.WriteString("Habitude formulée en intention d'implémentation (Gollwitzer) : Si [Déclencheur], Alors [Micro-Action].\n\n")

	sb.WriteString("### Journal des Victoires :\n")
	for _, entry := range h.History {
		sb.WriteString(entry + "\n")
	}

	return os.WriteFile(h.Path, []byte(sb.String()), 0644)
}

// ToggleHabitToday bascule la validation de l'habitude pour aujourd'hui.
// Retourne l'XP gagné (+XP) ou retiré (-XP).
func ToggleHabitToday(h *Habit) (int, error) {
	today := time.Now().Format("2006-01-02")

	if h.IsDoneToday() {
		// Déjà fait aujourd'hui : on annule (décochage)
		h.LastDoneDate = ""
		if h.Streak > 0 {
			h.Streak--
		}
		if h.TotalXP >= h.XPValue {
			h.TotalXP -= h.XPValue
		}

		// Retirer la dernière entrée si c'était aujourd'hui
		var newHistory []string
		for _, entry := range h.History {
			if !strings.Contains(entry, today) {
				newHistory = append(newHistory, entry)
			}
		}
		h.History = newHistory
		err := SaveHabitFile(h)
		return -h.XPValue, err
	}

	// Validation pour aujourd'hui
	h.LastDoneDate = today
	h.Streak++
	h.TotalXP += h.XPValue

	newEntry := fmt.Sprintf("- [x] %s : Validé (+%d XP)", today, h.XPValue)
	// Insérer en tête de l'historique
	h.History = append([]string{newEntry}, h.History...)

	err := SaveHabitFile(h)
	return h.XPValue, err
}
