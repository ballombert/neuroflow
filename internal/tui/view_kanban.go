package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"neuroflow/internal/kanban"
	"neuroflow/internal/pomodoro"
)

// RenderKanbanView affiche le tableau Kanban multi-colonnes.
func (a *App) RenderKanbanView() string {
	if a.Board == nil || len(a.Board.Columns) == 0 {
		return "Aucune colonne disponible dans le dossier Kanban."
	}

	numCols := len(a.Board.Columns)
	totalW := a.Width
	if totalW < 80 {
		totalW = 80
	}

	// Largeur par colonne calculée dynamiquement
	colWidth := (totalW - (numCols * 3)) / numCols
	if colWidth < 22 {
		colWidth = 22
	}

	var renderedColumns []string

	for cIdx, col := range a.Board.Columns {
		isColActive := cIdx == a.ActiveColIdx

		// Entête de colonne
		headerStyle := StyleColumnHeader
		if isColActive {
			headerStyle = StyleColumnActiveHeader
		}

		headerText := fmt.Sprintf("%s (%d)", col.DisplayName, len(col.Tasks))
		header := headerStyle.Width(colWidth).Render(headerText)

		// Liste des cartes de tâches
		var cards []string
		if len(col.Tasks) == 0 {
			emptyMsg := lipgloss.NewStyle().
				Foreground(ColorMuted).
				Italic(true).
				Padding(1, 1).
				Render("Aucune tâche ici")
			cards = append(cards, emptyMsg)
		} else {
			for tIdx, task := range col.Tasks {
				isTaskActive := isColActive && (tIdx == a.ActiveTaskIdx)
				card := a.renderTaskCard(task, isTaskActive, colWidth-4)
				cards = append(cards, card)
			}
		}

		columnContent := lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinVertical(lipgloss.Left, cards...))
		colBox := lipgloss.NewStyle().
			Width(colWidth).
			MarginRight(1).
			Render(columnContent)

		renderedColumns = append(renderedColumns, colBox)
	}

	boardView := lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...)

	// Barre d'état et aperçu de la tâche sélectionnée
	preview := a.renderSelectedTaskPreview(totalW)
	statusBar := a.renderKanbanStatusBar(totalW)

	return lipgloss.JoinVertical(lipgloss.Left, boardView, preview, statusBar)
}

// renderTaskCard génère le bloc visuel d'une tâche dans une colonne.
func (a *App) renderTaskCard(task *kanban.Task, isActive bool, width int) string {
	style := StyleCard
	if isActive {
		style = StyleActiveCard
	}

	// Titre tronqué si nécessaire
	title := task.Title
	if len(title) > width-2 && width > 5 {
		title = title[:width-5] + "..."
	}

	titleStyled := lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render(title)
	if isActive {
		titleStyled = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render("▶ " + title)
	}

	// Badge d'avancement des sous-tâches
	var progressLine string
	if len(task.Subtasks) > 0 {
		doneCount := 0
		for _, st := range task.Subtasks {
			if st.Done {
				doneCount++
			}
		}
		pct := int(float64(doneCount) / float64(len(task.Subtasks)) * 100)

		// Mini jauge à points ou carrés
		miniBar := renderMiniBar(doneCount, len(task.Subtasks))
		progressLine = fmt.Sprintf("%s %d%%", miniBar, pct)
		progressLine = StyleProgressBadge.Render(progressLine)
	} else {
		progressLine = lipgloss.NewStyle().Foreground(ColorMuted).Render("0 sous-tâche")
	}

	cardContent := fmt.Sprintf("%s\n%s", titleStyled, progressLine)
	return style.Width(width).Render(cardContent)
}

func renderMiniBar(done, total int) string {
	if total == 0 {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < total && i < 8; i++ {
		if i < done {
			sb.WriteString("■")
		} else {
			sb.WriteString("□")
		}
	}
	if total > 8 {
		sb.WriteString("…")
	}
	return sb.String()
}

// renderSelectedTaskPreview affiche les détails et sous-tâches de la tâche courante.
func (a *App) renderSelectedTaskPreview(totalW int) string {
	task := a.CurrentTask()
	if task == nil {
		return ""
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorMuted).
		Padding(0, 1).
		MarginTop(1).
		Width(totalW - 2)

	header := fmt.Sprintf("📄 Aperçu : %s (%s)", task.Title, task.ColumnDir)
	headerStyled := lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render(header)

	var subtaskList []string
	if len(task.Subtasks) == 0 {
		subtaskList = append(subtaskList, lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render("  (Aucune sous-tâche. Appuyez sur 'f' pour démarrer un focus ou 'n' pour en ajouter)"))
	} else {
		for i, st := range task.Subtasks {
			check := "[ ]"
			itemStyle := lipgloss.NewStyle().Foreground(ColorText)
			if st.Done {
				check = "[x]"
				itemStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Strikethrough(true)
			}
			isSubtaskFocused := i == a.ActiveSubtaskIdx
			prefix := "  "
			if isSubtaskFocused {
				prefix = "❯ "
				itemStyle = itemStyle.Bold(true).Foreground(ColorPrimary)
			}
			subtaskList = append(subtaskList, fmt.Sprintf("%s%s %s", prefix, check, itemStyle.Render(st.Text)))
		}
	}

	content := headerStyled + "\n" + strings.Join(subtaskList, "\n")
	return boxStyle.Render(content)
}

// renderKanbanStatusBar affiche la barre d'aide et l'état du timer en tâche de fond.
func (a *App) renderKanbanStatusBar(totalW int) string {
	var timerIndicator string
	if a.Timer.State != pomodoro.StateIdle {
		rem := a.Timer.RemainingDuration()
		taskName := ""
		if a.Timer.ActiveTask != nil {
			taskName = " - " + a.Timer.ActiveTask.Title
		}
		timerIndicator = fmt.Sprintf("[%s: %s%s] ", a.Timer.State.String(), pomodoro.FormatTimeDigital(rem), taskName)
		timerIndicator = lipgloss.NewStyle().Bold(true).Foreground(ColorWarning).Background(ColorSubtle).Render(timerIndicator)
	}

	// Message de notification flash (ex: "Sous-tâche cochée !")
	var flashMsg string
	if a.FlashMessage != "" {
		flashMsg = StyleNotification.Render("✨ " + a.FlashMessage + " ")
	}

	shortcuts := " [f] Focus  [s] Soft-Launch  [Tab] Habitudes  [x] Archiver  [X] Archives  [a] IA  [Espace] Cocher  [H/L] Déplacer  [n] Nouveau  [b] Dump  [m] Musique  [?] Aide"
	styledShortcuts := lipgloss.NewStyle().Foreground(ColorMuted).Render(shortcuts)

	fullLine := timerIndicator + flashMsg + styledShortcuts
	return lipgloss.NewStyle().MarginTop(1).Width(totalW).Render(fullLine)
}
