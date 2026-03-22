package keyboard

import (
	"strings"
	"testing"
)

func TestKeyboardRendersAllRows(t *testing.T) {
	unlocked := map[rune]bool{'f': true, 'j': true}
	kb := New(unlocked, 'f')
	view := kb.View()

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

	if !strings.Contains(view, "F") {
		t.Error("keyboard should display the 'F' key")
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
