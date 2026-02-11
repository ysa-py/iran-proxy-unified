package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// NetworkMetrics contains detailed network performance metrics
type NetworkMetrics struct {
	Latencies         []int64       // All latency measurements in ms
	PacketLoss        float64       // Packet loss percentage
	Jitter            float64       // Network jitter (latency variance)
	Throughput        float64       // Estimated throughput in Mbps
	DNSResolutionTime int64         // DNS lookup time in ms
	TCPHandshakeTime  int64         // TCP 3-way handshake time in ms
	TLSHandshakeTime  int64         // TLS handshake time in ms
	TTFBTime          int64         // Time to first byte in ms
	Timestamps        []time.Time   // Measurement timestamps
	ErrorRate         float64       // Connection error rate
	StabilityScore    float64       // Connection stability (0-100)
	DPIFingerprint    string        // Detected DPI fingerprint
}

// EnhancedHealthScorer provides advanced proxy health scoring
type EnhancedHealthScorer struct {
	metrics           *NetworkMetrics
	dpiEngine         *AIAntiDPIEngine
	historicalData    []NetworkMetrics
	maxHistorySize    int
	iranOptimized     bool
}

// NewEnhancedHealthScorer creates a new enhanced health scorer
func NewEnhancedHealthScorer(iranOptimized bool) *EnhancedHealthScorer {
	return &EnhancedHealthScorer{
		metrics: &NetworkMetrics{
			Latencies:  make([]int64, 0),
			Timestamps: make([]time.Time, 0),
		},
		dpiEngine:      NewAIAntiDPIEngine("advanced"),
		historicalData: make([]NetworkMetrics, 0),
		maxHistorySize: 50,
		iranOptimized:  iranOptimized,
	}
}

// RecordLatency records a latency measurement
func (ehs *EnhancedHealthScorer) RecordLatency(latency int64, timestamp time.Time) {
	ehs.metrics.Latencies = append(ehs.metrics.Latencies, latency)
	ehs.metrics.Timestamps = append(ehs.metrics.Timestamps, timestamp)
}

// RecordConnectionMetrics records detailed connection metrics
func (ehs *EnhancedHealthScorer) RecordConnectionMetrics(
	dnsTime, tcpTime, tlsTime, ttfbTime int64) {
	
	ehs.metrics.DNSResolutionTime = dnsTime
	ehs.metrics.TCPHandshakeTime = tcpTime
	ehs.metrics.TLSHandshakeTime = tlsTime
	ehs.metrics.TTFBTime = ttfbTime
}

// RecordError records a connection error
func (ehs *EnhancedHealthScorer) RecordError() {
	totalAttempts := len(ehs.metrics.Latencies) + 1
	errors := int(ehs.metrics.ErrorRate * float64(totalAttempts-1))
	errors++
	ehs.metrics.ErrorRate = float64(errors) / float64(totalAttempts)
}

// CalculateJitter calculates network jitter (latency variance)
func (ehs *EnhancedHealthScorer) CalculateJitter() float64 {
	if len(ehs.metrics.Latencies) < 2 {
		return 0.0
	}

	// Calculate standard deviation of latencies
	mean := ehs.calculateMeanLatency()
	variance := 0.0

	for _, latency := range ehs.metrics.Latencies {
		diff := float64(latency) - mean
		variance += diff * diff
	}

	variance /= float64(len(ehs.metrics.Latencies))
	stdDev := math.Sqrt(variance)

	ehs.metrics.Jitter = stdDev
	return stdDev
}

// calculateMeanLatency calculates average latency
func (ehs *EnhancedHealthScorer) calculateMeanLatency() float64 {
	if len(ehs.metrics.Latencies) == 0 {
		return 0.0
	}

	sum := int64(0)
	for _, lat := range ehs.metrics.Latencies {
		sum += lat
	}

	return float64(sum) / float64(len(ehs.metrics.Latencies))
}

