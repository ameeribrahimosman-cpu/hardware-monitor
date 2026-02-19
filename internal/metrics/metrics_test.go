package metrics

import (
	"testing"
)

func TestMockProvider(t *testing.T) {
	m := &MockProvider{}
	if err := m.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	for i := 0; i < 5; i++ {
		stats, err := m.GetStats()
		if err != nil {
			t.Fatalf("GetStats failed: %v", err)
		}

		if stats == nil {
			t.Fatal("Stats is nil")
		}

		if len(stats.Processes) != 50 {
			t.Errorf("Expected 50 processes, got %d", len(stats.Processes))
		}

		// GPU Process check
		// We expect 5 GPU processes based on mock logic
		if len(stats.GPU.Processes) != 5 {
			t.Errorf("Expected 5 GPU processes, got %d", len(stats.GPU.Processes))
		}

		// Historical Graph check
		if len(stats.GPU.HistoricalUtil) < 1 {
			t.Error("HistoricalUtil is empty")
		}
	}
}
