package ui

import (
	"strings"
	"testing"
)

func TestRenderPromptShowsText(t *testing.T) {
	result := RenderPrompt([]rune("fjfj"), []rune("fj"), 2, false)
	if !strings.Contains(result, "f") {
		t.Error("prompt should contain characters")
	}
}

func TestRenderPromptEmpty(t *testing.T) {
	result := RenderPrompt([]rune("fj"), []rune{}, 0, false)
	if result == "" {
		t.Error("prompt should render even with no input")
	}
}
