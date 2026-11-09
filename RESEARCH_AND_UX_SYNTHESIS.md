# Synthèse Scientifique : TDAH, Technique Pomodoro, Habitudes & Vie Professionnelle

> **Document de recherche théorique et empirique**  
> Destiné à la conception de l'application cross-plateforme Go : Pomodoro Adaptatif, Kanban Fichiers, Habitudes Gamifiées & Assistant Copilot IA.

---

## 1. Fondements Neurobiologiques et Cognitifs du TDAH

Le Trouble du Déficit de l'Attention avec ou sans Hyperactivité (TDAH / ADHD) est avant tout un **trouble neurodéveloppemental de l'autorégulation et des fonctions exécutives** (Barkley, 1997). Il ne s'agit pas d'un manque d'attention, mais d'une **difficulté de régulation et de distribution de l'attention**, intimement liée à un déficit de modulation dopaminergique et noradrénergique dans les circuits fronto-striataux.

### 1.1 La « Cécité Temporelle » (*Time Blindness*) & La Myopie Temporelle
* **Auteur clé : Dr Russell A. Barkley (1997, 2012)**
* **Mécanisme :** Les individus avec TDAH vivent dans une fenêtre temporelle extrêmement étroite : le **« MAINTENANT »** (*NOW*) vs **« PAS MAINTENANT »** (*NOT NOW*). Les échéances futures ou le passage abstrait du temps n'ont pas d'ancrage viscéral.
* **Conséquence UX :** Les horloges numériques traditionnelles (`14:23`) ou les décomptes numériques purs (`23:45... 23:44...`) sont abstraits et n'activent pas la conscience spatiale du temps. Le temps doit être **spatialisé**, rendu physique et visuel (principe du *Time Timer* de Jan Rogers).

### 1.2 Le Système Nerveux Fondé sur l'Intérêt (*Interest-Based Nervous System*)
* **Auteur clé : Dr William Dodson (2016)**
* **Mécanisme :** Le cerveau neurotypique priorise selon trois axes : *Importance*, *Priorité*, *Récompense différée*. Le cerveau TDAH est structurellement incapable d'exploiter durablement ces leviers. Il réagit à l'acronyme **ICNU** :
  - **I**nterest (Intérêt / Fascination)
  - **C**hallenge (Défi / Compétition ludique)
  - **N**ovelty (Nouveauté / Fraîcheur)
  - **U**rgency (Urgence / Échéance imminente)
* **Conséquence UX :** Les listes de tâches rigides et froides génèrent un blocage (paralysie d'initiation). L'outil doit introduire de la micro-nouveauté, du défi non culpabilisant et transformer l'amorce en micro-défi.

### 1.3 Théorie de la Stimulation Optimale & Modèle MBA
* **Auteurs clés : Sydney Zentall (1975) ; Sverker Sikström & Göran Söderlund (2007) - *Moderate Brain Arousal (MBA) model***
* **Mécanisme :** Le cortex préfrontal du TDAH est en état de **sous-arousal** basal. Pour fonctionner, le cerveau cherche inconsciemment une stimulation compensatoire (gigotements, pensées vagabondes, distraction).
* **Le Phénomène de Résonance Stochastique :** L'apport d'un stimulus sensoriel continu modéré (bruit blanc/brun, musique sans paroles, lo-fi, synthwave ou bandes-son de jeux vidéo à tempo régulier 120-130 BPM) rehausse le rapport signal/bruit neuronal. Contrairement aux neurotypiques chez qui le bruit perturbe l'attention, chez le TDAH, il stabilise le cortex préfrontal.
* **Conséquence UX :** Intégration native et fluide d'un flux audio stimulant (YouTube Music : Lo-Fi, Synthwave, Brown Noise) dès l'engagement de la session de travail pour amorcer la boucle dopaminergique.

---

## 2. Déconstruction Critique de la Méthode Pomodoro Classique

| Paramètre Pomodoro Classique | Problème pour le profil TDAH | Solution UX Neuro-Inclusive |
|---|---|---|
| **Bloc rigide de 25 minutes** | L'horizon de 25 minutes est perçu comme une montagne infranchissable en cas de paralysie d'initiation. | **Mode "Dopamine Ramp" (5 min)** : micro-session de démarrage pour briser la friction, avec auto-extension si la dynamique prend. |
| **Alarme abrupte à 25:00** | L'interruption brutale brise l'état d'**hyperfocus** (état précieux et difficile à atteindre) et provoque une réaction de frustration intense ou de panique. | **Mode "Ride the Wave" (Flow Extension)** : possibilité d'étendre la session à la volée (+5, +15 min) d'une seule touche sans sortir de l'état mental. |
| **Transition abrupte Travail ➔ Pause** | L'inertie cognitive empêche les bascules instantanées. Passer du code à la pause en 1 seconde est douloureux. | **Avertisseur de transition douce (Gentle Cue)** : avertissement visuel discret 60 secondes avant la fin pour préparer l'esprit au changement de contexte. |
| **Décompte numérique abstrait** | Regarder les chiffres s'égrener crée de l'anxiété temporelle (*timer anxiety*) sans matérialiser la durée restante. | **Jauge spatiale et chromatique** : barre de progression volumique, blocs textuels/visuels, changement graduel de nuance. |
| **Culpabilisation en cas d'abandon** | Les applications qui réinitialisent violemment le compteur ou affichent un échec induisent une spirale de honte (*shame spiral*). | **Philosophie non punitive** : toute minute passée est enregistrée comme une victoire cognitive ; boutons de pause douce sans jugement. |

---

## 3. Neurobiologie des Habitudes chez le TDAH & Gamification Non Punitive

### 3.1 Pourquoi les habitudes ne s'automatisent pas chez le TDAH
Chez une personne neurotypique, répéter un comportement pendant 21 à 66 jours transfère le contrôle des circuits corticaux préfrontaux vers les circuits sous-corticaux des ganglions de la base (striatum dorsal), créant un automatisme inconscient (Wendy Wood, 2019).
Chez le sujet TDAH, **le déficit de signalisation dopaminergique phasique altère la consolidation de cette automaticité** (Barkley, 2012). Même après six mois de pratique, se brosser les dents ou ranger son bureau demande toujours une décision exécutive consciente.

### 3.2 Intentions d'Implémentation & Micro-Habitudes
* **Peter Gollwitzer (1999) :** La formule **« Si [Déclencheur Contextuel], Alors [Micro-Action] »** décharge le cortex préfrontal en créant une liaison associative directe.
* **BJ Fogg (Tiny Habits, 2019) :** Pour contourner la résistance, le seuil d'amorce doit être minuscule (« Juste ouvrir l'application » ou « Écrire une ligne »).
* **Hors Kanban :** Les habitudes ne sont pas des tâches jetables à déplacer vers une colonne "Terminé", mais des **rituels d'ancrage permanents**. Elles ont leur propre espace physique (`./habits_workspace/`) avec leur historique et leurs déclencheurs.

