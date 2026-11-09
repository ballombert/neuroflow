package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderBrainDumpModal affiche la boîte de dialogue express pour décharger une pensée parasite.
func (a *App) RenderBrainDumpModal() string {
	totalW := a.Width
	if totalW < 70 {
		totalW = 70
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		Render("💡 BOUCLIER ANTI-INTERRUPTION (Brain Dump)")

	explanation := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render("Une pensée parasite vous distrait ? Déchargez votre mémoire de travail ici pour libérer votre esprit.\nL'idée sera enregistrée et prête à être traitée plus tard sans briser votre session.")

	inputView := a.BrainDumpInput.View()

	footer := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render("[Entrée] Sauvegarder & Reprendre le Focus   •   [Échap] Annuler")

	modalContent := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", header, explanation, inputView, footer)

	return lipgloss.Place(
		a.Width,
		a.Height,
		lipgloss.Center,
		lipgloss.Center,
		StyleModal.Width(64).Render(modalContent),
	)
}
