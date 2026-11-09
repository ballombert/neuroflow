package kanban

import (
	"fmt"
	"time"
)

// Subtask représente une case à cocher dans le fichier Markdown (- [ ] ou - [x]).
type Subtask struct {
	Index    int    // Index ordinal de la sous-tâche dans la liste
	LineNum  int    // Numéro de ligne dans le fichier Markdown d'origine
	Text     string // Libellé de la sous-tâche
	Done     bool   // État coché / non coché
}

// Task représente un fichier Markdown dans une colonne du Kanban.
type Task struct {
	FileName    string    // Ex: "implementer-auth.md"
	Path        string    // Chemin absolu ou relatif du fichier
	ColumnDir   string    // Nom du répertoire colonne parent (ex: "1_A_Faire")
	Title       string    // Titre extrait du premier '# ' ou déduit du nom de fichier
	BodyLines   []string  // Contenu complet ligne par ligne
	Subtasks    []Subtask // Liste des sous-tâches détectées
	ModifiedAt  time.Time // Horodatage de dernière modification
}

// ProgressPercentage retourne le pourcentage de complétion des sous-tâches (0 à 100).
func (t *Task) ProgressPercentage() int {
	if len(t.Subtasks) == 0 {
		return 0
	}
	doneCount := 0
	for _, st := range t.Subtasks {
		if st.Done {
			doneCount++
		}
	}
	return int(float64(doneCount) / float64(len(t.Subtasks)) * 100)
}

// ProgressText retourne une chaîne synthétique du type "[2/5] 40%".
func (t *Task) ProgressText() string {
	if len(t.Subtasks) == 0 {
		return "0 sous-tâche"
	}
	doneCount := 0
	for _, st := range t.Subtasks {
		if st.Done {
			doneCount++
		}
	}
	pct := int(float64(doneCount) / float64(len(t.Subtasks)) * 100)
	return fmt.Sprintf("[%d/%d] %d%%", doneCount, len(t.Subtasks), pct)
}

// Column représente un répertoire contenant des fichiers de tâches.
type Column struct {
	DirName     string  // Ex: "0_Inbox", "1_A_Faire", "2_En_Cours", "3_Termine"
	DisplayName string  // Ex: "📥 Boîte de Réception", "🔥 En Cours"
	Tasks       []*Task // Liste des tâches dans cette colonne
}

// Board représente l'espace de travail Kanban complet.
type Board struct {
	RootPath string
	Columns  []*Column
}