### 3.3 Gamification Non Punitive : Le Piège de la Spirale de Honte
Les trackers d'habitudes classiques qui détruisent une série de 30 jours à cause d'un oubli provoquent une **réaction de désespoir et d'abandon total** chez le TDAH (*all-or-nothing thinking*).
* **Solution NeuroFlow :**
  - **Points d'XP acquis à vie** : L'XP total ne diminue jamais.
  - **Progression en paliers de niveaux dopaminergiques** (*Niv. 1 Apprenti ➔ Niv. 2 Dompteur d'Inertie ➔ Niv. 3 Maître du Focus...*).
  - Récompense sensorielle immédiate lors de la validation (son de cloche, animation flash).

---

## 4. Organisation de la Vie Professionnelle avec un TDAH

La vie professionnelle soumet le cerveau TDAH à trois risques majeurs :
1. **La paralysie de priorisation :** Toutes les tâches urgentes des collègues et de la hiérarchie semblent d'égale importance, provoquant une surcharge cognitive.
   * *Solution :* **Kickoff Matinal & Règle des 3 Victoires Quotidiennes**. Choisir au maximum 3 livrables à fort impact pour la journée.
2. **Le trou noir des communications :** Ouvrir Slack ou ses emails sans cadre temporel absorbe toute l'énergie mentale de la journée.
   * *Solution :* **Habitude de Tri Express (15 min chrono)** sous compte à rebours strict.
3. **Le flou des frontières travail/vie privée & le burnout :** L'incapacité à "décrocher" fait que le travail non terminé hante la soirée.
   * *Solution :* **Rituel de Clôture Pro (Shutdown à 18h)** (Cal Newport, *Deep Work*). Ranger son bureau, fermer les onglets pros et noter explicitement la première micro-action du lendemain pour libérer la mémoire de travail.

---

## 5. L'Assistant IA Copilot comme « Body Double » et Décomposeur Atomique

Le concept de **Body Doubling** est l'une des techniques d'adaptation les plus puissantes validées pour le TDAH : la présence bienveillante d'une entité qui accompagne le démarrage réduit l'anxiété de la tâche.

L'intégration de **GitHub Copilot CLI** remplit trois rôles précis :
1. **Découpage atomique automatique :** Les tâches pro formulées de manière trop vague (ex: « Préparer la présentation client ») provoquent l'évitement. L'IA la déconstruit instantanément en 3 à 5 micro-étapes ridicules à très faible friction et les injecte sous forme de cases `- [ ]` dans le fichier Markdown de la tâche.
2. **Conseil d'amorce douce (Unblock) :** Fournit la première action de 60 secondes pour désamorcer la procrastination.
3. **Arbitrage matinal (Kickoff Pro) :** Analyse les fichiers du Kanban et identifie les 3 tâches prioritaires sans charge mentale pour l'utilisateur.

---

## 6. Références Bibliographiques Fondatrices

1. **Barkley, R. A. (1997).** *ADHD and the nature of self-control.* Guilford Press.
2. **Barkley, R. A. (2012).** *Executive Functions: What They Are, How They Work, and Why They Evolved.* Guilford Press.
3. **Dodson, W. (2016).** *Secrets of the ADHD Brain: Why we think, act, and feel the way we do.* ADDitude Magazine.
4. **Fogg, B. J. (2019).** *Tiny Habits: The Small Changes That Change Everything.* Houghton Mifflin Harcourt.
5. **Gollwitzer, P. M. (1999).** *Implementation intentions: Strong effects of simple plans.* American Psychologist, 54(7), 493–503.
6. **Newport, C. (2016).** *Deep Work: Rules for Focused Success in a Distracted World.* Grand Central Publishing.
7. **Sikström, S., & Söderlund, G. (2007).** *Stimulus-dependent dopamine release in ADHD: The moderate brain arousal model.* Psychological Review, 114(4), 1047–1075.
8. **Wood, W. (2019).** *Good Habits, Bad Habits: The Science of Making Positive Changes That Stick.* Farrar, Straus and Giroux.
9. **Zentall, S. S. (1975).** *Optimal stimulation as theoretical basis of hyperactivity.* American Journal of Orthopsychiatry, 45(4), 549–563.
