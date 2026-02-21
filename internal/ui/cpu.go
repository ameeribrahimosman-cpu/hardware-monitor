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
	stats  metrics.SystemStats // Holds all for summary
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

	// Calculate space for Cores
	// We need space for Memory and GPU summary at bottom?
	// The requirement says "Per-core CPU bars... memory breakdown, disk I/O, net RX/TX graphs" are in MIDDLE column?
	// Wait, the prompt says:
	// "Middle Column (~40% width): btop/htop hybrid – Full sortable process list..., bottom stacked: memory breakdown, disk I/O, net RX/TX graphs."
	// "Right Column (~30% width): Per-core CPU bars... load averages, quick GPU summary mini-graph."

	// So Memory/Disk/Net graphs should be in ProcessModel (Middle), not CPUModel (Right)?
	// But `ProcessModel` currently only has the table.
	// `RootModel` logic in `View` puts `process` in middle.
	// Currently `process.View` only renders the table.
	// I should probably move Memory/Disk/Net to Process column or a separate widget below it.
	// However, `RootModel` splits into 3 columns.
	// Left: GPU. Middle: Process. Right: CPU.
	// If the user wants Memory/Disk/Net in middle, I should add them to `ProcessModel` or split the middle column.
	// Given the complexity, I might stick to what I have or try to fit them.
	// But let's focus on CPUModel (Right Column) for now.
	// Requirement: Per-core bars, load averages, quick GPU summary.

	// Cores
	availHeight := m.height - 8 // Reserve for header, load, gpu summary
	if availHeight < 5 {
		availHeight = 5
	}

	cores := renderCores(m.stats.CPU.PerCoreUsage, m.stats.CPU.PerCoreTemp, m.width-4, availHeight)

	// GPU Summary Mini-Graph
	gpuSummary := ""
	if m.stats.GPU.Available {
		gpuSummary = RenderSparkline(m.stats.GPU.HistoricalUtil, m.width-4, 5, 100, fmt.Sprintf("GPU %d%%", m.stats.GPU.Utilization))
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
	// Target roughly 20 chars per column?
	colWidth := 20
	numCols := width / colWidth
	if numCols < 1 {
		numCols = 1
	}

	// Ensure we don't exceed height too much
	// Rows needed = ceil(count / cols)
	// If rows > height, increase cols if possible?
	// Or scroll? TUI usually doesn't scroll automatically without viewport.
	// Let's stick to simple grid.

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
