package kanban

import (
	"strings"
	"testing"
	"time"
)

func TestParseMarkdownTask(t *testing.T) {
	md := `# Ma Super Tâche TDAH

Ceci est une description de test.

### Sous-tâches :
- [ ] Première étape
- [x] Deuxième étape déjà terminée
- [ ] Troisième étape
`
	task := ParseMarkdownTask(md, "ma-tache.md", "1_A_Faire", "/tmp/ma-tache.md", time.Now())

	if task.Title != "Ma Super Tâche TDAH" {
		t.Fatalf("Titre attendu 'Ma Super Tâche TDAH', obtenu: '%s'", task.Title)
	}

	if len(task.Subtasks) != 3 {
		t.Fatalf("3 sous-tâches attendues, obtenu: %d", len(task.Subtasks))
	}

	if task.Subtasks[0].Done != false || task.Subtasks[0].Text != "Première étape" {
		t.Errorf("Sous-tâche 0 erronée: %+v", task.Subtasks[0])
	}
	if task.Subtasks[1].Done != true || task.Subtasks[1].Text != "Deuxième étape déjà terminée" {
		t.Errorf("Sous-tâche 1 erronée: %+v", task.Subtasks[1])
	}
	if task.Subtasks[2].Done != false || task.Subtasks[2].Text != "Troisième étape" {
		t.Errorf("Sous-tâche 2 erronée: %+v", task.Subtasks[2])
	}

	pct := task.ProgressPercentage()
	if pct != 33 {
		t.Errorf("Pourcentage attendu 33%%, obtenu: %d%%", pct)
	}
}

func TestToggleSubtask(t *testing.T) {
	md := `# Tâche

- [ ] Étape 1
- [x] Étape 2
`
	task := ParseMarkdownTask(md, "t.md", "1_A_Faire", "/tmp/t.md", time.Now())

	// Cocher l'étape 1
	if !task.ToggleSubtask(0) {
		t.Fatalf("Échec toggle sous-tâche 0")
	}
	if !task.Subtasks[0].Done {
		t.Errorf("Sous-tâche 0 devrait être cochée")
	}

	serialized := task.SerializeContent()
	if !strings.Contains(serialized, "- [x] Étape 1") {
		t.Errorf("Le Markdown sérialisé devrait contenir '- [x] Étape 1', contenu:\n%s", serialized)
	}

	// Décocher l'étape 2
	if !task.ToggleSubtask(1) {
		t.Fatalf("Échec toggle sous-tâche 1")
	}
	if task.Subtasks[1].Done {
		t.Errorf("Sous-tâche 1 devrait être décochée")
	}

	serialized = task.SerializeContent()
	if !strings.Contains(serialized, "- [ ] Étape 2") {
		t.Errorf("Le Markdown sérialisé devrait contenir '- [ ] Étape 2', contenu:\n%s", serialized)
	}
}

func TestAppendBrainDump(t *testing.T) {
	md := `# Tâche Simple
- [ ] Action
`
	task := ParseMarkdownTask(md, "t.md", "1_A_Faire", "/tmp/t.md", time.Now())

	task.AppendBrainDump("Commander du café")
	if len(task.Subtasks) != 2 {
		t.Fatalf("Devrait avoir 2 sous-tâches après le brain dump, obtenu: %d", len(task.Subtasks))
	}

	serialized := task.SerializeContent()
	if !strings.Contains(serialized, "Commander du café") || !strings.Contains(serialized, "Brain Dump") {
		t.Errorf("Le contenu sérialisé ne contient pas le brain dump:\n%s", serialized)
	}
}
