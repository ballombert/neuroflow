package kanban

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Mapping des noms de répertoires vers des titres avec émojis stimulants
var defaultColumnNames = map[string]string{
	"0_inbox":    "📥 Boîte d'Idées",
	"1_a_faire":  "📋 À Faire (Priorités)",
	"2_en_cours": "🔥 En Cours (Focus)",
	"3_termine":  "✅ Terminé (Victoires)",
}

// GetColumnDisplayName génère un libellé élégant pour un nom de répertoire.
func GetColumnDisplayName(dirName string) string {
	lower := strings.ToLower(dirName)
	if name, ok := defaultColumnNames[lower]; ok {
		return name
	}

	clean := strings.TrimPrefix(dirName, "0_")
	clean = strings.TrimPrefix(clean, "1_")
	clean = strings.TrimPrefix(clean, "2_")
	clean = strings.TrimPrefix(clean, "3_")
	clean = strings.TrimPrefix(clean, "4_")
	clean = strings.TrimPrefix(clean, "5_")
	clean = strings.ReplaceAll(clean, "_", " ")
	clean = strings.ReplaceAll(clean, "-", " ")
	return strings.Title(clean)
}

// InitDefaultWorkspace crée l'arborescence initiale et des tâches de bienvenue adaptées au TDAH.
func InitDefaultWorkspace(root string) error {
	cols := []string{"0_Inbox", "1_A_Faire", "2_En_Cours", "3_Termine"}
	for _, col := range cols {
		colPath := filepath.Join(root, col)
		if err := os.MkdirAll(colPath, 0755); err != nil {
			return err
		}
	}

	// Tâche 1 : En Cours
	taskEnCours := filepath.Join(root, "2_En_Cours", "decouvrir-neuroflow.md")
	if _, err := os.Stat(taskEnCours); os.IsNotExist(err) {
		content := `# Découvrir NeuroFlow et libérer son focus 🚀

Bienvenue dans votre nouvel allié neurodivergent !
Cette application est conçue pour respecter le fonctionnement unique de votre cerveau TDAH : zéro friction, visibilité spatiale du temps et dopamine positive.

### Vos premières micro-actions :
- [x] Lancer l'application NeuroFlow
- [ ] Tester le mode Focus sur cette tâche (Touche 'f')
- [ ] Lancer la musique YouTube Music pour ancrer l'attention (Touche 'm')
- [ ] Essayer le Soft-Launch de 5 minutes si vous ressentez une inertie (Touche 's')
- [ ] Décharger une idée parasite avec le Brain Dump (Touche 'b')

### Notes :
Ce fichier est stocké directement dans votre dossier './kanban_workspace/2_En_Cours/'.
Vous pouvez l'ouvrir avec n'importe quel éditeur de texte (Obsidian, VS Code, etc.).
`
		_ = os.WriteFile(taskEnCours, []byte(content), 0644)
	}

	// Tâche 2 : À Faire
	taskAFaire := filepath.Join(root, "1_A_Faire", "astuces-anti-procrastination.md")
	if _, err := os.Stat(taskAFaire); os.IsNotExist(err) {
		content := `# Astuces anti-procrastination TDAH 🧠

Le secret pour contourner la résistance du cerveau TDAH est d'abaisser l'énergie d'activation.

### Stratégies validées scientifiquement :
- [ ] Découper chaque projet en micro-étapes ridicules (ex: "ouvrir le fichier")
- [ ] Utiliser le Soft Launch de 5 min : s'engager seulement 5 minutes
- [ ] Garder la tâche visible en permanence pour contrer la cécité d'attention
- [ ] Noter immédiatement toute distraction dans le Brain Dump pour ne pas décrocher
`
		_ = os.WriteFile(taskAFaire, []byte(content), 0644)
	}

	// Tâche 3 : Inbox
	taskInbox := filepath.Join(root, "0_Inbox", "boite-a-idees-brutes.md")
	if _, err := os.Stat(taskInbox); os.IsNotExist(err) {
		content := `# Boîte à idées brutes 💡

Déposez ici toutes les idées, impulsions ou inspirations sans chercher à les structurer.

### Idées en vrac :
- [ ] Explorer les playlists Synthwave sur YouTube Music
- [ ] Synchroniser mon dossier Kanban avec Git
`
		_ = os.WriteFile(taskInbox, []byte(content), 0644)
	}

	return nil
}

