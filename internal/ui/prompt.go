package ui

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

func RenderPrompt(prompt []rune, position int, hasError bool) string {
	var sb strings.Builder

	// Render prompt line with spacing between characters
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDimText))
	for i, ch := range prompt {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(promptStyle.Render(string(unicode.ToUpper(ch))))
	}
	sb.WriteString("\n")

	// Render input line with same spacing
	for i := 0; i < len(prompt); i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		ch := string(unicode.ToUpper(prompt[i]))
		if i < position {
			sb.WriteString(CorrectStyle.Render(ch))
		} else if i == position {
			if hasError {
				sb.WriteString(ErrorStyle.Render(ch))
			} else {
				sb.WriteString(CursorStyle.Render(ch))
			}
		} else {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}
