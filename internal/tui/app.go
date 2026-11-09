package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"neuroflow/internal/assistant"
	"neuroflow/internal/config"
	"neuroflow/internal/habits"
	"neuroflow/internal/kanban"
	"neuroflow/internal/launcher"
	"neuroflow/internal/pomodoro"
)

// ViewMode définit la vue active affichée à l'écran.
type ViewMode int

const (
	ViewModeKanban ViewMode = iota
	ViewModeFocus
	ViewModeBrainDump
	ViewModeNewTask
	ViewModeHelp
	ViewModeHabits
	ViewModeAICoach
	ViewModeArchive
)

type tickMsg time.Time
type clearFlashMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearFlashCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return clearFlashMsg{}
	})
}

// App est le modèle central de l'application NeuroFlow.
type App struct {
	Config           *config.Config
	Board            *kanban.Board
	Timer            *pomodoro.Timer
	Habits           []*habits.Habit
	ArchivedTasks    []*kanban.Task
	Assistant        *assistant.Assistant
	CurrentMode      ViewMode
	PreviousMode     ViewMode

	ActiveColIdx     int
	ActiveTaskIdx    int
	ActiveSubtaskIdx int
	ActiveHabitIdx   int
	ActiveArchiveIdx int

	Width            int
	Height           int

	BrainDumpInput   textinput.Model
	NewTaskInput     textinput.Model
	FlashMessage     string

	AICoachResult    string
	AICoachLoading   bool
}

// NewApp instancie le modèle TUI principal.
func NewApp(cfg *config.Config, board *kanban.Board, habitList []*habits.Habit) *App {
	bdInput := textinput.New()
	bdInput.Placeholder = "Ex: Rappeler Julien pour le devis / Acheter du café..."
	bdInput.CharLimit = 150
	bdInput.Width = 50

	ntInput := textinput.New()
	ntInput.Placeholder = "Titre de la tâche (ex: Rédiger le rapport d'incident)"
	ntInput.CharLimit = 80
	ntInput.Width = 50

	timer := pomodoro.NewTimer()
	aiAssistant := assistant.NewAssistant(cfg.AIAssistant)

	app := &App{
		Config:         cfg,
		Board:          board,
		Timer:          timer,
		Habits:         habitList,
		Assistant:      aiAssistant,
		CurrentMode:    ViewModeKanban,
		PreviousMode:   ViewModeKanban,
		ActiveColIdx:   0,
		ActiveTaskIdx:  0,
		ActiveHabitIdx: 0,
		BrainDumpInput: bdInput,
		NewTaskInput:   ntInput,
		Width:          100,
		Height:         30,
	}

	// Si des tâches existent dans la colonne "En cours", on s'y positionne par défaut
	for i, col := range board.Columns {
		if col.DirName == "2_En_Cours" && len(col.Tasks) > 0 {
			app.ActiveColIdx = i
			break
		}
	}

	return app
}

// Init démarre les tickers de fond Bubble Tea.
func (a *App) Init() tea.Cmd {
	return tea.Batch(doTick(), textinput.Blink)
}

// CurrentTask renvoie la tâche actuellement sélectionnée dans la vue Kanban.
func (a *App) CurrentTask() *kanban.Task {
	if a.Board == nil || a.ActiveColIdx < 0 || a.ActiveColIdx >= len(a.Board.Columns) {
		return nil
	}
	col := a.Board.Columns[a.ActiveColIdx]
	if len(col.Tasks) == 0 || a.ActiveTaskIdx < 0 || a.ActiveTaskIdx >= len(col.Tasks) {
		return nil
	}
	return col.Tasks[a.ActiveTaskIdx]
}

