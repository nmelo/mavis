# Celebration System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add tiered visual celebrations (pass, perfect, level-up) with particle animations on drill/phase/level completion.

**Architecture:** New `celebrate.go` file in `internal/app/` handles particle system and celebration rendering. A new `stateCelebrating` state in the existing state machine triggers tick-based animation. The `onDrillComplete` method determines which tier fires based on precedence rules (level-up > phase pass > perfect drill).

**Tech Stack:** Go, Bubble Tea (tea.Tick for animation), Lipgloss (styling)

**Spec:** `docs/superpowers/specs/2026-03-21-celebration-system-design.md`

---

### Task 1: Celebration Data Types and Particle System

**Files:**
- Create: `internal/app/celebrate.go`
- Create: `internal/app/celebrate_test.go`

- [ ] **Step 1: Write tests for particle system**

Create `internal/app/celebrate_test.go`:

```go
package app

import "testing"

func TestSpawnPerfectCelebration(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	if len(c.particles) != 12 {
		t.Errorf("perfect tier particles = %d, want 12", len(c.particles))
	}
	if c.nextState != stateDrillComplete {
		t.Error("nextState should be stateDrillComplete")
	}
	if !c.active() {
		t.Error("celebration should be active after spawn")
	}
}

func TestSpawnLevelUpCelebration(t *testing.T) {
	c := spawnCelebration(tierLevelUp, stateLevelComplete, "LEVEL 3 COMPLETE!")
	if len(c.particles) != 20 {
		t.Errorf("level-up tier particles = %d, want 20", len(c.particles))
	}
}

func TestTickReducesLifetime(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	initialLife := c.particles[0].lifetime
	c.tick()
	if c.particles[0].lifetime != initialLife-1 {
		t.Errorf("lifetime after tick = %d, want %d", c.particles[0].lifetime, initialLife-1)
	}
}

func TestTickRemovesDeadParticles(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	// Set all particles to 1 tick remaining
	for i := range c.particles {
		c.particles[i].lifetime = 1
	}
	c.tick()
	if len(c.particles) != 0 {
		t.Errorf("particles after expiry = %d, want 0", len(c.particles))
	}
	if c.active() {
		t.Error("celebration should be inactive after all particles expire")
	}
}

func TestCelebrationRender(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	result := c.render(80)
	if result == "" {
		t.Error("render should produce output")
	}
}

func TestPassTierNoParticles(t *testing.T) {
	// Pass tier should not be used with spawnCelebration (it's just styled text)
	// but if called, it should produce 0 particles
	c := spawnCelebration(tierPass, statePhaseComplete, "Phase complete!")
	if len(c.particles) != 0 {
		t.Errorf("pass tier particles = %d, want 0", len(c.particles))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/app/ -v
```

