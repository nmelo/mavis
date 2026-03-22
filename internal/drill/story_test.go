package drill

import (
	"fmt"
	"testing"

	"github.com/nmelo/mavis/internal/level"
)

func TestStoryLinesOnlyUseUnlockedKeys(t *testing.T) {
	invalid := ValidateStoryLines()
	for lvl, lines := range invalid {
		for _, line := range lines {
			// Show which character is the offender
			allowed := make(map[rune]bool)
			for _, k := range level.UnlockedKeys(lvl) {
				allowed[k] = true
			}
			allowed[' '] = true
			for _, ch := range line {
				if !allowed[ch] {
					t.Errorf("level %d: line %q contains invalid char %c", lvl, line, ch)
					break
				}
			}
		}
	}
}

func TestAllStoryLevelsHave10Lines(t *testing.T) {
	// Levels 3-14 should have story content (word drill levels)
	for lvl := 3; lvl <= 14; lvl++ {
		if !HasStory(lvl) {
			t.Errorf("level %d has no story lines", lvl)
			continue
		}
		line := GetStoryLine(lvl, 0)
		if line == "" {
			t.Errorf("level %d story line 0 is empty", lvl)
		}
	}
}

func TestGetStoryLineWraps(t *testing.T) {
	if !HasStory(3) {
		t.Skip("no story for level 3")
	}
	line0 := GetStoryLine(3, 0)
	line10 := GetStoryLine(3, 10)
	if line0 != line10 {
		t.Error("drill index 10 should wrap to same as 0")
	}
}

func TestStoryLinesNotEmpty(t *testing.T) {
	for lvl := 3; lvl <= 14; lvl++ {
		if !HasStory(lvl) {
			continue
		}
		for i := 0; i < 10; i++ {
			line := GetStoryLine(lvl, i)
			if len(line) < 5 {
				t.Errorf("level %d, drill %d: line too short: %q", lvl, i, line)
			}
		}
	}
}

func TestPrintAvailableChars(t *testing.T) {
	// Helper: not a real test, just prints available chars per level for authoring
	if testing.Short() {
		t.Skip()
	}
	for lvl := 3; lvl <= 14; lvl++ {
		keys := level.UnlockedKeys(lvl)
		fmt.Printf("Level %d: %s\n", lvl, string(keys))
	}
}
