package progress

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewProgress(t *testing.T) {
	p := New()
	if p.CurrentLevel != 1 {
		t.Errorf("initial level = %d, want 1", p.CurrentLevel)
	}
	if p.CurrentPhase != PhaseChars {
		t.Errorf("initial phase = %q, want %q", p.CurrentPhase, PhaseChars)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "progress.json")

	p := New()
	p.CurrentLevel = 3
	p.CurrentPhase = PhaseWords

	err := Save(p, path)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.CurrentLevel != 3 {
		t.Errorf("loaded level = %d, want 3", loaded.CurrentLevel)
	}
	if loaded.CurrentPhase != PhaseWords {
		t.Errorf("loaded phase = %q, want %q", loaded.CurrentPhase, PhaseWords)
	}
}

func TestLoadNonExistentReturnsNew(t *testing.T) {
	p, err := Load("/nonexistent/path/progress.json")
	if err != nil {
		t.Fatalf("load should not error for missing file: %v", err)
	}
	if p.CurrentLevel != 1 {
		t.Errorf("should return fresh progress, got level %d", p.CurrentLevel)
	}
}

func TestRecordSession(t *testing.T) {
	p := New()
	session := Session{
		Date:     time.Now(),
		Duration: 300,
		Level:    1,
		Phase:    PhaseChars,
		Accuracy: 96.5,
		WPM:      22.0,
		KeyStats: map[string]KeyStat{
			"f": {Correct: 50, Incorrect: 2},
		},
	}

	p.RecordSession(session)

	if len(p.Sessions) != 1 {
		t.Fatalf("sessions count = %d, want 1", len(p.Sessions))
	}
	if p.Sessions[0].Accuracy != 96.5 {
		t.Errorf("session accuracy = %.1f, want 96.5", p.Sessions[0].Accuracy)
	}
}

func TestCompletePhaseSetsLevelBest(t *testing.T) {
	p := New()
	p.CompletePhase(1, PhaseChars, 97.3, 25.0)

	lvl, ok := p.Levels["1"]
	if !ok {
		t.Fatal("level 1 not recorded")
	}
	if lvl.CharDrills == nil {
		t.Fatal("char drills not recorded")
	}
	if !lvl.CharDrills.Completed {
		t.Error("char drills should be completed")
	}
	if lvl.CharDrills.BestAccuracy != 97.3 {
		t.Errorf("best accuracy = %.1f, want 97.3", lvl.CharDrills.BestAccuracy)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "progress.json")

	p := New()
	err := Save(p, path)
	if err != nil {
		t.Fatalf("save should create parent dirs: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file should exist after save")
	}
}