// LoadBoard charge l'ensemble des colonnes et des tâches Markdown du répertoire racine.
func LoadBoard(root string) (*Board, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := InitDefaultWorkspace(root); err != nil {
			return nil, err
		}
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var dirEntries []os.DirEntry
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			dirEntries = append(dirEntries, e)
		}
	}

	// Si le dossier existe mais ne contient aucune colonne, initialiser l'arborescence par défaut
	if len(dirEntries) == 0 {
		if err := InitDefaultWorkspace(root); err != nil {
			return nil, err
		}
		entries, err = os.ReadDir(root)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
				dirEntries = append(dirEntries, e)
			}
		}
	}

	// Tri des dossiers par ordre alphabétique (respecte 0_..., 1_..., 2_...)
	sort.Slice(dirEntries, func(i, j int) bool {
		return dirEntries[i].Name() < dirEntries[j].Name()
	})

	board := &Board{
		RootPath: root,
		Columns:  make([]*Column, 0),
	}

	for _, d := range dirEntries {
		dirName := d.Name()
		colPath := filepath.Join(root, dirName)

		col := &Column{
			DirName:     dirName,
			DisplayName: GetColumnDisplayName(dirName),
			Tasks:       make([]*Task, 0),
		}

		files, err := os.ReadDir(colPath)
		if err == nil {
			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".md") {
					filePath := filepath.Join(colPath, f.Name())
					info, err := f.Info()
					var modTime time.Time
					if err == nil {
						modTime = info.ModTime()
					}

					data, err := os.ReadFile(filePath)
					if err == nil {
						task := ParseMarkdownTask(string(data), f.Name(), dirName, filePath, modTime)
						col.Tasks = append(col.Tasks, task)
					}
				}
			}
		}

		board.Columns = append(board.Columns, col)
	}

	return board, nil
}

// sanitizeFileName assainit le nom pour respecter les contraintes strictes de Windows, macOS et Linux.
func sanitizeFileName(name string) string {
	// Caractères formellement interdits sous Windows : \ / : * ? " < > |
	forbidden := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	for _, ch := range forbidden {
		name = strings.ReplaceAll(name, ch, "-")
	}
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "'", "-")
	name = strings.Trim(name, ".- ")

	// Noms de périphériques réservés sous Windows (CON, PRN, AUX, NUL, COM1..9, LPT1..9)
	upper := strings.ToUpper(name)
	reserved := map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true,
		"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
	}
	if reserved[upper] {
		name = name + "-task"
	}
	if name == "" {
		name = "tache"
	}
	if len(name) > 35 {
		name = name[:35]
	}
	return strings.Trim(name, ".- ")
}

// MoveTask déplace physiquement une tâche d'un répertoire colonne vers un autre.
func (b *Board) MoveTask(task *Task, targetColDir string) error {
	if task.ColumnDir == targetColDir {
		return nil
	}

	oldPath := task.Path
	newPath := filepath.Join(b.RootPath, targetColDir, task.FileName)

	// Sous Windows, si le fichier cible existe déjà, os.Rename échoue : on supprime d'abord
	if _, err := os.Stat(newPath); err == nil && oldPath != newPath {
		_ = os.Remove(newPath)
	}

	// Déplacement atomique du fichier sur le système de fichiers
	if err := os.Rename(oldPath, newPath); err != nil {
		// Repli si déplacement entre deux disques/volumes différents sous Windows (EXDEV)
		if errCopy := copyAndRemoveFile(oldPath, newPath); errCopy != nil {
			return err
		}
	}

	// Retrait de l'ancienne colonne
	for _, col := range b.Columns {
		if col.DirName == task.ColumnDir {
			newTasks := make([]*Task, 0, len(col.Tasks))
			for _, t := range col.Tasks {
				if t.FileName != task.FileName {
					newTasks = append(newTasks, t)
				}
			}
			col.Tasks = newTasks
			break
		}
	}

	// Ajout à la nouvelle colonne
	task.Path = newPath
	task.ColumnDir = targetColDir
	task.ModifiedAt = time.Now()

	for _, col := range b.Columns {
		if col.DirName == targetColDir {
			col.Tasks = append(col.Tasks, task)
			break
		}
	}

	return nil
}

