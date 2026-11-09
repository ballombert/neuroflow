package launcher

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURL ouvre une URL dans le navigateur ou l'application par défaut du système.
// Totalement asynchrone et non-bloquant pour l'interface TUI, optimisé pour Windows, macOS et Linux.
func OpenURL(urlStr string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", urlStr)
	case "linux":
		cmd = exec.Command("xdg-open", urlStr)
	case "windows":
		// Sous Windows, rundll32 est léger et n'affiche pas de fenêtre de console noire parasite
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", urlStr)
	default:
		return fmt.Errorf("plateforme non supportée : %s", runtime.GOOS)
	}

	err := cmd.Start()
	if err != nil && runtime.GOOS == "windows" {
		// Repli Windows avec le verbe système 'start'
		fallbackCmd := exec.Command("cmd", "/c", "start", "", urlStr)
		return fallbackCmd.Start()
	}
	return err
}

// RingBell émet un signal sonore doux (terminal bell ASCII \a).
// Fonctionne sous Windows Terminal, PowerShell, CMD, macOS Terminal et Linux consoles.
func RingBell() {
	fmt.Print("\a")
}
