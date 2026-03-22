# Touch Typing Retraining TUI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a terminal-based touch typing retrainer with progressive key unlock, strict error correction, and full progress persistence.

**Architecture:** Single-screen Bubble Tea app. Pure data/logic layers (levels, drill engine, progress) are fully testable without TUI. UI components (keyboard widget, prompt, topbar) compose into a single app model. All state flows through the Bubble Tea Model/Update/View cycle.

**Tech Stack:** Go, Bubble Tea, Lipgloss, stdlib encoding/json for persistence, go:embed for word list.

**Spec:** `docs/superpowers/specs/2026-03-21-touch-typing-tui-design.md`

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `cmd/typer/main.go` (placeholder)
- Create: `internal/` directory tree

- [ ] **Step 1: Initialize Go module**

```bash
cd /Users/nmelo/Desktop/Projects/typer
go mod init github.com/nmelo/typer
```

- [ ] **Step 2: Install dependencies**

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
```

- [ ] **Step 3: Create directory structure**

```bash
mkdir -p cmd/typer
mkdir -p internal/{app,keyboard,drill,level,progress,ui}
mkdir -p data
```

- [ ] **Step 4: Create placeholder main.go**

Create `cmd/typer/main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("typer")
}
```

- [ ] **Step 5: Verify build**

```bash
go build ./cmd/typer
```

Expected: builds successfully, produces `typer` binary.

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum cmd/ internal/ data/
git commit -m "scaffold: initialize project structure and dependencies"
```

---

### Task 2: Level Definitions

**Files:**
- Create: `internal/level/levels.go`
- Create: `internal/level/levels_test.go`

This is pure data with helper functions. No dependencies on other packages.

- [ ] **Step 1: Write tests for level definitions**

Create `internal/level/levels_test.go`:

```go
package level

import "testing"

func TestAllLevelsExist(t *testing.T) {
	levels := All()
	if len(levels) != 18 {
		t.Fatalf("expected 18 levels, got %d", len(levels))
	}
}

func TestLevel1Keys(t *testing.T) {
	levels := All()
	l := levels[0]
	if l.Number != 1 {
		t.Errorf("first level number = %d, want 1", l.Number)
	}
	if string(l.NewKeys) != "fj" {
		t.Errorf("level 1 new keys = %q, want %q", string(l.NewKeys), "fj")
	}
	if l.Name != "Home Row: f j" {
		t.Errorf("level 1 name = %q, want %q", l.Name, "Home Row: f j")
	}
}

func TestUnlockedKeysAccumulate(t *testing.T) {
	levels := All()
	// After level 4 (a;), all home row keys should be unlocked
	unlocked := UnlockedKeys(4)
	expected := []rune{'f', 'j', 'd', 'k', 's', 'l', 'a', ';'}
	if len(unlocked) != len(expected) {
		t.Fatalf("unlocked keys at level 4: got %d, want %d", len(unlocked), len(expected))
	}
	for _, r := range expected {
		found := false
		for _, u := range unlocked {
			if u == r {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %c to be unlocked at level 4", r)
		}
	}
}

func TestLevelHasWordDrills(t *testing.T) {
	levels := All()
	// Levels 1-2 should NOT have word drills (fewer than 3 distinct chars)
	if levels[0].HasWordDrills {
		t.Error("level 1 should not have word drills")
	}
	if levels[1].HasWordDrills {
		t.Error("level 2 should not have word drills")
	}
	// Level 3 (s,l + f,j,d,k = 6 chars) should have word drills
	if !levels[2].HasWordDrills {
		t.Error("level 3 should have word drills")
	}
}

func TestLevelHasCodeDrills(t *testing.T) {
	levels := All()
	// Code drills unlock at level 10+
	for i := 0; i < 9; i++ {
		if levels[i].HasCodeDrills {
			t.Errorf("level %d should not have code drills", i+1)
		}
	}
	if !levels[9].HasCodeDrills {
		t.Error("level 10 should have code drills")
	}
}

func TestFingerAssignment(t *testing.T) {
	// Verify every new key in every level has a finger assignment
	levels := All()
	for _, l := range levels {
		for _, k := range l.NewKeys {
			finger := FingerForKey(k)
			if finger == "" {
				t.Errorf("level %d: key %c has no finger assignment", l.Number, k)
			}
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/level/ -v
```

