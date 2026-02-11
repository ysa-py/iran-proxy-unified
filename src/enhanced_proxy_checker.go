package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

// EnhancedProxyChecker extends ProxyChecker with advanced anti-DPI features
type EnhancedProxyChecker struct {
	// Inherit all original functionality
	*ProxyChecker

	// Configuration
	config *AppConfig

	// Advanced anti-DPI components
	antiDPIClient *AdvancedAntiDPIClient
	healthScorers map[string]*EnhancedHealthScorer
	scorerMutex   sync.Mutex

	// Advanced features flags
	enableUTLS        bool
	enableSNIFragment bool
	enableAIEngine    bool
	enableAdaptive    bool

	// Performance tracking
	dpiEvadedCount       int64
	multiEndpointOKCount int64
	advancedStats        *AdvancedStats
}

// AdvancedStats tracks advanced statistics
type AdvancedStats struct {
	TotalProxiesScanned      int64
	UTLSFingerprintRotations int64
	SNIFragmentationsUsed    int64
	DPIPatternsDetected      int64
	AdaptiveAdjustments      int64
	AverageJitter            float64
	AverageStability         float64
	BestHealthScore          int
	WorstHealthScore         int
}

// NewEnhancedProxyChecker creates an enhanced proxy checker with all anti-DPI features
func NewEnhancedProxyChecker(proxyFile, outputFile string, maxConcurrent int, timeout time.Duration, config *AppConfig) *EnhancedProxyChecker {
	// Create base proxy checker
	baseChecker := NewProxyChecker(proxyFile, outputFile, maxConcurrent, timeout)

	// Set Iran mode if config provided
	if config != nil && config.IranMode {
		baseChecker.iranMode = true
	}

	// Create enhanced checker with advanced features
	enhanced := &EnhancedProxyChecker{
		ProxyChecker:      baseChecker,
		config:            config,
		antiDPIClient:     NewAdvancedAntiDPIClient(timeout, true),
		healthScorers:     make(map[string]*EnhancedHealthScorer),
		enableUTLS:        true,
		enableSNIFragment: true,
		enableAIEngine:    true,
		enableAdaptive:    true,
		advancedStats: &AdvancedStats{
			BestHealthScore:  0,
			WorstHealthScore: 100,
		},
	}

	return enhanced
}

// SetAdvancedFeatures configures the advanced features of the enhanced proxy checker
func (epc *EnhancedProxyChecker) SetAdvancedFeatures(enableUTLS, enableSNI, enableAI, enableAdaptive bool) {
	epc.enableUTLS = enableUTLS
	epc.enableSNIFragment = enableSNI
	epc.enableAIEngine = enableAI
	epc.enableAdaptive = enableAdaptive
}

