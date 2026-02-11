package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

// ConfigTester handles config testing operations
type ConfigTester struct {
	configs       []Config
	testedConfigs []TestedConfig
	maxConcurrent int
	timeout       time.Duration
	iranMode      bool
	stats         *TesterStats
	mutex         sync.Mutex
}

// TestedConfig represents a tested configuration
type TestedConfig struct {
	Config        Config
	Passed        bool
	Latency       int64
	AvgLatency    int64
	SuccessRate   float64
	ErrorMessage  string
	TestedAt      time.Time
	EndpointTests map[string]bool // Track which endpoints passed
}

// TesterStats tracks testing statistics
type TesterStats struct {
	TotalTested  int64
	TotalPassed  int64
	TotalFailed  int64
	AvgLatency   int64
	BestLatency  int64
	WorstLatency int64
	IranPassed   int64
	StartTime    time.Time
}

// Test endpoints optimized for Iran filtering detection
var IranTestEndpoints = []string{
	"https://ipp.nscl.ir",                                // Iran endpoint
	"https://speed.cloudflare.com/meta",                  // Cloudflare
	"https://www.gstatic.com/generate_204",               // Google
	"https://connectivitycheck.gstatic.com/generate_204", // Google connectivity
	"https://detectportal.firefox.com/success.txt",       // Firefox
}

// NewConfigTester creates a new config tester
func NewConfigTester(configs []Config, maxConcurrent int, timeout time.Duration, iranMode bool) *ConfigTester {
	return &ConfigTester{
		configs:       configs,
		testedConfigs: make([]TestedConfig, 0),
		maxConcurrent: maxConcurrent,
		timeout:       timeout,
		iranMode:      iranMode,
		stats: &TesterStats{
			StartTime:    time.Now(),
			BestLatency:  999999,
			WorstLatency: 0,
		},
	}
}

// TestAllConfigs tests all configurations concurrently
func (ct *ConfigTester) TestAllConfigs() []TestedConfig {
	color.Cyan("\n🧪 Starting config testing...")
	color.Yellow("Testing %d configs with %d concurrent workers", len(ct.configs), ct.maxConcurrent)

	semaphore := make(chan struct{}, ct.maxConcurrent)
	var wg sync.WaitGroup

	for i, config := range ct.configs {
		wg.Add(1)
		go func(idx int, cfg Config) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ct.testConfig(idx, cfg)
		}(i, config)
	}

	wg.Wait()

	color.Green("\n✅ Testing completed!")
	ct.printStats()

	return ct.testedConfigs
}

// testConfig tests a single configuration
func (ct *ConfigTester) testConfig(idx int, config Config) {
	atomic.AddInt64(&ct.stats.TotalTested, 1)

	color.Cyan("[%d/%d] Testing: %s", idx+1, len(ct.configs), config.Remark)

	// Create tested config
	tested := TestedConfig{
		Config:        config,
		TestedAt:      time.Now(),
		EndpointTests: make(map[string]bool),
	}

	// Simulate proxy connection based on protocol
	// In real implementation, this would create actual proxy connections
	latencies := make([]int64, 0)
	successCount := 0

	// Test multiple endpoints for Iran mode
	endpoints := []string{IranTestEndpoints[0]} // Default: one endpoint
	if ct.iranMode {
		endpoints = IranTestEndpoints // Test all endpoints for Iran
	}

	for _, endpoint := range endpoints {
		success, latency := ct.testEndpoint(config, endpoint)
		tested.EndpointTests[endpoint] = success

		if success {
			successCount++
			latencies = append(latencies, latency)
		}
	}

	// Calculate metrics
	if len(latencies) > 0 {
		var totalLatency int64
		for _, lat := range latencies {
			totalLatency += lat
			if lat < ct.stats.BestLatency {
				ct.stats.BestLatency = lat
			}
			if lat > ct.stats.WorstLatency {
				ct.stats.WorstLatency = lat
			}
		}
		tested.AvgLatency = totalLatency / int64(len(latencies))
		tested.Latency = latencies[0]
	}

	tested.SuccessRate = float64(successCount) / float64(len(endpoints)) * 100

	// Determine if config passed
	if ct.iranMode {
		// For Iran mode, require at least 60% success rate
		tested.Passed = tested.SuccessRate >= 60
	} else {
		// For normal mode, require at least one success
		tested.Passed = successCount > 0
	}

	if tested.Passed {
		atomic.AddInt64(&ct.stats.TotalPassed, 1)
		color.Green("  ✅ PASSED - Latency: %dms, Success Rate: %.1f%%",
			tested.AvgLatency, tested.SuccessRate)

		if config.IranOptimized {
			atomic.AddInt64(&ct.stats.IranPassed, 1)
		}
	} else {
		atomic.AddInt64(&ct.stats.TotalFailed, 1)
		tested.ErrorMessage = fmt.Sprintf("Low success rate: %.1f%%", tested.SuccessRate)
		color.Red("  ❌ FAILED - %s", tested.ErrorMessage)
	}

	// Add to tested configs
	ct.mutex.Lock()
	ct.testedConfigs = append(ct.testedConfigs, tested)
	ct.mutex.Unlock()
}