// SaveTask écrit le contenu mis à jour de la tâche dans son fichier Markdown.
func (b *Board) SaveTask(task *Task) error {
	content := task.SerializeContent()
	return os.WriteFile(task.Path, []byte(content), 0644)
}

func copyAndRemoveFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	_ = os.Remove(src)
	return nil
}

// CreateTask crée un nouveau fichier Markdown de tâche dans la colonne spécifiée.
func (b *Board) CreateTask(colDir, title, description string, subtasks []string) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "Nouvelle Tâche"
	}

	// Génération d'un nom de fichier sécurisé pour tous les OS
	slug := sanitizeFileName(strings.ToLower(title))

	fileName := fmt.Sprintf("%s.md", slug)
	targetPath := filepath.Join(b.RootPath, colDir, fileName)

	// Éviter l'écrasement en cas de doublon de nom
	counter := 1
	for {
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			break
		}
		fileName = fmt.Sprintf("%s-%d.md", slug, counter)
		targetPath = filepath.Join(b.RootPath, colDir, fileName)
		counter++
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	if description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", description))
	}

	if len(subtasks) > 0 {
		sb.WriteString("### Sous-tâches :\n")
		for _, st := range subtasks {
			st = strings.TrimSpace(st)
			if st != "" {
				sb.WriteString(fmt.Sprintf("- [ ] %s\n", st))
			}
		}
	} else {
		sb.WriteString("### Sous-tâches :\n- [ ] Première étape atomique\n")
	}

	if err := os.WriteFile(targetPath, []byte(sb.String()), 0644); err != nil {
		return nil, err
	}

	task := ParseMarkdownTask(sb.String(), fileName, colDir, targetPath, time.Now())

	// Ajout à la colonne correspondante dans le Board en mémoire
	for _, col := range b.Columns {
		if col.DirName == colDir {
			col.Tasks = append(col.Tasks, task)
			break
		}
	}

	return task, nil
}

// DeleteTask supprime définitivement le fichier Markdown de la tâche.
func (b *Board) DeleteTask(task *Task) error {
	if err := os.Remove(task.Path); err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, col := range b.Columns {
		if col.DirName == task.ColumnDir {
			newTasks := make([]*Task, 0, len(col.Tasks))
			for _, t := range col.Tasks {
				if t.FileName != task.FileName {
					newTasks = append(newTasks, t)
				}
			}
			col.Tasks = newTasks
			break
		}
	}

	return nil
}

