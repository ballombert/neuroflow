package kanban

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorageMoveAndSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "neuroflow_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	board, err := LoadBoard(tmpDir)
	if err != nil {
		t.Fatalf("Erreur LoadBoard: %v", err)
	}

	if len(board.Columns) < 4 {
		t.Fatalf("Attendu au moins 4 colonnes, obtenu: %d", len(board.Columns))
	}

	// Créer une nouvelle tâche dans 0_Inbox
	task, err := board.CreateTask("0_Inbox", "Tâche Test", "Description", []string{"Sous-tâche A", "Sous-tâche B"})
	if err != nil {
		t.Fatalf("Erreur CreateTask: %v", err)
	}

	// Vérifier la présence physique du fichier
	if _, err := os.Stat(task.Path); os.IsNotExist(err) {
		t.Fatalf("Le fichier physique n'a pas été créé: %s", task.Path)
	}

	// Déplacer vers 2_En_Cours
	if err := board.MoveTask(task, "2_En_Cours"); err != nil {
		t.Fatalf("Erreur MoveTask: %v", err)
	}

	if task.ColumnDir != "2_En_Cours" {
		t.Errorf("Colonne de la tâche incorrecte: %s", task.ColumnDir)
	}

	if _, err := os.Stat(task.Path); os.IsNotExist(err) {
		t.Fatalf("Le fichier physique devrait exister dans le nouveau dossier: %s", task.Path)
	}

	// Vérifier que l'ancien fichier n'existe plus
	oldExpectedPath := filepath.Join(tmpDir, "0_Inbox", task.FileName)
	if _, err := os.Stat(oldExpectedPath); !os.IsNotExist(err) {
		t.Fatalf("L'ancien fichier ne devrait plus exister dans 0_Inbox: %s", oldExpectedPath)
	}
}

func TestArchiveAndRestoreTask(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "neuroflow_archive_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	board, err := LoadBoard(tmpDir)
	if err != nil {
		t.Fatalf("Erreur LoadBoard: %v", err)
	}

	archiveDir := filepath.Join(tmpDir, ".archive")

	// 1. Créer une tâche
	task, err := board.CreateTask("2_En_Cours", "Tâche à archiver", "Contenu", []string{"Action 1"})
	if err != nil {
		t.Fatalf("Erreur CreateTask: %v", err)
	}

	// 2. Archiver la tâche
	if err := board.ArchiveTask(task, archiveDir); err != nil {
		t.Fatalf("Erreur ArchiveTask: %v", err)
	}

	// Vérifier que la tâche n'est plus dans 2_En_Cours
	foundInBoard := false
	for _, col := range board.Columns {
		for _, tk := range col.Tasks {
			if tk.FileName == task.FileName {
				foundInBoard = true
			}
		}
	}
	if foundInBoard {
		t.Errorf("La tâche archivée ne devrait plus être présente dans les colonnes actives du tableau")
	}

	// Vérifier que le fichier est bien présent dans le dossier d'archives
	if _, err := os.Stat(task.Path); os.IsNotExist(err) {
		t.Fatalf("Le fichier archivé n'existe pas dans le dossier d'archive: %s", task.Path)
	}

	// 3. Charger les tâches archivées
	archivedList, err := LoadArchivedTasks(archiveDir)
	if err != nil {
		t.Fatalf("Erreur LoadArchivedTasks: %v", err)
	}
	if len(archivedList) != 1 {
		t.Fatalf("Attendu 1 tâche archivée, obtenu: %d", len(archivedList))
	}
	if archivedList[0].Title != "Tâche à archiver" {
		t.Errorf("Titre archivé attendu 'Tâche à archiver', obtenu: %s", archivedList[0].Title)
	}

	// 4. Restaurer la tâche vers 0_Inbox
	if err := board.RestoreTask(archivedList[0], "0_Inbox"); err != nil {
		t.Fatalf("Erreur RestoreTask: %v", err)
	}

	// Vérifier qu'elle est réapparue dans 0_Inbox
	foundRestored := false
	for _, tk := range board.Columns[0].Tasks {
		if tk.Title == "Tâche à archiver" {
			foundRestored = true
			break
		}
	}
	if !foundRestored {
		t.Errorf("La tâche restaurée n'a pas été trouvée dans 0_Inbox")
	}
}