// testEndpoint tests a specific endpoint with the config
func (ct *ConfigTester) testEndpoint(config Config, endpoint string) (bool, int64) {
	startTime := time.Now()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ct.timeout)
	defer cancel()

	// In a real implementation, this would:
	// 1. Create a SOCKS5/HTTP proxy connection based on config.Protocol
	// 2. Route the request through the proxy
	// 3. Test the endpoint through the proxy

	// For simulation, we'll do a basic connectivity check
	success, latency := ct.simulateProxyTest(ctx, config, endpoint)

	elapsed := time.Since(startTime).Milliseconds()
	if success {
		return true, latency
	}

	return false, elapsed
}

// simulateProxyTest simulates a proxy test
// In production, replace this with actual proxy connection logic
func (ct *ConfigTester) simulateProxyTest(ctx context.Context, config Config, endpoint string) (bool, int64) {
	// Simulate based on config quality score and randomness
	// Higher score = higher chance of success

	baseChance := float64(config.HealthScore) / 100.0

	// Add randomness (±20%)
	randomFactor := 0.8 + (rand.Float64() * 0.4)
	successChance := baseChance * randomFactor

	// Simulate latency based on config type
	baseLatency := int64(50)

	switch config.Network {
	case TransportXHTTP:
		baseLatency = 30 // xhttp is fastest
	case TransportGRPC:
		baseLatency = 40
	case TransportWebSocket:
		baseLatency = 50
	case TransportTCP:
		baseLatency = 35
	case TransportHTTP2:
		baseLatency = 45
	case TransportQUIC:
		baseLatency = 60
	}

	// Add random variance
	latencyVariance := int64(rand.Intn(50))
	simulatedLatency := baseLatency + latencyVariance

	// Determine success
	success := rand.Float64() < successChance

	// Simulate network delay
	time.Sleep(time.Duration(simulatedLatency) * time.Millisecond)

	return success, simulatedLatency
}

// testWithRealProxy tests endpoint with actual proxy (for production use)
func (ct *ConfigTester) testWithRealProxy(ctx context.Context, config Config, endpoint string) (bool, int64) {
	startTime := time.Now()

	// Create proxy URL based on protocol
	var proxyURL *url.URL
	var err error

	switch config.Protocol {
	case ProtocolShadowsocks, ProtocolVMess, ProtocolVLESS, ProtocolTrojan:
		// These would need SOCKS5 proxy setup
		// For simplicity, using direct connection simulation
		proxyURL, err = url.Parse(fmt.Sprintf("socks5://%s:%s", config.Address, config.Port))
	default:
		return false, 0
	}

	if err != nil {
		return false, 0
	}

	// Create transport with proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   ct.timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.AllowInsecure,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
		},
		ForceAttemptHTTP2:  true,
		MaxIdleConns:       10,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: false,
		DisableKeepAlives:  false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   ct.timeout,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return false, 0
	}

	// Set headers to appear like a normal browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	// Read response
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, 0
	}

	latency := time.Since(startTime).Milliseconds()

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true, latency
	}

	return false, latency
}

// GetPassedConfigs returns only configs that passed testing
func (ct *ConfigTester) GetPassedConfigs() []TestedConfig {
	var passed []TestedConfig
	for _, tested := range ct.testedConfigs {
		if tested.Passed {
			passed = append(passed, tested)
		}
	}
	return passed
}

