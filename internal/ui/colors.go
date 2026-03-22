package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/typer/internal/level"
)

const (
	ColorPinky  = "#FF6B9D"
	ColorRing   = "#C084FC"
	ColorMiddle = "#60A5FA"
	ColorIndex  = "#34D399"
	ColorThumb  = "#FBBF24"

	ColorCorrect = "#22C55E"
	ColorError   = "#EF4444"
	ColorLocked  = "#4B5563"
	ColorCursor  = "#F59E0B"
	ColorNextKey = "#FFFFFF"
	ColorDimText = "#6B7280"
)

var fingerColors = map[level.Finger]string{
	level.LPinky: ColorPinky,
	level.LRing:  ColorRing,
	level.LMid:   ColorMiddle,
	level.LIndex: ColorIndex,
	level.RIndex: ColorIndex,
	level.RMid:   ColorMiddle,
	level.RRing:  ColorRing,
	level.RPinky: ColorPinky,
	level.LThumb: ColorThumb,
	level.RThumb: ColorThumb,
}

func ColorForFinger(f level.Finger) string {
	return fingerColors[f]
}

func StyleForFinger(f level.Finger) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(fingerColors[f]))
}

var (
	CorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorCorrect))
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError))
	CursorStyle  = lipgloss.NewStyle().Background(lipgloss.Color(ColorCursor))
	LockedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorLocked))
	DimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDimText))
)
