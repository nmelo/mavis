package ui

import (
	"strings"
	"testing"
)

func TestTopBarShowsLevel(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "Home Row: f j") {
		t.Error("top bar should show level name")
	}
}

func TestTopBarShowsWPM(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "34") {
		t.Error("top bar should show WPM")
	}
}

func TestTopBarShowsAccuracy(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "96") {
		t.Error("top bar should show accuracy")
	}
}

func TestTopBarShowsDrillProgress(t *testing.T) {
	bar := RenderTopBar("Home Row: f j", 34.2, 96.5, 3, 10)
	if !strings.Contains(bar, "3/10") {
		t.Error("top bar should show drill progress")
	}
}
