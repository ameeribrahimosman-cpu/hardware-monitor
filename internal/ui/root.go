package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/omnitop/internal/config"
	"github.com/google/omnitop/internal/metrics"
)

type TickMsg time.Time

// tickRand is used to add jitter to polling intervals.
// Safe for use in this context as tick() is only called from
// the single-threaded Bubble Tea event loop.
var tickRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func tick() tea.Cmd {
	// Add jitter: 900-1100ms
	jitter := time.Duration(tickRand.Intn(200)-100) * time.Millisecond
	return tea.Tick(time.Second+jitter, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type RootModel struct {
	provider metrics.Provider
	config   *config.ProfileConfiguration

	// Sub-models
	gpu     GPUModel
	process ProcessModel
	cpu     CPUModel
	footer  FooterModel

	// Layout state
	width, height int
	col1Pct       float64 // Percentage of width for Left Column (GPU)
	col2Pct       float64 // Percentage of width for Middle Column (Process)
	// Right column takes remaining

	// Tooltip state
	mouseX, mouseY           int
	showTooltip              bool
	currentTooltipRegion     string
}

func NewRootModel(provider metrics.Provider, cfg *config.ProfileConfiguration) RootModel {
	// Use configuration values or defaults
	col1Pct := 0.30
	col2Pct := 0.40

	if cfg != nil {
		if val, ok := cfg.ColumnWidths["gpu"]; ok {
			col1Pct = val
		}
		if val, ok := cfg.ColumnWidths["process"]; ok {
			col2Pct = val
		}
	}

	return RootModel{
		provider: provider,
		config:   cfg,
		gpu:      NewGPUModel(),
		process:  NewProcessModel(),
		cpu:      NewCPUModel(),
		footer:   NewFooterModel(),
		col1Pct:  col1Pct,
		col2Pct:  col2Pct,
	}
}

func (m RootModel) Init() tea.Cmd {
	return tick()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "[": // Shrink Left Col
			m.col1Pct -= 0.05
			if m.col1Pct < 0.1 {
				m.col1Pct = 0.1
			}
			m.resizeModules()
		case "]": // Expand Left Col
			m.col1Pct += 0.05
			if m.col1Pct+m.col2Pct > 0.9 {
				m.col1Pct = 0.9 - m.col2Pct
			}
			m.resizeModules()
		case "{": // Shrink Middle Col (effectively expands Right)
			m.col2Pct -= 0.05
			if m.col2Pct < 0.1 {
				m.col2Pct = 0.1
			}
			m.resizeModules()
		case "}": // Expand Middle Col
			m.col2Pct += 0.05
			if m.col1Pct+m.col2Pct > 0.9 {
				m.col2Pct = 0.9 - m.col1Pct
			}
			m.resizeModules()
		}

		// Pass keys to sub-models if needed (e.g. process list scrolling)
		m.process, cmd = m.process.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeModules()

	case TickMsg:
		// Fetch metrics
		stats, err := m.provider.GetStats()
		if err == nil {
			m.gpu.SetStats(stats.GPU)
			m.process.SetStats(stats.Processes)
			m.cpu.SetStats(*stats)
		}
		// Continue tick
		cmds = append(cmds, tick())

	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y
		
		// Tooltip handling based on mouse position
		if m.config != nil && m.config.ShowTooltips {
			// Determine which UI region mouse is in
			region := m.determineMouseRegion()
			m.showTooltip = (region != "")
			
			// Store region for tooltip content generation
			if m.showTooltip {
				m.currentTooltipRegion = region
			}
		} else {
			m.showTooltip = false
			m.currentTooltipRegion = ""
		}

		// Pass mouse to sub-models
		m.process, cmd = m.process.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update other models if they have logic
	// gpu and cpu currently have no internal update logic needing msgs

	return m, tea.Batch(cmds...)
}

func (m *RootModel) resizeModules() {
	if m.width == 0 || m.height == 0 {
		return
	}

	// Calculate widths
	w1 := int(float64(m.width) * m.col1Pct)
	w2 := int(float64(m.width) * m.col2Pct)
	w3 := m.width - w1 - w2

	// Height available for columns (minus footer)
	h := m.height - 1
	if h < 1 {
		h = 1
	}

	m.gpu.SetSize(w1, h)
	m.process.SetSize(w2, h)
	m.cpu.SetSize(w3, h)
	m.footer.SetSize(m.width)
}

func (m RootModel) determineMouseRegion() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Calculate column widths
	w1 := int(float64(m.width) * m.col1Pct)
	w2 := int(float64(m.width) * m.col2Pct)
	
	// Column boundaries
	// Column 1: GPU (x: 0 to w1-1)
	// Column 2: Process (x: w1 to w1 + w2 - 1)
	// Column 3: CPU (x: w1 + w2 to m.width-1)
	
	// Footer is at bottom row (height - 1)
	footerRow := m.height - 1
	
	// Check if mouse is in footer
	if m.mouseY == footerRow {
		return "footer"
	}
	
	// Check if mouse is in GPU column
	if m.mouseX >= 0 && m.mouseX < w1 && m.mouseY < footerRow {
		return "gpu"
	}
	
	// Check if mouse is in Process column
	if m.mouseX >= w1 && m.mouseX < w1 + w2 && m.mouseY < footerRow {
		return "process"
	}
	
	// Check if mouse is in CPU column
	if m.mouseX >= w1 + w2 && m.mouseX < m.width && m.mouseY < footerRow {
		return "cpu"
	}
	
	return ""
}

