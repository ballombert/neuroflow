package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderNewTaskModal affiche le formulaire de création d'une nouvelle tâche Markdown.
func (a *App) RenderNewTaskModal() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render("✨ NOUVELLE TÂCHE (Markdown)")

	targetCol := "0_Inbox"
	if a.Board != nil && a.ActiveColIdx < len(a.Board.Columns) {
		targetCol = a.Board.Columns[a.ActiveColIdx].DisplayName
	}

	colInfo := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(fmt.Sprintf("Colonne cible : %s", targetCol))

	prompt := lipgloss.NewStyle().
		Foreground(ColorText).
		Render("Titre de la tâche ou micro-étape :")

	inputView := a.NewTaskInput.View()

	tip := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Italic(true).
		Render("💡 Conseil TDAH : Nommez une première action ridicule (ex: 'Créer le fichier' plutôt que 'Finir le projet').")

	footer := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render("[Entrée] Créer le fichier .md   •   [Échap] Annuler")

	modalContent := fmt.Sprintf("%s\n%s\n\n%s\n%s\n\n%s\n\n%s", header, colInfo, prompt, inputView, tip, footer)

	return lipgloss.Place(
		a.Width,
		a.Height,
		lipgloss.Center,
		lipgloss.Center,
		StyleModal.Width(66).Render(modalContent),
	)
}
