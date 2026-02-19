package metrics

import (
	"math/rand"
	"time"
)

type MockProvider struct {
	lastStats SystemStats
}

func (m *MockProvider) Init() error {
	m.lastStats = SystemStats{
		Timestamp: time.Now(),
		Uptime:    3600,
		CPU: CPUStats{
			PerCoreUsage: make([]float64, 8), // Simulate 8 cores
			PerCoreTemp:  make([]float64, 8),
		},
		GPU: GPUStats{
			Available:      true,
			Name:           "NVIDIA GeForce RTX 4090",
			MemoryTotal:    24576 * 1024 * 1024,
			HistoricalUtil: make([]float64, 100),
			Processes:      make([]GPUProcess, 0),
		},
		Processes: make([]ProcessInfo, 50),
		Disk: DiskStats{
			ReadSpeed:  1024 * 1024 * 10, // 10 MB/s
			WriteSpeed: 1024 * 1024 * 5,  // 5 MB/s
		},
		Net: NetStats{
			DownloadSpeed: 1024 * 1024 * 2, // 2 MB/s
			UploadSpeed:   1024 * 1024 * 1, // 1 MB/s
		},
	}
	// Fill historical graph with some noise
	for i := range m.lastStats.GPU.HistoricalUtil {
		m.lastStats.GPU.HistoricalUtil[i] = 10 + rand.Float64()*30
	}
	return nil
}

func (m *MockProvider) GetStats() (*SystemStats, error) {
	// Simulate metric updates with some randomness
	now := time.Now()
	m.lastStats.Timestamp = now
	m.lastStats.Uptime += 1

	// CPU
	m.lastStats.CPU.GlobalUsagePercent = 20 + rand.Float64()*10
	for i := range m.lastStats.CPU.PerCoreUsage {
		m.lastStats.CPU.PerCoreUsage[i] = 10 + rand.Float64()*30
		m.lastStats.CPU.PerCoreTemp[i] = 40 + rand.Float64()*10
	}
	m.lastStats.CPU.LoadAvg = [3]float64{1.5, 1.2, 0.8}

	// Memory
	m.lastStats.Memory.Total = 32 * 1024 * 1024 * 1024
	m.lastStats.Memory.Used = 12 * 1024 * 1024 * 1024 + uint64(rand.Int63n(1024*1024*1024))
	m.lastStats.Memory.Free = m.lastStats.Memory.Total - m.lastStats.Memory.Used
	m.lastStats.Memory.UsedPercent = float64(m.lastStats.Memory.Used) / float64(m.lastStats.Memory.Total) * 100
	m.lastStats.Memory.SwapTotal = 8 * 1024 * 1024 * 1024
	m.lastStats.Memory.SwapUsed = 1 * 1024 * 1024 * 1024
	m.lastStats.Memory.SwapPercent = 12.5

	// GPU
	m.lastStats.GPU.Utilization = uint32(50 + rand.Intn(30))
	m.lastStats.GPU.Temperature = uint32(60 + rand.Intn(10))
	m.lastStats.GPU.MemoryUsed = uint64(8 * 1024 * 1024 * 1024)
	m.lastStats.GPU.FanSpeed = uint32(40 + rand.Intn(10))
	m.lastStats.GPU.GraphicsClock = 2500
	m.lastStats.GPU.MemoryClock = 10500
	m.lastStats.GPU.PowerUsage = 150000 // mW
	m.lastStats.GPU.PowerLimit = 450000 // mW
	if m.lastStats.GPU.MemoryTotal > 0 {
		m.lastStats.GPU.MemoryUtil = uint32(float64(m.lastStats.GPU.MemoryUsed) / float64(m.lastStats.GPU.MemoryTotal) * 100.0)
	}

	// Historical Graph
	if len(m.lastStats.GPU.HistoricalUtil) > 0 {
		m.lastStats.GPU.HistoricalUtil = append(m.lastStats.GPU.HistoricalUtil[1:], float64(m.lastStats.GPU.Utilization))
	}

	// Disk & Net (Vary slightly)
	m.lastStats.Disk.ReadSpeed = uint64(float64(10*1024*1024) * (0.8 + rand.Float64()*0.4))
	m.lastStats.Disk.WriteSpeed = uint64(float64(5*1024*1024) * (0.8 + rand.Float64()*0.4))
	m.lastStats.Net.DownloadSpeed = uint64(float64(2*1024*1024) * (0.8 + rand.Float64()*0.4))
	m.lastStats.Net.UploadSpeed = uint64(float64(1*1024*1024) * (0.8 + rand.Float64()*0.4))


	// Fake Processes
	users := []string{"root", "jules", "systemd", "mysql"}
	cmds := []string{"chrome", "code", "go", "kworker", "bash", "python", "java"}
	m.lastStats.Processes = make([]ProcessInfo, 50)

	// Reset GPU processes for this tick
	m.lastStats.GPU.Processes = make([]GPUProcess, 0)

	for i := 0; i < len(m.lastStats.Processes); i++ {
		isGpu := i < 5
		pid := int32(1000 + i)
		name := cmds[rand.Intn(len(cmds))]
		m.lastStats.Processes[i] = ProcessInfo{
			PID:        pid,
			User:       users[rand.Intn(len(users))],
			Command:    name,
			State:      "R",
			CPUPercent: rand.Float64() * 5,
			MemPercent: rand.Float64() * 2,
			IsGPUUser:  isGpu,
			Threads:    int32(1 + rand.Intn(10)),
			Priority:   0,
		}

		if isGpu {
			m.lastStats.GPU.Processes = append(m.lastStats.GPU.Processes, GPUProcess{
				PID:        uint32(pid),
				Name:       name,
				MemoryUsed: uint64(rand.Int63n(1000) * 1024 * 1024),
			})
		}
	}

	// Simulate GPU process name resolution (mock already has it)

	// Return a COPY or POINTER?
	// If I modify m.lastStats next time, the caller might still be holding the pointer.
	// But in Bubble Tea loop, we usually consume stats immediately.
	// For safety, returning pointer to member is risky if caller modifies it, but here caller is read-only UI.
	// The problem is m.lastStats internal state (HistoricalUtil slice) is mutated.
	// That's fine.

	return &m.lastStats, nil
}

func (m *MockProvider) Shutdown() {}
