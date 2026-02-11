package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fatih/color"
)

// AdvancedAntiDPIClient represents the complete anti-DPI solution
type AdvancedAntiDPIClient struct {
	// Core components
	utlsDialer        *UTLSDialer
	sniFragmenter     *SNIFragmentDialer
	healthScorer      *EnhancedHealthScorer
	aiEngine          *AIAntiDPIEngine
	
	// Configuration
	timeout           time.Duration
	iranMode          bool
	adaptiveMode      bool
	
	// Statistics
	successfulConnections int64
	failedConnections     int64
	totalAttempts        int64
	mu                   sync.Mutex
	
	// Learning system
	connectionPatterns   map[string]*ConnectionPattern
	lastAdaptation       time.Time
	adaptationInterval   time.Duration
}

// ConnectionPattern stores learned patterns for successful connections
type ConnectionPattern struct {
	ProxyIP           string
	SuccessCount      int
	FailCount         int
	BestBrowser       string
	BestFragmentSize  int
	BestDelay         time.Duration
	AvgLatency        float64
	LastSuccess       time.Time
	DPISignature      string
}

// NewAdvancedAntiDPIClient creates a new advanced anti-DPI client
func NewAdvancedAntiDPIClient(timeout time.Duration, iranMode bool) *AdvancedAntiDPIClient {
	return &AdvancedAntiDPIClient{
		utlsDialer:         NewUTLSDialer(timeout),
		sniFragmenter:      NewSNIFragmentDialer(timeout),
		healthScorer:       NewEnhancedHealthScorer(iranMode),
		aiEngine:           NewAIAntiDPIEngine("advanced"),
		timeout:            timeout,
		iranMode:           iranMode,
		adaptiveMode:       true,
		connectionPatterns: make(map[string]*ConnectionPattern),
		lastAdaptation:     time.Now(),
		adaptationInterval: 5 * time.Minute,
	}
}

// CreateOptimizedTransport creates an HTTP transport with all anti-DPI features
func (ac *AdvancedAntiDPIClient) CreateOptimizedTransport(proxyURL string) (*http.Transport, error) {
	// Parse proxy URL
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %v", err)
	}

	// Get or create connection pattern for this proxy
	pattern := ac.getOrCreatePattern(parsedURL.Host)

	// Create custom dialer with all features
	customDialer := &AdvancedDialer{
		client:        ac,
		pattern:       pattern,
		utlsDialer:    ac.utlsDialer,
		sniFragmenter: ac.sniFragmenter,
		timeout:       ac.timeout,
	}

	// Create transport with custom dialer
	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedURL),
		DialContext: customDialer.DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
			// Use cipher suites from uTLS fingerprint
			CipherSuites: ac.utlsDialer.GetCurrentFingerprint().CipherSuites,
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   ac.timeout,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
		DisableKeepAlives:     false,
	}

	return transport, nil
}

// CreateOptimizedClient creates an HTTP client with all anti-DPI features
func (ac *AdvancedAntiDPIClient) CreateOptimizedClient(proxyURL string) (*http.Client, error) {
	transport, err := ac.CreateOptimizedTransport(proxyURL)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   ac.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	return client, nil
}

