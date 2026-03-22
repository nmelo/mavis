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
