package assistant

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"neuroflow/internal/config"
)

var subtaskLineRegex = regexp.MustCompile(`^[-*]\s*\[([ xX])\]\s*(.*)$`)

// Assistant gère les interactions avec GitHub Copilot CLI ou le moteur de coaching intégré.
type Assistant struct {
	cfg config.AIAssistantConfig
}

// NewAssistant instancie le service d'aide IA.
func NewAssistant(cfg config.AIAssistantConfig) *Assistant {
	return &Assistant{cfg: cfg}
}

// IsCopilotAvailable vérifie si le binaire copilot est présent dans le PATH ou dans les dossiers standards de l'OS.
func (a *Assistant) IsCopilotAvailable() bool {
	bin := a.cfg.CopilotBin
	if bin == "" {
		bin = "copilot"
	}

	if _, err := exec.LookPath(bin); err == nil {
		return true
	}

	if runtime.GOOS == "windows" {
		for _, ext := range []string{".cmd", ".exe", ".bat"} {
			if _, err := exec.LookPath(bin + ext); err == nil {
				return true
			}
		}
		appData := os.Getenv("APPDATA")
		if appData != "" {
			if _, err := os.Stat(filepath.Join(appData, "npm", "copilot.cmd")); err == nil {
				return true
			}
		}
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			if _, err := os.Stat(filepath.Join(localAppData, "Programs", "copilot", "copilot.exe")); err == nil {
				return true
			}
		}
	} else {
		if _, errHomebrew := os.Stat("/opt/homebrew/bin/copilot"); errHomebrew == nil {
			return true
		}
	}

	return false
}

func (a *Assistant) getBinaryPath() string {
	bin := a.cfg.CopilotBin
	if bin == "" {
		bin = "copilot"
	}

	if strings.Contains(bin, "/") || strings.Contains(bin, "\\") {
		return bin
	}

	if p, err := exec.LookPath(bin); err == nil {
		return p
	}

	if runtime.GOOS == "windows" {
		for _, ext := range []string{".cmd", ".exe", ".bat"} {
			if p, err := exec.LookPath(bin + ext); err == nil {
				return p
			}
		}
		appData := os.Getenv("APPDATA")
		if appData != "" {
			target := filepath.Join(appData, "npm", "copilot.cmd")
			if _, err := os.Stat(target); err == nil {
				return target
			}
		}
	} else {
		if _, err := os.Stat("/opt/homebrew/bin/copilot"); err == nil {
			return "/opt/homebrew/bin/copilot"
		}
	}

	return bin
}

func (a *Assistant) getSystemEnv() []string {
	if runtime.GOOS == "windows" {
		return os.Environ()
	}
	// Unix / macOS
	return append(os.Environ(), "PATH=/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:"+os.Getenv("PATH"))
}

func (a *Assistant) runCopilotPrompt(ctx context.Context, prompt string) ([]byte, error) {
	bin := a.getBinaryPath()
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" && (strings.HasSuffix(strings.ToLower(bin), ".cmd") || strings.HasSuffix(strings.ToLower(bin), ".bat")) {
		cmd = exec.CommandContext(ctx, "cmd", "/c", bin, "-s", "-p", prompt)
	} else {
		cmd = exec.CommandContext(ctx, bin, "-s", "-p", prompt)
	}

	cmd.Env = a.getSystemEnv()
	return cmd.CombinedOutput()
}

