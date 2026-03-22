# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test

```bash
go build ./cmd/mavis        # build binary
go test ./...               # run all tests
go test ./internal/drill/   # run tests for a single package
go install ./cmd/mavis      # install to $GOPATH/bin
```

No Makefile, no linter config, no external tooling. Standard Go commands only.

## Architecture

Mavis is a touch typing retrainer TUI built with Go and Bubble Tea. The module is `github.com/nmelo/mavis` (the directory is named `typer` for historical reasons).

**Bubble Tea Model/Update/View pattern:** `internal/app/app.go` is the central hub. It owns a 4-state state machine (`stateDrilling` > `stateDrillComplete` > `statePhaseComplete` > `stateLevelComplete`) and composes all other packages.

**Package dependency flow:**
```
cmd/mavis/main.go
  └── internal/app      (imports everything below)
        ├── drill        (engine.go: DrillState, input eval, WPM/accuracy)
        │                (content.go: word list loading, filtering, code snippets)
        ├── level        (levels.go: 18 level definitions, finger mappings)
        │                (promotion.go: 95% accuracy gate, phase progression)
        ├── keyboard     (keyboard.go: virtual keyboard with background-colored keys)
        ├── progress     (store.go: JSON persistence to ~/.config/typer/progress.json)
        └── ui           (colors.go, topbar.go, prompt.go: styling and rendering)
```

**Data layer:** `data/embed.go` embeds `data/words.txt` (3000 English words) via `//go:embed`. `data/code_snippets.go` has curated Go/Python snippets.

## Key Design Decisions

- **Error blocks progress:** Wrong keypress sets `HasError`, cursor cannot advance until backspace clears it. This is intentional for retraining muscle memory.
- **Space excluded from char drills:** `GenerateCharDrill` skips space. Space is accepted from level 5+ in word/code drills only. The `Update` method silently ignores space during char phase.
- **Phase retry on failure:** `advancePhase()` keeps `currentPhase` unchanged if accuracy < 95%, causing a retry with fresh drills.
- **Content filtering:** Word drills and code snippets are filtered to only include characters the user has unlocked. Word drills are skipped for levels 1-2 (too few characters for real words).
- **Centering uses Align, not Place:** `lipgloss.Place` centers each line independently, which breaks the keyboard stagger. The app uses `lipgloss.NewStyle().Align().AlignVertical()` instead.
- **Background-colored keys, not borders:** Box-drawing border characters render inconsistently across terminals. Keys use solid background blocks instead.
