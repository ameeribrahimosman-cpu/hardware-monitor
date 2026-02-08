# Implementation Plan

Implement comprehensive fixes for the OmniTop hardware monitor repository to address identified issues including missing configuration system, build artifacts, documentation discrepancies, and code improvements.

This implementation addresses five main issues identified in the OmniTop repository: 1) Empty profiles.json configuration file not integrated with the application, 2) Missing omnitop.png icon file required for AppImage builds, 3) Unresolved TODO comment for tooltip functionality, 4) README repository origin mismatch causing confusion, and 5) Verification of the complete build process from source to AppImage. The plan follows a phased approach starting with critical build infrastructure fixes, then configuration system implementation, followed by code improvements and documentation updates to ensure the repository is fully functional and properly documented for users and developers.

[Types]  
Implement configuration types, icon asset handling, and repository metadata structures.

The implementation requires the following type system changes:
1. **ProfileConfiguration** - Struct in `internal/config/types.go` defining user-configurable settings:
   - `Theme` (string): "lich-king", "default", or "custom"
   - `ColumnWidths` (map[string]float64): Percentage widths for GPU, Process, CPU columns
   - `RefreshInterval` (int): Milliseconds between updates (default: 1000)
   - `MaxProcesses` (int): Maximum processes to display (default: 200)
   - `GPUHistoryLength` (int): Number of historical GPU utilization points (default: 100)
   - `ShowTooltips` (bool): Whether to display hover tooltips (default: true)

2. **IconAsset** - Implicit handling for 256x256 PNG icon with proper licensing metadata

3. **RepositoryMetadata** - Updates to README.md and documentation to reflect actual repository fork

[Files]
Create and modify files to implement configuration system, add missing assets, and fix documentation.

Detailed breakdown:
- **New files to be created**:
  - `omnitop.png` (256x256): Application icon in root directory with OmniTop branding using Lich King theme colors (midnight black, ice blue accents)
  - `internal/config/config.go`: Configuration loading and parsing logic with JSON unmarshaling
  - `internal/config/types.go`: Type definitions for ProfileConfiguration and validation methods
  - `example_profiles.json`: Example configuration file with commented settings for users

- **Existing files to be modified**:
  - `README.md`: Update repository URLs from "google/omnitop" to "ameeribrahimosman-cpu/hardware-monitor", add configuration documentation, clarify build instructions
  - `profiles.json`: Populate with default configuration that matches current hardcoded settings
  - `internal/ui/root.go`: Replace TODO comment with basic tooltip implementation using mouse coordinates
  - `build_appimage.sh`: Add verification step for icon existence with better error messaging
  - `cmd/omnitop/main.go`: Integrate configuration loading with fallback to defaults
  - `internal/ui/root.go`: Use configuration values for column percentages and other settings

- **Files to be deleted or moved**: None

- **Configuration file updates**:
  - `go.mod`: No changes required (dependencies already sufficient)
  - `COORDINATION.md`: Update with implementation progress and new coordination procedures

[Functions]
Add configuration loading functions, tooltip implementation, and build verification.

Detailed breakdown:
- **New functions**:
  - `LoadConfig(path string) (*config.ProfileConfiguration, error)` in `internal/config/config.go`: Loads and validates JSON configuration with sensible defaults
  - `DefaultConfig() *config.ProfileConfiguration` in `internal/config/config.go`: Returns hardcoded defaults matching current behavior
  - `RenderTooltip(x, y int, content string) string` in `internal/ui/root.go`: Basic tooltip rendering using lipgloss styles
  - `VerifyBuildEnvironment() error` in `build_appimage.sh`: Check for required tools (go, wget, appimagetool)

- **Modified functions**:
  - `NewRootModel(provider metrics.Provider)` in `internal/ui/root.go`: Accept configuration parameter, use config values for col1Pct, col2Pct
  - `main()` in `cmd/omnitop/main.go`: Load configuration before creating root model, handle config errors gracefully
  - `resizeModules()` in `internal/ui/root.go`: Use configuration values instead of hardcoded 0.30, 0.40

- **Removed functions**: None

[Classes]
No new classes required; extend existing structures with configuration support.

Detailed breakdown:
- **New classes**: None (Go uses structs, not classes)
- **Modified structs**:
  - `RootModel` in `internal/ui/root.go`: Add `config *config.ProfileConfiguration` field
  - `RealProvider` and `MockProvider` in `internal/metrics/`: No changes needed
- **Removed classes**: None

[Dependencies]
No new external dependencies; use existing Go standard library JSON packages.

Details of new packages, version changes, and integration requirements:
- **Standard library packages**: `encoding/json` for configuration parsing (already available)
- **No new external dependencies**: Configuration system uses only Go standard library
- **Integration requirements**: Configuration must be backward compatible - missing profiles.json should use defaults without errors
- **Build dependencies**: `appimagetool` for AppImage creation (already handled in build script)

[Testing]
Implement unit tests for configuration system and integration tests for build process.

Test file requirements, existing test modifications, and validation strategies:
1. **Unit Tests**: Create `internal/config/config_test.go` with:
   - Test loading valid configuration
   - Test loading missing file (should use defaults)
   - Test invalid JSON handling
   - Test validation of out-of-range values

2. **Integration Tests**:
   - Verify binary builds successfully with `go build`
   - Test mock mode execution for basic functionality
   - Verify AppImage build process (can be manual/scripted)

3. **Existing Test Modifications**:
   - `internal/metrics/metrics_test.go`: No changes needed
   - Update any tests that depend on hardcoded layout values

4. **Validation Strategies**:
   - Manual verification of tooltip functionality
   - Configuration round-trip test (save config, load back, compare)
   - Build script verification in clean environment

[Implementation Order]
Sequential implementation starting with critical build fixes, then configuration, then code improvements.

Numbered steps showing the logical order of changes to minimize conflicts and ensure successful integration:
1. **Step 1 - Create missing icon file**: Generate 256x256 omnitop.png using basic graphic design (can be simple placeholder with text)
2. **Step 2 - Fix README repository references**: Update all "google/omnitop" references to actual fork URL
3. **Step 3 - Implement configuration types**: Create internal/config/types.go with ProfileConfiguration struct
4. **Step 4 - Implement configuration loading**: Create internal/config/config.go with LoadConfig and DefaultConfig functions
5. **Step 5 - Integrate configuration into main**: Modify cmd/omnitop/main.go to load profiles.json
6. **Step 6 - Update root model for configuration**: Modify internal/ui/root.go to use config values
7. **Step 7 - Populate profiles.json**: Add default configuration matching current behavior
8. **Step 8 - Create example configuration**: Add example_profiles.json with commented options
9. **Step 9 - Implement basic tooltips**: Replace TODO in root.go with simple tooltip rendering
10. **Step 10 - Enhance build script**: Add verification steps and better error messages
11. **Step 11 - Test complete build process**: Verify go build, mock execution, and AppImage creation
12. **Step 12 - Update coordination document**: Add implementation details to COORDINATION.md