// CalculateStabilityScore calculates connection stability (0-100)
func (ehs *EnhancedHealthScorer) CalculateStabilityScore() float64 {
	if len(ehs.metrics.Latencies) < 3 {
		return 50.0 // Not enough data
	}

	// Factors affecting stability:
	// 1. Jitter (lower is better)
	// 2. Error rate (lower is better)
	// 3. Latency consistency (lower variance is better)
	// 4. Connection timing consistency

	jitter := ehs.CalculateJitter()
	mean := ehs.calculateMeanLatency()

	// Normalize jitter (0-100 scale, inverted)
	jitterScore := 100.0
	if mean > 0 {
		jitterRatio := jitter / mean
		if jitterRatio > 1.0 {
			jitterScore = 0.0
		} else {
			jitterScore = 100.0 * (1.0 - jitterRatio)
		}
	}

	// Error rate score (0-100 scale, inverted)
	errorScore := 100.0 * (1.0 - ehs.metrics.ErrorRate)

	// Latency consistency score
	consistencyScore := ehs.calculateConsistencyScore()

	// Combined stability score (weighted average)
	stabilityScore := (jitterScore * 0.4) + (errorScore * 0.4) + (consistencyScore * 0.2)

	ehs.metrics.StabilityScore = stabilityScore
	return stabilityScore
}

// calculateConsistencyScore calculates how consistent the latencies are
func (ehs *EnhancedHealthScorer) calculateConsistencyScore() float64 {
	if len(ehs.metrics.Latencies) < 3 {
		return 50.0
	}

	// Calculate coefficient of variation (CV)
	mean := ehs.calculateMeanLatency()
	stdDev := ehs.metrics.Jitter

	if mean == 0 {
		return 0.0
	}

	cv := stdDev / mean

	// Convert CV to score (lower CV = higher score)
	// CV > 1.0 is very inconsistent (score near 0)
	// CV < 0.1 is very consistent (score near 100)
	
	if cv > 1.0 {
		return 0.0
	}

	return 100.0 * (1.0 - cv)
}

// CalculateAdvancedHealthScore calculates comprehensive health score (0-100)
func (ehs *EnhancedHealthScorer) CalculateAdvancedHealthScore() int {
	if len(ehs.metrics.Latencies) == 0 {
		return 0
	}

	// Component scores (0-100 each)
	latencyScore := ehs.calculateLatencyScore()
	stabilityScore := ehs.CalculateStabilityScore()
	reliabilityScore := 100.0 * (1.0 - ehs.metrics.ErrorRate)
	dpiEvadScore := ehs.calculateDPIEvasionScore()
	performanceScore := ehs.calculatePerformanceScore()

	// Weighted combination
	weights := map[string]float64{
		"latency":     0.25,
		"stability":   0.25,
		"reliability": 0.20,
		"dpi_evasion": 0.20, // Critical for Iran
		"performance": 0.10,
	}

	// For Iran-optimized mode, increase DPI evasion weight
	if ehs.iranOptimized {
		weights["dpi_evasion"] = 0.30
		weights["latency"] = 0.20
		weights["stability"] = 0.20
	}

	finalScore := (latencyScore * weights["latency"]) +
		(stabilityScore * weights["stability"]) +
		(reliabilityScore * weights["reliability"]) +
		(dpiEvadScore * weights["dpi_evasion"]) +
		(performanceScore * weights["performance"])

	// Apply bonus/penalty modifiers
	finalScore = ehs.applyModifiers(finalScore)

	// Ensure score is in valid range
	if finalScore < 0 {
		finalScore = 0
	}
	if finalScore > 100 {
		finalScore = 100
	}

	return int(math.Round(finalScore))
}