Expected: compilation errors (types don't exist yet).

- [ ] **Step 3: Implement level definitions**

Create `internal/level/levels.go`:

```go
package level

// Finger represents which finger should press a key.
type Finger string

const (
	LPinky  Finger = "L.pinky"
	LRing   Finger = "L.ring"
	LMid    Finger = "L.mid"
	LIndex  Finger = "L.index"
	RIndex  Finger = "R.index"
	RMid    Finger = "R.mid"
	RRing   Finger = "R.ring"
	RPinky  Finger = "R.pinky"
	LThumb  Finger = "L.thumb"
	RThumb  Finger = "R.thumb"
)

// Level defines a progression level with the keys it introduces.
type Level struct {
	Number        int
	Name          string
	NewKeys       []rune
	HasWordDrills bool
	HasCodeDrills bool
}

var fingerMap = map[rune]Finger{
	// Home row
	'a': LPinky, 's': LRing, 'd': LMid, 'f': LIndex,
	'j': RIndex, 'k': RMid, 'l': RRing, ';': RPinky,
	'g': LIndex, 'h': RIndex,
	// Top row
	'q': LPinky, 'w': LRing, 'e': LMid, 'r': LIndex, 't': LIndex,
	'y': RIndex, 'u': RIndex, 'i': RMid, 'o': RRing, 'p': RPinky,
	// Bottom row
	'z': LPinky, 'x': LRing, 'c': LMid, 'v': LIndex, 'b': LIndex,
	'n': RIndex, 'm': RIndex, ',': RMid, '.': RRing, '/': RPinky,
	// Space
	' ': RThumb,
	// Numbers
	'1': LPinky, '2': LRing, '3': LMid, '4': LIndex, '5': LIndex,
	'6': RIndex, '7': RIndex, '8': RMid, '9': RRing, '0': RPinky,
	// Symbols
	'-': RPinky, '=': RPinky, '[': RPinky, ']': RPinky,
	'\\': RPinky, '\'': RPinky,
}

var levels = []Level{
	{1, "Home Row: f j", []rune{'f', 'j'}, false, false},
	{2, "Home Row: d k", []rune{'d', 'k'}, false, false},
	{3, "Home Row: s l", []rune{'s', 'l'}, true, false},
	{4, "Home Row: a ;", []rune{'a', ';'}, true, false},
	{5, "Home Row: g h", []rune{'g', 'h'}, true, false},
	{6, "Top Row: t y", []rune{'t', 'y'}, true, false},
	{7, "Top Row: r u", []rune{'r', 'u'}, true, false},
	{8, "Top Row: e i", []rune{'e', 'i'}, true, false},
	{9, "Top Row: w o", []rune{'w', 'o'}, true, false},
	{10, "Top Row: q p", []rune{'q', 'p'}, true, true},
	{11, "Bottom Row: v m", []rune{'v', 'm'}, true, true},
	{12, "Bottom Row: c ,", []rune{'c', ','}, true, true},
	{13, "Bottom Row: x .", []rune{'x', '.'}, true, true},
	{14, "Bottom Row: z /", []rune{'z', '/'}, true, true},
	{15, "Space", []rune{' '}, true, true},
	{16, "Shift + Keys", []rune{}, true, true},
	{17, "Number Row", []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}, true, true},
	{18, "Symbols", []rune{'-', '=', '[', ']', '\\', '\''}, true, true},
}

// All returns all level definitions in order.
func All() []Level {
	return levels
}

// Get returns the level definition for the given level number (1-indexed).
func Get(n int) Level {
	return levels[n-1]
}

// UnlockedKeys returns all keys unlocked up to and including the given level.
func UnlockedKeys(levelNum int) []rune {
	var keys []rune
	for i := 0; i < levelNum && i < len(levels); i++ {
		keys = append(keys, levels[i].NewKeys...)
	}
	return keys
}

// FingerForKey returns the finger assignment for a given key.
func FingerForKey(k rune) Finger {
	return fingerMap[k]
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/level/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/level/
git commit -m "feat: add level definitions with finger mappings and progression data"
```

---

### Task 3: Drill Engine - Character Drill Generation

**Files:**
- Create: `internal/drill/engine.go`
- Create: `internal/drill/engine_test.go`

- [ ] **Step 1: Write tests for character drill generation**

Create `internal/drill/engine_test.go`:

```go
package drill

import (
	"testing"

	"github.com/nmelo/typer/internal/level"
)

func TestGenerateCharDrill(t *testing.T) {
	keys := level.UnlockedKeys(1) // f, j
	drill := GenerateCharDrill(keys, level.Get(1).NewKeys, 25)

	if len(drill) != 25 {
		t.Fatalf("drill length = %d, want 25", len(drill))
	}

	allowed := map[rune]bool{'f': true, 'j': true}
	for i, ch := range drill {
		if !allowed[ch] && ch != ' ' {
			t.Errorf("position %d: unexpected character %c", i, ch)
		}
	}
}

func TestCharDrillWeightsNewKeys(t *testing.T) {
	// Level 5: new keys are g, h; unlocked includes f,j,d,k,s,l,a,;,g,h
	keys := level.UnlockedKeys(5)
	newKeys := level.Get(5).NewKeys

	counts := map[rune]int{}
	// Generate many drills to check distribution
	for i := 0; i < 100; i++ {
		drill := GenerateCharDrill(keys, newKeys, 100)
		for _, ch := range drill {
			counts[ch]++
		}
	}

	// New keys (g, h) should appear more frequently than any single review key
	newKeyCount := counts['g'] + counts['h']
	reviewKeyCount := counts['f'] + counts['j'] + counts['d'] + counts['k'] +
		counts['s'] + counts['l'] + counts['a'] + counts[';']

	// New keys are 2 of 10 total keys but should be ~60% of chars
	newKeyRatio := float64(newKeyCount) / float64(newKeyCount+reviewKeyCount)
	if newKeyRatio < 0.45 || newKeyRatio > 0.75 {
		t.Errorf("new key ratio = %.2f, want ~0.60 (between 0.45 and 0.75)", newKeyRatio)
	}
}

func TestGenerateCharDrillMinLength(t *testing.T) {
	keys := level.UnlockedKeys(1)
	drill := GenerateCharDrill(keys, level.Get(1).NewKeys, 1)
	if len(drill) < 1 {
		t.Error("drill should have at least 1 character")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/drill/ -v
```

Expected: compilation errors.

- [ ] **Step 3: Implement character drill generation**

Create `internal/drill/engine.go`:

```go
package drill

import (
	"math/rand"
)

// GenerateCharDrill creates a character drill of the given length.
// newKeys are weighted at ~60%, review keys at ~40%.
func GenerateCharDrill(allKeys []rune, newKeys []rune, length int) []rune {
	if len(allKeys) == 0 {
		return nil
	}

	// Build weighted pool: new keys get 3x weight
	// Exclude space from char drills (space is for word/code drills only per spec)
	var pool []rune
	newSet := make(map[rune]bool)
	for _, k := range newKeys {
		newSet[k] = true
	}

	for _, k := range allKeys {
		if k == ' ' {
			continue
		}
		if newSet[k] {
			pool = append(pool, k, k, k) // 3x weight for new keys
		} else {
			pool = append(pool, k)
		}
	}

	if len(pool) == 0 {
		return nil
	}

	drill := make([]rune, length)
	for i := range drill {
		drill[i] = pool[rand.Intn(len(pool))]
	}
	return drill
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/drill/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/drill/
git commit -m "feat: add character drill generation with weighted key distribution"
```

---

### Task 4: Drill Engine - Input Evaluation

**Files:**
- Modify: `internal/drill/engine.go`
- Modify: `internal/drill/engine_test.go`

The core typing logic: evaluate keypresses, track accuracy, calculate WPM.

- [ ] **Step 1: Write tests for input evaluation**

Append to `internal/drill/engine_test.go`:

```go
func TestNewDrillState(t *testing.T) {
	prompt := []rune("fjfj")
	state := NewDrillState(prompt)

	if state.Position != 0 {
		t.Errorf("initial position = %d, want 0", state.Position)
	}
	if state.TotalKeystrokes != 0 {
		t.Errorf("initial keystrokes = %d, want 0", state.TotalKeystrokes)
	}
	if state.Errors != 0 {
		t.Errorf("initial errors = %d, want 0", state.Errors)
	}
	if state.HasError {
		t.Error("should not start with error")
	}
	if state.Complete {
		t.Error("should not start complete")
	}
}

func TestCorrectKeypress(t *testing.T) {
	state := NewDrillState([]rune("fj"))
	state.HandleKey('f')

	if state.Position != 1 {
		t.Errorf("position after correct key = %d, want 1", state.Position)
	}
	if state.TotalKeystrokes != 1 {
		t.Errorf("keystrokes = %d, want 1", state.TotalKeystrokes)
	}
	if state.HasError {
		t.Error("should not have error after correct key")
	}
}

func TestWrongKeypress(t *testing.T) {
	state := NewDrillState([]rune("fj"))
	state.HandleKey('j') // wrong, expected 'f'

	if state.Position != 0 {
		t.Errorf("position after wrong key = %d, want 0 (should not advance)", state.Position)
	}
	if state.TotalKeystrokes != 1 {
		t.Errorf("keystrokes = %d, want 1", state.TotalKeystrokes)
	}
	if state.Errors != 1 {
		t.Errorf("errors = %d, want 1", state.Errors)
	}
	if !state.HasError {
		t.Error("should have error flag set")
	}
}

func TestBackspaceOnError(t *testing.T) {
	state := NewDrillState([]rune("fj"))
	state.HandleKey('j') // wrong
	state.HandleBackspace()

	if state.HasError {
		t.Error("error should be cleared after backspace")
	}
	if state.Position != 0 {
		t.Errorf("position after backspace = %d, want 0", state.Position)
	}
}

func TestBackspaceOnCorrectDoesNothing(t *testing.T) {
	state := NewDrillState([]rune("fj"))
	state.HandleKey('f') // correct
	state.HandleBackspace()

	// Cannot backspace through correct characters
	if state.Position != 1 {
		t.Errorf("position = %d, want 1 (backspace should not undo correct chars)", state.Position)
	}
}

func TestDrillCompletion(t *testing.T) {
	state := NewDrillState([]rune("fj"))
	state.HandleKey('f')
	state.HandleKey('j')

	if !state.Complete {
		t.Error("drill should be complete")
	}
}

func TestAccuracy(t *testing.T) {
	state := NewDrillState([]rune("fj"))
	state.HandleKey('f') // correct
	state.HandleKey('k') // wrong
	state.HandleBackspace()
	state.HandleKey('j') // correct

	// 2 correct out of 3 total keystrokes
	acc := state.Accuracy()
	expected := (2.0 / 3.0) * 100
	if acc < expected-0.1 || acc > expected+0.1 {
		t.Errorf("accuracy = %.1f, want %.1f", acc, expected)
	}
}

func TestKeyStats(t *testing.T) {
	state := NewDrillState([]rune("ff"))
	state.HandleKey('f') // correct
	state.HandleKey('j') // wrong - expected 'f', pressed 'j'
	state.HandleBackspace()
	state.HandleKey('f') // correct

	stats := state.KeyStats()
	if stats['f'].Correct != 2 {
		t.Errorf("f correct = %d, want 2", stats['f'].Correct)
	}
	// 'j' was the wrong key pressed - error is tracked on the key that was pressed
	if stats['j'].Incorrect != 1 {
		t.Errorf("j incorrect = %d, want 1", stats['j'].Incorrect)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/drill/ -v
```

Expected: compilation errors (DrillState type doesn't exist).

- [ ] **Step 3: Implement DrillState**

Add to `internal/drill/engine.go`:

```go
import "time"

// KeyStat tracks correct and incorrect presses for a key.
type KeyStat struct {
	Correct   int
	Incorrect int
}

// DrillState tracks the state of an active drill.
type DrillState struct {
	Prompt          []rune
	Position        int
	TotalKeystrokes int
	Errors          int
	HasError        bool
	Complete        bool
	StartTime       time.Time
	keyCorrect      map[rune]int
	keyIncorrect    map[rune]int
}

// NewDrillState creates a new drill state for the given prompt.
func NewDrillState(prompt []rune) *DrillState {
	return &DrillState{
		Prompt:       prompt,
		StartTime:    time.Now(),
		keyCorrect:   make(map[rune]int),
		keyIncorrect: make(map[rune]int),
	}
}

// HandleKey processes a keypress against the current expected character.
func (d *DrillState) HandleKey(ch rune) {
	if d.Complete || d.HasError {
		return
	}

	d.TotalKeystrokes++
	expected := d.Prompt[d.Position]

	if ch == expected {
		d.keyCorrect[expected]++
		d.Position++
		if d.Position >= len(d.Prompt) {
			d.Complete = true
		}
	} else {
		d.Errors++
		d.keyIncorrect[ch]++
		d.HasError = true
	}
}

// HandleBackspace clears the current error state.
func (d *DrillState) HandleBackspace() {
	if d.HasError {
		d.HasError = false
	}
}

// Accuracy returns the accuracy percentage.
func (d *DrillState) Accuracy() float64 {
	if d.TotalKeystrokes == 0 {
		return 100.0
	}
	correct := d.TotalKeystrokes - d.Errors
	return float64(correct) / float64(d.TotalKeystrokes) * 100
}

// WPM returns the current words-per-minute based on correct characters typed.
func (d *DrillState) WPM() float64 {
	elapsed := time.Since(d.StartTime).Minutes()
	if elapsed < 0.001 {
		return 0
	}
	return (float64(d.Position) / 5.0) / elapsed
}

// KeyStats returns per-key accuracy statistics.
func (d *DrillState) KeyStats() map[rune]KeyStat {
	stats := make(map[rune]KeyStat)
	for k, v := range d.keyCorrect {
		s := stats[k]
		s.Correct = v
		stats[k] = s
	}
	for k, v := range d.keyIncorrect {
		s := stats[k]
		s.Incorrect = v
		stats[k] = s
	}
	return stats
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/drill/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/drill/
git commit -m "feat: add drill state with input evaluation, accuracy, and WPM tracking"
```

---

### Task 5: Word List and Word Drill Generation

**Files:**
- Create: `data/words.txt`
- Create: `data/embed.go`
- Create: `internal/drill/content.go`
- Create: `internal/drill/content_test.go`

- [ ] **Step 1: Write tests for word filtering**

Create `internal/drill/content_test.go`:

```go
package drill

import "testing"

func TestFilterWordsForKeys(t *testing.T) {
	words := []string{"fall", "salad", "flask", "hello", "dad", "ask"}
	keys := map[rune]bool{'f': true, 'j': true, 'd': true, 'k': true, 's': true, 'l': true, 'a': true, ';': true}

	filtered := FilterWords(words, keys)

	// "hello" has 'h', 'e', 'o' which are not in the key set
	for _, w := range filtered {
		for _, ch := range w {
			if !keys[ch] {
				t.Errorf("word %q contains unlocked char %c", w, ch)
			}
		}
	}

	// Should include words like "fall", "salad", "flask", "dad", "ask"
	if len(filtered) < 4 {
		t.Errorf("expected at least 4 filtered words, got %d", len(filtered))
	}
}

func TestFilterWordsEmptyResult(t *testing.T) {
	words := []string{"hello", "world"}
	keys := map[rune]bool{'f': true, 'j': true}

	filtered := FilterWords(words, keys)
	if len(filtered) != 0 {
		t.Errorf("expected 0 words, got %d", len(filtered))
	}
}

func TestGenerateWordDrill(t *testing.T) {
	words := []string{"fall", "dad", "ask", "flask", "salad"}
	drill := GenerateWordDrill(words, 5)

	if len(drill) != 5 {
		t.Fatalf("expected 5 words in drill, got %d", len(drill))
	}

	for _, w := range drill {
		found := false
		for _, src := range words {
			if w == src {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("drill word %q not in source list", w)
		}
	}
}

func TestLoadWordList(t *testing.T) {
	words := LoadWordList()
	if len(words) < 100 {
		t.Errorf("word list too small: %d words", len(words))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/drill/ -v
```

Expected: compilation errors.

- [ ] **Step 3: Create word list data file**

Create `data/words.txt` with a curated English word list (common words, one per line). Should contain at least 2000 words. Source from a standard English frequency list. Include only lowercase words. Example (truncated, the actual file should have 2000+ entries):

```
the
of
and
a
to
in
is
you
that
it
he
was
for
on
are
...
```

Use this command to generate the initial word list from a standard source:

```bash
curl -sL "https://raw.githubusercontent.com/first20hours/google-10000-english/master/google-10000-english-no-swears.txt" | head -3000 > data/words.txt
```

- [ ] **Step 4: Create embed file**

Create `data/embed.go`:

```go
package data

import _ "embed"

//go:embed words.txt
var WordList string
```

- [ ] **Step 5: Implement word filtering and drill generation**

Create `internal/drill/content.go`:

```go
package drill

import (
	"math/rand"
	"strings"

	"github.com/nmelo/typer/data"
)

// LoadWordList returns the embedded word list as a slice of strings.
func LoadWordList() []string {
	var words []string
	for _, w := range strings.Split(data.WordList, "\n") {
		w = strings.TrimSpace(w)
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}

// FilterWords returns only words whose characters are all in the allowed key set.
func FilterWords(words []string, allowedKeys map[rune]bool) []string {
	var filtered []string
	for _, w := range words {
		ok := true
		for _, ch := range w {
			if !allowedKeys[ch] {
				ok = false
				break
			}
		}
		if ok && len(w) > 0 {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

// GenerateWordDrill picks n random words from the provided word list.
func GenerateWordDrill(words []string, n int) []string {
	if len(words) == 0 {
		return nil
	}
	drill := make([]string, n)
	for i := range drill {
		drill[i] = words[rand.Intn(len(words))]
	}
	return drill
}
```

- [ ] **Step 6: Run tests to verify they pass**

```bash
go test ./internal/drill/ -v
```

Expected: all tests PASS.

- [ ] **Step 7: Commit**

```bash
git add data/ internal/drill/content.go internal/drill/content_test.go
git commit -m "feat: add word list embedding and word drill generation with key filtering"
```

---

### Task 6: Progress Persistence

**Files:**
- Create: `internal/progress/store.go`
- Create: `internal/progress/store_test.go`

- [ ] **Step 1: Write tests for progress store**

Create `internal/progress/store_test.go`:

```go
package progress

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewProgress(t *testing.T) {
	p := New()
	if p.CurrentLevel != 1 {
		t.Errorf("initial level = %d, want 1", p.CurrentLevel)
	}
	if p.CurrentPhase != PhaseChars {
		t.Errorf("initial phase = %q, want %q", p.CurrentPhase, PhaseChars)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	p := New()
	p.CurrentLevel = 3
	p.CurrentPhase = PhaseWords

	err := Save(p, path)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.CurrentLevel != 3 {
		t.Errorf("loaded level = %d, want 3", loaded.CurrentLevel)
	}
	if loaded.CurrentPhase != PhaseWords {
		t.Errorf("loaded phase = %q, want %q", loaded.CurrentPhase, PhaseWords)
	}
}

func TestLoadNonExistentReturnsNew(t *testing.T) {
	p, err := Load("/nonexistent/path/progress.json")
	if err != nil {
		t.Fatalf("load should not error for missing file: %v", err)
	}
	if p.CurrentLevel != 1 {
		t.Errorf("should return fresh progress, got level %d", p.CurrentLevel)
	}
}

func TestRecordSession(t *testing.T) {
	p := New()
	session := Session{
		Date:     time.Now(),
		Duration: 300,
		Level:    1,
		Phase:    PhaseChars,
		Accuracy: 96.5,
		WPM:      22.0,
		KeyStats: map[string]KeyStat{
			"f": {Correct: 50, Incorrect: 2},
		},
	}

	p.RecordSession(session)

	if len(p.Sessions) != 1 {
		t.Fatalf("sessions count = %d, want 1", len(p.Sessions))
	}
	if p.Sessions[0].Accuracy != 96.5 {
		t.Errorf("session accuracy = %.1f, want 96.5", p.Sessions[0].Accuracy)
	}
}

func TestCompletePhaseSetsLevelBest(t *testing.T) {
	p := New()
	p.CompletePhase(1, PhaseChars, 97.3, 25.0)

	lvl, ok := p.Levels["1"]
	if !ok {
		t.Fatal("level 1 not recorded")
	}
	if lvl.CharDrills == nil {
		t.Fatal("char drills not recorded")
	}
	if !lvl.CharDrills.Completed {
		t.Error("char drills should be completed")
	}
	if lvl.CharDrills.BestAccuracy != 97.3 {
		t.Errorf("best accuracy = %.1f, want 97.3", lvl.CharDrills.BestAccuracy)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "progress.json")

	p := New()
	err := Save(p, path)
	if err != nil {
		t.Fatalf("save should create parent dirs: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file should exist after save")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/progress/ -v
```

Expected: compilation errors.

- [ ] **Step 3: Implement progress store**

Create `internal/progress/store.go`:

```go
package progress

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	PhaseChars = "chars"
	PhaseWords = "words"
	PhaseCode  = "code"
)

type KeyStat struct {
	Correct   int `json:"correct"`
	Incorrect int `json:"incorrect"`
}

type Session struct {
	Date     time.Time          `json:"date"`
	Duration int                `json:"duration_sec"`
	Level    int                `json:"level"`
	Phase    string             `json:"phase"`
	Accuracy float64            `json:"accuracy"`
	WPM      float64            `json:"wpm"`
	KeyStats map[string]KeyStat `json:"key_stats"`
}

type PhaseResult struct {
	Completed    bool    `json:"completed"`
	BestAccuracy float64 `json:"best_accuracy"`
	BestWPM      float64 `json:"best_wpm"`
}

type LevelProgress struct {
	CharDrills *PhaseResult `json:"char_drills"`
	WordDrills *PhaseResult `json:"word_drills"`
	CodeDrills *PhaseResult `json:"code_drills"`
}

type Progress struct {
	CurrentLevel int                      `json:"current_level"`
	CurrentPhase string                   `json:"current_phase"`
	Levels       map[string]LevelProgress `json:"levels"`
	Sessions     []Session                `json:"sessions"`
}

// New returns a fresh progress state starting at level 1.
func New() *Progress {
	return &Progress{
		CurrentLevel: 1,
		CurrentPhase: PhaseChars,
		Levels:       make(map[string]LevelProgress),
	}
}

// Load reads progress from the given file path. Returns fresh progress if file doesn't exist.
func Load(path string) (*Progress, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return nil, err
	}
	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.Levels == nil {
		p.Levels = make(map[string]LevelProgress)
	}
	return &p, nil
}

// Save writes progress to the given file path, creating parent directories.
func Save(p *Progress, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// RecordSession appends a session to the progress.
func (p *Progress) RecordSession(s Session) {
	p.Sessions = append(p.Sessions, s)
}

// CompletePhase marks a phase as completed with its best stats.
func (p *Progress) CompletePhase(levelNum int, phase string, accuracy, wpm float64) {
	key := strconv.Itoa(levelNum)

	lvl := p.Levels[key]
	result := &PhaseResult{
		Completed:    true,
		BestAccuracy: accuracy,
		BestWPM:      wpm,
	}

	switch phase {
	case PhaseChars:
		lvl.CharDrills = result
	case PhaseWords:
		lvl.WordDrills = result
	case PhaseCode:
		lvl.CodeDrills = result
	}

	p.Levels[key] = lvl
}

// DefaultPath returns the default progress file path.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "typer", "progress.json")
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/progress/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/progress/
git commit -m "feat: add progress persistence with session recording and phase completion"
```

---

### Task 7: Promotion Logic

**Files:**
- Create: `internal/level/promotion.go`
- Create: `internal/level/promotion_test.go`

- [ ] **Step 1: Write tests for promotion logic**

Create `internal/level/promotion_test.go`:

```go
package level

import "testing"

func TestPhasePassesAt95(t *testing.T) {
	if !PhasePassed(95.0) {
		t.Error("95.0% should pass")
	}
	if !PhasePassed(100.0) {
		t.Error("100.0% should pass")
	}
}

func TestPhaseFailsBelow95(t *testing.T) {
	if PhasePassed(94.9) {
		t.Error("94.9% should not pass")
	}
	if PhasePassed(0) {
		t.Error("0% should not pass")
	}
}

func TestNextPhaseAfterChars(t *testing.T) {
	lvl := All()[2] // Level 3: has word drills, no code drills
	next := NextPhase("chars", lvl)
	if next != "words" {
		t.Errorf("next phase after chars = %q, want %q", next, "words")
	}
}

func TestNextPhaseAfterCharsNoWords(t *testing.T) {
	lvl := All()[0] // Level 1: no word drills, no code drills
	next := NextPhase("chars", lvl)
	if next != "" {
		t.Errorf("next phase after chars (no words) = %q, want empty (level complete)", next)
	}
}

func TestNextPhaseAfterWordsWithCode(t *testing.T) {
	lvl := All()[9] // Level 10: has code drills
	next := NextPhase("words", lvl)
	if next != "code" {
		t.Errorf("next phase after words = %q, want %q", next, "code")
	}
}

func TestNextPhaseAfterWordsNoCode(t *testing.T) {
	lvl := All()[2] // Level 3: no code drills
	next := NextPhase("words", lvl)
	if next != "" {
		t.Errorf("next phase after words (no code) = %q, want empty (level complete)", next)
	}
}

func TestNextPhaseAfterCode(t *testing.T) {
	lvl := All()[9] // Level 10
	next := NextPhase("code", lvl)
	if next != "" {
		t.Errorf("next phase after code = %q, want empty (level complete)", next)
	}
}

func TestPhasesForLevel(t *testing.T) {
	// Level 1: chars only
	phases := PhasesFor(All()[0])
	if len(phases) != 1 || phases[0] != "chars" {
		t.Errorf("level 1 phases = %v, want [chars]", phases)
	}

	// Level 3: chars + words
	phases = PhasesFor(All()[2])
	if len(phases) != 2 {
		t.Errorf("level 3 phases = %v, want [chars, words]", phases)
	}

	// Level 10: chars + words + code
	phases = PhasesFor(All()[9])
	if len(phases) != 3 {
		t.Errorf("level 10 phases = %v, want [chars, words, code]", phases)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/level/ -v -run "Phase|Promotion|Next"
```

Expected: compilation errors.

- [ ] **Step 3: Implement promotion logic**

Create `internal/level/promotion.go`:

```go
package level

const AccuracyThreshold = 95.0

// PhasePassed returns true if the given accuracy meets the promotion threshold.
func PhasePassed(accuracy float64) bool {
	return accuracy >= AccuracyThreshold
}

// NextPhase returns the next phase for a level, or "" if the level is complete.
func NextPhase(currentPhase string, lvl Level) string {
	switch currentPhase {
	case "chars":
		if lvl.HasWordDrills {
			return "words"
		}
		if lvl.HasCodeDrills {
			return "code"
		}
		return ""
	case "words":
		if lvl.HasCodeDrills {
			return "code"
		}
		return ""
	case "code":
		return ""
	}
	return ""
}

// PhasesFor returns the list of phases a level requires.
func PhasesFor(lvl Level) []string {
	phases := []string{"chars"}
	if lvl.HasWordDrills {
		phases = append(phases, "words")
	}
	if lvl.HasCodeDrills {
		phases = append(phases, "code")
	}
	return phases
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/level/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/level/promotion.go internal/level/promotion_test.go
git commit -m "feat: add phase promotion logic with 95% accuracy gate"
```

---

### Task 8: UI Colors and Finger Zone Scheme

**Files:**
- Create: `internal/ui/colors.go`
- Create: `internal/ui/colors_test.go`

- [ ] **Step 1: Write tests for color assignments**

Create `internal/ui/colors_test.go`:

```go
package ui

import (
	"testing"

	"github.com/nmelo/typer/internal/level"
)

func TestEveryFingerHasColor(t *testing.T) {
	fingers := []level.Finger{
		level.LPinky, level.LRing, level.LMid, level.LIndex,
		level.RIndex, level.RMid, level.RRing, level.RPinky,
		level.LThumb, level.RThumb,
	}
	for _, f := range fingers {
		c := ColorForFinger(f)
		if c == "" {
			t.Errorf("finger %q has no color", f)
		}
	}
}

func TestMatchingFingersShareColor(t *testing.T) {
	// Left and right index should be the same color
	if ColorForFinger(level.LIndex) != ColorForFinger(level.RIndex) {
		t.Error("L.index and R.index should share color")
	}
}

func TestDistinctColorsForDifferentFingers(t *testing.T) {
	colors := map[string]bool{}
	fingers := []level.Finger{level.LPinky, level.LRing, level.LMid, level.LIndex, level.LThumb}
	for _, f := range fingers {
		colors[ColorForFinger(f)] = true
	}
	if len(colors) < 4 {
		t.Errorf("expected at least 4 distinct colors, got %d", len(colors))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/ui/ -v
```

Expected: compilation errors.

- [ ] **Step 3: Implement color scheme**

Create `internal/ui/colors.go`:

```go
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/typer/internal/level"
)

// Hex color constants for finger zones.
const (
	ColorPinky  = "#FF6B9D" // pink
	ColorRing   = "#C084FC" // purple
	ColorMiddle = "#60A5FA" // blue
	ColorIndex  = "#34D399" // green
	ColorThumb  = "#FBBF24" // yellow

	ColorCorrect  = "#22C55E" // green
	ColorError    = "#EF4444" // red
	ColorLocked   = "#4B5563" // dim gray
	ColorCursor   = "#F59E0B" // amber
	ColorNextKey  = "#FFFFFF" // bright white
	ColorDimText  = "#6B7280" // gray
)

var fingerColors = map[level.Finger]string{
	level.LPinky:  ColorPinky,
	level.LRing:   ColorRing,
	level.LMid:    ColorMiddle,
	level.LIndex:  ColorIndex,
	level.RIndex:  ColorIndex,
	level.RMid:    ColorMiddle,
	level.RRing:   ColorRing,
	level.RPinky:  ColorPinky,
	level.LThumb:  ColorThumb,
	level.RThumb:  ColorThumb,
}

// ColorForFinger returns the hex color string for a given finger.
func ColorForFinger(f level.Finger) string {
	return fingerColors[f]
}

// StyleForFinger returns a lipgloss style for the given finger's color.
func StyleForFinger(f level.Finger) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(fingerColors[f]))
}

// Styles for typing feedback.
var (
	CorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorCorrect))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError))
	CursorStyle  = lipgloss.NewStyle().Background(lipgloss.Color(ColorCursor))
	LockedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorLocked))
	DimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDimText))
)
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/ui/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/ui/colors.go internal/ui/colors_test.go
git commit -m "feat: add finger-zone color scheme and typing feedback styles"
```

---

### Task 9: Virtual Keyboard Widget

**Files:**
- Create: `internal/keyboard/keyboard.go`
- Create: `internal/keyboard/keyboard_test.go`

- [ ] **Step 1: Write tests for keyboard rendering**

Create `internal/keyboard/keyboard_test.go`:

```go
package keyboard

import (
	"strings"
	"testing"
)

func TestKeyboardRendersAllRows(t *testing.T) {
	unlocked := map[rune]bool{'f': true, 'j': true}
	kb := New(unlocked, 'f')
	view := kb.View()

	// Should contain at least 4 rows of keys
	lines := strings.Split(view, "\n")
	nonEmpty := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonEmpty++
		}
	}
	if nonEmpty < 4 {
		t.Errorf("keyboard has %d non-empty lines, want at least 4", nonEmpty)
	}
}

func TestKeyboardHighlightsNextKey(t *testing.T) {
	unlocked := map[rune]bool{'f': true, 'j': true}
	kb := New(unlocked, 'f')
	view := kb.View()

	// The view should contain 'f' somewhere (it's the next key)
	if !strings.Contains(view, "f") {
		t.Error("keyboard should display the 'f' key")
	}
}

func TestKeyboardShowsLegend(t *testing.T) {
	unlocked := map[rune]bool{'f': true, 'j': true}
	kb := New(unlocked, 'f')
	view := kb.View()

	if !strings.Contains(view, "index") || !strings.Contains(view, "pinky") {
		t.Error("keyboard should show finger legend")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/keyboard/ -v
```

Expected: compilation errors.

- [ ] **Step 3: Implement keyboard widget**

Create `internal/keyboard/keyboard.go`:

```go
package keyboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/typer/internal/level"
	"github.com/nmelo/typer/internal/ui"
)

var rows = [][]rune{
	{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '='},
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/'},
}

// Keyboard renders a virtual keyboard with finger-zone coloring.
type Keyboard struct {
	unlocked map[rune]bool
	nextKey  rune
}

// New creates a keyboard widget.
func New(unlocked map[rune]bool, nextKey rune) *Keyboard {
	return &Keyboard{
		unlocked: unlocked,
		nextKey:  nextKey,
	}
}

// Update changes the next expected key.
func (k *Keyboard) Update(nextKey rune) {
	k.nextKey = nextKey
}

// View renders the keyboard as a string.
func (k *Keyboard) View() string {
	var sb strings.Builder

	for rowIdx, row := range rows {
		// Indent for staggered layout
		indent := strings.Repeat(" ", rowIdx+1)
		sb.WriteString(indent)

		for _, key := range row {
			styled := k.renderKey(key)
			sb.WriteString(styled)
		}
		sb.WriteString("\n")
	}

	// Space bar
	sb.WriteString("      ")
	sb.WriteString(k.renderSpaceBar())
	sb.WriteString("\n\n")

	// Legend
	sb.WriteString(k.renderLegend())

	return sb.String()
}

func (k *Keyboard) renderKey(key rune) string {
	label := string(key)
	if !k.unlocked[key] {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorLocked)).
			Render("[" + label + "]")
	}

	finger := level.FingerForKey(key)
	color := ui.ColorForFinger(finger)

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	if key == k.nextKey {
		style = style.Bold(true).Background(lipgloss.Color(color)).
			Foreground(lipgloss.Color("#000000"))
	}

	return style.Render("[" + label + "]")
}

func (k *Keyboard) renderSpaceBar() string {
	label := "       space       "
	if !k.unlocked[' '] {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorLocked)).
			Render("[" + label + "]")
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(ui.ColorThumb))
	if k.nextKey == ' ' {
		style = style.Bold(true).Background(lipgloss.Color(ui.ColorThumb)).
			Foreground(lipgloss.Color("#000000"))
	}
	return style.Render("[" + label + "]")
}

func (k *Keyboard) renderLegend() string {
	entries := []struct {
		label string
		color string
	}{
		{"pinky", ui.ColorPinky},
		{"ring", ui.ColorRing},
		{"mid", ui.ColorMiddle},
		{"index", ui.ColorIndex},
		{"thumb", ui.ColorThumb},
	}

	var parts []string
	for _, e := range entries {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(e.color))
		parts = append(parts, style.Render(e.label))
	}

	return "  " + strings.Join(parts, "  ")
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/keyboard/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/keyboard/
git commit -m "feat: add virtual keyboard widget with finger-zone coloring and next-key highlight"
```

---

### Task 10: Top Bar and Prompt UI Components

**Files:**
- Create: `internal/ui/topbar.go`
- Create: `internal/ui/prompt.go`
- Create: `internal/ui/topbar_test.go`
- Create: `internal/ui/prompt_test.go`

- [ ] **Step 1: Write tests for top bar**

Create `internal/ui/topbar_test.go`:

```go
package ui

import (
	"strings"
	"testing"
)

func TestTopBarShowsLevel(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "Home Row: f j") {
		t.Error("top bar should show level name")
	}
}

func TestTopBarShowsWPM(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "34") {
		t.Error("top bar should show WPM")
	}
}

func TestTopBarShowsAccuracy(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "96") {
		t.Error("top bar should show accuracy")
	}
}

func TestTopBarShowsDrillProgress(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "3/10") {
		t.Error("top bar should show drill progress")
	}
}
```

- [ ] **Step 2: Write tests for prompt rendering**

Create `internal/ui/prompt_test.go`:

```go
package ui

import (
	"strings"
	"testing"
)

func TestRenderPromptShowsText(t *testing.T) {
	result := RenderPrompt([]rune("fjfj"), []rune("fj"), 2, false)
	// Should contain both the prompt and input characters
	if !strings.Contains(result, "f") {
		t.Error("prompt should contain characters")
	}
}

func TestRenderPromptEmpty(t *testing.T) {
	result := RenderPrompt([]rune("fj"), []rune{}, 0, false)
	if result == "" {
		t.Error("prompt should render even with no input")
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
go test ./internal/ui/ -v
```

Expected: compilation errors.

- [ ] **Step 4: Implement top bar**

Create `internal/ui/topbar.go`:

```go
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var topBarStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 1)

// RenderTopBar renders the status bar with level, WPM, accuracy, and drill progress.
func RenderTopBar(levelName string, wpm, accuracy float64, drillNum, totalDrills int) string {
	level := fmt.Sprintf("Level: %s", levelName)
	wpmStr := fmt.Sprintf("WPM: %.0f", wpm)
	accStr := fmt.Sprintf("Accuracy: %.0f%%", accuracy)
	progress := fmt.Sprintf("%d/%d", drillNum, totalDrills)

	return topBarStyle.Render(
		fmt.Sprintf("  %s    %s    %s    %s", level, wpmStr, accStr, progress),
	)
}
```

- [ ] **Step 5: Implement prompt rendering**

Create `internal/ui/prompt.go`:

```go
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderPrompt renders the prompt text and user input with color feedback.
// position is the current cursor position, hasError indicates if the last keypress was wrong.
func RenderPrompt(prompt []rune, input []rune, position int, hasError bool) string {
	var sb strings.Builder

	// Render the prompt line (what to type)
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDimText))
	sb.WriteString(promptStyle.Render(string(prompt)))
	sb.WriteString("\n")

	// Render the input line with color feedback
	for i := 0; i < len(prompt); i++ {
		ch := string(prompt[i])
		if i < position {
			// Already typed correctly
			sb.WriteString(CorrectStyle.Render(ch))
		} else if i == position {
			if hasError {
				sb.WriteString(ErrorStyle.Render(ch))
			} else {
				sb.WriteString(CursorStyle.Render(ch))
			}
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}
```

- [ ] **Step 6: Run tests to verify they pass**

```bash
go test ./internal/ui/ -v
```

Expected: all tests PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/ui/topbar.go internal/ui/topbar_test.go internal/ui/prompt.go internal/ui/prompt_test.go
git commit -m "feat: add top bar and prompt UI components"
```

---

### Task 11: App Model - Bubble Tea Integration

**Files:**
- Create: `internal/app/app.go`
- Modify: `cmd/typer/main.go`

This task ties all components together into the Bubble Tea Model/Update/View cycle. This is primarily integration code; the logic is already tested in individual packages.

- [ ] **Step 1: Implement the app model**

Create `internal/app/app.go`:

```go
package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/typer/internal/drill"
	"github.com/nmelo/typer/internal/keyboard"
	"github.com/nmelo/typer/internal/level"
	"github.com/nmelo/typer/internal/progress"
	"github.com/nmelo/typer/internal/ui"
)

const drillsPerPhase = 10
const drillLength = 25

type state int

const (
	stateDrilling state = iota
	stateDrillComplete
	statePhaseComplete
	stateLevelComplete
)

// Model is the top-level Bubble Tea model.
type Model struct {
	progress     *progress.Progress
	progressPath string
	words        []string

	currentLevel level.Level
	currentPhase string
	drillNum     int // 1-indexed, current drill within phase
	drillState   *drill.DrillState
	kb           *keyboard.Keyboard
	state        state

	// Aggregate stats for current phase
	phaseKeystrokes int
	phaseErrors     int

	// Message to show between drills
	message string
}

// New creates a new app model.
func New(prog *progress.Progress, progressPath string, words []string) Model {
	lvl := level.Get(prog.CurrentLevel)
	unlocked := unlockedKeySet(prog.CurrentLevel)

	m := Model{
		progress:     prog,
		progressPath: progressPath,
		words:        words,
		currentLevel: lvl,
		currentPhase: prog.CurrentPhase,
		drillNum:     1,
		state:        stateDrilling,
	}

	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.drillState.Prompt[0])

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.saveProgress()
			return m, tea.Quit

		case tea.KeyBackspace:
			if m.state == stateDrilling {
				m.drillState.HandleBackspace()
				m.updateKeyboard()
			}

		case tea.KeyEnter:
			switch m.state {
			case stateDrillComplete:
				m.advanceDrill()
			case statePhaseComplete:
				m.advancePhase()
			case stateLevelComplete:
				m.advanceLevel()
			}

		case tea.KeyRunes:
			if m.state == stateDrilling && len(msg.Runes) > 0 {
				m.drillState.HandleKey(msg.Runes[0])
				m.updateKeyboard()

				if m.drillState.Complete {
					m.onDrillComplete()
				}
			}

		case tea.KeySpace:
			if m.state == stateDrilling {
				m.drillState.HandleKey(' ')
				m.updateKeyboard()

				if m.drillState.Complete {
					m.onDrillComplete()
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	// Top bar
	wpm := m.drillState.WPM()
	acc := m.drillState.Accuracy()
	sb.WriteString(ui.RenderTopBar(m.currentLevel.Name, wpm, acc, m.drillNum, drillsPerPhase))
	sb.WriteString("\n\n")

	// Prompt area
	sb.WriteString(ui.RenderPrompt(
		m.drillState.Prompt,
		nil,
		m.drillState.Position,
		m.drillState.HasError,
	))
	sb.WriteString("\n\n")

	// Message area (between drills)
	if m.message != "" {
		msgStyle := lipgloss.NewStyle().Bold(true)
		sb.WriteString(msgStyle.Render(m.message))
		sb.WriteString("\n\n")
	}

	// Keyboard
	sb.WriteString(m.kb.View())
	sb.WriteString("\n")

	// Help
	sb.WriteString(ui.DimStyle.Render("  esc/ctrl+c to quit"))

	return sb.String()
}

func (m *Model) generateDrill(unlocked map[rune]bool) *drill.DrillState {
	var prompt []rune

	switch m.currentPhase {
	case progress.PhaseChars:
		prompt = drill.GenerateCharDrill(
			level.UnlockedKeys(m.currentLevel.Number),
			m.currentLevel.NewKeys,
			drillLength,
		)
	case progress.PhaseWords:
		filtered := drill.FilterWords(m.words, unlocked)
		words := drill.GenerateWordDrill(filtered, 5)
		prompt = []rune(strings.Join(words, " "))
	case progress.PhaseCode:
		// Code drills: use character drills as fallback for now
		prompt = drill.GenerateCharDrill(
			level.UnlockedKeys(m.currentLevel.Number),
			m.currentLevel.NewKeys,
			drillLength,
		)
	}

	return drill.NewDrillState(prompt)
}

func (m *Model) onDrillComplete() {
	m.phaseKeystrokes += m.drillState.TotalKeystrokes
	m.phaseErrors += m.drillState.Errors

	// Record session
	session := progress.Session{
		Date:     time.Now(),
		Duration: int(time.Since(m.drillState.StartTime).Seconds()),
		Level:    m.currentLevel.Number,
		Phase:    m.currentPhase,
		Accuracy: m.drillState.Accuracy(),
		WPM:      m.drillState.WPM(),
		KeyStats: convertKeyStats(m.drillState.KeyStats()),
	}
	m.progress.RecordSession(session)

	if m.drillNum >= drillsPerPhase {
		phaseAcc := phaseAccuracy(m.phaseKeystrokes, m.phaseErrors)
		if level.PhasePassed(phaseAcc) {
			m.progress.CompletePhase(m.currentLevel.Number, m.currentPhase, phaseAcc, m.drillState.WPM())
			next := level.NextPhase(m.currentPhase, m.currentLevel)
			if next == "" {
				m.state = stateLevelComplete
				m.message = fmt.Sprintf("Level %d complete! Accuracy: %.1f%%. Press Enter to continue.",
					m.currentLevel.Number, phaseAcc)
			} else {
				m.state = statePhaseComplete
				m.message = fmt.Sprintf("Phase complete! Accuracy: %.1f%%. Press Enter for %s drills.",
					phaseAcc, next)
			}
		} else {
			m.state = statePhaseComplete
			m.message = fmt.Sprintf("Phase accuracy: %.1f%% (need 95%%). Press Enter to retry.",
				phaseAcc)
		}
	} else {
		m.state = stateDrillComplete
		m.message = fmt.Sprintf("Drill %d/%d done. Accuracy: %.1f%%. Press Enter to continue.",
			m.drillNum, drillsPerPhase, m.drillState.Accuracy())
	}

	m.saveProgress()
}

func (m *Model) advanceDrill() {
	m.drillNum++
	unlocked := unlockedKeySet(m.currentLevel.Number)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.drillState.Prompt[0])
	m.state = stateDrilling
	m.message = ""
}

func (m *Model) advancePhase() {
	phaseAcc := phaseAccuracy(m.phaseKeystrokes, m.phaseErrors)
	if level.PhasePassed(phaseAcc) {
		// Phase passed: move to next phase
		next := level.NextPhase(m.currentPhase, m.currentLevel)
		m.currentPhase = next
	}
	// If phase failed, m.currentPhase stays the same (retry)

	// Reset phase counters for fresh attempt
	m.phaseKeystrokes = 0
	m.phaseErrors = 0
	m.drillNum = 1
	m.progress.CurrentPhase = m.currentPhase

	unlocked := unlockedKeySet(m.currentLevel.Number)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.drillState.Prompt[0])
	m.state = stateDrilling
	m.message = ""
}

