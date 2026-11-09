package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderAICoachModal affiche la boîte de dialogue de l'assistant d'aide IA Copilot.
func (a *App) RenderAICoachModal() string {
	totalW := a.Width
	if totalW < 75 {
		totalW = 75
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render("🤖 ASSISTANT D'AIDE COPILOT & COACH TDAH")

	task := a.CurrentTask()
	taskTitle := "Aucune tâche sélectionnée"
	if task != nil {
		taskTitle = task.Title
	}

	taskBanner := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		Render(fmt.Sprintf("Tâche ciblée : ▶ %s", taskTitle))

	var body string
	if a.AICoachResult != "" {
		// Affichage du résultat obtenu
		resultHeader := lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Render("✨ Résultat de l'Assistant :")
		resultBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Padding(1, 2).
			Width(68).
			Render(a.AICoachResult)

		footer := lipgloss.NewStyle().
			Foreground(ColorMuted).
			Render("[Entrée] ou [Échap] Fermer et revenir au tableau")

		body = fmt.Sprintf("%s\n\n%s\n\n%s", resultHeader, resultBox, footer)
	} else if a.AICoachLoading {
		// En cours de réflexion
		loadingMsg := lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true).
			Render("⏳ Consultation de Copilot en cours... Décomposition atomique...")
		body = loadingMsg
	} else {
		// Menu des actions proposées par l'assistant
		menuTitle := lipgloss.NewStyle().Foreground(ColorText).Render("Choisissez une intervention de l'assistant :")

		options := []string{
			" [1] ✂️  Découper cette tâche en micro-actions avec Copilot (injecte les - [ ] dans le .md)",
			" [2] 🌟 Body Doubling Virtuel (Conseil d'amorce douce pour vaincre l'inertie)",
			" [3] 💼 Kickoff Pro (Règle des 3 : analyser le tableau et choisir les 3 victoires du jour)",
		}

		styledOptions := lipgloss.NewStyle().Foreground(ColorText).Render(strings.Join(options, "\n\n"))

		copilotStatus := "GitHub Copilot CLI : Détecté & Actif 🟢"
		if !a.Assistant.IsCopilotAvailable() {
			copilotStatus = "GitHub Copilot CLI : Non détecté (mode de repli heuristique TDAH actif) 🟡"
		}
		statusStyled := lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render(copilotStatus)

		footer := lipgloss.NewStyle().
			Foreground(ColorMuted).
			Render("[1/2/3] Choisir une option   •   [Échap] Annuler")

		body = fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", menuTitle, styledOptions, statusStyled, footer)
	}

	modalContent := fmt.Sprintf("%s\n%s\n\n%s", header, taskBanner, body)

	return lipgloss.Place(
		a.Width,
		a.Height,
		lipgloss.Center,
		lipgloss.Center,
		StyleModal.Width(72).Render(modalContent),
	)
}