// ClampIndices s'assure que les pointeurs de sélection restent dans les bornes.
func (a *App) ClampIndices() {
	if a.Board == nil || len(a.Board.Columns) == 0 {
		a.ActiveColIdx = 0
		a.ActiveTaskIdx = 0
		a.ActiveSubtaskIdx = 0
		return
	}

	if a.ActiveColIdx < 0 {
		a.ActiveColIdx = 0
	}
	if a.ActiveColIdx >= len(a.Board.Columns) {
		a.ActiveColIdx = len(a.Board.Columns) - 1
	}

	col := a.Board.Columns[a.ActiveColIdx]
	if len(col.Tasks) == 0 {
		a.ActiveTaskIdx = 0
		a.ActiveSubtaskIdx = 0
		return
	}

	if a.ActiveTaskIdx < 0 {
		a.ActiveTaskIdx = 0
	}
	if a.ActiveTaskIdx >= len(col.Tasks) {
		a.ActiveTaskIdx = len(col.Tasks) - 1
	}

	task := col.Tasks[a.ActiveTaskIdx]
	if len(task.Subtasks) == 0 {
		a.ActiveSubtaskIdx = 0
		return
	}
	if a.ActiveSubtaskIdx < 0 {
		a.ActiveSubtaskIdx = 0
	}
	if a.ActiveSubtaskIdx >= len(task.Subtasks) {
		a.ActiveSubtaskIdx = len(task.Subtasks) - 1
	}

	if a.ActiveHabitIdx < 0 {
		a.ActiveHabitIdx = 0
	}
	if len(a.Habits) > 0 && a.ActiveHabitIdx >= len(a.Habits) {
		a.ActiveHabitIdx = len(a.Habits) - 1
	}
}

