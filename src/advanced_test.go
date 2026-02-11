package main

import (
	"testing"
	"time"
)

func TestUTLSDialer(t *testing.T) {
	dialer := NewUTLSDialer(5 * time.Second)

	if dialer == nil {
		t.Fatal("Failed to create uTLS dialer")
	}

	if dialer.timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", dialer.timeout)
	}

	if len(dialer.availableBrowsers) == 0 {
		t.Error("No browser fingerprints available")
	}
}

func TestSNIFragmenter(t *testing.T) {
	fragmenter := NewSNIFragmenter()

	if fragmenter == nil {
		t.Fatal("Failed to create SNI fragmenter")
	}

	if fragmenter.fragmentSize <= 0 {
		t.Error("Fragment size must be positive")
	}

	if fragmenter.delayBetween < 0 {
		t.Error("Delay must be non-negative")
	}
}

func TestEnhancedHealthScorer(t *testing.T) {
	scorer := NewEnhancedHealthScorer(true)

	if scorer == nil {
		t.Fatal("Failed to create health scorer")
	}

	scorer.RecordLatency(100, time.Now())
	scorer.RecordLatency(110, time.Now())
	scorer.RecordLatency(105, time.Now())

	jitter := scorer.CalculateJitter()
	if jitter < 0 {
		t.Error("Jitter cannot be negative")
	}

	health := scorer.CalculateAdvancedHealthScore()
	if health < 0 || health > 100 {
		t.Errorf("Health score %d out of range 0-100", health)
	}
}

func TestAdvancedAntiDPIClient(t *testing.T) {
	client := NewAdvancedAntiDPIClient(5*time.Second, true)

	if client == nil {
		t.Fatal("Failed to create advanced client")
	}

	if client.timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", client.timeout)
	}

	if !client.iranMode {
		t.Error("Iran mode should be enabled")
	}

	stats := client.GetStatistics()
	if stats == nil {
		t.Error("Statistics should not be nil")
	}
}

func TestBrowserFingerprintDatabase(t *testing.T) {
	browsers := []string{"chrome120", "firefox121", "edge120", "safari17"}

	for _, browser := range browsers {
		fp, exists := FingerprintDatabase[browser]
		if !exists {
			t.Errorf("Browser fingerprint %s not found", browser)
			continue
		}

		if len(fp.CipherSuites) == 0 {
			t.Errorf("Browser %s has no cipher suites", browser)
		}

		if len(fp.SupportedCurves) == 0 {
			t.Errorf("Browser %s has no supported curves", browser)
		}
	}
}

func BenchmarkJitterCalculation(b *testing.B) {
	scorer := NewEnhancedHealthScorer(true)

	for i := 0; i < 100; i++ {
		scorer.RecordLatency(int64(100+i%50), time.Now())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scorer.CalculateJitter()
	}
}

func BenchmarkHealthScoreCalculation(b *testing.B) {
	scorer := NewEnhancedHealthScorer(true)

	for i := 0; i < 100; i++ {
		scorer.RecordLatency(int64(100+i%50), time.Now())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scorer.CalculateAdvancedHealthScore()
	}
}
