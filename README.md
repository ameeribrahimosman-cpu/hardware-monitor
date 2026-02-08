# OmniTop - The Unified System Monitor

OmniTop merges the best features of `nvtop`, `htop`, and `btop` into a single, cohesive TUI dashboard. Inspired by the "Wrath of the Lich King" aesthetic, it provides high-density system metrics with a focus on GPU telemetry, process management, and per-core CPU visualization.

## Features

- **Unified Dashboard**: 3-column layout replicating your multi-window workflow.
  - **Left**: GPU History & Telemetry (NVTop style).
  - **Middle**: Process list (HTop style).
  - **Right**: Per-core CPU bars & Load Averages (BTop style).
- **GPU First**: Native NVIDIA GPU monitoring via NVML (temps, fans, clocks, power).
- **Lich King Theme**: Midnight Black, Ice Blue, and Blood Crimson aesthetics.
- **Mock Mode**: Run without hardware sensors for testing/demo purposes.
- **Keyboard Resizing**: Adjust column widths dynamically.

## Installation

### From Source

Requirements: Go 1.24.3+

```bash
git clone https://github.com/ameeribrahimosman-cpu/hardware-monitor.git
cd hardware-monitor
go build -o omnitop ./cmd/omnitop
```

### Running

```bash
# Run with real sensors (requires NVIDIA GPU for GPU metrics)
./omnitop

# Run in Mock Mode (simulated data)
./omnitop --mock
```

## Key Bindings

| Key | Action |
|---|---|
| `q` / `Ctrl+C` | Quit |
| `[` | Shrink Left Column (GPU) |
| `]` | Expand Left Column (GPU) |
| `{` | Shrink Middle Column (Process) |
| `}` | Expand Middle Column (Process) |
| `Up` / `Down` | Navigate Process List |

## Building AppImage

To create a portable AppImage:

1. Ensure `appimagetool` is installed.
2. Run the build script:
   ```bash
   ./build_appimage.sh
   ```

## Configuration

OmniTop now features a comprehensive configuration system. The application reads settings from `profiles.json` in the current working directory, with sensible defaults provided if the file is missing or invalid.

### Configuration Options

Example `profiles.json`:
```json
{
  "theme": "lich-king",
  "columnWidths": {
    "gpu": 0.30,
    "process": 0.40,
    "cpu": 0.30
  },
  "refreshInterval": 1000,
  "maxProcesses": 200,
  "gpuHistoryLength": 100,
  "showTooltips": true
}
```

### Available Settings:

- **theme**: UI theme selection (currently only 'lich-king' supported)
- **columnWidths**: GPU, Process, and CPU panel width percentages (must sum to ≤ 0.9)
- **refreshInterval**: Polling interval in milliseconds (100-5000ms)
- **maxProcesses**: Maximum number of processes to display (10-500)
- **gpuHistoryLength**: GPU metrics history length for trend visualization (50-500)
- **showTooltips**: Enable/disable UI tooltips based on mouse position

### Configuration Management

- Configuration is loaded automatically on startup
- Missing or invalid fields fall back to sensible defaults
- Configuration validation ensures values are within safe ranges
- The system maintains backward compatibility with previous versions

## Troubleshooting

### GPU Monitoring Issues

Based on analysis of system dmesg logs, OmniTop includes robust error handling for common GPU monitoring scenarios:

**NVIDIA GPU Xid 56 Errors**: 
- These errors indicate "GPU stopped processing commands" (hardware or driver issue)
- OmniTop will gracefully handle GPU read failures and continue monitoring other metrics
- Consider checking: GPU overheating, driver updates, hardware stability

**General GPU Monitoring Recommendations**:
1. Ensure NVIDIA drivers are up to date
2. Verify GPU temperature stays within safe operating range
3. Check system logs for Xid errors using `dmesg | grep -i nvidia`
4. Test with `nvidia-smi` to verify GPU is accessible

### Common Issues

1. **TTY Error in Non-Interactive Environments**:
   - OmniTop requires a terminal with TTY support
   - This is expected behavior for TUI applications
   - Run in a terminal emulator or with proper TTY allocation

2. **Missing GPU Metrics**:
   - Ensure NVIDIA drivers and NVML library are installed
   - Run with `--mock` flag to test without GPU hardware
   - Verify permissions to access GPU device files

3. **Configuration Issues**:
   - Invalid `profiles.json` will be ignored with warnings
   - Check JSON syntax if configuration isn't loading
   - Default settings will be used for any invalid values

### Testing with Mock Mode

When GPU hardware is unavailable or for testing purposes:
```bash
./omnitop --mock
```
Mock mode provides simulated data for all metrics, allowing full testing of the UI and configuration system.

## Development

### Build Requirements
- Go 1.24.3 or later
- NVIDIA Management Library (NVML) for GPU monitoring
- Bubble Tea TUI framework dependencies

### Project Structure
- `cmd/omnitop/` - Main application entry point
- `internal/config/` - Configuration loading and validation
- `internal/ui/` - TUI components and layout
- `internal/metrics/` - Hardware metric collection
- `profiles.json` - User configuration file

### Testing the Build Process
```bash
# Basic compilation test
go build ./cmd/omnitop

# AppImage creation test
./build_appimage.sh
```