// DecomposeTask découpe une tâche en micro-actions atomiques via Copilot ou par heuristique TDAH.
func (a *Assistant) DecomposeTask(title, description string) ([]string, string, error) {
	if !a.cfg.Enabled {
		return a.fallbackDecomposition(title), "Généré via moteur TDAH intégré (IA désactivée dans config)", nil
	}

	prompt := fmt.Sprintf(
		"Tu es un coach expert TDAH. Découpe la tâche suivante en exactement 3 à 5 micro-étapes atomiques, ultra-concrètes et à friction minimale (seuil d'amorce ridicule). Réponds UNIQUEMENT avec des lignes de cases à cocher au format '- [ ] Libellé'.\nTâche : %s\nContexte : %s",
		title, description,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	out, err := a.runCopilotPrompt(ctx, prompt)
	outputStr := string(out)

	if err == nil {
		subtasks := parseSubtasksFromText(outputStr)
		if len(subtasks) > 0 {
			return subtasks, "Découpé avec succès par GitHub Copilot ✨", nil
		}
	}

	// Si Copilot n'est pas authentifié ou indisponible, bascule douce vers l'heuristique TDAH
	fallback := a.fallbackDecomposition(title)
	sourceMsg := "Généré via l'heuristique TDAH intégrée (Copilot hors ligne ou non authentifié)"
	return fallback, sourceMsg, nil
}

// GetDailyKickoffAdvice analyse les tâches du jour et propose les 3 victoires prioritaires.
func (a *Assistant) GetDailyKickoffAdvice(taskTitles []string) (string, error) {
	if len(taskTitles) == 0 {
		return "Votre tableau est vide ! Créez votre première micro-tâche avec la touche 'n' pour amorcer la journée.", nil
	}

	if a.cfg.Enabled && a.IsCopilotAvailable() {
		tasksList := strings.Join(taskTitles, ", ")
		prompt := fmt.Sprintf(
			"Tu es un coach professionnel TDAH. Voici les tâches professionnelles de mon tableau : [%s]. Applique la 'Règle des 3' : choisis les 3 tâches prioritaires à plus fort impact pour aujourd'hui et donne un court conseil bienveillant de 2 phrases pour vaincre l'inertie matinale.",
			tasksList,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		out, err := a.runCopilotPrompt(ctx, prompt)
		if err == nil && len(strings.TrimSpace(string(out))) > 20 {
			return strings.TrimSpace(string(out)), nil
		}
	}

	// Heuristique intégrée
	var sb strings.Builder
	sb.WriteString("🎯 Règle des 3 Victoires Pro (Conseil TDAH) :\n\n")
	limit := 3
	if len(taskTitles) < limit {
		limit = len(taskTitles)
	}
	for i := 0; i < limit; i++ {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, taskTitles[i]))
	}
	sb.WriteString("\n💡 Conseil : Ne touchez à rien d'autre tant que la victoire n°1 n'est pas amorcée avec le Soft-Launch (5 min) !")
	return sb.String(), nil
}

// GetUnblockAdvice génère un conseil de déblocage immédiat (body doubling).
func (a *Assistant) GetUnblockAdvice(taskTitle string) string {
	return fmt.Sprintf(
		"🌟 Body Doubling Virtuel pour '%s' :\n\n"+
			"1. Réduisez l'enjeu : vous n'avez pas à finir cette tâche maintenant.\n"+
			"2. Lancez le Soft-Launch de 5 min (Touche 's').\n"+
			"3. Votre seul but pour les 5 prochaines minutes : ouvrir le fichier et lire 3 lignes sans jugement.",
		taskTitle,
	)
}

func parseSubtasksFromText(text string) []string {
	var results []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		line = strings.TrimSpace(line)
		if m := subtaskLineRegex.FindStringSubmatch(line); len(m) == 3 {
			content := strings.TrimSpace(m[2])
			if content != "" {
				results = append(results, content)
			}
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			content := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "))
			if content != "" && !strings.HasPrefix(content, "#") {
				results = append(results, content)
			}
		}
	}
	return results
}

func (a *Assistant) fallbackDecomposition(title string) []string {
	return []string{
		fmt.Sprintf("Ouvrir l'espace ou le document lié à '%s'", title),
		"Écrire un premier jet imparfait ou lister 2 points clés (principe 80/20)",
		"Traiter le premier bloc de 15 minutes sans consulter d'autres onglets",
		"Valider l'avancement et célébrer la première micro-victoire",
	}
}