func (m *Model) advanceLevel() {
	if m.currentLevel.Number < len(level.All()) {
		m.progress.CurrentLevel = m.currentLevel.Number + 1
		m.progress.CurrentPhase = progress.PhaseChars
		m.currentLevel = level.Get(m.progress.CurrentLevel)
		m.currentPhase = progress.PhaseChars
	}

	m.phaseKeystrokes = 0
	m.phaseErrors = 0
	m.drillNum = 1

	unlocked := unlockedKeySet(m.currentLevel.Number)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.drillState.Prompt[0])
	m.state = stateDrilling
	m.message = ""
}

func (m *Model) updateKeyboard() {
	if !m.drillState.Complete && m.drillState.Position < len(m.drillState.Prompt) {
		m.kb.Update(m.drillState.Prompt[m.drillState.Position])
	}
}

func (m *Model) saveProgress() {
	_ = progress.Save(m.progress, m.progressPath)
}

func unlockedKeySet(levelNum int) map[rune]bool {
	keys := level.UnlockedKeys(levelNum)
	set := make(map[rune]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	// Space is implicitly available from level 5+
	if levelNum >= 5 {
		set[' '] = true
	}
	return set
}

func phaseAccuracy(keystrokes, errors int) float64 {
	if keystrokes == 0 {
		return 100.0
	}
	correct := keystrokes - errors
	return float64(correct) / float64(keystrokes) * 100
}

func convertKeyStats(drillStats map[rune]drill.KeyStat) map[string]progress.KeyStat {
	result := make(map[string]progress.KeyStat)
	for k, v := range drillStats {
		result[string(k)] = progress.KeyStat{
			Correct:   v.Correct,
			Incorrect: v.Incorrect,
		}
	}
	return result
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./internal/app/
```

Expected: compiles without errors.

- [ ] **Step 3: Update main.go**

Replace `cmd/typer/main.go`:

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nmelo/typer/internal/app"
	"github.com/nmelo/typer/internal/drill"
	"github.com/nmelo/typer/internal/progress"
)

func main() {
	path := progress.DefaultPath()
	prog, err := progress.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading progress: %v\n", err)
		os.Exit(1)
	}

	words := drill.LoadWordList()
	model := app.New(prog, path, words)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Verify full build**

```bash
go build ./cmd/typer
```

Expected: produces `typer` binary.

- [ ] **Step 5: Commit**

```bash
git add internal/app/ cmd/typer/main.go
git commit -m "feat: integrate all components into Bubble Tea app model with main entry point"
```

---

### Task 12: Manual Smoke Test and Polish

**Files:**
- Possibly modify: `internal/app/app.go`, `internal/keyboard/keyboard.go`, `internal/ui/prompt.go`

- [ ] **Step 1: Run the app**

```bash
./typer
```

Verify:
- Top bar shows Level 1 info, WPM 0, Accuracy 100%, 1/10
- Prompt shows a random sequence of 'f' and 'j' characters
- Virtual keyboard shows with 'f' and 'j' colored, all other keys dimmed
- Next key is highlighted on keyboard
- Typing correct keys advances the cursor (green)
- Typing wrong key shows error (red), blocks progress
- Backspace clears error
- After completing a drill, message shows accuracy and "Press Enter"
- After 10 drills with 95%+ accuracy, phase/level advances
- Esc quits cleanly

- [ ] **Step 2: Fix any visual issues found during smoke test**

Adjust spacing, alignment, or rendering based on what looks off in the terminal.

- [ ] **Step 3: Run all tests**

```bash
go test ./... -v
```

Expected: all tests PASS.

- [ ] **Step 4: Commit any fixes**

```bash
git add -A
git commit -m "fix: polish UI rendering from smoke test"
```

---

### Task 13: Code Snippets (Stretch)

**Files:**
- Create: `data/code_snippets.go`
- Modify: `internal/drill/content.go`
- Modify: `internal/drill/content_test.go`
- Modify: `internal/app/app.go`

Code drill content for levels 10+. This can use the character drill fallback for v1 and be backfilled later, so it's lower priority than getting the core loop working.

- [ ] **Step 1: Write tests for code snippet selection**

Add to `internal/drill/content_test.go`:

```go
func TestFilterCodeSnippets(t *testing.T) {
	snippets := LoadCodeSnippets()
	if len(snippets) == 0 {
		t.Fatal("should have at least some code snippets")
	}

	// Level 10 unlocks q,p + all previous
	keys := make(map[rune]bool)
	for _, r := range []rune{'f', 'j', 'd', 'k', 's', 'l', 'a', ';', 'g', 'h', 't', 'y', 'r', 'u', 'e', 'i', 'w', 'o', 'q', 'p'} {
		keys[r] = true
	}

	filtered := FilterCodeSnippets(snippets, keys)
	for _, s := range filtered {
		for _, ch := range s.Code {
			if ch != ' ' && ch != '\n' && !keys[ch] {
				t.Errorf("snippet %q contains locked char %c", s.Code[:20], ch)
			}
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/drill/ -v -run CodeSnippet
```

Expected: compilation error.

- [ ] **Step 3: Create code snippets data**

Create `data/code_snippets.go`:

```go
package data

// CodeSnippet is a curated code snippet with its required character set.
type CodeSnippet struct {
	Code         string
	RequiredKeys string // all unique chars needed (excluding space)
	Language     string
}

// Snippets contains curated code snippets for typing practice.
var Snippets = []CodeSnippet{
	{
		Code:         "if true { return }",
		RequiredKeys: "iftrue{rn}",
		Language:     "go",
	},
	{
		Code:         "for i := 0; i < 10; i++ {",
		RequiredKeys: "fori:=0;<1+{",
		Language:     "go",
	},
	{
		Code:         "def split(words):",
		RequiredKeys: "defsplit(wor):",
		Language:     "python",
	},
	{
		Code:         "result := 0",
		RequiredKeys: "result:=0",
		Language:     "go",
	},
	{
		Code:         "err != nil",
		RequiredKeys: "er!=nil",
		Language:     "go",
	},
	{
		Code:         "type status struct {}",
		RequiredKeys: "typesauc{}",
		Language:     "go",
	},
	{
		Code:         "for k, v := range items {",
		RequiredKeys: "fork,v:=angeitms{",
		Language:     "go",
	},
	{
		Code:         "import os",
		RequiredKeys: "importas",
		Language:     "python",
	},
	{
		Code:         "while true:",
		RequiredKeys: "whlietru:",
		Language:     "python",
	},
	{
		Code:         "print(sorted(list))",
		RequiredKeys: "pint(sored(ls))",
		Language:     "python",
	},
}
```

- [ ] **Step 4: Implement code snippet filtering**

Add to `internal/drill/content.go`:

```go
import "github.com/nmelo/typer/data"

// FilterCodeSnippets returns snippets whose characters are all in the allowed key set.
func FilterCodeSnippets(snippets []data.CodeSnippet, allowedKeys map[rune]bool) []data.CodeSnippet {
	var filtered []data.CodeSnippet
	for _, s := range snippets {
		ok := true
		for _, ch := range s.Code {
			if ch != ' ' && ch != '\n' && !allowedKeys[ch] {
				ok = false
				break
			}
		}
		if ok {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// LoadCodeSnippets returns all curated code snippets.
func LoadCodeSnippets() []data.CodeSnippet {
	return data.Snippets
}
```

- [ ] **Step 5: Update app model to use code snippets**

In `internal/app/app.go`, update the `generateDrill` method's `PhaseCode` case:

```go
case progress.PhaseCode:
    allSnippets := drill.LoadCodeSnippets()
    filtered := drill.FilterCodeSnippets(allSnippets, unlocked)
    if len(filtered) > 0 {
        snippet := filtered[rand.Intn(len(filtered))]
        prompt = []rune(snippet.Code)
    } else {
        // Fallback to char drill if no snippets match
        prompt = drill.GenerateCharDrill(
            level.UnlockedKeys(m.currentLevel.Number),
            m.currentLevel.NewKeys,
            drillLength,
        )
    }
```

- [ ] **Step 6: Run all tests**

```bash
go test ./... -v
```

Expected: all tests PASS.

- [ ] **Step 7: Commit**

```bash
git add data/code_snippets.go internal/drill/content.go internal/drill/content_test.go internal/app/app.go
git commit -m "feat: add curated code snippets for code drill phases"
```