// Update réagit aux événements clavier et d'horloge.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.Width = msg.Width
		a.Height = msg.Height

	case clearFlashMsg:
		a.FlashMessage = ""

	case tickMsg:
		warning, done := a.Timer.TickSecond()
		if warning {
			if a.Config.SoundBell {
				launcher.RingBell()
			}
			a.FlashMessage = "🔔 Transition douce : plus que 60 secondes pour clore en douceur."
			cmds = append(cmds, clearFlashCmd())
		}
		if done {
			if a.Config.SoundBell {
				launcher.RingBell()
			}
			if a.Timer.State == pomodoro.StateFocus || a.Timer.State == pomodoro.StateFlowExtension || a.Timer.State == pomodoro.StateSoftLaunch {
				isLong := (a.Timer.CompletedSessions % a.Config.LongBreakInterval) == 0
				dur := a.Config.ShortBreakDuration
				if isLong {
					dur = a.Config.LongBreakDuration
				}
				a.Timer.StartBreak(isLong, dur)
				a.FlashMessage = "🎉 Session terminée ! Prenez une pause bienfaisante."
			} else {
				a.Timer.Reset()
				a.FlashMessage = "☕ Pause terminée. Prêt pour un nouvel élan ?"
			}
			cmds = append(cmds, clearFlashCmd())
		}
		cmds = append(cmds, doTick())

	case tea.KeyMsg:
		switch a.CurrentMode {
		case ViewModeBrainDump:
			cmd := a.handleBrainDumpKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ViewModeNewTask:
			cmd := a.handleNewTaskKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ViewModeAICoach:
			cmd := a.handleAICoachKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ViewModeHabits:
			cmd := a.handleHabitsKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ViewModeArchive:
			cmd := a.handleArchiveKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ViewModeHelp:
			if msg.String() == "esc" || msg.String() == "?" || msg.String() == "q" {
				a.CurrentMode = a.PreviousMode
			}
		case ViewModeFocus:
			cmd := a.handleFocusKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case ViewModeKanban:
			cmd := a.handleKanbanKeys(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKanbanKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit

	case "tab":
		// Bascule vers l'écran de suivi des habitudes TDAH
		a.PreviousMode = ViewModeKanban
		a.CurrentMode = ViewModeHabits
		return nil

	case "?", "f1":
		a.PreviousMode = ViewModeKanban
		a.CurrentMode = ViewModeHelp

	case "a", "c":
		// Ouvrir l'Assistant d'aide Copilot / Coach TDAH
		a.PreviousMode = ViewModeKanban
		a.CurrentMode = ViewModeAICoach
		a.AICoachResult = ""
		a.AICoachLoading = false

	case "left", "h":
		if a.ActiveColIdx > 0 {
			a.ActiveColIdx--
			a.ClampIndices()
		}

	case "right", "l":
		if a.ActiveColIdx < len(a.Board.Columns)-1 {
			a.ActiveColIdx++
			a.ClampIndices()
		}

	case "up", "k":
		if a.ActiveSubtaskIdx > 0 {
			a.ActiveSubtaskIdx--
		} else if a.ActiveTaskIdx > 0 {
			a.ActiveTaskIdx--
			a.ActiveSubtaskIdx = 0
			a.ClampIndices()
		}

	case "down", "j":
		task := a.CurrentTask()
		if task != nil && a.ActiveSubtaskIdx < len(task.Subtasks)-1 {
			a.ActiveSubtaskIdx++
		} else if a.ActiveColIdx < len(a.Board.Columns) && a.ActiveTaskIdx < len(a.Board.Columns[a.ActiveColIdx].Tasks)-1 {
			a.ActiveTaskIdx++
			a.ActiveSubtaskIdx = 0
			a.ClampIndices()
		}

	case "shift+left", "H":
		// Déplacer physiquement la tâche vers la colonne de gauche (os.Rename)
		task := a.CurrentTask()
		if task != nil && a.ActiveColIdx > 0 {
			targetCol := a.Board.Columns[a.ActiveColIdx-1].DirName
			if err := a.Board.MoveTask(task, targetCol); err == nil {
				a.ActiveColIdx--
				a.ClampIndices()
				a.FlashMessage = fmt.Sprintf("Tâche déplacée vers %s", targetCol)
				cmds = append(cmds, clearFlashCmd())
			}
		}

	case "shift+right", "L":
		// Déplacer physiquement la tâche vers la colonne de droite (os.Rename)
		task := a.CurrentTask()
		if task != nil && a.ActiveColIdx < len(a.Board.Columns)-1 {
			targetCol := a.Board.Columns[a.ActiveColIdx+1].DirName
			if err := a.Board.MoveTask(task, targetCol); err == nil {
				a.ActiveColIdx++
				a.ClampIndices()
				a.FlashMessage = fmt.Sprintf("Tâche déplacée vers %s", targetCol)
				cmds = append(cmds, clearFlashCmd())
			}
		}

	case " ":
		// Cocher / Décocher la sous-tâche courante
		task := a.CurrentTask()
		if task != nil && len(task.Subtasks) > 0 && a.ActiveSubtaskIdx < len(task.Subtasks) {
			task.ToggleSubtask(a.ActiveSubtaskIdx)
			_ = a.Board.SaveTask(task)
			st := task.Subtasks[a.ActiveSubtaskIdx]
			if st.Done {
				a.FlashMessage = "✨ Sous-tâche accomplie ! (+1 élan dopaminergique)"
			} else {
				a.FlashMessage = "Sous-tâche rouverte."
			}
			cmds = append(cmds, clearFlashCmd())
		}

	case "f":
		// Lancer le focus classique sur la tâche active
		task := a.CurrentTask()
		a.Timer.StartFocus(task, a.Config.WorkDuration)
		if a.Config.AutoLaunchMusic {
			_ = launcher.OpenURL(a.Config.YouTubeMusicURL)
		}
		a.CurrentMode = ViewModeFocus

	case "s":
		// Lancer le Soft-Launch (5 min) pour amorcer l'action
		task := a.CurrentTask()
		a.Timer.StartSoftLaunch(task, a.Config.SoftLaunchDuration)
		if a.Config.AutoLaunchMusic {
			_ = launcher.OpenURL(a.Config.YouTubeMusicURL)
		}
		a.CurrentMode = ViewModeFocus

	case "b":
		// Ouvrir le Brain Dump
		a.PreviousMode = ViewModeKanban
		a.CurrentMode = ViewModeBrainDump
		a.BrainDumpInput.SetValue("")
		a.BrainDumpInput.Focus()

	case "m":
		// Lancer YouTube Music
		_ = launcher.OpenURL(a.Config.YouTubeMusicURL)
		a.FlashMessage = "🎵 YouTube Music lancé ! Ancrage audio actif."
		cmds = append(cmds, clearFlashCmd())

	case "n":
		// Créer une nouvelle tâche
		a.PreviousMode = ViewModeKanban
		a.CurrentMode = ViewModeNewTask
		a.NewTaskInput.SetValue("")
		a.NewTaskInput.Focus()

	case "x":
		// Archiver la tâche active (désencombrement sans perte d'information)
		task := a.CurrentTask()
		if task != nil {
			if err := a.Board.ArchiveTask(task, a.Config.ArchiveDir); err == nil {
				a.ClampIndices()
				a.FlashMessage = "📦 Tâche archivée ! (Consultez vos archives avec Maj+X)"
				cmds = append(cmds, clearFlashCmd())
			}
		}

	case "X":
		// Consulter les archives
		archived, err := kanban.LoadArchivedTasks(a.Config.ArchiveDir)
		if err == nil {
			a.ArchivedTasks = archived
			a.ActiveArchiveIdx = 0
			a.PreviousMode = ViewModeKanban
			a.CurrentMode = ViewModeArchive
		}

	case "d":
		// Supprimer la tâche active
		task := a.CurrentTask()
		if task != nil {
			_ = a.Board.DeleteTask(task)
			a.ClampIndices()
			a.FlashMessage = "Tâche supprimée."
			cmds = append(cmds, clearFlashCmd())
		}
	}

	return tea.Batch(cmds...)
}

