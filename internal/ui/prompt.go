package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderPrompt(prompt []rune, position int, hasError bool) string {
	var sb strings.Builder

	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDimText))
	sb.WriteString(promptStyle.Render(string(prompt)))
	sb.WriteString("\n")

	for i := 0; i < len(prompt); i++ {
		ch := string(prompt[i])
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
