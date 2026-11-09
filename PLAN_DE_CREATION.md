# Plan de Création Détaillé : NeuroFlow (Go Cross-Platform)

> **Application Pomodoro Adaptatif, Kanban Système de Fichiers & Scaffolding TDAH**  
> *Conçue pour transformer l'énergie cognitive en élan durable par l'externalisation du temps et l'élimination des frictions.*

---

## 1. Vision et Identité du Projet

- **Nom :** **NeuroFlow** (ou `neuroflow`)
- **Langage & Écosystème :** Go (Golang) 1.22+
- **Architecture d'Interface :** TUI (Terminal User Interface) avancée via l'écosystème Charmbracelet (**Bubble Tea**, **Lipgloss**, **Bubbles**).
- **Cible Cross-Platform :** macOS (Apple Silicon & Intel), Linux (toutes distributions), Windows 10/11 (Windows Terminal, PowerShell, CMD).
- **Stockage de Données :** 100% Système de fichiers local transparent (Markdown & Dossiers), sans base de données opaque. Compatible Git, Obsidian, Logseq, VS Code.

---

## 2. Architecture Modulaire

```
adhd-pomodoro-kanban/
├── cmd/
│   └── neuroflow/
│       └── main.go               # Point d'entrée, gestion des flags CLI et initialisation
├── internal/
│   ├── config/
│   │   └── config.go             # Gestion de config.json, valeurs par défaut neuro-inclusives
│   ├── kanban/
│   │   ├── model.go              # Structures Column, Task, Subtask
│   │   ├── storage.go            # Scan de répertoires, lecture/écriture Markdown, os.Rename
│   │   └── parser.go             # Analyseur Markdown des cases à cocher (- [ ] / - [x])
│   ├── pomodoro/
│   │   ├── timer.go              # Machine à états (Idle, SoftLaunch, Focus, FlowExtension, Break)
│   │   └── stats.go              # Calcul des micro-victoires, streaks et temps de travail
│   ├── launcher/
│   │   └── opener.go             # Lanceur multiplateforme d'URL (YouTube Music, navigateur OS)
│   └── tui/
│       ├── app.go                # Modèle Bubble Tea principal, orchestration des vues
│       ├── view_kanban.go        # Rendu multi-colonnes Kanban avec cartes et badges
│       ├── view_focus.go         # Vue focus immersive (Timer spatial + Ancre attentionnelle)
│       ├── view_braindump.go     # Modal express de capture de pensées parasites
│       ├── view_help.go          # Guide des raccourcis et philosophie TDAH
│       └── styles.go             # Thèmes visuels Lipgloss apaisants (Tokyo Night, Catppuccin)
├── kanban_workspace/             # Espace Kanban exemple pré-configuré
│   ├── 0_Inbox/
│   ├── 1_A_Faire/
│   ├── 2_En_Cours/
│   └── 3_Termine/
├── RESEARCH_AND_UX_SYNTHESIS.md  # Synthèse de la littérature scientifique
├── PLAN_DE_CREATION.md           # Ce présent document
├── Makefile                      # Automatisation du build, test et run
├── go.mod
└── README.md                     # Manuel utilisateur et explications scientifiques
```

---

## 3. Spécifications Fonctionnelles Détaillées

### 3.1 Système Kanban "Dossiers & Fichiers Markdown"
- **Structure physique :**
  - Un dossier racine (ex. `./kanban_workspace`).
  - Chaque sous-dossier correspond à une colonne du Kanban (ex: `0_Inbox`, `1_A_Faire`, `2_En_Cours`, `3_Termine`). L'ordre est déterminé par le préfixe numérique ou l'ordre alphabétique.
  - Chaque fichier `.md` dans une colonne représente une tâche.
  - Le titre est extrait du premier `# Titre` ou du nom de fichier.
  - Le corps contient les notes, descriptions et la liste des sous-tâches sous forme de cases à cocher :
    ```markdown
    # Implémenter le parseur Markdown
    Description ou contexte...

    - [x] Détecter les lignes avec - [ ]
    - [ ] Permettre la bascule d'état (toggle)
    - [ ] Calculer le ratio d'avancement
    ```
- **Opérations atomiques :**
  - **Déplacement :** Déplacer une tâche vers la colonne de droite (`L` ou `➔`) ou de gauche (`H` ou `⬅`) effectue un appel atomique `os.Rename(src, dst)` sur le système de fichiers.
  - **Bascule de sous-tâche :** Appuyer sur `Espace` ou `Enter` sur une sous-tâche modifie directement le fichier `.md` sans altérer le reste du contenu.
  - **Nouvelle tâche :** Raccourci `n` dans une colonne pour créer immédiatement un nouveau fichier Markdown avec un modèle prêt à l'emploi.

