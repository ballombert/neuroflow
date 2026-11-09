package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"neuroflow/internal/pomodoro"
)

// RenderFocusView affiche la vue immersive de concentration anti-cécité temporelle.
func (a *App) RenderFocusView() string {
	totalW := a.Width
	if totalW < 80 {
		totalW = 80
	}

	timer := a.Timer
	task := timer.ActiveTask

	// 1. ANCRE ATTENTIONNELLE (Résout la perte de mémoire de travail immédiate)
	anchorTitle := "Session de Focus Libre"
	if task != nil {
		anchorTitle = task.Title
	}

	anchorBanner := StyleAttentionAnchor.Width(totalW - 4).Render(
		fmt.Sprintf("🎯 ANCRE ATTENTIONNELLE : Qu'est-ce que je fais MAINTENANT ?\n   ▶ %s", anchorTitle),
	)

	// 2. BADGE D'ÉTAT DU CYCLE
	var stateBadge string
	switch timer.State {
	case pomodoro.StateSoftLaunch:
		stateBadge = StyleSoftLaunchBadge.Render(" 🌱 RAMPE DOPAMINERGIQUE (Amorce 5 min sans pression) ")
	case pomodoro.StateFlowExtension:
		stateBadge = StyleFlowBadge.Render(fmt.Sprintf(" 🌊 HYPERFOCUS PROTÉGÉ (+%d min ajoutées) ", timer.ExtensionCount*5))
	case pomodoro.StateShortBreak:
		stateBadge = StyleBreakBadge.Render(" ☕ PAUSE COURTE (Respirer, s'étirer, boire de l'eau) ")
	case pomodoro.StateLongBreak:
		stateBadge = StyleBreakBadge.Render(" 🌴 PAUSE LONGUE RÉGÉNÉRATRICE ")
	default:
		stateBadge = StyleFocusBadge.Render(" 🔥 SESSION DE FOCUS INTENSE ")
	}

	if timer.IsPaused {
		pauseBadge := lipgloss.NewStyle().Bold(true).Foreground(ColorDanger).Render(" [⏸ EN PAUSE SANS JUGEMENT] ")
		stateBadge = stateBadge + " " + pauseBadge
	}

	// 3. MINUTEUR NUMÉRIQUE & PROGRESSION
	rem := timer.RemainingDuration()
	elap := timer.Elapsed
	pct := int(timer.ProgressRatio() * 100)

	digitalTime := pomodoro.FormatTimeDigital(rem)
	timeStats := fmt.Sprintf("%s restant   •   %s écoulé   •   %d%% accompli", digitalTime, pomodoro.FormatTimeDigital(elap), pct)
	styledDigital := lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(digitalTime)
	styledStats := lipgloss.NewStyle().Foreground(ColorMuted).Render(timeStats)

	// 4. JAUGE SPATIALE VISUELLE (Time Timer numérique)
	gaugeWidth := totalW - 10
	if gaugeWidth > 70 {
		gaugeWidth = 70
	}
	spatialGauge := timer.SpatialGauge(gaugeWidth)

	// Couleurs de la jauge adaptées selon le mode
	gaugeColor := ColorPrimary
	if timer.State == pomodoro.StateSoftLaunch {
		gaugeColor = ColorSuccess
	} else if timer.State == pomodoro.StateFlowExtension {
		gaugeColor = ColorSecondary
	} else if timer.State == pomodoro.StateShortBreak || timer.State == pomodoro.StateLongBreak {
		gaugeColor = ColorWarning
	}

	styledGauge := lipgloss.NewStyle().Foreground(gaugeColor).Bold(true).Render(spatialGauge)

	// 5. LISTE DES SOUS-TÂCHES INTERACTIVE DIRECTEMENT EN VUE FOCUS
	var subtasksBox string
	if task != nil && len(task.Subtasks) > 0 {
		var stLines []string
		stHeader := lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render("📋 Micro-actions atomiques (Espace pour cocher) :")
		stLines = append(stLines, stHeader)

		for i, st := range task.Subtasks {
			check := "[ ]"
			itemStyle := lipgloss.NewStyle().Foreground(ColorText)
			if st.Done {
				check = "[x]"
				itemStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Strikethrough(true)
			}

			isFocused := i == a.ActiveSubtaskIdx
			prefix := "  "
			if isFocused {
				prefix = "❯ "
				itemStyle = itemStyle.Bold(true).Foreground(ColorPrimary)
			}
			stLines = append(stLines, fmt.Sprintf("%s%s %s", prefix, check, itemStyle.Render(st.Text)))
		}

		subtasksBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Padding(1, 2).
			MarginTop(1).
			Width(totalW - 6).
			Render(strings.Join(stLines, "\n"))
	}

	// 6. SÉRIES & STREAKS DE DOPAMINE
	streakMsg := fmt.Sprintf("⭐ Victoires accomplies aujourd'hui : %d sessions focus", timer.CompletedSessions)
	styledStreak := lipgloss.NewStyle().Foreground(ColorSuccess).Italic(true).Render(streakMsg)

	// Notification flash (ex: Micro-célébration de sous-tâche validée)
	var flashLine string
	if a.FlashMessage != "" {
		flashLine = StyleNotification.Render("✨ " + a.FlashMessage) + "\n"
	}

	// 7. BARRE DE COMMANDES TDAH
	commands := []string{
		"[+] +5m Flow",
		"[f] +15m Wave",
		"[p] Pause",
		"[b] Brain Dump",
		"[m] YouTube Music",
		"[Esc] Kanban",
	}
	styledCmds := lipgloss.NewStyle().Foreground(ColorMuted).Render(strings.Join(commands, "   "))

	// Assemblage global centré et aéré
	timerBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		stateBadge,
		"",
		styledDigital,
		styledStats,
		"",
		styledGauge,
		"",
		styledStreak,
	)

	viewContent := lipgloss.JoinVertical(
		lipgloss.Left,
		anchorBanner,
		timerBlock,
		flashLine,
		subtasksBox,
		"",
		styledCmds,
	)

	return lipgloss.NewStyle().Padding(1, 2).Render(viewContent)
}