func (a *App) handleArchiveKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit

	case "esc", "X", "tab":
		a.CurrentMode = ViewModeKanban
		return nil

	case "up", "k":
		if a.ActiveArchiveIdx > 0 {
			a.ActiveArchiveIdx--
		}

	case "down", "j":
		if a.ActiveArchiveIdx < len(a.ArchivedTasks)-1 {
			a.ActiveArchiveIdx++
		}

	case "enter", "u":
		// Restaurer la tâche sélectionnée dans le Kanban
		if len(a.ArchivedTasks) > 0 && a.ActiveArchiveIdx < len(a.ArchivedTasks) {
			task := a.ArchivedTasks[a.ActiveArchiveIdx]
			targetCol := "0_Inbox"
			if len(a.Board.Columns) > 0 {
				targetCol = a.Board.Columns[0].DirName
			}
			if err := a.Board.RestoreTask(task, targetCol); err == nil {
				a.FlashMessage = fmt.Sprintf("✨ Tâche '%s' restaurée dans le Kanban !", task.Title)
				a.ArchivedTasks, _ = kanban.LoadArchivedTasks(a.Config.ArchiveDir)
				if a.ActiveArchiveIdx >= len(a.ArchivedTasks) && a.ActiveArchiveIdx > 0 {
					a.ActiveArchiveIdx = len(a.ArchivedTasks) - 1
				}
				cmds = append(cmds, clearFlashCmd())
			}
		}

	case "d":
		// Supprimer définitivement la tâche archivée
		if len(a.ArchivedTasks) > 0 && a.ActiveArchiveIdx < len(a.ArchivedTasks) {
			task := a.ArchivedTasks[a.ActiveArchiveIdx]
			_ = os.Remove(task.Path)
			a.ArchivedTasks, _ = kanban.LoadArchivedTasks(a.Config.ArchiveDir)
			if a.ActiveArchiveIdx >= len(a.ArchivedTasks) && a.ActiveArchiveIdx > 0 {
				a.ActiveArchiveIdx = len(a.ArchivedTasks) - 1
			}
			a.FlashMessage = "Tâche archivée supprimée définitivement."
			cmds = append(cmds, clearFlashCmd())
		}
	}

	return tea.Batch(cmds...)
}