Expected: compilation errors (types don't exist yet).

- [ ] **Step 3: Implement celebration system**

Create `internal/app/celebrate.go`:

```go
package app

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/mavis/internal/ui"
)

type tier int

const (
	tierPass    tier = iota
	tierPerfect
	tierLevelUp
)

var sparkleColors = []string{
	ui.ColorPinky,
	ui.ColorRing,
	ui.ColorMiddle,
	ui.ColorIndex,
	ui.ColorThumb,
}

var sparkleChars = []rune{'*', '+', '.'}

type particle struct {
	x, y     int
	char     rune
	colorIdx int
	lifetime int
}

type celebration struct {
	particles []particle
	message   string
	msgStyle  lipgloss.Style
	nextState state
	tickCount int
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func spawnCelebration(t tier, nextState state, message string) *celebration {
	c := &celebration{
		message:   message,
		nextState: nextState,
	}

	switch t {
	case tierPass:
		c.msgStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ui.ColorCorrect))
	case tierPerfect:
		c.msgStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
		c.particles = makeParticles(12, 20, 3)
	case tierLevelUp:
		c.msgStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
		c.particles = makeParticles(20, 30, 4)
	}

	return c
}

func makeParticles(count, spreadX, spreadY int) []particle {
	particles := make([]particle, count)
	for i := range particles {
		particles[i] = particle{
			x:        rand.Intn(spreadX*2+1) - spreadX,
			y:        rand.Intn(spreadY*2+1) - spreadY,
			char:     sparkleChars[rand.Intn(len(sparkleChars))],
			colorIdx: rand.Intn(len(sparkleColors)),
			lifetime: 10 + rand.Intn(6), // 10-15 ticks = 1.0-1.5 seconds
		}
	}
	return particles
}

func (c *celebration) active() bool {
	return len(c.particles) > 0
}

func (c *celebration) tick() {
	c.tickCount++
	alive := c.particles[:0]
	for i := range c.particles {
		c.particles[i].lifetime--
		c.particles[i].colorIdx = (c.particles[i].colorIdx + 1) % len(sparkleColors)
		if c.particles[i].lifetime > 0 {
			alive = append(alive, c.particles[i])
		}
	}
	c.particles = alive
}

func (c *celebration) render(width int) string {
	centerX := width / 2
	centerY := 0

	// Build a sparse grid of particles
	grid := map[[2]int]particle{}
	for _, p := range c.particles {
		grid[[2]int{centerX + p.x, centerY + p.y}] = p
	}

	var sb strings.Builder

	// Render particle rows above the message
	for row := centerY - 4; row <= centerY+4; row++ {
		var line strings.Builder
		hasContent := false
		for col := centerX - 30; col <= centerX+30; col++ {
			if p, ok := grid[[2]int{col, row}]; ok {
				color := sparkleColors[p.colorIdx]
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				line.WriteString(style.Render(string(p.char)))
				hasContent = true
			} else {
				line.WriteString(" ")
			}
		}
		if hasContent {
			sb.WriteString(line.String())
		}
		sb.WriteString("\n")
	}

	// Render the celebration message
	sb.WriteString(c.msgStyle.Render(c.message))

	if !c.active() {
		sb.WriteString(fmt.Sprintf("\n\n%s", lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorDimText)).
			Render("Press Enter to continue.")))
	}

	return sb.String()
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/app/ -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/app/celebrate.go internal/app/celebrate_test.go
git commit -m "feat: add celebration particle system with tiered animations"
```

---

### Task 2: Integrate Celebrations into App State Machine

**Files:**
- Modify: `internal/app/app.go`

- [ ] **Step 1: Add stateCelebrating to state enum**

In `internal/app/app.go`, change:

```go
const (
	stateDrilling state = iota
	stateDrillComplete
	statePhaseComplete
	stateLevelComplete
)
```

to:

```go
const (
	stateDrilling state = iota
	stateDrillComplete
	statePhaseComplete
	stateLevelComplete
	stateCelebrating
)
```

- [ ] **Step 2: Add celebration field to Model**

Add to the Model struct:

```go
	celebration *celebration
```

- [ ] **Step 3: Handle tickMsg in Update**

Add a new case in the `Update` method, before the `tea.KeyMsg` case:

```go
	case tickMsg:
		if m.state == stateCelebrating && m.celebration != nil {
			m.celebration.tick()
			if !m.celebration.active() {
				return m, nil
			}
			return m, tickCmd()
		}
		return m, nil
```

- [ ] **Step 4: Allow Enter to skip celebration**

In the `tea.KeyEnter` handler, add a case for `stateCelebrating`:

```go
		case tea.KeyEnter:
			switch m.state {
			case stateCelebrating:
				if m.celebration != nil {
					m.state = m.celebration.nextState
					m.celebration = nil
				}
			case stateDrillComplete:
				m.advanceDrill()
			// ... rest unchanged
```

- [ ] **Step 5: Update onDrillComplete with celebration triggers**

Replace the `onDrillComplete` method. The key changes:
- After determining the next state (drillComplete/phaseComplete/levelComplete), check celebration tiers in precedence order
- Level-up: spawn tierLevelUp, set stateCelebrating
- Phase pass (>= 95%): style message green with checkmark (no animation)
- Perfect drill (100%): spawn tierPerfect, set stateCelebrating
- Return a `tickCmd` when entering stateCelebrating

Change the method signature to return a `tea.Cmd`:

```go
func (m *Model) onDrillComplete() tea.Cmd {
	m.phaseKeystrokes += m.drillState.TotalKeystrokes
	m.phaseErrors += m.drillState.Errors

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

	drillAcc := m.drillState.Accuracy()

	if m.drillNum >= drillsPerPhase {
		phaseAcc := phaseAccuracy(m.phaseKeystrokes, m.phaseErrors)
		if level.PhasePassed(phaseAcc) {
			m.progress.CompletePhase(m.currentLevel.Number, m.currentPhase, phaseAcc, m.drillState.WPM())
			next := level.NextPhase(m.currentPhase, m.currentLevel)
			if next == "" {
				// Level complete - highest tier
				msg := fmt.Sprintf("Level %d complete! Accuracy: %.1f%%.", m.currentLevel.Number, phaseAcc)
				m.celebration = spawnCelebration(tierLevelUp, stateLevelComplete, msg)
				m.state = stateCelebrating
				m.message = ""
				m.saveProgress()
				return tickCmd()
			}
			// Phase passed
			msg := fmt.Sprintf("Phase complete! Accuracy: %.1f%%. Press Enter for %s drills.", phaseAcc, next)
			m.state = statePhaseComplete
			m.message = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ui.ColorCorrect)).Render("+ " + msg)
		} else {
			m.state = statePhaseComplete
			m.message = fmt.Sprintf("Phase accuracy: %.1f%% (need 95%%). Press Enter to retry.", phaseAcc)
		}
	} else {
		// Individual drill complete
		if drillAcc == 100.0 {
			msg := fmt.Sprintf("PERFECT! Drill %d/%d. Press Enter to continue.", m.drillNum, drillsPerPhase)
			m.celebration = spawnCelebration(tierPerfect, stateDrillComplete, msg)
			m.state = stateCelebrating
			m.message = ""
			m.saveProgress()
			return tickCmd()
		}
		m.state = stateDrillComplete
		m.message = fmt.Sprintf("Drill %d/%d done. Accuracy: %.1f%%. Press Enter to continue.",
			m.drillNum, drillsPerPhase, drillAcc)
	}

	m.saveProgress()
	return nil
}
```

- [ ] **Step 6: Update callers of onDrillComplete to use returned Cmd**

In the `Update` method, everywhere `m.onDrillComplete()` is called, change to return the cmd:

```go
			if m.drillState.Complete {
				cmd := m.onDrillComplete()
				return m, cmd
			}
```

There are two call sites (KeyRunes and KeySpace handlers).

- [ ] **Step 7: Render celebration in View**

In the `View` method, replace the message rendering block:

```go
	if m.state == stateCelebrating && m.celebration != nil {
		sb.WriteString(m.celebration.render(m.width))
		sb.WriteString("\n\n")
	} else if m.message != "" {
		msgStyle := lipgloss.NewStyle().Bold(true)
		sb.WriteString(msgStyle.Render(m.message))
		sb.WriteString("\n\n")
	}
```

- [ ] **Step 8: Verify build and tests**

```bash
go build ./cmd/mavis && go test ./... -v
```

Expected: all tests PASS, binary builds.

- [ ] **Step 9: Commit**

```bash
git add internal/app/app.go
git commit -m "feat: integrate celebration state machine with tiered triggers"
```

---

### Task 3: Manual Test and Polish

**Files:**
- Possibly modify: `internal/app/celebrate.go`, `internal/app/app.go`

- [ ] **Step 1: Build and install**

```bash
go build ./cmd/mavis && go install ./cmd/mavis
```

- [ ] **Step 2: Test in tmux**

Restart mavis in tmux. Test each tier:
- Type a drill with some errors. Verify normal message (no celebration).
- Type a drill with 100% accuracy. Verify sparkle animation plays, then "Press Enter" appears.
- Press Enter during animation to verify skip works.
- Use ctrl+n to skip to level 1, complete all drills with 95%+ to verify level-up celebration.
- Verify 95% phase pass shows green-styled message (no animation).

- [ ] **Step 3: Fix any visual issues**

Adjust particle spread, colors, timing, or message styling based on what looks right.

- [ ] **Step 4: Run all tests**

```bash
go test ./... -v
```

- [ ] **Step 5: Commit any fixes**

```bash
git add -A
git commit -m "fix: polish celebration animations"
```