// calculateLatencyScore scores based on latency (0-100)
func (ehs *EnhancedHealthScorer) calculateLatencyScore() float64 {
	mean := ehs.calculateMeanLatency()

	// Latency scoring (exponential decay)
	// < 50ms = 100 points
	// 100ms = ~82 points
	// 200ms = ~60 points
	// 500ms = ~30 points
	// 1000ms = ~15 points
	// > 2000ms = ~5 points

	if mean < 50 {
		return 100.0
	}

	score := 100.0 * math.Exp(-mean/500.0)
	return score
}

// calculateDPIEvasionScore scores DPI evasion capability
func (ehs *EnhancedHealthScorer) calculateDPIEvasionScore() float64 {
	// Factors:
	// 1. TLS handshake time (longer may indicate DPI inspection)
	// 2. Connection stability (unstable = likely being throttled)
	// 3. Packet timing patterns (detecting active DPI)

	score := 100.0

	// Penalize unusually long TLS handshakes (sign of DPI inspection)
	if ehs.metrics.TLSHandshakeTime > 0 {
		if ehs.metrics.TLSHandshakeTime > 3000 {
			score -= 30.0 // Very slow TLS = likely blocked/inspected
		} else if ehs.metrics.TLSHandshakeTime > 1500 {
			score -= 15.0
		} else if ehs.metrics.TLSHandshakeTime > 1000 {
			score -= 5.0
		}
	}

	// Bonus for consistent connection times (evading DPI successfully)
	if ehs.metrics.StabilityScore > 80 {
		score += 10.0
	}

	// Penalize high jitter (sign of active throttling)
	mean := ehs.calculateMeanLatency()
	if mean > 0 && ehs.metrics.Jitter > 0 {
		jitterRatio := ehs.metrics.Jitter / mean
		if jitterRatio > 0.5 {
			score -= 20.0 * jitterRatio
		}
	}

	// Analyze packet timing patterns for DPI signatures
	if ehs.detectDPIPatterns() {
		score -= 25.0 // DPI signature detected
		ehs.metrics.DPIFingerprint = "suspected"
	}

	if score < 0 {
		score = 0
	}
	return score
}

// detectDPIPatterns detects DPI interference patterns
func (ehs *EnhancedHealthScorer) detectDPIPatterns() bool {
	if len(ehs.metrics.Latencies) < 5 {
		return false
	}

	// Pattern 1: Sudden latency spikes (sign of active inspection)
	mean := ehs.calculateMeanLatency()
	for _, lat := range ehs.metrics.Latencies {
		if float64(lat) > mean*3.0 {
			return true // Suspicious spike
		}
	}

	// Pattern 2: Regular throttling pattern (every N packets)
	if len(ehs.metrics.Latencies) >= 10 {
		slow := 0
		for _, lat := range ehs.metrics.Latencies {
			if float64(lat) > mean*1.5 {
				slow++
			}
		}
		throttlingRatio := float64(slow) / float64(len(ehs.metrics.Latencies))
		if throttlingRatio > 0.3 && throttlingRatio < 0.7 {
			return true // Selective throttling pattern
		}
	}

	// Pattern 3: Timing-based DPI (consistent delay at specific intervals)
	if ehs.detectTimingBasedDPI() {
		return true
	}

	return false
}

// detectTimingBasedDPI detects timing-based DPI patterns
func (ehs *EnhancedHealthScorer) detectTimingBasedDPI() bool {
	if len(ehs.metrics.Timestamps) < 5 {
		return false
	}

	// Calculate inter-arrival times
	intervals := make([]float64, 0)
	for i := 1; i < len(ehs.metrics.Timestamps); i++ {
		interval := ehs.metrics.Timestamps[i].Sub(ehs.metrics.Timestamps[i-1]).Seconds()
		intervals = append(intervals, interval)
	}

	if len(intervals) < 3 {
		return false
	}

	// Check for suspicious regularity (sign of rate limiting)
	mean := 0.0
	for _, iv := range intervals {
		mean += iv
	}
	mean /= float64(len(intervals))

	variance := 0.0
	for _, iv := range intervals {
		diff := iv - mean
		variance += diff * diff
	}
	variance /= float64(len(intervals))

	// Very low variance in timing = suspicious (natural traffic is more random)
	if variance < 0.001 && mean > 0.1 {
		return true
	}

	return false
}

