package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
)

// SystemMonitor tracks system metrics and performance
type SystemMonitor struct {
	startTime   time.Time
	events      []MonitorEvent
	metrics     map[string]float64
	mu          sync.RWMutex
	enabled     bool
	reportFile  string
}

// MonitorEvent represents a tracked event
type MonitorEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Details   string    `json:"details,omitempty"`
}

// MetricsReport contains performance metrics
type MetricsReport struct {
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Duration     float64                `json:"duration_seconds"`
	Events       []MonitorEvent         `json:"events"`
	Metrics      map[string]float64     `json:"metrics"`
	SystemInfo   SystemInfo             `json:"system_info"`
}

// SystemInfo contains system information
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"num_cpu"`
	GOMAXPROCS   int    `json:"gomaxprocs"`
}

// NewSystemMonitor creates a new system monitor
func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{
		startTime:  time.Now(),
		events:     make([]MonitorEvent, 0),
		metrics:    make(map[string]float64),
		enabled:    true,
		reportFile: "results/monitoring-report.json",
	}
}

// Start initializes the monitoring system
func (sm *SystemMonitor) Start() {
	if !sm.enabled {
		return
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.startTime = time.Now()
	sm.RecordEvent("monitoring_started")
	
	color.Green("📊 System monitoring started")
}

// Stop finalizes the monitoring session
func (sm *SystemMonitor) Stop() {
	if !sm.enabled {
		return
	}
	
	sm.RecordEvent("monitoring_stopped")
	color.Green("📊 System monitoring stopped")
}

// RecordEvent logs an event
func (sm *SystemMonitor) RecordEvent(event string) {
	if !sm.enabled {
		return
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.events = append(sm.events, MonitorEvent{
		Timestamp: time.Now(),
		Event:     event,
	})
}

// RecordEventWithDetails logs an event with details
func (sm *SystemMonitor) RecordEventWithDetails(event, details string) {
	if !sm.enabled {
		return
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.events = append(sm.events, MonitorEvent{
		Timestamp: time.Now(),
		Event:     event,
		Details:   details,
	})
}

// RecordMetric records a metric value
func (sm *SystemMonitor) RecordMetric(name string, value float64) {
	if !sm.enabled {
		return
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.metrics[name] = value
}

// IncrementMetric increments a metric
func (sm *SystemMonitor) IncrementMetric(name string, delta float64) {
	if !sm.enabled {
		return
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.metrics[name] += delta
}

// GetMetric retrieves a metric value
func (sm *SystemMonitor) GetMetric(name string) float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	return sm.metrics[name]
}

// GenerateReport creates a comprehensive monitoring report
func (sm *SystemMonitor) GenerateReport() error {
	if !sm.enabled {
		return nil
	}
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	endTime := time.Now()
	duration := endTime.Sub(sm.startTime).Seconds()
	
	report := MetricsReport{
		StartTime: sm.startTime,
		EndTime:   endTime,
		Duration:  duration,
		Events:    sm.events,
		Metrics:   sm.metrics,
		SystemInfo: SystemInfo{
			GoVersion:  runtime.Version(),
			GOOS:       runtime.GOOS,
			GOARCH:     runtime.GOARCH,
			NumCPU:     runtime.NumCPU(),
			GOMAXPROCS: runtime.GOMAXPROCS(0),
		},
	}
	
	// Ensure results directory exists
	os.MkdirAll("results", 0755)
	
	// Write JSON report
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	
	if err := os.WriteFile(sm.reportFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}
	
	// Print summary
	sm.printReportSummary(report)
	
	return nil
}

func (sm *SystemMonitor) printReportSummary(report MetricsReport) {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                  📊 MONITORING REPORT                         ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Duration:            %-40.2fs║", report.Duration)
	color.Yellow("║  Total Events:        %-40d║", len(report.Events))
	color.Yellow("║  Total Metrics:       %-40d║", len(report.Metrics))
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	
	// Show key metrics
	if len(report.Metrics) > 0 {
		color.Cyan("║                      KEY METRICS                              ║")
		color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
		
		keyMetrics := []string{
			"proxy_check_duration",
			"config_gen_duration",
			"config_test_duration",
			"configs_generated",
			"configs_passed",
		}
		
		for _, key := range keyMetrics {
			if val, ok := report.Metrics[key]; ok {
				color.Yellow("║  %-30s: %-25.2f  ║", key, val)
			}
		}
	}
	
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Report saved to: %-43s║", truncate(sm.reportFile, 43))
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// GetEventCount returns the number of events recorded
func (sm *SystemMonitor) GetEventCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	return len(sm.events)
}

// GetMetricCount returns the number of metrics recorded
func (sm *SystemMonitor) GetMetricCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	return len(sm.metrics)
}

// GetDuration returns the monitoring duration
func (sm *SystemMonitor) GetDuration() time.Duration {
	return time.Since(sm.startTime)
}
