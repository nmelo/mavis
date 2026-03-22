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
			lifetime: 10 + rand.Intn(6),
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

	grid := map[[2]int]particle{}
	for _, p := range c.particles {
		grid[[2]int{centerX + p.x, centerY + p.y}] = p
	}

	var sb strings.Builder

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

	sb.WriteString(c.msgStyle.Render(c.message))

	if !c.active() {
		sb.WriteString(fmt.Sprintf("\n\n%s", lipgloss.NewStyle().
			Foreground(lipgloss.Color(ui.ColorDimText)).
			Render("Press Enter to continue.")))
	}

	return sb.String()
}