func (m RootModel) getTooltipContent(region string) string {
	switch region {
	case "gpu":
		gpuStats := m.gpu.stats
		base := "GPU Panel: Shows GPU utilization, memory, temperature, and health status. Click to expand/collapse details."
		
		if gpuStats.Available {
			// Add health status
			healthInfo := ""
			switch gpuStats.HealthStatus {
			case metrics.GPUHealthHealthy:
				healthInfo = "Health: ✅ Healthy"
			case metrics.GPUHealthDegraded:
				healthInfo = "Health: ⚠️ Degraded"
			case metrics.GPUHealthFailed:
				healthInfo = "Health: 🔴 Failed"
			}
			
			// Add temperature warning if needed
			tempWarning := ""
			if gpuStats.Temperature > 90 {
				tempWarning = " | 🔥 CRITICAL TEMPERATURE!"
			} else if gpuStats.Temperature > 80 {
				tempWarning = " | ⚠️ High Temperature"
			}
			
			// Add fan warning
			fanWarning := ""
			if gpuStats.Temperature > 80 && gpuStats.FanSpeed < 30 {
				fanWarning = " | 💨 Low Fan Speed"
			}
			
			// Add error info if any
			errorInfo := ""
			if gpuStats.ErrorCount > 0 {
				errorInfo = fmt.Sprintf(" | Errors: %d", gpuStats.ErrorCount)
			}
			
			return fmt.Sprintf("%s\n\n%s%s%s%s", base, healthInfo, tempWarning, fanWarning, errorInfo)
		}
		return base
		
	case "process":
		return "Process List: Shows running processes with highest GPU/CPU usage. Use Up/Down arrows to navigate."
	case "cpu":
		return "CPU Panel: Shows CPU utilization, temperature, frequency, and core usage. Click to toggle core view."
	case "footer":
		return "Footer: Shows hotkeys and current time. [ and ] resize GPU column, { and } resize Process column."
	default:
		return ""
	}
}

func (m RootModel) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Render columns
	cols := lipgloss.JoinHorizontal(lipgloss.Top,
		m.gpu.View(),
		m.process.View(),
		m.cpu.View(),
	)

	// Add Footer
	mainView := lipgloss.JoinVertical(lipgloss.Left,
		cols,
		m.footer.View(),
	)
	
	// Add tooltip if enabled and mouse is in a region
	if m.showTooltip && m.currentTooltipRegion != "" {
		tooltipContent := m.getTooltipContent(m.currentTooltipRegion)
		if tooltipContent != "" {
			tooltipStyle := lipgloss.NewStyle().
				Background(lipgloss.Color(ColorSteelGray)).
				Foreground(lipgloss.Color(ColorMidnightBlack)).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(ColorIceBlue))
			
			tooltip := tooltipStyle.Render(tooltipContent)
			tooltipWidth := lipgloss.Width(tooltip)
			
			// Position tooltip horizontally near mouse cursor
			// Try to align left edge with mouse X, but keep within terminal bounds
			desiredX := m.mouseX
			if desiredX + tooltipWidth > m.width {
				// Shift left to keep tooltip on screen
				desiredX = m.width - tooltipWidth
			}
			if desiredX < 0 {
				desiredX = 0
			}
			
			// Add spaces before tooltip to position it horizontally
			spaces := ""
			if desiredX > 0 {
				spaces = strings.Repeat(" ", desiredX)
			}
			
			// Add tooltip on new line after main view
			return mainView + "\n" + spaces + tooltip
		}
	}

	return mainView
}
