package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"neuroflow/internal/habits"
)

// RenderHabitsView affiche l'écran de suivi des habitudes et de gamification TDAH.
func (a *App) RenderHabitsView() string {
	totalW := a.Width
	if totalW < 80 {
		totalW = 80
	}

	profile := habits.CalculateProfile(a.Habits)

	// 1. BANNIÈRE PROFIL JOUEUR & DOPAMINE
	rankBanner := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#414868")).
		Padding(0, 2).
		MarginBottom(1).
		Width(totalW - 4).
		Render(fmt.Sprintf("%s (Niveau %d)   •   XP Total : %d   •   Prochain palier : %d XP",
			profile.Title, profile.Level, profile.TotalXP, profile.NextXP,
		))

	xpBar := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true).
		Render(profile.XPBar(35))

	headerBlock := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render("🌱 RITUELS & HABITUDES TDAH (Fichiers Markdown hors Kanban)"),
		rankBanner,
		"Progression vers le niveau suivant : "+xpBar,
		lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render("💡 Philosophie TDAH : Vos XP sont acquis à vie. Pas de remise à zéro punitive en cas d'oubli."),
		"",
	)

	// 2. LISTE DES HABITUDES
	var habitCards []string
	if len(a.Habits) == 0 {
		emptyMsg := lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render("Aucune habitude trouvée dans " + a.Config.HabitsDir)
		habitCards = append(habitCards, emptyMsg)
	} else {
		for i, h := range a.Habits {
			isSelected := i == a.ActiveHabitIdx
			card := a.renderHabitCard(h, isSelected, totalW-6)
			habitCards = append(habitCards, card)
		}
	}

	habitsContent := lipgloss.JoinVertical(lipgloss.Left, habitCards...)

	// 3. BARRE D'ÉTAT HABITUDES
	var flashLine string
	if a.FlashMessage != "" {
		flashLine = StyleNotification.Render("✨ " + a.FlashMessage) + "\n"
	}

	shortcuts := " [Espace] Cocher/Valider aujourd'hui  [Tab] Retour Kanban  [j/k] Naviguer  [m] Musique  [?] Aide  [q] Quitter"
	statusBar := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(ColorMuted).
		Render(shortcuts)

	viewContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBlock,
		habitsContent,
		flashLine,
		statusBar,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(viewContent)
}

func (a *App) renderHabitCard(h *habits.Habit, isSelected bool, width int) string {
	style := StyleCard
	if isSelected {
		style = StyleActiveCard
	}

	// Case à cocher pour aujourd'hui
	check := "[ ]"
	statusStyle := lipgloss.NewStyle().Foreground(ColorText)
	if h.IsDoneToday() {
		check = "[x]"
		statusStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	}

	prefix := "  "
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	if isSelected {
		prefix = "❯ "
		titleStyle = titleStyle.Foreground(ColorPrimary)
	}

	firstLine := fmt.Sprintf("%s%s %s  %s", prefix, check, h.CategoryEmoji(), titleStyle.Render(h.Title))

	// Formule d'intention d'implémentation (Gollwitzer) : Si... Alors...
	intention := fmt.Sprintf("    ↳ Si : %s\n    ↳ Alors : %s", h.Trigger, h.MicroAction)
	styledIntention := lipgloss.NewStyle().Foreground(ColorMuted).Render(intention)

	// Métriques de streak et d'XP
	streakStr := fmt.Sprintf("🔥 Série : %d jours", h.Streak)
	xpStr := fmt.Sprintf("⭐ +%d XP", h.XPValue)
	metaLine := fmt.Sprintf("    %s   •   %s   •   Fichier : %s",
		lipgloss.NewStyle().Foreground(ColorWarning).Render(streakStr),
		lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render(xpStr),
		lipgloss.NewStyle().Foreground(ColorMuted).Render(h.FileName),
	)

	content := fmt.Sprintf("%s\n%s\n%s", statusStyle.Render(firstLine), styledIntention, metaLine)
	return style.Width(width).Render(content)
}
