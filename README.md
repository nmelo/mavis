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

## How it plays

You launch `mavis`. The screen shows a virtual keyboard with only F and J lit up. Everything else is dimmed. A prompt appears above a warm input box with a cursor blinking on the first character.

**Level 1 is just two keys.** You drill F and J with your index fingers. Ten drills of 25 characters each. Pure muscle memory. The keyboard highlights which key to press next. Wrong key? Red flash, you're stuck until you backspace. Hit 95% accuracy across all 10 drills and you unlock Level 2 (D and K, middle fingers).

**Levels 1-2 are character drills only.** No English words exist with just those letters.

**Level 3 unlocks word drills and the story begins.** After passing character drills, you enter a second phase: a progressive cyberpunk narrative inspired by Vernor Vinge's *True Names*. The catch is that every sentence uses only the keys you've unlocked so far. Early levels are cryptic fragments. Later levels become full narrative.

The story builds as your keyboard grows:

- **Level 4** (+ A): `alas a sad flask` / `all lads ask dad`
- **Level 5** (+ G H): `flash ash gash hall` / `a glad lad had a dash`
- **Level 8** (+ E I): `the digital shield is hid` / `fire strikes their dark field`
- **Level 10** (+ Q P): `speak the password to the portal` / `quest through opaque pathways`
- **Level 12** (+ C): `the cult circles, code is core` / `cosmic data cascades from the rift`
- **Level 14** (+ Z): `the wizard gazed at a froze maze` / `fuzzy daze grips the dizzied quiz`

**Hit 100% on a drill** and sparkles burst across the screen with a golden PERFECT message. **Complete a level** and shooting stars streak across with a bigger celebration. Fail a phase (below 95%) and you retry with fresh content.

**Code drills** unlock at Level 10, with curated Go and Python snippets filtered to your available keys.

By Level 14 you have all lowercase letters. Levels 15-18 add space, shift, numbers, and symbols.

## Controls

| Key | Action |
|-----|--------|
| Type | Match the characters shown |
| Backspace | Clear an error (errors block progress) |
| Enter | Advance to next drill |
| Esc | Quit (progress saves automatically) |
| Ctrl+L | Open level selector |
| Ctrl+N | Skip to next level |
| Ctrl+B | Go back one level |

## Level selector

Press Ctrl+L to open a two-pane modal. Left side lists all 18 levels grouped by keyboard region (home, top, bottom). Right side shows the phases (char/word/code drills) for the highlighted level. Navigate with j/k, Tab to switch panes, Enter to jump in.

## Progress

Everything saves to `~/.config/typer/progress.json`. Per-drill accuracy, WPM, per-key error rates, session history. Pick up where you left off.

## License

MIT
