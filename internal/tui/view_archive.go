package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"neuroflow/internal/kanban"
)

// RenderArchiveView affiche l'écran de consultation et restauration des tâches archivées.
func (a *App) RenderArchiveView() string {
	totalW := a.Width
	if totalW < 80 {
		totalW = 80
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWarning).
		Render("📦 ARCHIVES & JOURNAL DES VICTOIRES PASSÉES")

	subHeader := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render("Les tâches archivées sont conservées en toute sécurité dans votre dossier .archive sans encombrer votre champ visuel.\nConsultez vos accomplissements ou restaurez une tâche à tout moment.")

	var taskCards []string
	if len(a.ArchivedTasks) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true).
			Padding(2, 2).
			Render("Aucune tâche archivée pour le moment. Dans le Kanban, appuyez sur 'x' sur une tâche pour l'archiver.")
		taskCards = append(taskCards, emptyMsg)
	} else {
		for i, task := range a.ArchivedTasks {
			isSelected := i == a.ActiveArchiveIdx
			card := a.renderArchivedCard(task, isSelected, totalW-6)
			taskCards = append(taskCards, card)
		}
	}

	tasksContent := lipgloss.JoinVertical(lipgloss.Left, taskCards...)

	// Barre de commandes
	var flashLine string
	if a.FlashMessage != "" {
		flashLine = StyleNotification.Render("✨ " + a.FlashMessage) + "\n"
	}

	shortcuts := " [Entrée/u] Restaurer dans le Kanban  [d] Supprimer définitivement  [j/k] Naviguer  [Échap/X] Retour"
	statusBar := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(ColorMuted).
		Render(shortcuts)

	viewContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		subHeader,
		"",
		tasksContent,
		flashLine,
		statusBar,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(viewContent)
}

func (a *App) renderArchivedCard(task *kanban.Task, isSelected bool, width int) string {
	style := StyleCard
	if isSelected {
		style = StyleActiveCard
	}

	prefix := "  "
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	if isSelected {
		prefix = "❯ "
		titleStyle = titleStyle.Foreground(ColorPrimary)
	}

	titleLine := fmt.Sprintf("%s📦 %s", prefix, titleStyle.Render(task.Title))

	progressStr := task.ProgressText()
	metaLine := fmt.Sprintf("    Avancement : %s   •   Fichier : %s   •   Date : %s",
		lipgloss.NewStyle().Foreground(ColorSuccess).Render(progressStr),
		lipgloss.NewStyle().Foreground(ColorMuted).Render(task.FileName),
		lipgloss.NewStyle().Foreground(ColorWarning).Render(task.ModifiedAt.Format("2006-01-02 15:04")),
	)

	content := fmt.Sprintf("%s\n%s", titleLine, metaLine)
	return style.Width(width).Render(content)
}
