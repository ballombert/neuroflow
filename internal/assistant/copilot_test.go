package assistant

import (
	"testing"

	"neuroflow/internal/config"
)

func TestParseSubtasksFromText(t *testing.T) {
	sampleOutput := `Voici les étapes :
- [ ] Ouvrir le dépôt Git et vérifier la branche
- [x] Regarder les logs d'erreurs
- [ ] Rédiger le correctif
* [ ] Valider les tests unitaires
`
	subtasks := parseSubtasksFromText(sampleOutput)

	if len(subtasks) != 4 {
		t.Fatalf("Attendu 4 sous-tâches, obtenu: %d", len(subtasks))
	}

	if subtasks[0] != "Ouvrir le dépôt Git et vérifier la branche" {
		t.Errorf("Sous-tâche 0 inattendue: %s", subtasks[0])
	}
	if subtasks[3] != "Valider les tests unitaires" {
		t.Errorf("Sous-tâche 3 inattendue: %s", subtasks[3])
	}
}

func TestFallbackDecomposition(t *testing.T) {
	cfg := config.AIAssistantConfig{Enabled: false}
	a := NewAssistant(cfg)

	subtasks, source, err := a.DecomposeTask("Refactoriser le service d'authentification", "Description")
	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if len(subtasks) < 3 {
		t.Errorf("Attendu au moins 3 sous-tâches de secours, obtenu: %d", len(subtasks))
	}

	if source == "" {
		t.Errorf("La source du message ne devrait pas être vide")
	}
}
