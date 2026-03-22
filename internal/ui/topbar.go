package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var topBarStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(1, 3)

func RenderTopBar(levelName string, wpm, accuracy float64, drillNum, totalDrills int) string {
	level := fmt.Sprintf("Level: %s", levelName)
	wpmStr := fmt.Sprintf("WPM: %.0f", wpm)
	accStr := fmt.Sprintf("Accuracy: %.0f%%", accuracy)
	progress := fmt.Sprintf("%d/%d", drillNum, totalDrills)

	return topBarStyle.Render(
		fmt.Sprintf("  %s      %s      %s      %s", level, wpmStr, accStr, progress),
	)
}
