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
	keys := level.UnlockedKeys(5)
	newKeys := level.Get(5).NewKeys

	counts := map[rune]int{}
	for i := 0; i < 100; i++ {
		drill := GenerateCharDrill(keys, newKeys, 100)
		for _, ch := range drill {
			counts[ch]++
		}
	}

	newKeyCount := counts['g'] + counts['h']
	reviewKeyCount := counts['f'] + counts['j'] + counts['d'] + counts['k'] +
		counts['s'] + counts['l'] + counts['a'] + counts[';']

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
	if stats['j'].Incorrect != 1 {
		t.Errorf("j incorrect = %d, want 1", stats['j'].Incorrect)
	}
}