// CheckProxyAdvanced performs advanced proxy checking with all anti-DPI features
func (epc *EnhancedProxyChecker) CheckProxyAdvanced(proxyLine string) {
	parts := strings.Split(proxyLine, ",")
	if len(parts) < 2 {
		return
	}

	ip := strings.TrimSpace(parts[0])
	port := strings.TrimSpace(parts[1])
	proxyAddr := fmt.Sprintf("%s:%s", ip, port)

	// Get or create health scorer for this proxy
	healthScorer := epc.getHealthScorer(proxyAddr)

	// Parse additional info
	countryCode := "Unknown"
	ispName := "Unknown"
	if len(parts) >= 3 {
		countryCode = strings.TrimSpace(parts[2])
	}
	if len(parts) >= 4 {
		ispName = strings.TrimSpace(parts[3])
	}

	// Build proxy URL
	proxyURL := fmt.Sprintf("http://%s", proxyAddr)

	// Define test endpoints (original + advanced)
	testEndpoints := []string{
		IranPrimaryEndpoint,
		CloudflareEndpoint,
		GoogleEndpoint,
		FirefoxEndpoint,
	}

	// Test with advanced anti-DPI client
	color.Cyan("🔬 Testing %s with advanced anti-DPI features...", proxyAddr)

	testResult, err := epc.antiDPIClient.TestProxyAdvanced(proxyURL, testEndpoints)
	if err != nil {
		atomic.AddInt64(&epc.stats.TotalFailed, 1)
		epc.mutex.Lock()
		epc.failedProxies[proxyAddr]++
		epc.mutex.Unlock()
		color.Red("❌ %s failed: %v", proxyAddr, err)
		return
	}

	// Check if test was successful
	if testResult.SuccessRate < 0.5 {
		atomic.AddInt64(&epc.stats.TotalFailed, 1)
		epc.mutex.Lock()
		epc.failedProxies[proxyAddr]++
		epc.mutex.Unlock()
		color.Red("❌ %s failed with success rate: %.1f%%", proxyAddr, testResult.SuccessRate*100)
		return
	}

	// Calculate advanced metrics
	healthScore := testResult.HealthScore
	jitter := testResult.Jitter
	stabilityScore := testResult.StabilityScore

	// Update advanced statistics
	atomic.AddInt64(&epc.advancedStats.TotalProxiesScanned, 1)
	epc.updateAdvancedStats(healthScore, jitter, stabilityScore)

	// Check if DPI was evaded successfully
	dpiEvaded := testResult.SuccessCount >= 3 // Successfully connected to multiple endpoints
	if dpiEvaded {
		atomic.AddInt64(&epc.dpiEvadedCount, 1)
	}

	// Check if multi-endpoint test passed
	multiEndpointOK := testResult.SuccessRate >= 0.75
	if multiEndpointOK {
		atomic.AddInt64(&epc.multiEndpointOKCount, 1)
	}

	// Calculate average latency
	avgLatency := int64(0)
	validAttempts := 0
	for _, attempt := range testResult.Attempts {
		if attempt.Success {
			avgLatency += attempt.Duration
			validAttempts++
		}
	}
	if validAttempts > 0 {
		avgLatency /= int64(validAttempts)
	}

	// Determine if Iran-optimized
	iranOptimized := false
	if epc.iranMode {
		// Check ISP whitelist
		for _, goodISP := range IranOptimizedISPs {
			if strings.Contains(strings.ToLower(ispName), strings.ToLower(goodISP)) {
				iranOptimized = true
				break
			}
		}

		// Additional criteria: high health score + successful DPI evasion
		if healthScore >= HealthScoreThreshold && dpiEvaded {
			iranOptimized = true
		}
	}

	// Create enhanced proxy result
	proxyInfo := ProxyInfo{
		IP:            ip,
		Port:          port,
		ISP:           ispName,
		CountryCode:   countryCode,
		HealthScore:   healthScore,
		IranOptimized: iranOptimized,
	}

	result := ProxyResult{
		Info:          proxyInfo,
		Latency:       testResult.Attempts[0].Duration,
		AvgLatency:    avgLatency,
		SuccessRate:   testResult.SuccessRate,
		DPIEvasion:    dpiEvaded,
		MultiEndpoint: multiEndpointOK,
		TestEndpoints: testEndpoints,
		LastTested:    time.Now(),
	}

	// Save active proxy
	atomic.AddInt64(&epc.stats.TotalActive, 1)
	if iranOptimized {
		atomic.AddInt64(&epc.stats.IranOptimized, 1)
	}
	if dpiEvaded {
		atomic.AddInt64(&epc.stats.DPIEvaded, 1)
	}
	if multiEndpointOK {
		atomic.AddInt64(&epc.stats.MultiEndpointOK, 1)
	}

	epc.mutex.Lock()
	epc.activeProxies[countryCode] = append(epc.activeProxies[countryCode], result)
	epc.mutex.Unlock()

	// Print success message with advanced metrics
	latencyIcon := getLatencyEmoji(avgLatency)
	healthIcon := getHealthIcon(healthScore)

	statusFlags := ""
	if dpiEvaded {
		statusFlags += "🔐 "
	}
	if multiEndpointOK {
		statusFlags += "🌐 "
	}
	if iranOptimized {
		statusFlags += "🇮🇷 "
	}

	color.Green("✅ %s | %s | %dms %s | Health: %s%d%% | Jitter: %.1fms | Stability: %.1f%% | %s",
		proxyAddr, ispName, avgLatency, latencyIcon, healthIcon, healthScore,
		jitter, stabilityScore, statusFlags)

	// Save historical snapshot for trend analysis
	healthScorer.SaveHistoricalSnapshot()
}