// TestProxyAdvanced tests a proxy with all anti-DPI features enabled
func (ac *AdvancedAntiDPIClient) TestProxyAdvanced(proxyURL string, testEndpoints []string) (*ProxyTestResult, error) {
	startTime := time.Now()
	
	// Create optimized client
	client, err := ac.CreateOptimizedClient(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	result := &ProxyTestResult{
		ProxyURL:      proxyURL,
		StartTime:     startTime,
		TestEndpoints: testEndpoints,
		Attempts:      make([]AttemptResult, 0),
	}

	// Test multiple endpoints
	for _, endpoint := range testEndpoints {
		attemptStart := time.Now()
		
		// Record pre-connection timing
		dnsStart := time.Now()
		
		// Make request with timing measurements
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			ac.recordFailure()
			continue
		}

		// Add realistic headers using AI engine
		ac.addOptimizedHeaders(req, endpoint)

		// Perform request
		resp, err := client.Do(req)
		
		attemptDuration := time.Since(attemptStart).Milliseconds()
		
		attempt := AttemptResult{
			Endpoint:  endpoint,
			StartTime: attemptStart,
			Duration:  attemptDuration,
			Success:   err == nil && resp != nil,
		}

		if err == nil && resp != nil {
			attempt.StatusCode = resp.StatusCode
			attempt.Success = resp.StatusCode >= 200 && resp.StatusCode < 400
			resp.Body.Close()
			
			if attempt.Success {
				ac.recordSuccess()
				ac.healthScorer.RecordLatency(attemptDuration, attemptStart)
			} else {
				ac.recordFailure()
			}
		} else {
			ac.recordFailure()
			ac.healthScorer.RecordError()
			attempt.Error = err
		}

		result.Attempts = append(result.Attempts, attempt)
		
		// Record detailed timing
		dnsTime := time.Since(dnsStart).Milliseconds()
		ac.healthScorer.RecordConnectionMetrics(dnsTime, 0, 0, attemptDuration)

		// Small delay between tests to appear more natural
		if ac.iranMode {
			time.Sleep(ac.aiEngine.ApplyTimingObfuscation())
		}
	}

	// Calculate final results
	result.EndTime = time.Now()
	result.TotalDuration = result.EndTime.Sub(startTime)
	result.SuccessCount = 0
	result.FailCount = 0

	for _, attempt := range result.Attempts {
		if attempt.Success {
			result.SuccessCount++
		} else {
			result.FailCount++
		}
	}

	result.SuccessRate = float64(result.SuccessCount) / float64(len(result.Attempts))
	result.HealthScore = ac.healthScorer.CalculateAdvancedHealthScore()
	result.Jitter = ac.healthScorer.CalculateJitter()
	result.StabilityScore = ac.healthScorer.CalculateStabilityScore()

	// Update connection pattern
	ac.updateConnectionPattern(proxyURL, result)

	// Adaptive learning
	if ac.adaptiveMode && time.Since(ac.lastAdaptation) > ac.adaptationInterval {
		ac.adaptStrategies()
	}

	return result, nil
}

// addOptimizedHeaders adds realistic headers based on target endpoint
func (ac *AdvancedAntiDPIClient) addOptimizedHeaders(req *http.Request, endpoint string) {
	// Determine target service
	targetService := "generic"
	if contains(endpoint, "cloudflare") {
		targetService = "cloudflare"
	} else if contains(endpoint, "google") {
		targetService = "google"
	} else if contains(endpoint, "microsoft") {
		targetService = "microsoft"
	}

	// Get AI-generated headers
	mimicry := ac.aiEngine.GenerateTrafficMimicry(targetService)

	// Apply headers
	for key, value := range mimicry {
		if strValue, ok := value.(string); ok {
			req.Header.Set(key, strValue)
		}
	}

	// Add standard headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}
	if req.Header.Get("Cache-Control") == "" {
		req.Header.Set("Cache-Control", "max-age=0")
	}
	if req.Header.Get("DNT") == "" {
		req.Header.Set("DNT", "1")
	}
	if req.Header.Get("Upgrade-Insecure-Requests") == "" {
		req.Header.Set("Upgrade-Insecure-Requests", "1")
	}
}

// getOrCreatePattern retrieves or creates a connection pattern
func (ac *AdvancedAntiDPIClient) getOrCreatePattern(proxyHost string) *ConnectionPattern {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if pattern, exists := ac.connectionPatterns[proxyHost]; exists {
		return pattern
	}

	pattern := &ConnectionPattern{
		ProxyIP:          proxyHost,
		BestBrowser:      "chrome120",
		BestFragmentSize: 5,
		BestDelay:        2 * time.Millisecond,
	}

	ac.connectionPatterns[proxyHost] = pattern
	return pattern
}

// updateConnectionPattern updates the learned pattern for a proxy
func (ac *AdvancedAntiDPIClient) updateConnectionPattern(proxyURL string, result *ProxyTestResult) {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return
	}

	pattern := ac.getOrCreatePattern(parsedURL.Host)
	
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if result.SuccessRate > 0.5 {
		pattern.SuccessCount++
		pattern.LastSuccess = time.Now()
		
		// Update best configuration
		pattern.BestBrowser = ac.utlsDialer.selectedBrowser
		
		// Calculate average latency
		totalLatency := int64(0)
		validAttempts := 0
		for _, attempt := range result.Attempts {
			if attempt.Success {
				totalLatency += attempt.Duration
				validAttempts++
			}
		}
		if validAttempts > 0 {
			pattern.AvgLatency = float64(totalLatency) / float64(validAttempts)
		}
	} else {
		pattern.FailCount++
	}

	// Update DPI signature if detected
	if ac.healthScorer.metrics.DPIFingerprint != "" {
		pattern.DPISignature = ac.healthScorer.metrics.DPIFingerprint
	}
}

