package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Palette de couleurs neuro-apaisantes et hautement contrastées (Tokyo Night / Nord)
var (
	ColorPrimary   = lipgloss.Color("#7AA2F7") // Bleu doux attentif
	ColorSecondary = lipgloss.Color("#BB9AF7") // Lavande / Flow
	ColorSuccess   = lipgloss.Color("#9ECE6A") // Vert sauge dopaminergique
	ColorWarning   = lipgloss.Color("#E0AF68") // Ambre apaisant
	ColorDanger    = lipgloss.Color("#F7768E") // Corail doux
	ColorMuted     = lipgloss.Color("#565F89") // Gris bleuté discret
	ColorText      = lipgloss.Color("#C0CAF5") // Blanc doux lisible
	ColorSubtle    = lipgloss.Color("#24283B") // Fond sombre profond
	ColorCardBg    = lipgloss.Color("#1A1B26") // Fond de carte
	ColorActiveCard= lipgloss.Color("#2E3440") // Fond carte sélectionnée
)

// Styles Lipgloss réutilisables
var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1)

	StyleColumnHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(ColorMuted).
				Padding(0, 1)

	StyleColumnActiveHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSecondary).
				BorderStyle(lipgloss.ThickBorder()).
				BorderBottom(true).
				BorderForeground(ColorSecondary).
				Padding(0, 1)

	StyleCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Background(ColorCardBg).
			Padding(0, 1).
			MarginBottom(1)

	StyleActiveCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Background(ColorActiveCard).
			Padding(0, 1).
			MarginBottom(1)

	StyleProgressBadge = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)

	StyleAttentionAnchor = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#3D59A1")).
				Padding(0, 2).
				MarginBottom(1)

	StyleTimerBig = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(1, 2)

	StyleSoftLaunchBadge = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSuccess).
				Background(ColorSubtle).
				Padding(0, 1)

	StyleFocusBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(ColorSubtle).
			Padding(0, 1)

	StyleFlowBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(ColorSubtle).
			Padding(0, 1)

	StyleBreakBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning).
			Background(ColorSubtle).
			Padding(0, 1)

	StyleKeyBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorMuted).
			Padding(0, 1)

	StyleHelpKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	StyleHelpDesc = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleNotification = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSuccess).
				Background(ColorSubtle).
				Padding(0, 1)

	StyleModal = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorSecondary).
			Background(ColorCardBg).
			Padding(1, 2)
)