func (a *App) handleHabitsKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit

	case "tab", "esc":
		a.CurrentMode = ViewModeKanban
		return nil

	case "?":
		a.PreviousMode = ViewModeHabits
		a.CurrentMode = ViewModeHelp

	case "up", "k":
		if a.ActiveHabitIdx > 0 {
			a.ActiveHabitIdx--
		}

	case "down", "j":
		if a.ActiveHabitIdx < len(a.Habits)-1 {
			a.ActiveHabitIdx++
		}

	case " ":
		// Cocher / Décocher l'habitude courante pour aujourd'hui
		if len(a.Habits) > 0 && a.ActiveHabitIdx < len(a.Habits) {
			h := a.Habits[a.ActiveHabitIdx]
			gainedXP, err := habits.ToggleHabitToday(h)
			if err == nil {
				if gainedXP > 0 {
					if a.Config.SoundBell {
						launcher.RingBell()
					}
					a.FlashMessage = fmt.Sprintf("🎉 +%d XP ! Habitude validée ! Votre niveau augmente.", gainedXP)
				} else {
					a.FlashMessage = "Habitude décochée."
				}
				cmds = append(cmds, clearFlashCmd())
			}
		}

	case "m":
		_ = launcher.OpenURL(a.Config.YouTubeMusicURL)
		a.FlashMessage = "🎵 YouTube Music relancé."
		cmds = append(cmds, clearFlashCmd())
	}

	return tea.Batch(cmds...)
}

func (a *App) handleAICoachKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit

	case "esc":
		a.CurrentMode = a.PreviousMode
		a.AICoachResult = ""
		a.AICoachLoading = false
		return nil

	case "enter":
		if a.AICoachResult != "" {
			a.CurrentMode = a.PreviousMode
			a.AICoachResult = ""
			a.AICoachLoading = false
			return nil
		}

	case "1":
		// Découper la tâche en cours
		task := a.CurrentTask()
		if task == nil {
			a.AICoachResult = "Veuillez d'abord sélectionner une tâche dans le Kanban."
			return nil
		}

		desc := strings.Join(task.BodyLines, "\n")
		subtasks, source, _ := a.Assistant.DecomposeTask(task.Title, desc)

		for _, st := range subtasks {
			task.AddSubtask(st)
		}
		_ = a.Board.SaveTask(task)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%s\n\n", source))
		sb.WriteString("Micro-actions atomiques ajoutées au fichier .md :\n")
		for _, st := range subtasks {
			sb.WriteString(fmt.Sprintf("  • %s\n", st))
		}
		a.AICoachResult = sb.String()
		a.FlashMessage = "✨ Micro-tâches injectées avec succès !"
		cmds = append(cmds, clearFlashCmd())

	case "2":
		// Body Doubling / Déblocage
		task := a.CurrentTask()
		title := "votre tâche"
		if task != nil {
			title = task.Title
		}
		a.AICoachResult = a.Assistant.GetUnblockAdvice(title)

	case "3":
		// Kickoff Pro : Règle des 3
		var titles []string
		for _, col := range a.Board.Columns {
			for _, t := range col.Tasks {
				titles = append(titles, t.Title)
			}
		}
		advice, _ := a.Assistant.GetDailyKickoffAdvice(titles)
		a.AICoachResult = advice
	}

	return tea.Batch(cmds...)
}