### 3.2 Moteur Pomodoro Adaptatif TDAH
- **Machine à états :**
  1. `StateIdle` : En attente de session.
  2. `StateSoftLaunch` (Dopamine Ramp) : 5 minutes pour démarrer sans pression cognitive. Si l'utilisateur accroche, bascule fluide vers le Focus.
  3. `StateFocus` : Session principale (configurable : 15, 20, 25, 45 min).
  4. `StateFlowExtension` (Ride the Wave) : Permet d'ajouter +5 min ou +15 min avec une seule touche (`+`), préservant le précieux état d'hyperfocus sans alarme punitive.
  5. `StateShortBreak` : Pause courte (ex. 5 min).
  6. `StateLongBreak` : Pause longue (ex. 15-20 min) toutes les $N$ sessions.
- **Ancre Visuelle d'Attention (*Anchoring*) :**
  - Pendant le décompte, la tâche en cours et ses sous-tâches restent affichées en grand. L'utilisateur qui relève la tête sait instantanément sur quoi porte son effort sans devoir reconstituer sa mémoire de travail.
- **Représentation Spatiale du Temps (*Spatial Gauge*) :**
  - Remplacement de la simple horloge numérique par une jauge visuelle continue en blocs Unicode (`████████░░░░░░░░`).
  - Indicateur de progression en pourcentage et temps restant/écoulé.

### 3.3 Amorceur Sensoriel YouTube Music (Cross-Platform)
- **Déclenchement :**
  - Automatique lors du début d'une session de focus (option activable/désactivable dans `config.json`).
  - Manuel à tout moment via la touche `m` (Music).
- **Implémentation multiplateforme robuste :**
  - Détection automatique du système d'exploitation via `runtime.GOOS`.
  - Exécution du lanceur adéquat (`open` sous Darwin/macOS, `xdg-open` sous Linux, `rundll32` ou `cmd /c start` sous Windows).
  - URL par défaut : Playlist de concentration YouTube Music (Lo-Fi, Synthwave, Bruit Brun, ou l'URL personnalisée de l'utilisateur).

### 3.4 Bouclier Anti-Interruption (*Brain Dump Buffer*)
- **Problème résolu :** L'angoisse de perdre une idée parasite (« Payer la facture », « Envoyer un mail à Paul ») pousse la personne TDAH à interrompre sa tâche en cours.
- **Solution UX :**
  - Touche `b` (Brain dump) à tout moment.
  - Fenêtre contextuelle légère pour taper la pensée.
  - Au choix : ajout instantané en tant que sous-tâche de la tâche active, ou création d'une note rapide dans `0_Inbox`.
  - Fermeture et retour direct au focus en moins de 3 secondes, sans altérer le chronomètre.

---

## 4. Plan de Réalisation Étape par Étape

1. **Étape 1 : Initialisation de l'environnement Go & Dépendances**
   - Initialisation du module `go.mod` (`neuroflow`).
   - Importation des dépendances Charmbracelet (`bubbletea`, `lipgloss`, `bubbles`).
2. **Étape 2 : Moteur de Fichiers & Analyseur Markdown Kanban**
   - Scan des dossiers et ordonnancement des colonnes.
   - Parseur de tâches `.md` et détection des sous-tâches `[ ]` / `[x]`.
   - Manipulation atomique des fichiers (`os.Rename`, écriture préservant les retours à la ligne).
3. **Étape 3 : Lanceur d'Audio YouTube Music Cross-Platform**
   - Implémentation du paquet `launcher` pour ouvrir les URLs sur macOS, Linux et Windows de façon non bloquante.
4. **Étape 4 : Moteur Pomodoro Neuro-Inclusif**
   - Implémentation de la machine à états : Soft Launch (5m), Focus ajustable, Flow Extension (+5m/+15m), Pauses.
   - Système de calcul de jauge spatiale.
5. **Étape 5 : Conception et Assemblage de la TUI (Bubble Tea)**
   - Écran Kanban interactif (colonnes, cartes, déplacement, sélection de tâche).
   - Écran Focus immersif avec ancre attentionnelle et jauge spatiale.
   - Modal Brain Dump.
   - Écran de configuration / raccourcis.
6. **Étape 6 : Tests, Validation et Préparation de l'Espace Kanban Exemple**
   - Création de l'arborescence `./kanban_workspace` avec des colonnes et des tâches guides rédigées avec bienveillance pour le TDAH.
   - Compilation et vérification des binaires.
7. **Étape 7 : Documentation Complète et Guide d'Utilisation**
   - Rédaction du `README.md` exhaustif avec guide d'installation, raccourcis, explications des principes TDAH.
