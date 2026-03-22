package ui

import (
	"strings"
	"testing"
)

func TestRenderPromptShowsText(t *testing.T) {
	result := RenderPrompt([]rune("fjfj"), 2, false)
	if !strings.Contains(result, "f") {
		t.Error("prompt should contain characters")
	}
}

func TestRenderPromptEmpty(t *testing.T) {
	result := RenderPrompt([]rune("fj"), 0, false)
	if result == "" {
		t.Error("prompt should render even with no input")
	}
}