func (a *App) handleFocusKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit

	case "esc", "k":
		// Retour au Kanban, le timer continue en arrière-plan
		a.CurrentMode = ViewModeKanban

	case "?":
		a.PreviousMode = ViewModeFocus
		a.CurrentMode = ViewModeHelp

	case "+":
		// Extension de 5 minutes
		a.Timer.ExtendFlow(5 * time.Minute)
		a.FlashMessage = "🌊 +5 minutes ajoutées sans briser votre élan !"
		cmds = append(cmds, clearFlashCmd())

	case "f":
		// Extension de 15 minutes (Ride the wave)
		a.Timer.ExtendFlow(15 * time.Minute)
		a.FlashMessage = "🌊 +15 minutes d'hyperfocus préservées !"
		cmds = append(cmds, clearFlashCmd())

	case "p":
		a.Timer.TogglePause()
		if a.Timer.IsPaused {
			a.FlashMessage = "⏸ En pause (respirez, reprenez quand vous voulez)."
		} else {
			a.FlashMessage = "▶ C'est reparti !"
		}
		cmds = append(cmds, clearFlashCmd())

	case "b":
		a.PreviousMode = ViewModeFocus
		a.CurrentMode = ViewModeBrainDump
		a.BrainDumpInput.SetValue("")
		a.BrainDumpInput.Focus()

	case "m":
		_ = launcher.OpenURL(a.Config.YouTubeMusicURL)
		a.FlashMessage = "🎵 YouTube Music relancé."
		cmds = append(cmds, clearFlashCmd())

	case "up":
		if a.ActiveSubtaskIdx > 0 {
			a.ActiveSubtaskIdx--
		}

	case "down":
		task := a.Timer.ActiveTask
		if task != nil && a.ActiveSubtaskIdx < len(task.Subtasks)-1 {
			a.ActiveSubtaskIdx++
		}

	case " ":
		task := a.Timer.ActiveTask
		if task != nil && len(task.Subtasks) > 0 && a.ActiveSubtaskIdx < len(task.Subtasks) {
			task.ToggleSubtask(a.ActiveSubtaskIdx)
			_ = a.Board.SaveTask(task)
			st := task.Subtasks[a.ActiveSubtaskIdx]
			if st.Done {
				a.FlashMessage = "✨ Bravo ! Une micro-étape de franchie (+dopamine)."
			} else {
				a.FlashMessage = "Sous-tâche réouverte."
			}
			cmds = append(cmds, clearFlashCmd())
		}
	}

	return tea.Batch(cmds...)
}

func (a *App) handleBrainDumpKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		a.CurrentMode = a.PreviousMode
		return nil

	case "enter":
		thought := a.BrainDumpInput.Value()
		if thought != "" {
			task := a.Timer.ActiveTask
			if task == nil {
				task = a.CurrentTask()
			}

			if task != nil {
				task.AppendBrainDump(thought)
				_ = a.Board.SaveTask(task)
			} else {
				_ = a.Board.SaveQuickThought(thought)
			}
			a.FlashMessage = "💡 Pensée parasite déchargée ! Votre esprit est libre."
			cmds = append(cmds, clearFlashCmd())
		}
		a.CurrentMode = a.PreviousMode
		return tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	a.BrainDumpInput, cmd = a.BrainDumpInput.Update(msg)
	return cmd
}

func (a *App) handleNewTaskKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		a.CurrentMode = ViewModeKanban
		return nil

	case "enter":
		title := a.NewTaskInput.Value()
		if title != "" {
			targetCol := "0_Inbox"
			if a.ActiveColIdx < len(a.Board.Columns) {
				targetCol = a.Board.Columns[a.ActiveColIdx].DirName
			}
			task, err := a.Board.CreateTask(targetCol, title, "", []string{"Première micro-action"})
			if err == nil && task != nil {
				col := a.Board.Columns[a.ActiveColIdx]
				a.ActiveTaskIdx = len(col.Tasks) - 1
				a.ActiveSubtaskIdx = 0
				a.FlashMessage = "✨ Tâche créée dans le Kanban !"
				cmds = append(cmds, clearFlashCmd())
			}
		}
		a.CurrentMode = ViewModeKanban
		return tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	a.NewTaskInput, cmd = a.NewTaskInput.Update(msg)
	return cmd
}

// View affiche la vue correspondant au mode actif.
func (a *App) View() string {
	switch a.CurrentMode {
	case ViewModeFocus:
		return a.RenderFocusView()
	case ViewModeBrainDump:
		return a.RenderBrainDumpModal()
	case ViewModeNewTask:
		return a.RenderNewTaskModal()
	case ViewModeHelp:
		return a.RenderHelpView()
	case ViewModeHabits:
		return a.RenderHabitsView()
	case ViewModeAICoach:
		return a.RenderAICoachModal()
	case ViewModeArchive:
		return a.RenderArchiveView()
	default:
		return a.RenderKanbanView()
	}
}