// adaptStrategies adapts anti-DPI strategies based on learned patterns
func (ac *AdvancedAntiDPIClient) adaptStrategies() {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	totalSuccess := int64(0)
	totalFail := int64(0)

	for _, pattern := range ac.connectionPatterns {
		totalSuccess += int64(pattern.SuccessCount)
		totalFail += int64(pattern.FailCount)
	}

	if totalSuccess+totalFail < 10 {
		return // Not enough data
	}

	successRate := float64(totalSuccess) / float64(totalSuccess+totalFail)

	// Adapt based on overall success rate
	if successRate < 0.5 {
		// Low success - increase aggression
		color.Yellow("⚠️  Low success rate (%.1f%%) - Increasing anti-DPI aggression", successRate*100)
		ac.sniFragmenter.SetFragmentSize(3) // Smaller fragments
		ac.sniFragmenter.SetDelay(5 * time.Millisecond) // Longer delays
	} else if successRate > 0.85 {
		// High success - can optimize for performance
		color.Green("✅ High success rate (%.1f%%) - Optimizing for performance", successRate*100)
		ac.sniFragmenter.SetFragmentSize(6) // Larger fragments
		ac.sniFragmenter.SetDelay(1 * time.Millisecond) // Shorter delays
	}

	ac.lastAdaptation = time.Now()
}

// recordSuccess records a successful connection
func (ac *AdvancedAntiDPIClient) recordSuccess() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.successfulConnections++
	ac.totalAttempts++
}

// recordFailure records a failed connection
func (ac *AdvancedAntiDPIClient) recordFailure() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.failedConnections++
	ac.totalAttempts++
}

// GetStatistics returns comprehensive statistics
func (ac *AdvancedAntiDPIClient) GetStatistics() map[string]interface{} {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	successRate := 0.0
	if ac.totalAttempts > 0 {
		successRate = float64(ac.successfulConnections) / float64(ac.totalAttempts)
	}

	return map[string]interface{}{
		"total_attempts":         ac.totalAttempts,
		"successful_connections": ac.successfulConnections,
		"failed_connections":     ac.failedConnections,
		"success_rate":           successRate,
		"learned_patterns":       len(ac.connectionPatterns),
		"iran_mode":              ac.iranMode,
		"adaptive_mode":          ac.adaptiveMode,
		"utls_fingerprint":       ac.utlsDialer.GetFingerprintInfo(),
		"sni_fragmentation":      ac.sniFragmenter.GetFragmenter().GetStats(),
	}
}

// PrintStatistics prints detailed statistics
func (ac *AdvancedAntiDPIClient) PrintStatistics() {
	stats := ac.GetStatistics()
	
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║        ADVANCED ANTI-DPI CLIENT STATISTICS                    ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Total Attempts:         %-40d║", stats["total_attempts"])
	color.Green("║  Successful:             %-40d║", stats["successful_connections"])
	color.Green("║  Failed:                 %-40d║", stats["failed_connections"])
	color.Green("║  Success Rate:           %.1f%%%-37s║", stats["success_rate"].(float64)*100, "")
	color.Green("║  Learned Patterns:       %-40d║", stats["learned_patterns"])
	color.Green("║  Iran Mode:              %-40v║", stats["iran_mode"])
	color.Green("║  Adaptive Mode:          %-40v║", stats["adaptive_mode"])
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Browser Fingerprint:    %-40s║", stats["utls_fingerprint"].(string))
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
}

// ProxyTestResult contains comprehensive test results
type ProxyTestResult struct {
	ProxyURL        string
	StartTime       time.Time
	EndTime         time.Time
	TotalDuration   time.Duration
	TestEndpoints   []string
	Attempts        []AttemptResult
	SuccessCount    int
	FailCount       int
	SuccessRate     float64
	HealthScore     int
	Jitter          float64
	StabilityScore  float64
}

// AttemptResult contains results for a single test attempt
type AttemptResult struct {
	Endpoint   string
	StartTime  time.Time
	Duration   int64
	StatusCode int
	Success    bool
	Error      error
}

// AdvancedDialer combines all dialing features
type AdvancedDialer struct {
	client        *AdvancedAntiDPIClient
	pattern       *ConnectionPattern
	utlsDialer    *UTLSDialer
	sniFragmenter *SNIFragmentDialer
	timeout       time.Duration
}

// DialContext implements custom dialing with all features
func (ad *AdvancedDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// Apply learned optimizations from pattern
	if ad.pattern.BestBrowser != "" {
		ad.utlsDialer.SetBrowser(ad.pattern.BestBrowser)
	}

	if ad.pattern.BestFragmentSize > 0 {
		ad.sniFragmenter.SetFragmentSize(ad.pattern.BestFragmentSize)
	}

	if ad.pattern.BestDelay > 0 {
		ad.sniFragmenter.SetDelay(ad.pattern.BestDelay)
	}

	// Use SNI fragmentation dialer
	return ad.sniFragmenter.DialContext(ctx, network, addr)
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
		 (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		  findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