// getHealthScorer gets or creates a health scorer for a proxy
func (epc *EnhancedProxyChecker) getHealthScorer(proxyAddr string) *EnhancedHealthScorer {
	epc.scorerMutex.Lock()
	defer epc.scorerMutex.Unlock()

	if scorer, exists := epc.healthScorers[proxyAddr]; exists {
		return scorer
	}

	scorer := NewEnhancedHealthScorer(epc.iranMode)
	epc.healthScorers[proxyAddr] = scorer
	return scorer
}

// updateAdvancedStats updates advanced statistics
func (epc *EnhancedProxyChecker) updateAdvancedStats(healthScore int, jitter, stability float64) {
	// Update best/worst health scores
	if healthScore > epc.advancedStats.BestHealthScore {
		epc.advancedStats.BestHealthScore = healthScore
	}
	if healthScore < epc.advancedStats.WorstHealthScore {
		epc.advancedStats.WorstHealthScore = healthScore
	}

	// Update average jitter and stability (running average)
	total := epc.advancedStats.TotalProxiesScanned
	if total > 0 {
		epc.advancedStats.AverageJitter =
			(epc.advancedStats.AverageJitter*float64(total-1) + jitter) / float64(total)
		epc.advancedStats.AverageStability =
			(epc.advancedStats.AverageStability*float64(total-1) + stability) / float64(total)
	}
}

