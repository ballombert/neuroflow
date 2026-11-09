package pomodoro

import (
	"fmt"
	"strings"
	"time"

	"neuroflow/internal/kanban"
)

// State définit l'état du cycle attentionnel.
type State int

const (
	StateIdle State = iota
	StateSoftLaunch      // Rampe dopaminergique (5 min sans pression)
	StateFocus           // Session de focus principal
	StateFlowExtension   // Préservation de l'hyperfocus ("Ride the Wave")
	StateShortBreak      // Pause courte régénératrice
	StateLongBreak       // Pause longue
)

func (s State) String() string {
	switch s {
	case StateSoftLaunch:
		return "🌱 Rampe Dopaminergique (5 min)"
	case StateFocus:
		return "🔥 Session de Focus"
	case StateFlowExtension:
		return "🌊 Hyperfocus Préserve (+Ext)"
	case StateShortBreak:
		return "☕ Pause Courte"
	case StateLongBreak:
		return "🌴 Pause Longue"
	default:
		return "💤 En Attente"
	}
}

// Timer gère le temps, l'externalisation spatiale et la transition douce.
type Timer struct {
	State                  State
	ActiveTask             *kanban.Task
	TotalDuration          time.Duration
	Elapsed                time.Duration
	IsPaused               bool
	CompletedSessions      int
	ExtensionCount         int
	TransitionWarningEmitted bool
}

// NewTimer instancie le gestionnaire de pomodoro neuro-adapté.
func NewTimer() *Timer {
	return &Timer{
		State: StateIdle,
	}
}

// StartFocus démarre une session de travail standard.
func (t *Timer) StartFocus(task *kanban.Task, duration time.Duration) {
	t.ActiveTask = task
	t.State = StateFocus
	t.TotalDuration = duration
	t.Elapsed = 0
	t.IsPaused = false
	t.ExtensionCount = 0
	t.TransitionWarningEmitted = false
}

// StartSoftLaunch lance le défi d'initiation douce de 5 minutes.
func (t *Timer) StartSoftLaunch(task *kanban.Task, duration time.Duration) {
	t.ActiveTask = task
	t.State = StateSoftLaunch
	t.TotalDuration = duration
	t.Elapsed = 0
	t.IsPaused = false
	t.ExtensionCount = 0
	t.TransitionWarningEmitted = false
}

// ExtendFlow ajoute du temps à la volée sans casser l'état d'hyperfocus.
func (t *Timer) ExtendFlow(extra time.Duration) {
	if t.State == StateFocus || t.State == StateFlowExtension || t.State == StateSoftLaunch {
		t.State = StateFlowExtension
		t.TotalDuration += extra
		t.ExtensionCount++
		t.TransitionWarningEmitted = false
	}
}

// StartBreak bascule vers une phase de repos.
func (t *Timer) StartBreak(isLong bool, duration time.Duration) {
	if isLong {
		t.State = StateLongBreak
	} else {
		t.State = StateShortBreak
	}
	t.TotalDuration = duration
	t.Elapsed = 0
	t.IsPaused = false
	t.TransitionWarningEmitted = false
}

// TogglePause suspend ou reprend le décompte sans jugement.
func (t *Timer) TogglePause() {
	if t.State != StateIdle {
		t.IsPaused = !t.IsPaused
	}
}

// Reset réinitialise le chronomètre à l'état de repos.
func (t *Timer) Reset() {
	t.State = StateIdle
	t.TotalDuration = 0
	t.Elapsed = 0
	t.IsPaused = false
	t.TransitionWarningEmitted = false
}

// RemainingDuration calcule le temps restant.
func (t *Timer) RemainingDuration() time.Duration {
	rem := t.TotalDuration - t.Elapsed
	if rem < 0 {
		return 0
	}
	return rem
}

// ProgressRatio renvoie la proportion de temps écoulé entre 0.0 et 1.0.
func (t *Timer) ProgressRatio() float64 {
	if t.TotalDuration <= 0 {
		return 0.0
	}
	r := float64(t.Elapsed) / float64(t.TotalDuration)
	if r > 1.0 {
		return 1.0
	}
	return r
}

// TickSecond avance le compteur d'une seconde.
// Renvoie warning=true si on arrive à 60s de la fin (préparation de transition).
// Renvoie done=true si la session se termine.
func (t *Timer) TickSecond() (warning bool, done bool) {
	if t.State == StateIdle || t.IsPaused {
		return false, false
	}

	t.Elapsed += time.Second

	rem := t.RemainingDuration()
	// Avertissement de transition douce à 60 secondes de l'échéance
	if rem <= 60*time.Second && rem > 58*time.Second && !t.TransitionWarningEmitted {
		t.TransitionWarningEmitted = true
		warning = true
	}

	if t.Elapsed >= t.TotalDuration {
		if t.State == StateFocus || t.State == StateFlowExtension || t.State == StateSoftLaunch {
			t.CompletedSessions++
		}
		done = true
		return warning, done
	}

	return warning, false
}

// SpatialGauge génère une représentation visuelle volumique du temps écoulé en blocs Unicode.
// Principe clé de Russell Barkley : externaliser le temps pour le rendre spatialement visible.
func (t *Timer) SpatialGauge(width int) string {
	if width < 10 {
		width = 10
	}

	ratio := t.ProgressRatio()
	filledCount := int(ratio * float64(width))
	if filledCount > width {
		filledCount = width
	}
	emptyCount := width - filledCount

	filled := strings.Repeat("█", filledCount)
	empty := strings.Repeat("░", emptyCount)

	return filled + empty
}

// FormatTimeDigital renvoie "MM:SS".
func FormatTimeDigital(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
