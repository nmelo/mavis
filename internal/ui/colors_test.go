package ui

import (
	"testing"

	"github.com/nmelo/mavis/internal/level"
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
