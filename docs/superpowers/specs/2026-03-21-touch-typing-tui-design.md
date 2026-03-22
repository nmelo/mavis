# Touch Typing Retraining TUI

## Context

A terminal-based touch typing trainer designed specifically for retraining, not learning from scratch. The user has 15+ years of typing with incorrect finger habits (3 fingers) and wants to rebuild from the ground up with proper home-row touch typing technique using all 10 fingers.

The key design constraint: every decision optimizes for breaking old muscle memory and building correct new patterns. This means strict error handling, progressive key unlocking (so old habits can't take over on familiar keys), and constant visual reinforcement of which finger owns which key.

## Technology

- Language: Go
- TUI framework: Bubble Tea (bubbletea) + Lipgloss for styling
- No external dependencies beyond the Charm ecosystem
- Single binary, no runtime requirements

## Screen Layout

Single-screen TUI with three vertical zones:

```
+------------------------------------------------------------------+
|  Level: Home Row (asdfjkl;)    WPM: 34    Accuracy: 96%    3/10 |
+------------------------------------------------------------------+
|                                                                    |
|                     f j f k d l s ; a j                           |
|                     f j f k _                                     |
|                                                                    |
+------------------------------------------------------------------+
|        [1][2][3][4][5][6][7][8][9][0][-][=]                      |
|         [q][w][e][r][t][y][u][i][o][p]                            |
|          [a][s][d][f][g][h][j][k][l][;]                           |
|           [z][x][c][v][b][n][m][,][.][/]                          |
|                  [       space       ]                             |
|                                                                    |
|   L.pinky  L.ring  L.mid  L.index  R.index  R.mid  R.ring  R.pinky |
+------------------------------------------------------------------+
```

**Top bar:** Current level name, live WPM, live accuracy percentage, exercise progress (e.g. 3/10 drills completed in current phase).

**Middle zone:** The prompt line shows what to type. The input line shows what has been typed. Correct characters render green. The cursor position is highlighted. Errors render red and block forward progress until corrected with backspace.

**Bottom zone:** Virtual keyboard with finger-zone coloring. Each finger gets a distinct color. Unlocked keys show in their finger color. Locked keys are dimmed/greyed. The key that should be pressed next is highlighted brightly, providing constant reinforcement of correct finger assignment.

A legend below the keyboard maps colors to finger names.

## Level Progression System

### Level Sequence

Each level introduces 2-4 new keys, paired ergonomically (matching fingers on opposite hands when possible):

| Level | Keys | Fingers |
|-------|------|---------|
| 1 | f j | Index, home |
| 2 | d k | Middle |
| 3 | s l | Ring |
| 4 | a ; | Pinky |
| 5 | g h | Index, reach inward |
| 6 | t y | Index, top row |
| 7 | r u | Index, top row |
| 8 | e i | Middle, top row |
| 9 | w o | Ring, top row |
| 10 | q p | Pinky, top row |
| 11 | v m | Index, bottom row |
| 12 | c , | Middle, bottom row |
| 13 | x . | Ring, bottom row |
| 14 | z / | Pinky, bottom row |
| 15 | space | Thumbs (formally introduced; accepted as input from level 5 onward for word/code drills but not shown in character drills or on the keyboard until level 15) |
| 16 | shift + keys | Capitalization |
| 17 | 1-0 | Number row |
| 18 | symbols | - = [ ] etc. |

### Phases Within Each Level

Each level has three sequential phases:

1. **Character drills:** Raw repetition sequences built from new keys + all unlocked keys. New keys weighted at ~60%, review keys at ~40%. Short bursts of 20-30 characters, 10 drills per phase.

2. **Word drills:** Real English words filtered to only contain unlocked characters. Drawn from an embedded word list. 10 drills per phase. Word drills are skipped for levels where fewer than ~3 distinct characters are unlocked (levels 1-2), since no meaningful English words exist with only those letters. These levels run character drills only.

3. **Code drills:** Go and Python snippets using only unlocked characters. Unlocks around level 10+ when enough characters are available. Pre-curated snippets tagged with required character sets. 10 drills per phase.

### Promotion

- Each phase requires 95%+ aggregate accuracy across its 10 drills to pass.
- Failing a phase means repeating it with freshly generated drills.
- Passing a phase is permanent; you don't lose progress on completed phases.
- All three phases (or two, if code drills aren't yet available) must pass to unlock the next level.

## Core Typing Engine

### Input Handling

Bubble Tea captures raw key events. Each keypress is evaluated against the expected character at the current cursor position.

- **Correct key:** Character renders green, cursor advances.
- **Wrong key:** Character renders red, cursor does NOT advance. A short visual flash (red border or background pulse) signals the error. User must press backspace to clear the error before trying again.
- **Backspace:** Removes the last error. Cannot backspace through correct characters, only back to the last correct position.

### Metrics

- **WPM:** Standard calculation: (characters typed / 5) / minutes elapsed. Only counts correct forward progress, not error corrections. Updated live in the top bar.
- **Accuracy:** correct_keypresses / total_keypresses * 100. Every wrong keypress counts against accuracy, even if corrected. Incentivizes getting it right the first time.

### Drill Generation

- **Character drills:** Randomizer builds sequences from current + unlocked keys, weighted toward new keys.
- **Word drills:** Words selected from embedded word list, filtered to contain only unlocked characters.
- **Code drills:** Pre-curated snippets tagged with required character sets, selected when all required characters are unlocked.

### Phase Completion Flow

A phase consists of 10 drills. After each drill, that drill's accuracy is shown. After all 10, the aggregate phase accuracy determines pass/fail against the 95% gate. On failure, the phase repeats with fresh drills.

## Data Model and Persistence

### Progress File

Stored at `~/.config/typer/progress.json`:

```json
{
  "current_level": 4,
  "current_phase": "words",
  "levels": {
    "1": {
      "char_drills": {"completed": true, "best_accuracy": 98.2, "best_wpm": 22},
      "word_drills": {"completed": true, "best_accuracy": 96.1, "best_wpm": 18},
      "code_drills": null
    }
  },
  "sessions": [
    {
      "date": "2026-03-21T20:30:00Z",
      "duration_sec": 600,
      "level": 3,
      "phase": "chars",
      "accuracy": 94.8,
      "wpm": 19,
      "key_stats": {
        "d": {"correct": 45, "incorrect": 3},
        "k": {"correct": 42, "incorrect": 1}
      }
    }
  ]
}
```

### What's Tracked

- **Current position:** Level and phase, so you pick up where you left off.
- **Per-level bests:** Best accuracy and WPM for each completed phase.
- **Session log:** Append-only. Each practice session records date, duration, level, phase, accuracy, WPM, and per-key breakdown.
- **Per-key stats:** Correct/incorrect counts per key per session. Enables identifying problem keys over time.

### Level Definitions

Static Go code defining the level sequence, key introductions, and content generation rules. Not external config.

### Word List

English word list embedded via `//go:embed` from a text file.

### Code Snippets

Curated Go and Python snippets embedded in a Go source file, tagged with required character sets.

## Project Structure

```
typer/
  cmd/
    typer/
      main.go              # Entry point, Bubble Tea program init
  internal/
    app/
      app.go               # Top-level Bubble Tea model, Update/View
    keyboard/
      keyboard.go          # Virtual keyboard widget, finger-zone colors
    drill/
      engine.go            # Drill generation, input evaluation, stats
      content.go           # Word list filtering, code snippet selection
    level/
      levels.go            # Level definitions, progression logic
      promotion.go         # Phase/level completion, accuracy gating
    progress/
      store.go             # JSON read/write to ~/.config/typer/
      stats.go             # Session recording, per-key tracking
    ui/
      prompt.go            # Text prompt and input line rendering
      topbar.go            # WPM, accuracy, level display
      colors.go            # Finger-zone color scheme, error/success colors
  data/
    words.txt              # English word list (embedded)
    code_snippets.go       # Curated code snippets (embedded)
  go.mod
  go.sum
```

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling and layout
- `github.com/charmbracelet/bubbles` - Reusable components (if needed)

No other external dependencies. Progress file handled with stdlib `encoding/json` and `os`.

## Build

`go build ./cmd/typer` produces a single binary.

## Scope Boundaries

**In scope for v1:**
- Single-screen drill interface with virtual keyboard
- 18-level progressive unlock system
- Three-phase content (chars, words, code)
- Strict error correction (must fix before continuing)
- 95% accuracy gate for promotion
- Full progress persistence with session history and per-key stats
- Live WPM and accuracy display

**Explicitly out of scope:**
- Multi-screen navigation / menus / dashboards
- Stats visualization (WPM trends, heatmaps)
- Gamification (XP, streaks, achievements)
- Custom key layouts (Dvorak, Colemak)
- Multiplayer or online features
- Sound effects
- `typer stats` CLI subcommand (future addition)
