package level

import "testing"

func TestPhasePassesAt95(t *testing.T) {
	if !PhasePassed(95.0) {
		t.Error("95.0% should pass")
	}
	if !PhasePassed(100.0) {
		t.Error("100.0% should pass")
	}
}

func TestPhaseFailsBelow95(t *testing.T) {
	if PhasePassed(94.9) {
		t.Error("94.9% should not pass")
	}
	if PhasePassed(0) {
		t.Error("0% should not pass")
	}
}

func TestNextPhaseAfterChars(t *testing.T) {
	lvl := All()[2] // Level 3: has word drills, no code drills
	next := NextPhase("chars", lvl)
	if next != "words" {
		t.Errorf("next phase after chars = %q, want %q", next, "words")
	}
}

func TestNextPhaseAfterCharsNoWords(t *testing.T) {
	lvl := All()[0] // Level 1: no word drills, no code drills
	next := NextPhase("chars", lvl)
	if next != "" {
		t.Errorf("next phase after chars (no words) = %q, want empty (level complete)", next)
	}
}

func TestNextPhaseAfterWordsWithCode(t *testing.T) {
	lvl := All()[9] // Level 10: has code drills
	next := NextPhase("words", lvl)
	if next != "code" {
		t.Errorf("next phase after words = %q, want %q", next, "code")
	}
}

func TestNextPhaseAfterWordsNoCode(t *testing.T) {
	lvl := All()[2] // Level 3: no code drills
	next := NextPhase("words", lvl)
	if next != "" {
		t.Errorf("next phase after words (no code) = %q, want empty (level complete)", next)
	}
}

func TestNextPhaseAfterCode(t *testing.T) {
	lvl := All()[9] // Level 10
	next := NextPhase("code", lvl)
	if next != "" {
		t.Errorf("next phase after code = %q, want empty (level complete)", next)
	}
}

func TestPhasesForLevel(t *testing.T) {
	phases := PhasesFor(All()[0])
	if len(phases) != 1 || phases[0] != "chars" {
		t.Errorf("level 1 phases = %v, want [chars]", phases)
	}

	phases = PhasesFor(All()[2])
	if len(phases) != 2 {
		t.Errorf("level 3 phases = %v, want [chars, words]", phases)
	}

	phases = PhasesFor(All()[9])
	if len(phases) != 3 {
		t.Errorf("level 10 phases = %v, want [chars, words, code]", phases)
	}
}
