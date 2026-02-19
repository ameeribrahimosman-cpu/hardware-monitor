package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/omnitop/internal/metrics"
)

type CPUModel struct {
	width  int
	height int
	stats  metrics.SystemStats
	Alert  bool
}

func NewCPUModel() CPUModel {
	return CPUModel{}
}

func (m CPUModel) Init() tea.Cmd {
	return nil
}

func (m CPUModel) Update(msg tea.Msg) (CPUModel, tea.Cmd) {
	return m, nil
}

func (m *CPUModel) SetStats(stats metrics.SystemStats) {
	m.stats = stats
}

func (m *CPUModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m CPUModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	style := PanelStyle
	if m.Alert {
		style = AlertPanelStyle
	}
	style = style.Copy().Width(m.width - 2).Height(m.height - 2)

	// Uptime
	uptimeDuration := m.stats.Uptime
	days := uptimeDuration / 86400
	hours := (uptimeDuration % 86400) / 3600
	mins := (uptimeDuration % 3600) / 60
	uptimeStr := fmt.Sprintf("Up: %dd %02dh %02dm", days, hours, mins)

	// CPU Header
	// We need to calculate available width for spacer
	availSpacer := m.width - 20 - len(uptimeStr) - 2 // -2 border/padding
	if availSpacer < 0 { availSpacer = 0 }

	cpuHeader := lipgloss.JoinHorizontal(lipgloss.Left,
		TitleStyle.Render(fmt.Sprintf("CPU: %.1f%%", m.stats.CPU.GlobalUsagePercent)),
		lipgloss.PlaceHorizontal(availSpacer, lipgloss.Right, " "),
		MetricLabelStyle.Render(uptimeStr),
	)

	// Load Average
	loadStr := fmt.Sprintf("Load: %.2f %.2f %.2f", m.stats.CPU.LoadAvg[0], m.stats.CPU.LoadAvg[1], m.stats.CPU.LoadAvg[2])
	load := MetricLabelStyle.Render(loadStr)

	// Calculate space for Cores
	// Header(1) + Load(1) + \n(1) + \n(1) + GPU Header(1) + GPU Bar(1) = 6 lines
	// Available height for cores:
	contentH := m.height - 2
	availHeight := contentH - 6
	if availHeight < 2 {
		availHeight = 2
	}

	cores := renderCores(m.stats.CPU.PerCoreUsage, m.stats.CPU.PerCoreTemp, m.width-4, availHeight)

	// GPU Summary Mini-Graph
	gpuSummary := ""
	if m.stats.GPU.Available {
		gpuSummary = renderBar(int(m.stats.GPU.Utilization), 100, m.width-4, fmt.Sprintf("GPU %d%%", m.stats.GPU.Utilization))
	} else {
		gpuSummary = MetricLabelStyle.Render("GPU: N/A")
	}

	// Combine
	content := lipgloss.JoinVertical(lipgloss.Left,
		cpuHeader,
		load,
		"\n",
		cores,
		"\n",
		TitleStyle.Render("GPU Summary"),
		gpuSummary,
	)

	return style.Render(content)
}

func renderCores(usage []float64, temps []float64, width, height int) string {
	if len(usage) == 0 {
		return "No CPU Data"
	}

	// Dynamic columns based on width and count
	// Target roughly 15-20 chars per column
	// If width is small, maybe 1 column
	// If width is large, maybe 2 or 3
	colWidth := 20
	if width < 25 {
		colWidth = width // Single column if very narrow
	}

	numCols := width / colWidth
	if numCols < 1 { numCols = 1 }

	var sb strings.Builder

	// Calculate how many rows we can fit
	rows := (len(usage) + numCols - 1) / numCols

	// If rows exceed height, we might need more columns or truncation
	// For MVP, just truncate if too tall, but try to fit

	for r := 0; r < rows; r++ {
		if r >= height {
			// sb.WriteString(MetricLabelStyle.Render("..."))
			break
		}

		rowStr := ""
		for c := 0; c < numCols; c++ {
			idx := r*numCols + c
			if idx >= len(usage) {
				break
			}

			// Render individual core bar
			label := fmt.Sprintf("%2d", idx)
			u := usage[idx]

			// Calculate width for this column item
			// itemWidth := width / numCols
			// Padding between columns?
			itemWidth := (width - ((numCols - 1) * 2)) / numCols
			if itemWidth < 5 { itemWidth = 5 }

			bar := renderBarCompact(int(u), 100, itemWidth, label)

			if c < numCols - 1 {
				rowStr += bar + "  "
			} else {
				rowStr += bar
			}
		}
		sb.WriteString(rowStr + "\n")
	}

	return sb.String()
}

func renderBarCompact(value, max, width int, label string) string {
	// [Label ||||| ]
	// Label usually 2 chars " 0", " 1"
	labelLen := len(label)

	// bar brackets [ ] take 2 chars
	// space take 1 char
	// total overhead = labelLen + 3

	barLen := width - labelLen - 3
	if barLen < 1 {
		// Just text if too small
		return fmt.Sprintf("%s%d", label, value)
	}

	filled := int(float64(value) / float64(max) * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	empty := barLen - filled

	bar := strings.Repeat("|", filled) + strings.Repeat(" ", empty)

	style := BarStyle
	if value > 80 {
		style = AlertBarStyle
	}

	return fmt.Sprintf("%s [%s]", label, style.Render(bar))
}
