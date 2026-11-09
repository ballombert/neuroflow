package kanban

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	// Regex pour détecter les cases à cocher Markdown : - [ ] ou - [x]
	subtaskRegex = regexp.MustCompile(`^(\s*[-*]\s*\[)([ xX])(\]\s+)(.*)$`)
	titleRegex   = regexp.MustCompile(`^#\s+(.+)$`)
)

// ParseMarkdownTask lit le contenu Markdown et génère la structure Task.
func ParseMarkdownTask(content string, fileName, columnDir, fullPath string, modTime time.Time) *Task {
	task := &Task{
		FileName:   fileName,
		Path:       fullPath,
		ColumnDir:  columnDir,
		ModifiedAt: modTime,
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	subtaskIdx := 0

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		task.BodyLines = append(task.BodyLines, line)

		// Détection du titre principal
		if task.Title == "" {
			if matches := titleRegex.FindStringSubmatch(line); len(matches) > 1 {
				task.Title = strings.TrimSpace(matches[1])
			}
		}

		// Détection d'une sous-tâche
		if matches := subtaskRegex.FindStringSubmatch(line); len(matches) == 5 {
			isDone := strings.ToLower(matches[2]) == "x"
			subtaskText := strings.TrimSpace(matches[4])

			task.Subtasks = append(task.Subtasks, Subtask{
				Index:   subtaskIdx,
				LineNum: lineNum,
				Text:    subtaskText,
				Done:    isDone,
			})
			subtaskIdx++
		}

		lineNum++
	}

	// Si aucun titre n'a été spécifié via '# ', on déduit le titre du nom de fichier
	if task.Title == "" {
		cleanName := strings.TrimSuffix(fileName, ".md")
		cleanName = strings.ReplaceAll(cleanName, "-", " ")
		cleanName = strings.ReplaceAll(cleanName, "_", " ")
		task.Title = strings.Title(cleanName)
	}

	return task
}

// ToggleSubtask bascule l'état d'une sous-tâche par son index ordinal.
func (t *Task) ToggleSubtask(index int) bool {
	if index < 0 || index >= len(t.Subtasks) {
		return false
	}

	st := &t.Subtasks[index]
	st.Done = !st.Done

	// Mise à jour de la ligne correspondante dans BodyLines
	if st.LineNum < len(t.BodyLines) {
		line := t.BodyLines[st.LineNum]
		if st.Done {
			line = regexp.MustCompile(`^(\s*[-*]\s*\[)[ ](\]\s+)`).ReplaceAllString(line, "${1}x${2}")
		} else {
			line = regexp.MustCompile(`^(\s*[-*]\s*\[)[xX](\]\s+)`).ReplaceAllString(line, "${1} ${2}")
		}
		t.BodyLines[st.LineNum] = line
	}

	t.ModifiedAt = time.Now()
	return true
}

// AddSubtask ajoute une nouvelle sous-tâche à la fin du fichier Markdown.
func (t *Task) AddSubtask(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	line := fmt.Sprintf("- [ ] %s", text)
	t.BodyLines = append(t.BodyLines, line)
	t.Subtasks = append(t.Subtasks, Subtask{
		Index:   len(t.Subtasks),
		LineNum: len(t.BodyLines) - 1,
		Text:    text,
		Done:    false,
	})
	t.ModifiedAt = time.Now()
}

// AppendBrainDump injecte une pensée parasite dans la section dédiée du fichier.
func (t *Task) AppendBrainDump(thought string) {
	thought = strings.TrimSpace(thought)
	if thought == "" {
		return
	}

	hasHeader := false
	for _, l := range t.BodyLines {
		if strings.Contains(l, "Pensées Parasites") || strings.Contains(l, "Brain Dump") {
			hasHeader = true
			break
		}
	}

	if !hasHeader {
		t.BodyLines = append(t.BodyLines, "", "### 💡 Pensées Parasites (Brain Dump)")
	}

	timestamp := time.Now().Format("15:04")
	thoughtLine := fmt.Sprintf("- [ ] 💡 [%s] %s", timestamp, thought)
	t.BodyLines = append(t.BodyLines, thoughtLine)
	t.Subtasks = append(t.Subtasks, Subtask{
		Index:   len(t.Subtasks),
		LineNum: len(t.BodyLines) - 1,
		Text:    fmt.Sprintf("💡 [%s] %s", timestamp, thought),
		Done:    false,
	})
	t.ModifiedAt = time.Now()
}

// SerializeContent recompose le contenu Markdown textuel complet.
func (t *Task) SerializeContent() string {
	return strings.Join(t.BodyLines, "\n") + "\n"
}