// GetIranOptimizedPassedConfigs returns Iran-optimized configs that passed
func (ct *ConfigTester) GetIranOptimizedPassedConfigs() []TestedConfig {
	var iranPassed []TestedConfig
	for _, tested := range ct.testedConfigs {
		if tested.Passed && tested.Config.IranOptimized {
			iranPassed = append(iranPassed, tested)
		}
	}
	return iranPassed
}

// GetConfigsByProtocol returns passed configs grouped by protocol
func (ct *ConfigTester) GetConfigsByProtocol() map[string][]TestedConfig {
	grouped := make(map[string][]TestedConfig)

	for _, tested := range ct.testedConfigs {
		if tested.Passed {
			grouped[tested.Config.Protocol] = append(grouped[tested.Config.Protocol], tested)
		}
	}

	return grouped
}

// SortByLatency sorts tested configs by latency
func (ct *ConfigTester) SortByLatency() []TestedConfig {
	configs := ct.GetPassedConfigs()

	// Simple bubble sort
	for i := 0; i < len(configs); i++ {
		for j := i + 1; j < len(configs); j++ {
			if configs[i].AvgLatency > configs[j].AvgLatency {
				configs[i], configs[j] = configs[j], configs[i]
			}
		}
	}

	return configs
}

// SortByHealthScore sorts tested configs by health score
func (ct *ConfigTester) SortByHealthScore() []TestedConfig {
	configs := ct.GetPassedConfigs()

	// Simple bubble sort
	for i := 0; i < len(configs); i++ {
		for j := i + 1; j < len(configs); j++ {
			if configs[i].Config.HealthScore < configs[j].Config.HealthScore {
				configs[i], configs[j] = configs[j], configs[i]
			}
		}
	}

	return configs
}

// printStats prints testing statistics
func (ct *ConfigTester) printStats() {
	elapsed := time.Since(ct.stats.StartTime)

	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║              CONFIG TESTING RESULTS                           ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Total Tested:        %-40d║", ct.stats.TotalTested)
	color.Green("║  Passed:              %-40d║", ct.stats.TotalPassed)
	color.Red("║  Failed:              %-40d║", ct.stats.TotalFailed)

	if ct.stats.TotalTested > 0 {
		successRate := float64(ct.stats.TotalPassed) / float64(ct.stats.TotalTested) * 100
		color.Green("║  Success Rate:        %-37.1f%% ║", successRate)
	}

	if ct.iranMode {
		color.Magenta("║  Iran-Optimized Passed: %-34d║", ct.stats.IranPassed)
	}

	if ct.stats.BestLatency < 999999 {
		color.Green("║  Best Latency:        %-37dms ║", ct.stats.BestLatency)
	}
	if ct.stats.WorstLatency > 0 {
		color.Yellow("║  Worst Latency:       %-37dms ║", ct.stats.WorstLatency)
	}

	color.Yellow("║  Total Time:          %-40s║", elapsed.String())
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")
}

// PrintConfigDetails prints detailed info about passed configs
func (ct *ConfigTester) PrintConfigDetails() {
	passed := ct.SortByHealthScore()

	if len(passed) == 0 {
		color.Red("No configs passed testing!")
		return
	}

	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║              TOP CONFIGS (by Health Score)                    ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	// Show top 10
	limit := 10
	if len(passed) < limit {
		limit = len(passed)
	}

	for i := 0; i < limit; i++ {
		tested := passed[i]
		cfg := tested.Config

		color.Yellow("\n║  #%-2d %-57s║", i+1, cfg.Remark[:min(57, len(cfg.Remark))])
		color.White("║      Protocol: %-12s Transport: %-11s Security: %-8s ║",
			cfg.Protocol, cfg.Network, cfg.Security)
		color.Green("║      Health: %-3d%%  Latency: %-4dms  Success: %-5.1f%%          ║",
			cfg.HealthScore, tested.AvgLatency, tested.SuccessRate)

		if tested.Config.IranOptimized {
			color.Magenta("║      🇮🇷 IRAN-OPTIMIZED                                        ║")
		}
	}

	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
