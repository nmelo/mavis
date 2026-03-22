package keyboard

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/mavis/internal/level"
	"github.com/nmelo/mavis/internal/ui"
)

var rows = [][]rune{
	{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '='},
	{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},
	{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';'},
	{'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/'},
}

var (
	keyBase = lipgloss.NewStyle().
		Width(3).
		Align(lipgloss.Center).
		Padding(0, 1).
		MarginRight(1)

	lockedKey = keyBase.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ui.ColorLocked)).
			Foreground(lipgloss.Color(ui.ColorLocked))

	spaceBase = lipgloss.NewStyle().
			Width(30).
			Align(lipgloss.Center).
			Padding(0, 1)
)

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
	var rowViews []string

	for rowIdx, row := range rows {
		var keys []string
		for _, key := range row {
			keys = append(keys, k.renderKey(key))
		}
		indent := strings.Repeat("  ", rowIdx)
		rowView := indent + lipgloss.JoinHorizontal(lipgloss.Top, keys...)
		rowViews = append(rowViews, rowView)
	}

	kb := strings.Join(rowViews, "\n")

	space := k.renderSpaceBar()
	legend := k.renderLegend()

	return kb + "\n" + space + "\n\n" + legend
}

func (k *Keyboard) renderKey(key rune) string {
	label := string(unicode.ToUpper(key))

	if !k.unlocked[key] {
		return lockedKey.Render(label)
	}

	finger := level.FingerForKey(key)
	color := lipgloss.Color(ui.ColorForFinger(finger))

	style := keyBase.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Foreground(color)

	if key == k.nextKey {
		style = style.
			Bold(true).
			Background(color).
			Foreground(lipgloss.Color("#000000"))
	}

	return style.Render(label)
}

func (k *Keyboard) renderSpaceBar() string {
	label := "SPACE"

	if !k.unlocked[' '] {
		return spaceBase.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ui.ColorLocked)).
			Foreground(lipgloss.Color(ui.ColorLocked)).
			Render(label)
	}

	color := lipgloss.Color(ui.ColorThumb)
	style := spaceBase.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Foreground(color)

	if k.nextKey == ' ' {
		style = style.
			Bold(true).
			Background(color).
			Foreground(lipgloss.Color("#000000"))
	}

	return style.Render(label)
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
		dot := lipgloss.NewStyle().Foreground(lipgloss.Color(e.color)).Render("●")
		lbl := lipgloss.NewStyle().Foreground(lipgloss.Color(e.color)).Render(e.label)
		parts = append(parts, dot+" "+lbl)
	}

	return "  " + strings.Join(parts, "    ")
}
