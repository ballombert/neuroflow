package pomodoro

import (
	"strings"
	"testing"
	"time"
)

func TestPomodoroFlowAndGauge(t *testing.T) {
	timer := NewTimer()

	if timer.State != StateIdle {
		t.Errorf("État initial attendu StateIdle, obtenu %v", timer.State)
	}

	// 1. Démarrage Soft Launch
	timer.StartSoftLaunch(nil, 5*time.Minute)
	if timer.State != StateSoftLaunch {
		t.Errorf("État attendu StateSoftLaunch, obtenu %v", timer.State)
	}
	if timer.RemainingDuration() != 5*time.Minute {
		t.Errorf("Durée restante attendue 5 min, obtenu %v", timer.RemainingDuration())
	}

	// 2. Vérification de la jauge spatiale au début (vide)
	gauge := timer.SpatialGauge(20)
	if !strings.HasPrefix(gauge, "░") {
		t.Errorf("Au départ la jauge doit être vide, obtenu: %s", gauge)
	}

	// 3. Simuler du temps écoulé
	timer.Elapsed = 2 * time.Minute + 30*time.Second // 50%
	ratio := timer.ProgressRatio()
	if ratio != 0.5 {
		t.Errorf("Ratio attendu 0.5, obtenu %f", ratio)
	}

	gauge = timer.SpatialGauge(20)
	if strings.Count(gauge, "█") != 10 {
		t.Errorf("Attendu 10 blocs pleins pour 50%% sur 20 de large, obtenu: %s", gauge)
	}

	// 4. Flow Extension (+15 min)
	timer.ExtendFlow(15 * time.Minute)
	if timer.State != StateFlowExtension {
		t.Errorf("État attendu StateFlowExtension, obtenu %v", timer.State)
	}
	if timer.TotalDuration != 20*time.Minute {
		t.Errorf("Total attendu 20 min, obtenu %v", timer.TotalDuration)
	}

	// 5. Avertissement doux de transition à 60s
	timer.Elapsed = 19 * time.Minute // Reste exactement 60s
	warning, done := timer.TickSecond()
	if !warning {
		t.Errorf("L'avertissement de transition douce à 60s aurait dû être émis")
	}
	if done {
		t.Errorf("La session ne devrait pas être terminée à 60s de la fin")
	}
}
