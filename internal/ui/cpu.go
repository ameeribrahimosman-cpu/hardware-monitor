package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/omnitop/internal/metrics"
)

type CPUModel struct {
	width        int
	height       int
	stats        metrics.SystemStats // Holds all for summary
	Alert        bool
	highlightPID int
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

func (m *CPUModel) SetHighlight(pid int) {
	m.highlightPID = pid
}

func (m CPUModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	style := PanelStyle
	if m.Alert {
		style = AlertPanelStyle
	}
	style = style.Copy().Width(m.width).Height(m.height)

	// Uptime
	uptimeDuration := m.stats.Uptime
	days := uptimeDuration / 86400
	hours := (uptimeDuration % 86400) / 3600
	mins := (uptimeDuration % 3600) / 60
	uptimeStr := fmt.Sprintf("Up: %dd %02dh %02dm", days, hours, mins)

	// CPU Header
	cpuHeader := lipgloss.JoinHorizontal(lipgloss.Left,
		TitleStyle.Render(fmt.Sprintf("CPU: %.1f%%", m.stats.CPU.GlobalUsagePercent)),
		lipgloss.PlaceHorizontal(m.width-20-len(uptimeStr), lipgloss.Right, " "),
		MetricLabelStyle.Render(uptimeStr),
	)

	// Load Average
	loadStr := fmt.Sprintf("Load: %.2f %.2f %.2f", m.stats.CPU.LoadAvg[0], m.stats.CPU.LoadAvg[1], m.stats.CPU.LoadAvg[2])
	load := MetricLabelStyle.Render(loadStr)

	// Cores
	availHeight := m.height - 8 // Reserve for header, load, gpu summary
	if availHeight < 5 {
		availHeight = 5
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
	colWidth := 20
	numCols := width / colWidth
	if numCols < 1 {
		numCols = 1
	}

	var sb strings.Builder

	rows := (len(usage) + numCols - 1) / numCols

	for r := 0; r < rows; r++ {
		if r >= height {
			sb.WriteString(MetricLabelStyle.Render("..."))
			break
		}

		rowStr := ""
		for c := 0; c < numCols; c++ {
			idx := r*numCols + c
			if idx >= len(usage) {
				break
			}

			// Render individual core bar
			// [ 0] ||||| 50%
			label := fmt.Sprintf("%2d", idx)
			// Compact bar
			u := usage[idx]
			// We have colWidth - padding
			w := (width / numCols) - 2
			if w < 5 {
				w = 5
			}

			bar := renderBarCompact(int(u), 100, w, label)
			rowStr += bar + "  "
		}
		sb.WriteString(rowStr + "\n")
	}

	return sb.String()
}

func renderBarCompact(value, max, width int, label string) string {
	// [Label ||||| ]
	labelLen := len(label)
	barLen := width - labelLen - 3 // [ ] and space
	if barLen < 5 {
		// Just text if too small
		return fmt.Sprintf("%s %d%%", label, value)
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
