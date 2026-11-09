# 🧠 NeuroFlow v1.1.0

> **Application Pomodoro Adaptatif, Kanban Fichiers, Rituels d'Habitudes & Assistant Copilot pour Esprits Neurodivergents (TDAH)**  
> *Développée en Go — 100% Cross-Platform (Windows, macOS, Linux) — Architecture 100% Système de Fichiers (Markdown & YAML).*

---

## 🪟 Compatibilité Windows Complète

NeuroFlow a été spécialement optimisé pour fonctionner de manière native et fluide sous **Windows 10 et Windows 11** :

1. **Rendu UTF-8 Automatique (`kernel32.dll`)** :  
   À son lancement sous Windows, NeuroFlow configure automatiquement la console en page de code UTF-8 (`SetConsoleOutputCP(65001)`). Tous les emojis (🔥, 🧠, 💼, ⚡) et les blocs graphiques de progression (█, ░) s'affichent impeccablement dans **Windows Terminal**, **PowerShell** et **l'invite de commandes (CMD)**.
2. **Gestion des Fichiers & Noms Réservés Windows** :  
   - Assainissement automatique des caractères interdits sous Windows (`\ / : * ? " < > |`).
   - Protection contre les noms de périphériques réservés DOS (`CON`, `PRN`, `AUX`, `NUL`, etc.).
   - Support des retours chariot Windows (`\r\n` / CRLF) transparent sans altération des fichiers Markdown.
   - Gestion des déplacements de fichiers atomiques (`os.Rename`) même si le fichier cible existe déjà ou en cas de déplacement inter-disques (`copyAndRemoveFile`).
3. **Lancement Audio YouTube Music sous Windows** :  
   Ouverture automatique et silencieuse via l'association système (`rundll32 url.dll,FileProtocolHandler` avec repli `start ""`) sans faire clignoter de fenêtre noire parasite.
4. **Assistant Copilot CLI sous Windows** :  
   Détection automatique des scripts Windows (`copilot.cmd`, `copilot.bat`, `copilot.exe`) dans le PATH, `%APPDATA%\npm`, et `%LOCALAPPDATA%`.

---

## 🚀 Lancement sous Windows

### Méthode 1 : Binaire Direct (`neuroflow.exe`)
Le binaire pré-compilé pour Windows **`neuroflow.exe`** est inclus à la racine :
```powershell
# Dans Windows Terminal ou PowerShell :
.\neuroflow.exe

# Ou dans CMD :
neuroflow.exe
```

### Méthode 2 : Double-Clic avec `run.bat` ou `run.ps1`
* Double-cliquez simplement sur **`run.bat`** dans l'Explorateur Windows.
* Ou exécutez le script PowerShell :
```powershell
.\run.ps1
```

> **Conseil d'affichage :** Pour la meilleure expérience visuelle sous Windows, utilisez **Windows Terminal** (installé par défaut sur Windows 11, disponible sur le Microsoft Store pour Windows 10).

---

## 🍎 Lancement sous macOS & 🐧 Linux

```bash
cd adhd-pomodoro-kanban

# Lancer sous macOS / Linux :
./neuroflow

# Ou compiler :
make build
```

---

## 📁 Architecture Système de Fichiers

Tout est lisible et modifiable avec Obsidian, VS Code, Notepad++ ou Vim :

