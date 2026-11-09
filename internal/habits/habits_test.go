package habits

import (
	"os"
	"testing"
)

func TestHabitsParseAndToggle(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "neuroflow_habits_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := InitDefaultHabits(tmpDir); err != nil {
		t.Fatalf("Erreur InitDefaultHabits: %v", err)
	}

	habitsList, err := LoadHabits(tmpDir)
	if err != nil {
		t.Fatalf("Erreur LoadHabits: %v", err)
	}

	if len(habitsList) < 5 {
		t.Fatalf("Attendu au moins 5 habitudes par défaut, obtenu: %d", len(habitsList))
	}

	// Sélectionner la première habitude (ex: Kickoff pro)
	h := habitsList[0]
	if h.IsDoneToday() {
		t.Errorf("L'habitude ne devrait pas être marquée terminée avant validation")
	}

	initialXP := h.TotalXP
	initialStreak := h.Streak

	// Valider pour aujourd'hui
	gained, err := ToggleHabitToday(h)
	if err != nil {
		t.Fatalf("Erreur ToggleHabitToday: %v", err)
	}
	if gained <= 0 {
		t.Errorf("L'XP gagné devrait être positif, obtenu: %d", gained)
	}
	if !h.IsDoneToday() {
		t.Errorf("L'habitude devrait maintenant être terminée aujourd'hui")
	}
	if h.Streak != initialStreak+1 {
		t.Errorf("Le streak devrait avoir augmenté de 1")
	}
	if h.TotalXP != initialXP+gained {
		t.Errorf("Le TotalXP devrait être %d, obtenu %d", initialXP+gained, h.TotalXP)
	}

	// Recharger depuis le fichier pour valider la persistance Markdown sur disque
	reloaded, err := LoadHabits(tmpDir)
	if err != nil {
		t.Fatalf("Erreur rechargement: %v", err)
	}
	found := false
	for _, rh := range reloaded {
		if rh.FileName == h.FileName {
			found = true
			if !rh.IsDoneToday() {
				t.Errorf("L'habitude rechargée devrait être marquée comme terminée aujourd'hui")
			}
			if rh.TotalXP != h.TotalXP {
				t.Errorf("XP rechargé incorrect: %d vs %d", rh.TotalXP, h.TotalXP)
			}
		}
	}
	if !found {
		t.Fatalf("Fichier d'habitude %s non trouvé après rechargement", h.FileName)
	}

	// Décocher
	cancelled, err := ToggleHabitToday(h)
	if err != nil {
		t.Fatalf("Erreur dé-validation: %v", err)
	}
	if cancelled >= 0 {
		t.Errorf("L'XP de retrait devrait être négatif, obtenu: %d", cancelled)
	}
	if h.IsDoneToday() {
		t.Errorf("L'habitude ne devrait plus être terminée")
	}
}

func TestPlayerProfileCalculation(t *testing.T) {
	testHabits := []*Habit{
		{TotalXP: 80},
		{TotalXP: 140},
	}
	profile := CalculateProfile(testHabits)

	if profile.TotalXP != 220 {
		t.Errorf("TotalXP attendu 220, obtenu %d", profile.TotalXP)
	}
	if profile.Level != 3 {
		t.Errorf("Niveau attendu 3, obtenu %d", profile.Level)
	}
	if profile.Title != "🔥 Maître du Focus" {
		t.Errorf("Titre attendu '🔥 Maître du Focus', obtenu '%s'", profile.Title)
	}
}
