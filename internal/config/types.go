package config

// ProfileConfiguration defines the user-configurable settings.
type ProfileConfiguration struct {
	Theme            string             `json:"theme"`
	ColumnWidths     map[string]float64 `json:"column_widths"`
	RefreshInterval  int                `json:"refresh_interval"` // Milliseconds
	MaxProcesses     int                `json:"max_processes"`
	GPUHistoryLength int                `json:"gpu_history_length"`
	ShowTooltips     bool               `json:"show_tooltips"`
	AlertThresholds  AlertThresholds    `json:"alert_thresholds"`
}

// AlertThresholds defines the limits for triggering alerts.
type AlertThresholds struct {
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	CPUTempCelsius     float64 `json:"cpu_temp_celsius"`
	GPUUsagePercent    float64 `json:"gpu_usage_percent"`
	GPUTempCelsius     float64 `json:"gpu_temp_celsius"`
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	DiskUsagePercent   float64 `json:"disk_usage_percent"`
}