// calculatePerformanceScore scores overall performance
func (ehs *EnhancedHealthScorer) calculatePerformanceScore() float64 {
	score := 100.0

	// DNS resolution time
	if ehs.metrics.DNSResolutionTime > 500 {
		score -= 15.0
	} else if ehs.metrics.DNSResolutionTime > 200 {
		score -= 5.0
	}

	// TCP handshake time
	if ehs.metrics.TCPHandshakeTime > 1000 {
		score -= 15.0
	} else if ehs.metrics.TCPHandshakeTime > 500 {
		score -= 5.0
	}

	// Time to first byte
	if ehs.metrics.TTFBTime > 2000 {
		score -= 20.0
	} else if ehs.metrics.TTFBTime > 1000 {
		score -= 10.0
	}

	if score < 0 {
		score = 0
	}
	return score
}

// applyModifiers applies bonus/penalty modifiers to final score
func (ehs *EnhancedHealthScorer) applyModifiers(baseScore float64) float64 {
	score := baseScore

	// Bonus for very low latency and high stability
	mean := ehs.calculateMeanLatency()
	if mean < 100 && ehs.metrics.StabilityScore > 90 {
		score += 5.0 // Excellent connection bonus
	}

	// Penalty for high error rate
	if ehs.metrics.ErrorRate > 0.2 {
		score *= 0.7 // 30% penalty for unreliable connection
	}

	// Iran-specific modifiers
	if ehs.iranOptimized {
		// Bonus for successful DPI evasion indicators
		if ehs.metrics.DPIFingerprint == "" && ehs.metrics.StabilityScore > 75 {
			score += 10.0 // Likely evading DPI successfully
		}

		// Penalty for DPI detection
		if ehs.metrics.DPIFingerprint == "suspected" {
			score *= 0.8 // 20% penalty for suspected DPI interference
		}
	}

	return score
}

// GetMetricsSummary returns a formatted summary of all metrics
func (ehs *EnhancedHealthScorer) GetMetricsSummary() map[string]interface{} {
	mean := ehs.calculateMeanLatency()
	median := ehs.calculateMedianLatency()
	p95 := ehs.calculatePercentile(95)

	return map[string]interface{}{
		"mean_latency_ms":     mean,
		"median_latency_ms":   median,
		"p95_latency_ms":      p95,
		"jitter_ms":           ehs.metrics.Jitter,
		"packet_loss_pct":     ehs.metrics.PacketLoss,
		"error_rate_pct":      ehs.metrics.ErrorRate * 100,
		"stability_score":     ehs.metrics.StabilityScore,
		"dns_time_ms":         ehs.metrics.DNSResolutionTime,
		"tcp_handshake_ms":    ehs.metrics.TCPHandshakeTime,
		"tls_handshake_ms":    ehs.metrics.TLSHandshakeTime,
		"ttfb_ms":             ehs.metrics.TTFBTime,
		"measurements":        len(ehs.metrics.Latencies),
		"dpi_fingerprint":     ehs.metrics.DPIFingerprint,
	}
}

// calculateMedianLatency calculates median latency
func (ehs *EnhancedHealthScorer) calculateMedianLatency() float64 {
	if len(ehs.metrics.Latencies) == 0 {
		return 0.0
	}

	sorted := make([]int64, len(ehs.metrics.Latencies))
	copy(sorted, ehs.metrics.Latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return float64(sorted[mid-1]+sorted[mid]) / 2.0
	}
	return float64(sorted[mid])
}

