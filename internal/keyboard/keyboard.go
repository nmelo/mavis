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

type Keyboard struct {
	unlocked map[rune]bool
	nextKey  rune
}

func New(unlocked map[rune]bool, nextKey rune) *Keyboard {
	return &Keyboard{unlocked: unlocked, nextKey: nextKey}
}

func (k *Keyboard) Update(nextKey rune) {
	k.nextKey = nextKey
}

func (k *Keyboard) View() string {
	var sb strings.Builder

	for rowIdx, row := range rows {
		indent := strings.Repeat(" ", rowIdx+1)
		sb.WriteString(indent)

		for _, key := range row {
			sb.WriteString(k.renderKey(key))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("      ")
	sb.WriteString(k.renderSpaceBar())
	sb.WriteString("\n\n")

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