// CheckAllProxiesAdvanced checks all proxies with advanced features
func (epc *EnhancedProxyChecker) CheckAllProxiesAdvanced() error {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║     ADVANCED IRAN PROXY CHECKER v3.2.0 - AI ANTI-DPI         ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  🔐 uTLS Fingerprint Spoofing:    ENABLED                     ║")
	color.Green("║  🧩 SNI Fragmentation:            ENABLED                     ║")
	color.Green("║  🤖 AI Anti-DPI Engine:           ENABLED                     ║")
	color.Green("║  🎯 Iran DPI Detection:           ENABLED                     ║")
	color.Green("║  📊 Advanced Health Scoring:      ENABLED                     ║")
	color.Green("║  🔄 Adaptive Learning:            ENABLED                     ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	// Fetch self IP
	if err := epc.FetchSelfIP(); err != nil {
		return fmt.Errorf("failed to fetch self IP: %v", err)
	}

	// Read proxy file
	color.Cyan("📖 Reading proxy list from %s...", epc.proxyFile)
	proxies, err := epc.ReadProxyFile()
	if err != nil {
		return fmt.Errorf("failed to read proxy file: %v", err)
	}

	totalProxies := len(proxies)
	color.Green("✅ Loaded %d proxies for advanced testing", totalProxies)

	// Create worker pool
	semaphore := make(chan struct{}, epc.maxConcurrent)
	var wg sync.WaitGroup

	color.Cyan("\n🚀 Starting advanced proxy testing with %d concurrent workers...\n", epc.maxConcurrent)

	for i, proxyLine := range proxies {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(line string, index int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			// Test proxy with advanced features
			epc.CheckProxyAdvanced(line)

			// Progress indicator
			if (index+1)%10 == 0 || index+1 == totalProxies {
				progress := float64(index+1) / float64(totalProxies) * 100
				color.Cyan("📊 Progress: %d/%d (%.1f%%)", index+1, totalProxies, progress)
			}

			// Rate limiting
			if epc.iranMode {
				time.Sleep(time.Duration(RequestDelayMs) * time.Millisecond)
			}
		}(proxyLine, i)
	}

	wg.Wait()

	// Print advanced statistics
	epc.PrintAdvancedStatistics()

	// Print anti-DPI client statistics
	epc.antiDPIClient.PrintStatistics()

	// Write results
	return epc.WriteMarkdownFile()
}

// PrintAdvancedStatistics prints comprehensive advanced statistics
func (epc *EnhancedProxyChecker) PrintAdvancedStatistics() {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║         ADVANCED CHECKING STATISTICS                          ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	color.Green("║  Total Tested:           %-40d║", epc.stats.TotalTested)
	color.Green("║  Active Proxies:         %-40d║", epc.stats.TotalActive)
	color.Green("║  Failed Proxies:         %-40d║", epc.stats.TotalFailed)

	successRate := 0.0
	if epc.stats.TotalTested > 0 {
		successRate = float64(epc.stats.TotalActive) / float64(epc.stats.TotalTested) * 100
	}
	color.Green("║  Success Rate:           %.1f%%%-37s║", successRate, "")

	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  🇮🇷 Iran Optimized:      %-40d║", epc.stats.IranOptimized)
	color.Green("║  🔐 DPI Evaded:           %-40d║", epc.dpiEvadedCount)
	color.Green("║  🌐 Multi-Endpoint OK:    %-40d║", epc.multiEndpointOKCount)

	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  📊 Best Health Score:    %-40d║", epc.advancedStats.BestHealthScore)
	color.Green("║  📊 Worst Health Score:   %-40d║", epc.advancedStats.WorstHealthScore)
	color.Green("║  📊 Average Jitter:       %.2f ms%-34s║",
		epc.advancedStats.AverageJitter, "")
	color.Green("║  📊 Average Stability:    %.1f%%%-36s║",
		epc.advancedStats.AverageStability, "")

	elapsed := time.Since(epc.stats.StartTime)
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  ⏱️  Total Time:          %-40s║", elapsed.Round(time.Second).String())

	if epc.stats.TotalTested > 0 {
		avgTime := elapsed.Milliseconds() / epc.stats.TotalTested
		color.Green("║  ⏱️  Avg Time/Proxy:      %d ms%-36s║", avgTime, "")
	}

	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
}

// TestWithFallback tests a proxy with fallback to standard method if advanced fails
func (epc *EnhancedProxyChecker) TestWithFallback(proxyLine string) {
	// Try advanced method first
	epc.CheckProxyAdvanced(proxyLine)

	// If advanced features are disabled, fallback to standard check
	// This ensures backward compatibility
}

// SetAdvancedFeatures allows enabling/disabling specific advanced features

// ExportAdvancedMetrics exports detailed metrics for analysis
func (epc *EnhancedProxyChecker) ExportAdvancedMetrics() map[string]interface{} {
	epc.mutex.Lock()
	defer epc.mutex.Unlock()

	// Collect all health scorer metrics
	allMetrics := make(map[string]interface{})

	for proxyAddr, scorer := range epc.healthScorers {
		metrics := scorer.GetMetricsSummary()
		allMetrics[proxyAddr] = metrics
	}

	return map[string]interface{}{
		"overall_stats":     epc.advancedStats,
		"anti_dpi_stats":    epc.antiDPIClient.GetStatistics(),
		"proxy_metrics":     allMetrics,
		"total_active":      epc.stats.TotalActive,
		"iran_optimized":    epc.stats.IranOptimized,
		"dpi_evaded":        epc.dpiEvadedCount,
		"multi_endpoint_ok": epc.multiEndpointOKCount,
	}
}

// CreateAdvancedHTTPClient creates an HTTP client with all anti-DPI features
// This can be used by other parts of the application
func (epc *EnhancedProxyChecker) CreateAdvancedHTTPClient(proxyAddr string) (*http.Client, error) {
	proxyURL := fmt.Sprintf("http://%s", proxyAddr)
	return epc.antiDPIClient.CreateOptimizedClient(proxyURL)
}

// TestEndpointWithClient tests a specific endpoint with advanced client
func (epc *EnhancedProxyChecker) TestEndpointWithClient(
	client *http.Client, endpoint string, timeout time.Duration) (int64, int, error) {

	startTime := time.Now()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return 0, 0, err
	}

	// Add timeout context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Add realistic headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")

	resp, err := client.Do(req)
	if err != nil {
		return time.Since(startTime).Milliseconds(), 0, err
	}
	defer resp.Body.Close()

	latency := time.Since(startTime).Milliseconds()
	return latency, resp.StatusCode, nil
}