// calculatePercentile calculates the specified percentile of latencies
func (ehs *EnhancedHealthScorer) calculatePercentile(percentile int) float64 {
	if len(ehs.metrics.Latencies) == 0 {
		return 0.0
	}

	sorted := make([]int64, len(ehs.metrics.Latencies))
	copy(sorted, ehs.metrics.Latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	index := int(math.Ceil(float64(len(sorted)) * float64(percentile) / 100.0)) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return float64(sorted[index])
}

// SaveHistoricalSnapshot saves current metrics to history
func (ehs *EnhancedHealthScorer) SaveHistoricalSnapshot() {
	snapshot := *ehs.metrics // Create a copy
	ehs.historicalData = append(ehs.historicalData, snapshot)

	// Keep only last N snapshots
	if len(ehs.historicalData) > ehs.maxHistorySize {
		ehs.historicalData = ehs.historicalData[1:]
	}
}

// AnalyzeTrend analyzes performance trends over time
func (ehs *EnhancedHealthScorer) AnalyzeTrend() string {
	if len(ehs.historicalData) < 3 {
		return "insufficient_data"
	}

	// Calculate trend in average latency
	recentMean := 0.0
	olderMean := 0.0
	
	halfPoint := len(ehs.historicalData) / 2
	
	for i := 0; i < halfPoint; i++ {
		if len(ehs.historicalData[i].Latencies) > 0 {
			sum := int64(0)
			for _, lat := range ehs.historicalData[i].Latencies {
				sum += lat
			}
			olderMean += float64(sum) / float64(len(ehs.historicalData[i].Latencies))
		}
	}
	olderMean /= float64(halfPoint)
	
	for i := halfPoint; i < len(ehs.historicalData); i++ {
		if len(ehs.historicalData[i].Latencies) > 0 {
			sum := int64(0)
			for _, lat := range ehs.historicalData[i].Latencies {
				sum += lat
			}
			recentMean += float64(sum) / float64(len(ehs.historicalData[i].Latencies))
		}
	}
	recentMean /= float64(len(ehs.historicalData) - halfPoint)
	
	// Determine trend
	if recentMean < olderMean*0.9 {
		return "improving"
	} else if recentMean > olderMean*1.1 {
		return "degrading"
	}
	
	return "stable"
}

// PrintDetailedMetrics prints comprehensive metrics report
func (ehs *EnhancedHealthScorer) PrintDetailedMetrics() {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           ADVANCED HEALTH METRICS REPORT                      ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	
	mean := ehs.calculateMeanLatency()
	median := ehs.calculateMedianLatency()
	p95 := ehs.calculatePercentile(95)
	jitter := ehs.CalculateJitter()
	stability := ehs.CalculateStabilityScore()
	health := ehs.CalculateAdvancedHealthScore()
	
	fmt.Printf("║  Latency (Mean/Median/P95):  %.1f / %.1f / %.1f ms\n", mean, median, p95)
	fmt.Printf("║  Jitter:                     %.2f ms\n", jitter)
	fmt.Printf("║  Stability Score:            %.1f%%\n", stability)
	fmt.Printf("║  Error Rate:                 %.2f%%\n", ehs.metrics.ErrorRate*100)
	fmt.Printf("║  Health Score:               %d/100\n", health)
	
	if ehs.metrics.DNSResolutionTime > 0 {
		fmt.Printf("║  DNS Resolution:             %d ms\n", ehs.metrics.DNSResolutionTime)
	}
	if ehs.metrics.TCPHandshakeTime > 0 {
		fmt.Printf("║  TCP Handshake:              %d ms\n", ehs.metrics.TCPHandshakeTime)
	}
	if ehs.metrics.TLSHandshakeTime > 0 {
		fmt.Printf("║  TLS Handshake:              %d ms\n", ehs.metrics.TLSHandshakeTime)
	}
	
	if ehs.metrics.DPIFingerprint != "" {
		fmt.Printf("║  DPI Status:                 %s\n", ehs.metrics.DPIFingerprint)
	}
	
	trend := ehs.AnalyzeTrend()
	fmt.Printf("║  Performance Trend:          %s\n", trend)
	
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}
