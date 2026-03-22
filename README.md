# mavis

A terminal-based touch typing retrainer. Not for learning from scratch, but for breaking bad habits and rebuilding with proper home-row technique.

Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Install

```bash
go install github.com/nmelo/mavis/cmd/mavis@latest
```

Or build from source:

```bash
git clone https://github.com/nmelo/mavis.git
cd mavis
go build ./cmd/mavis
./mavis
```

## How it works

You start at Level 1 with just `f` and `j`. Each level introduces 2-4 new keys, paired by finger. You must hit 95% accuracy across 10 drills to unlock the next level.

**18 levels** from home row through numbers and symbols. Three drill phases per level:

- **Character drills** -- raw key repetition, weighted toward new keys
- **Word drills** -- real English words using only your unlocked keys
- **Code drills** -- Go and Python snippets (unlocks at level 10)

Mistakes block progress. You must backspace and correct before continuing. The virtual keyboard shows which finger should hit each key.

Progress saves automatically to `~/.config/typer/progress.json`.

## Controls

- Type the characters shown on screen
- Backspace to clear errors
- Enter to advance between drills
- Esc to quit (progress is saved)
