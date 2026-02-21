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

// RenderSparkline renders a multi-line ASCII sparkline graph
func RenderSparkline(data []float64, width, height int, maxVal float64, label string) string {
	if width < 1 || height < 1 {
		return ""
	}

	// If maxVal is 0, find max in data
	if maxVal <= 0 {
		for _, v := range data {
			if v > maxVal {
				maxVal = v
			}
		}
	}
	if maxVal <= 0 {
		maxVal = 1.0 // Avoid div by zero
	}

	// Use last N points that fit width
	window := data
	if len(window) > width {
		window = window[len(window)-width:]
	}
	startIdx := width - len(window)

	// Symbols for graph:   ▂▃▄▅▆▇█
	symbols := []rune{' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	// Grid
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for x, val := range window {
		// Normalized height (0 to height)
		normH := (val / maxVal) * float64(height)
		fullBlocks := int(normH)
		remainder := normH - float64(fullBlocks)

		colIdx := startIdx + x
		if colIdx >= width {
			continue
		}

		// Draw from bottom up
		for y := 0; y < fullBlocks; y++ {
			rowIdx := height - 1 - y
			if rowIdx >= 0 {
				grid[rowIdx][colIdx] = '█'
			}
		}

		// Partial block
		if fullBlocks < height {
			symIdx := int(remainder * 8)
			if symIdx > 8 {
				symIdx = 8
			}
			if symIdx < 0 {
				symIdx = 0
			}
			rowIdx := height - 1 - fullBlocks
			if rowIdx >= 0 {
				grid[rowIdx][colIdx] = symbols[symIdx]
			}
		}
	}

	var sb strings.Builder
	if label != "" {
		sb.WriteString(MetricLabelStyle.Render(label) + "\n")
	}

	for _, row := range grid {
		sb.WriteString(BarStyle.Render(string(row)) + "\n")
	}

	return sb.String()
}