```
adhd-pomodoro-kanban/
├── config.yaml                      # Configuration principale en YAML
├── run.bat                          # Lanceur Windows UTF-8 (double-clic)
├── run.ps1                          # Lanceur PowerShell Windows
├── neuroflow.exe                    # Binaire exécutable autonome Windows (AMD64)
├── neuroflow                        # Binaire exécutable autonome macOS/Linux
├── kanban_workspace/                # TABLEAU KANBAN (chaque dossier = une colonne)
│   ├── 0_Inbox/                     # Tâches brutes
│   ├── 1_A_Faire/                   # Tâches planifiées
│   ├── 2_En_Cours/                  # Tâche active de focus
│   └── 3_Termine/                   # Tâches accomplies
├── habits_workspace/                # HABITUDES HORS KANBAN (Markdown)
│   ├── kickoff-matinal-pro.md       # Règle des 3 victoires
│   ├── tri-inbox-chrono-pro.md      # 15 min chrono
│   ├── rituel-shutdown-pro.md       # Clôture à 18h
│   ├── amorce-soft-launch.md        # 5 min anti-paralysie
│   └── pause-eau-mouvement.md       # Hydratation & posture
├── cmd/neuroflow/                   # Code source
│   ├── main.go                      # Point d'entrée
│   ├── console_windows.go           # UTF-8 Windows Console (CP 65001)
│   └── console_other.go             # No-op Unix
├── internal/
│   ├── assistant/copilot.go         # Intégration Copilot CLI & Coach TDAH
│   ├── config/config.go             # Chargeur YAML / JSON
│   ├── habits/                      # Moteur d'habitudes et gamification XP
│   ├── kanban/                      # Moteur Kanban dossiers & fichiers .md
│   ├── launcher/opener.go           # Lanceur cross-platform YouTube Music
│   ├── pomodoro/timer.go            # Pomodoro adaptatif & jauge spatiale
│   └── tui/                         # Vues Bubble Tea
└── Makefile                         # make build, make build-windows, make test
```

---

## ⌨️ Raccourcis Clavier Complets

| Touche | Vue | Action |
|---|---|---|
| **`Tab`** | Partout | **Basculer entre Kanban et l'écran des Habitudes & XP** |
| **`a` / `c`** | Kanban | **Ouvrir l'Assistant Copilot IA** (découper la tâche, kickoff, déblocage) |
| **`f`** | Kanban | Lancer une **Session de Focus** (25 min par défaut) |
| **`s`** | Kanban | Lancer le **Soft-Launch (5 min)** : rampe dopaminergique |
| **`Espace`** | Kanban / Focus | Cocher / décocher la sous-tâche sélectionnée (écrit dans le `.md`) |
| **`Espace`** | Habitudes | **Valider l'habitude du jour (+XP et streak)** |
| `h` / `l` ou `⬅` / `➔` | Kanban | Naviguer entre les colonnes |
| `j` / `k` ou `⬆` / `⬇` | Partout | Naviguer entre les tâches / sous-tâches / habitudes |
| `H` / `L` | Kanban | **Déplacer physiquement** la tâche de colonne (`os.Rename`) |
| `+` / `f` | Focus | **Flow Extension (+5m ou +15m)** : préserver l'hyperfocus |
| `p` | Focus | Mettre en pause ou reprendre sans culpabilité |
| `b` | Partout | **Brain Dump express** : capturer une pensée parasite sans quitter le flux |
| `m` | Partout | Lancer / relancer **YouTube Music** dans l'app par défaut |
| `n` | Kanban | Créer une nouvelle tâche Markdown |
| `d` | Kanban | Supprimer la tâche active |
| `?` | Partout | Guide complet des touches et fondements neuroscientifiques |
| `q` | Partout | Quitter NeuroFlow |

---

## ⚙️ Configuration YAML (`config.yaml`)

```yaml
# Emplacements des données
kanban_dir: "./kanban_workspace"
habits_dir: "./habits_workspace"

# Paramètres du cycle Pomodoro (en minutes)
work_duration_minutes: 25
soft_launch_minutes: 5
short_break_minutes: 5
long_break_minutes: 15
long_break_interval: 4

# Musique & Sensory Priming
youtube_music_url: "https://music.youtube.com/watch?v=jfKfPfyJRdk"
auto_launch_music: false
sound_bell: true
theme: "tokyo-night"

# Assistant d'Aide IA Copilot
ai_assistant:
  enabled: true
  provider: "copilot"
  copilot_bin: "copilot"
  role_prompt: >
    Tu es un coach expert en TDAH et en productivité bienveillante.
    Tu aides à décomposer les tâches professionnelles intimidantes en micro-étapes
    atomiques ridicules, sans culpabilité ni jargon superflu.

# Organisation & Protection de la vie professionnelle
work_life:
  daily_wins_goal: 3             # 3 victoires du jour au kickoff
  shutdown_time: "18:00"         # Clôture formelle de la journée de travail
  inbox_limit_min: 15            # Limite du tri des messages pros
```
