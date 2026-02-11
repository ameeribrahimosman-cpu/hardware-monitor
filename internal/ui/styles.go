package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme colors based on "Wrath of the Lich King" palette
const (
	ColorMidnightBlack = "#0A001F" // Background
	ColorIceBlue       = "#81A1C1" // Primary UI/Text
	ColorSteelGray     = "#4C566A" // Panels/Borders
	ColorPaleBlue      = "#8FBCBB" // Graphs/Normal Metrics
	ColorBloodCrimson  = "#C41E3A" // Alerts/Errors
)

var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(ColorMidnightBlack)).
			Foreground(lipgloss.Color(ColorIceBlue))

	// Panel styles
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorSteelGray)).
			Padding(0, 1)

	AlertPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorBloodCrimson)).
			Padding(0, 1)

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorIceBlue)).
			Bold(true)

	TextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorIceBlue))

	MetricLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSteelGray))

	MetricValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorPaleBlue))

	// Alert styles
	AlertStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBloodCrimson)).
			Bold(true)

	// Bar styles
	BarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPaleBlue))

	AlertBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBloodCrimson))
)

// renderBar renders a simple progress bar
func renderBar(value, max, width int, label string) string {
	if width < 10 {
		return label
	}
	barWidth := width - lipgloss.Width(label) - 2
	if barWidth < 0 {
		barWidth = 0
	}

	filled := int(float64(value) / float64(max) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	style := BarStyle
	if value > 80 {
		style = AlertBarStyle
	}

	return fmt.Sprintf("%s %s", label, style.Render(bar))
}