// SaveQuickThought décharge une pensée parasite dans le fichier _pensees_rapides.md de la première colonne.
func (b *Board) SaveQuickThought(thought string) error {
	if len(b.Columns) == 0 {
		return fmt.Errorf("aucune colonne disponible")
	}

	firstCol := b.Columns[0].DirName
	filePath := filepath.Join(b.RootPath, firstCol, "_pensees_rapides.md")

	timestamp := time.Now().Format("2006-01-02 15:04")
	entry := fmt.Sprintf("- [ ] 💡 [%s] %s\n", timestamp, strings.TrimSpace(thought))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		header := "# 💡 Pensées Parasites & Idées Éphémères (Brain Dump)\n\n"
		return os.WriteFile(filePath, []byte(header+entry), 0644)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

// ArchiveTask déplace une tâche vers le dossier des archives pour désencombrer l'espace visuel sans perte d'information.
func (b *Board) ArchiveTask(task *Task, archiveDir string) error {
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}

	// Ajouter une note d'archivage dans le contenu du fichier
	archiveTime := time.Now().Format("2006-01-02 15:04")
	archiveBadge := fmt.Sprintf("> **📦 Archivé le :** %s (depuis %s)", archiveTime, task.ColumnDir)
	task.BodyLines = append([]string{archiveBadge, ""}, task.BodyLines...)
	_ = b.SaveTask(task)

	dstPath := filepath.Join(archiveDir, task.FileName)

	// Éviter l'écrasement en cas de doublon
	if _, err := os.Stat(dstPath); err == nil && task.Path != dstPath {
		timeTag := time.Now().Format("20060102-150405")
		nameWithoutExt := strings.TrimSuffix(task.FileName, ".md")
		dstPath = filepath.Join(archiveDir, fmt.Sprintf("%s-%s.md", nameWithoutExt, timeTag))
	}

	if err := os.Rename(task.Path, dstPath); err != nil {
		if errCopy := copyAndRemoveFile(task.Path, dstPath); errCopy != nil {
			return err
		}
	}

	// Retirer la tâche de sa colonne en mémoire
	for _, col := range b.Columns {
		if col.DirName == task.ColumnDir {
			newTasks := make([]*Task, 0, len(col.Tasks))
			for _, t := range col.Tasks {
				if t.FileName != task.FileName {
					newTasks = append(newTasks, t)
				}
			}
			col.Tasks = newTasks
			break
		}
	}

	task.Path = dstPath
	task.ColumnDir = "archive"
	task.ModifiedAt = time.Now()
	return nil
}

// LoadArchivedTasks lit l'ensemble des tâches archivées.
func LoadArchivedTasks(archiveDir string) ([]*Task, error) {
	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		return []*Task{}, nil
	}

	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil, err
	}

	var list []*Task
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			fullPath := filepath.Join(archiveDir, e.Name())
			info, _ := e.Info()
			var modTime time.Time
			if info != nil {
				modTime = info.ModTime()
			}
			data, err := os.ReadFile(fullPath)
			if err == nil {
				task := ParseMarkdownTask(string(data), e.Name(), "archive", fullPath, modTime)
				list = append(list, task)
			}
		}
	}

	// Trier par date de modification décroissante (les plus récemment archivées en premier)
	sort.Slice(list, func(i, j int) bool {
		return list[i].ModifiedAt.After(list[j].ModifiedAt)
	})

	return list, nil
}

// RestoreTask restaure une tâche archivée vers une colonne du Kanban.
func (b *Board) RestoreTask(task *Task, targetColDir string) error {
	if len(b.Columns) == 0 {
		return fmt.Errorf("aucune colonne disponible pour la restauration")
	}

	targetColExists := false
	for _, col := range b.Columns {
		if col.DirName == targetColDir {
			targetColExists = true
			break
		}
	}
	if !targetColExists {
		targetColDir = b.Columns[0].DirName
	}

	dstPath := filepath.Join(b.RootPath, targetColDir, task.FileName)
	// Éviter l'écrasement
	counter := 1
	slug := strings.TrimSuffix(task.FileName, ".md")
	for {
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			break
		}
		task.FileName = fmt.Sprintf("%s-restaure-%d.md", slug, counter)
		dstPath = filepath.Join(b.RootPath, targetColDir, task.FileName)
		counter++
	}

	if err := os.Rename(task.Path, dstPath); err != nil {
		if errCopy := copyAndRemoveFile(task.Path, dstPath); errCopy != nil {
			return err
		}
	}

	task.Path = dstPath
	task.ColumnDir = targetColDir
	task.ModifiedAt = time.Now()

	for _, col := range b.Columns {
		if col.DirName == targetColDir {
			col.Tasks = append(col.Tasks, task)
			break
		}
	}

	return nil
}
