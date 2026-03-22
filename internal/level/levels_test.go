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
	_ = levels
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
	if levels[0].HasWordDrills {
		t.Error("level 1 should not have word drills")
	}
	if levels[1].HasWordDrills {
		t.Error("level 2 should not have word drills")
	}
	if !levels[2].HasWordDrills {
		t.Error("level 3 should have word drills")
	}
}

func TestLevelHasCodeDrills(t *testing.T) {
	levels := All()
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
