package progress

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	PhaseChars = "chars"
	PhaseWords = "words"
	PhaseCode  = "code"
)

type KeyStat struct {
	Correct   int `json:"correct"`
	Incorrect int `json:"incorrect"`
}

type Session struct {
	Date     time.Time          `json:"date"`
	Duration int                `json:"duration_sec"`
	Level    int                `json:"level"`
	Phase    string             `json:"phase"`
	Accuracy float64            `json:"accuracy"`
	WPM      float64            `json:"wpm"`
	KeyStats map[string]KeyStat `json:"key_stats"`
}

type PhaseResult struct {
	Completed    bool    `json:"completed"`
	BestAccuracy float64 `json:"best_accuracy"`
	BestWPM      float64 `json:"best_wpm"`
}

type LevelProgress struct {
	CharDrills *PhaseResult `json:"char_drills"`
	WordDrills *PhaseResult `json:"word_drills"`
	CodeDrills *PhaseResult `json:"code_drills"`
}

type Progress struct {
	CurrentLevel int                      `json:"current_level"`
	CurrentPhase string                   `json:"current_phase"`
	Levels       map[string]LevelProgress `json:"levels"`
	Sessions     []Session                `json:"sessions"`
}

func New() *Progress {
	return &Progress{
		CurrentLevel: 1,
		CurrentPhase: PhaseChars,
		Levels:       make(map[string]LevelProgress),
	}
}

func Load(path string) (*Progress, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return nil, err
	}
	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.Levels == nil {
		p.Levels = make(map[string]LevelProgress)
	}
	return &p, nil
}

func Save(p *Progress, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (p *Progress) RecordSession(s Session) {
	p.Sessions = append(p.Sessions, s)
}

func (p *Progress) CompletePhase(levelNum int, phase string, accuracy, wpm float64) {
	key := strconv.Itoa(levelNum)

	lvl := p.Levels[key]
	result := &PhaseResult{
		Completed:    true,
		BestAccuracy: accuracy,
		BestWPM:      wpm,
	}

	switch phase {
	case PhaseChars:
		lvl.CharDrills = result
	case PhaseWords:
		lvl.WordDrills = result
	case PhaseCode:
		lvl.CodeDrills = result
	}

	p.Levels[key] = lvl
}

func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "typer", "progress.json")
}
