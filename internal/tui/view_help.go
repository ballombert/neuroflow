package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHelpView affiche le panneau d'aide des raccourcis et la philosophie scientifique TDAH.
func (a *App) RenderHelpView() string {
	totalW := a.Width
	if totalW < 80 {
		totalW = 80
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render("📖 GUIDE, RACCOURCIS & ASSISTANCE TDAH (NeuroFlow)")

	// Table des raccourcis
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"Tab", "Basculer entre le tableau Kanban et l'écran des Habitudes (XP & Gamification)"},
		{"a / c", "Ouvrir l'Assistant IA Copilot (découpage atomique de tâche, kickoff pro, déblocage)"},
		{"f", "Lancer une session de Focus sur la tâche sélectionnée (par défaut 25 min)"},
		{"s", "Lancer le Soft-Launch (5 min) : défi d'amorce douce pour vaincre la paralysie"},
		{"Espace", "Cocher / décocher la sous-tâche (ou l'habitude du jour dans l'écran Habitudes)"},
		{"h / l ou ⬅ / ➔", "Naviguer entre les colonnes du Kanban"},
		{"j / k ou ⬆ / ⬇", "Naviguer entre les tâches ou sous-tâches"},
		{"H / L ou Maj+Flèches", "Déplacer physiquement la tâche vers la colonne gauche / droite (os.Rename)"},
		{"+ / f (en focus)", "Flow Extension (+5m ou +15m) : protège l'état d'hyperfocus"},
		{"p (en focus)", "Mettre en pause ou reprendre la session sans culpabilité"},
		{"b", "Brain Dump : décharger instantanément une idée parasite sans perte de contexte"},
		{"m", "Lancer YouTube Music (ouvre le flux audio dans votre navigateur/app)"},
		{"n", "Créer une nouvelle tâche Markdown dans la colonne active"},
		{"x", "Archiver la tâche sélectionnée vers le dossier .archive (espace visuel allégé)"},
		{"X / Maj+X", "Consulter les Archives et restaurer d'anciennes tâches dans le Kanban"},
		{"d", "Supprimer définitivement le fichier Markdown de la tâche"},
		{"Esc / k", "Basculer entre la vue Focus et la vue Kanban"},
		{"? / F1", "Ouvrir ou fermer ce guide d'aide"},
		{"q / Ctrl+C", "Quitter NeuroFlow"},
	}

	var sb strings.Builder
	for _, sc := range shortcuts {
		k := lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Width(24).Render(sc.key)
		d := lipgloss.NewStyle().Foreground(ColorText).Render(sc.desc)
		sb.WriteString(fmt.Sprintf("  %s %s\n", k, d))
	}

	// Synthèse neuroscientifique
	scienceHeader := lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Render("\n🧠 FONDEMENTS SCIENTIFIQUES TDAH, PROCRASTINATION & VIE PRO :")
	scienceText := lipgloss.NewStyle().Foreground(ColorMuted).Render(`
  • Formation des Habitudes & TDAH (Dr Russell Barkley, BJ Fogg, James Clear) :
    Les automatismes striataux étant défaillants, une habitude n'est jamais 'automatique'.
    Elle nécessite des intentions d'implémentation : Si [Déclencheur], Alors [Micro-action].
  • Gamification Non Punitive : Les ruptures de séries brutales créent une spirale de honte.
    Dans NeuroFlow, les XP sont cumulatifs et acquis à vie (pas de remise à zéro).
  • Gestion de la Vie Professionnelle :
    - Règle des 3 : Choisir 3 victoires quotidiennes au kickoff matinal.
    - Shutdown Ritual : Clôture formelle de la journée à 18h pour préserver la vie perso.
  • Assistant Copilot IA : Découpe les tâches pro intimidantes en micro-étapes atomiques
    ridicules pour désamorcer la résistance de l'amygdale.
`)

	footer := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render("\n[Échap] ou [?] pour fermer cette fenêtre d'aide.")

	content := fmt.Sprintf("%s\n\n%s%s%s%s", title, sb.String(), scienceHeader, scienceText, footer)

	return lipgloss.Place(
		a.Width,
		a.Height,
		lipgloss.Center,
		lipgloss.Center,
		StyleModal.Width(totalW-6).Render(content),
	)
}
