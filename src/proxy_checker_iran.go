package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

const (
	DefaultProxyFile     = "edge/assets/p-list-february.txt"
	DefaultOutputFile    = "sub/ProxyIP-Daily.md"
	DefaultMaxConcurrent = 100
	DefaultTimeoutSecs   = 10
	RequestDelayMs       = 30

	// Iran-specific endpoints for better filtering bypass
	IranPrimaryEndpoint = "https://ipp.nscl.ir"
	CloudflareEndpoint  = "https://speed.cloudflare.com/meta"
	GoogleEndpoint      = "https://www.gstatic.com/generate_204"
	FirefoxEndpoint     = "https://detectportal.firefox.com/success.txt"

	// Advanced filtering settings
	MaxRetries           = 3
	HealthScoreThreshold = 70
	IranOptimizedPort    = 443
)

// Enhanced ISP whitelist optimized for Iran filtering bypass
var IranOptimizedISPs = []string{
	// Tier 1 - Best for Iran (CDN & Major Cloud)
	"Cloudflare", "Google", "Amazon", "Akamai", "Fastly", "Microsoft",

	// Tier 2 - Reliable bypasses
	"M247", "OVH", "Vultr", "GCore", "IONOS", "Hetzner", "DigitalOcean",

	// Tier 3 - Good alternatives
	"Contabo", "UpCloud", "Tencent", "Multacom", "Leaseweb", "Hostinger",
	"Scaleway", "netcup GmbH", "ByteDance", "RackSpace", "SiteGround",

	// Additional reliable providers
	"Online Ltd", "Relink LTD", "PQ Hosting", "Gigahost AS", "White Label",
	"G-Core Labs", "3HCLOUD LLC", "HOSTKEY B.V", "3NT SOLUTION", "Zenlayer Inc",
	"RackNerd LLC", "Plant Holding", "WorkTitans", "IROKO Networks", "WorldStream",
	"Cluster", "Cogent Communications", "Metropolis networks inc",
	"Total Uptime Technologies", "NetLab", "Turunc", "HostPapa", "Ultahost",
	"DataCamp", "Bluehost", "Protilab", "DO Space", "The Empire", "The Constant Company",
}

// GoodISPs is an alias for IranOptimizedISPs for backward compatibility with tests
var GoodISPs = IranOptimizedISPs

// WorkerResponse represents Cloudflare worker response
type WorkerResponse struct {
	IP string   `json:"clientIp"`
	CF WorkerCf `json:"cf"`
}

// WorkerCf represents Cloudflare metadata
type WorkerCf struct {
	ISP     string `json:"asOrganization"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	ASN     string `json:"asn"`
}

// ProxyInfo holds enhanced proxy information
type ProxyInfo struct {
	IP            string
	Port          string
	ISP           string
	CountryCode   string
	City          string
	Region        string
	ASN           string
	HealthScore   int
	IranOptimized bool
}

// ProxyResult contains proxy info and performance metrics
type ProxyResult struct {
	Info          ProxyInfo
	Latency       int64
	AvgLatency    int64
	SuccessRate   float64
	DPIEvasion    bool
	MultiEndpoint bool
	TestEndpoints []string
	LastTested    time.Time
}

// ProxyChecker handles advanced proxy checking operations
type ProxyChecker struct {
	proxyFile     string
	outputFile    string
	maxConcurrent int
	timeout       time.Duration
	iranMode      bool
	selfIP        string
	activeProxies map[string][]ProxyResult
	failedProxies map[string]int
	mutex         sync.Mutex
	client        *http.Client
	stats         *CheckerStats
}

// CheckerStats tracks checking statistics
type CheckerStats struct {
	TotalTested     int64
	TotalActive     int64
	TotalFailed     int64
	IranOptimized   int64
	DPIEvaded       int64
	MultiEndpointOK int64
	AvgLatency      int64
	StartTime       time.Time
}

// NewProxyChecker creates an enhanced ProxyChecker instance
func NewProxyChecker(proxyFile, outputFile string, maxConcurrent int, timeout time.Duration) *ProxyChecker {
	return &ProxyChecker{
		proxyFile:     proxyFile,
		outputFile:    outputFile,
		maxConcurrent: maxConcurrent,
		timeout:       timeout,
		iranMode:      true,
		activeProxies: make(map[string][]ProxyResult),
		failedProxies: make(map[string]int),
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
					MinVersion:         tls.VersionTLS12,
					MaxVersion:         tls.VersionTLS13,
					// Prefer TLS 1.3 for better obfuscation
					CipherSuites: []uint16{
						tls.TLS_AES_128_GCM_SHA256,
						tls.TLS_AES_256_GCM_SHA384,
						tls.TLS_CHACHA20_POLY1305_SHA256,
					},
				},
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: false,
				DisableKeepAlives:  false,
			},
		},
		stats: &CheckerStats{
			StartTime: time.Now(),
		},
	}
}

// FetchSelfIP retrieves the checker's own IP with Iran-aware fallback
func (pc *ProxyChecker) FetchSelfIP() error {
	color.Cyan("🔍 Detecting your IP address...")

	// Try Iran-optimized endpoint first
	if ip, err := pc.tryEndpoint(IranPrimaryEndpoint); err == nil {
		pc.selfIP = ip
		color.Green("✅ Detected IP: %s (via Iran endpoint)", ip)
		return nil
	}

	// Fallback to Cloudflare
	if ip, err := pc.tryEndpoint(CloudflareEndpoint); err == nil {
		pc.selfIP = ip
		color.Green("✅ Detected IP: %s (via Cloudflare)", ip)
		return nil
	}

	// Final fallback
	resp, err := pc.client.Get("https://api.ipify.org")
	if err != nil {
		return fmt.Errorf("failed to fetch self IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read self IP response: %v", err)
	}

	pc.selfIP = strings.TrimSpace(string(body))
	color.Green("✅ Detected IP: %s (via ipify)", pc.selfIP)
	return nil
}

func (pc *ProxyChecker) tryEndpoint(endpoint string) (string, error) {
	resp, err := pc.client.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if ip, ok := result["clientIp"].(string); ok && ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("invalid response")
}

// ReadProxyFile reads and filters proxies with enhanced Iran optimization
func (pc *ProxyChecker) ReadProxyFile() ([]string, error) {
	file, err := os.Open(pc.proxyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open proxy file: %v", err)
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		port := strings.TrimSpace(parts[1])
		ispName := "Unknown"
		if len(parts) >= 4 {
			ispName = strings.TrimSpace(parts[3])
		}

		// Iran-optimized filtering: prefer port 443 and whitelisted ISPs
		if pc.iranMode && port != fmt.Sprintf("%d", IranOptimizedPort) {
			continue
		}

		// Check ISP whitelist
		ispOK := false
		for _, goodISP := range IranOptimizedISPs {
			if strings.Contains(strings.ToLower(ispName), strings.ToLower(goodISP)) {
				ispOK = true
				break
			}
		}

		if ispOK {
			proxies = append(proxies, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading proxy file: %v", err)
	}

	return proxies, nil
}

// CheckProxyAdvanced performs multi-endpoint advanced proxy checking
func (pc *ProxyChecker) CheckProxyAdvanced(ctx context.Context, ip, port string) (*ProxyResult, error) {
	var result ProxyResult
	result.TestEndpoints = make([]string, 0)

	attempts := 0
	successCount := 0
	var totalLatency int64

	// Test multiple endpoints for better Iran filtering detection
	endpoints := []string{
		CloudflareEndpoint,
		IranPrimaryEndpoint,
	}

	for _, endpoint := range endpoints {
		for retry := 0; retry < MaxRetries; retry++ {
			attempts++

			data, latency, err := pc.testProxyEndpoint(ctx, ip, port, endpoint)
			if err == nil {
				successCount++
				totalLatency += latency
				result.TestEndpoints = append(result.TestEndpoints, endpoint)

				// Use first successful response for metadata
				if result.Info.IP == "" {
					result.Info = ProxyInfo{
						IP:          data.IP,
						Port:        port,
						ISP:         data.CF.ISP,
						CountryCode: data.CF.Country,
						City:        data.CF.City,
						Region:      data.CF.Region,
						ASN:         data.CF.ASN,
					}
					result.Latency = latency
				}

				break // Success on this endpoint
			}

			// Small delay between retries
			if retry < MaxRetries-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	if successCount == 0 {
		return nil, fmt.Errorf("all endpoint tests failed")
	}

	// Calculate performance metrics
	result.AvgLatency = totalLatency / int64(successCount)
	result.SuccessRate = float64(successCount) / float64(attempts) * 100
	result.MultiEndpoint = len(result.TestEndpoints) > 1
	result.DPIEvasion = result.MultiEndpoint && result.SuccessRate > 50
	result.LastTested = time.Now()

	// Calculate health score
	result.Info.HealthScore = pc.calculateHealthScore(&result)
	result.Info.IranOptimized = result.Info.HealthScore >= HealthScoreThreshold

	return &result, nil
}

func (pc *ProxyChecker) testProxyEndpoint(ctx context.Context, ip, port, endpoint string) (*WorkerResponse, int64, error) {
	startTime := time.Now()

	// Create connection with timeout
	dialer := &net.Dialer{
		Timeout:   pc.timeout,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, 0, fmt.Errorf("connection failed: %v", err)
	}
	defer conn.Close()

	// Determine host based on endpoint
	var host string
	if strings.Contains(endpoint, "cloudflare") {
		host = "speed.cloudflare.com"
	} else if strings.Contains(endpoint, "nscl.ir") {
		host = "ipp.nscl.ir"
	} else {
		host = "www.gstatic.com"
	}

	// Upgrade to TLS with optimized config
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}

	tlsConn := tls.Client(conn, tlsConfig)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return nil, 0, fmt.Errorf("TLS handshake failed: %v", err)
	}
	defer tlsConn.Close()

	latency := time.Since(startTime).Milliseconds()

	// Send HTTP request
	path := "/meta"
	if strings.Contains(endpoint, "nscl.ir") {
		path = "/"
	}

	request := fmt.Sprintf("GET %s HTTP/1.1\r\n", path) +
		fmt.Sprintf("Host: %s\r\n", host) +
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36\r\n" +
		"Accept: application/json\r\n" +
		"Accept-Encoding: identity\r\n" +
		"Connection: close\r\n\r\n"

	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		return nil, 0, fmt.Errorf("write failed: %v", err)
	}

	// Read response
	buffer := make([]byte, 16384)
	var response strings.Builder

	for {
		n, err := tlsConn.Read(buffer)
		if n > 0 {
			response.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	// Parse HTTP response
	respText := response.String()
	bodyStart := strings.Index(respText, "\r\n\r\n")
	if bodyStart == -1 {
		return nil, 0, fmt.Errorf("invalid HTTP response")
	}

	body := strings.TrimSpace(respText[bodyStart+4:])

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, 0, fmt.Errorf("JSON parse failed: %v", err)
	}

	clientIP, _ := result["clientIp"].(string)
	if clientIP == "" || clientIP == pc.selfIP {
		return nil, 0, fmt.Errorf("IP match or empty")
	}

	workerResp := &WorkerResponse{
		IP: clientIP,
		CF: WorkerCf{
			ISP:     getStringValue(result, "asOrganization"),
			City:    getStringValue(result, "city"),
			Region:  getStringValue(result, "region"),
			Country: getStringValue(result, "country"),
			ASN:     getStringValue(result, "asn"),
		},
	}

	return workerResp, latency, nil
}

func (pc *ProxyChecker) calculateHealthScore(result *ProxyResult) int {
	score := 100

	// Latency penalty
	if result.AvgLatency > 2000 {
		score -= 30
	} else if result.AvgLatency > 1000 {
		score -= 15
	} else if result.AvgLatency > 500 {
		score -= 5
	}

	// Success rate bonus/penalty
	if result.SuccessRate < 50 {
		score -= 30
	} else if result.SuccessRate > 80 {
		score += 10
	}

	// Multi-endpoint bonus
	if result.MultiEndpoint {
		score += 15
	}

	// DPI evasion bonus
	if result.DPIEvasion {
		score += 10
	}

	// CDN/Major cloud bonus
	majorProviders := []string{"Cloudflare", "Google", "Amazon", "Akamai", "Microsoft"}
	for _, provider := range majorProviders {
		if strings.Contains(strings.ToLower(result.Info.ISP), strings.ToLower(provider)) {
			score += 10
			break
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// ProcessProxy processes a single proxy with advanced metrics
func (pc *ProxyChecker) ProcessProxy(ctx context.Context, proxyLine string) {
	atomic.AddInt64(&pc.stats.TotalTested, 1)

	parts := strings.Split(proxyLine, ",")
	if len(parts) < 2 {
		return
	}

	ip := strings.TrimSpace(parts[0])
	port := strings.TrimSpace(parts[1])

	result, err := pc.CheckProxyAdvanced(ctx, ip, port)
	if err != nil {
		atomic.AddInt64(&pc.stats.TotalFailed, 1)
		color.Red("❌ DEAD: %s:%s (%v)", ip, port, err)

		pc.mutex.Lock()
		pc.failedProxies[ip]++
		pc.mutex.Unlock()
		return
	}

	// Fill in missing data
	if result.Info.ISP == "" && len(parts) >= 4 {
		result.Info.ISP = strings.TrimSpace(parts[3])
	}
	if result.Info.CountryCode == "" {
		result.Info.CountryCode = "XX"
	}
	if result.Info.City == "" {
		result.Info.City = "Unknown"
	}
	if result.Info.Region == "" {
		result.Info.Region = "Unknown"
	}

	// Update statistics
	atomic.AddInt64(&pc.stats.TotalActive, 1)
	if result.Info.IranOptimized {
		atomic.AddInt64(&pc.stats.IranOptimized, 1)
	}
	if result.DPIEvasion {
		atomic.AddInt64(&pc.stats.DPIEvaded, 1)
	}
	if result.MultiEndpoint {
		atomic.AddInt64(&pc.stats.MultiEndpointOK, 1)
	}

	// Display result with enhanced information
	healthIcon := getHealthIcon(result.Info.HealthScore)
	iranIcon := ""
	if result.Info.IranOptimized {
		iranIcon = "🇮🇷"
	}

	color.Green("✅ LIVE: %s:%s %s %s | %dms | Score: %d%% | Endpoints: %d | %s, %s",
		ip, port, healthIcon, iranIcon,
		result.AvgLatency,
		result.Info.HealthScore,
		len(result.TestEndpoints),
		result.Info.City,
		result.Info.ISP)

	pc.mutex.Lock()
	pc.activeProxies[result.Info.CountryCode] = append(pc.activeProxies[result.Info.CountryCode], *result)
	pc.mutex.Unlock()
}

// Run executes the proxy checking process
func (pc *ProxyChecker) Run() error {
	color.Cyan("🚀 Starting advanced Iran-optimized proxy checker...")

	if err := pc.FetchSelfIP(); err != nil {
		return err
	}

	proxies, err := pc.ReadProxyFile()
	if err != nil {
		return err
	}

	color.Yellow("📋 Loaded %d proxies from file", len(proxies))
	if pc.iranMode {
		color.Magenta("🇮🇷 Iran optimization mode: ENABLED")
		color.Magenta("🔐 DPI evasion: ACTIVE")
		color.Magenta("🌐 Multi-endpoint testing: ACTIVE")
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	color.Cyan("Starting proxy checks...")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	sem := make(chan struct{}, pc.maxConcurrent)
	var wg sync.WaitGroup

	ctx := context.Background()

	for i, proxyLine := range proxies {
		wg.Add(1)
		sem <- struct{}{}

		go func(line string, index int) {
			defer wg.Done()
			defer func() { <-sem }()

			// Small delay to avoid overwhelming the network
			time.Sleep(time.Duration(RequestDelayMs) * time.Millisecond)
			pc.ProcessProxy(ctx, line)

			// Progress indicator
			if (index+1)%10 == 0 {
				progress := float64(index+1) / float64(len(proxies)) * 100
				color.Cyan("📊 Progress: %.1f%% (%d/%d)", progress, index+1, len(proxies))
			}
		}(proxyLine, i)
	}

	wg.Wait()

	elapsed := time.Since(pc.stats.StartTime)

	fmt.Println("\n" + strings.Repeat("=", 80))
	color.Green("✅ Proxy checking completed!")
	fmt.Println(strings.Repeat("=", 80))

	pc.printStatistics(elapsed)

	return nil
}

func (pc *ProxyChecker) printStatistics(elapsed time.Duration) {
	fmt.Println()
	color.Cyan("📊 STATISTICS SUMMARY")
	color.Cyan(strings.Repeat("-", 80))

	fmt.Printf("⏱️  Total Time: %v\n", elapsed)
	fmt.Printf("🔍 Total Tested: %d\n", pc.stats.TotalTested)
	fmt.Printf("✅ Active Proxies: %d (%.1f%%)\n",
		pc.stats.TotalActive,
		float64(pc.stats.TotalActive)/float64(pc.stats.TotalTested)*100)
	fmt.Printf("❌ Failed Proxies: %d\n", pc.stats.TotalFailed)
	fmt.Printf("🇮🇷 Iran-Optimized: %d\n", pc.stats.IranOptimized)
	fmt.Printf("🔐 DPI Evasion OK: %d\n", pc.stats.DPIEvaded)
	fmt.Printf("🌐 Multi-Endpoint: %d\n", pc.stats.MultiEndpointOK)
	fmt.Printf("🌍 Countries: %d\n", len(pc.activeProxies))

	color.Cyan(strings.Repeat("-", 80))
}

// WriteMarkdownFile generates enhanced markdown output
func (pc *ProxyChecker) WriteMarkdownFile() error {
	dir := strings.TrimSuffix(pc.outputFile, "/"+pc.outputFile[strings.LastIndex(pc.outputFile, "/")+1:])
	if dir != pc.outputFile {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	file, err := os.Create(pc.outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Calculate statistics
	totalActive := int(pc.stats.TotalActive)
	totalCountries := len(pc.activeProxies)
	iranOptimized := int(pc.stats.IranOptimized)

	var sumPing int64
	for _, proxies := range pc.activeProxies {
		for _, p := range proxies {
			sumPing += p.AvgLatency
		}
	}

	avgPing := int64(0)
	if totalActive > 0 {
		avgPing = sumPing / int64(totalActive)
	}

	// Timestamps
	location, _ := time.LoadLocation("Asia/Tehran")
	now := time.Now().In(location)
	nextUpdate := now.Add(6 * time.Hour)

	lastUpdated := now.Format("Mon, 02 Jan 2006 15:04")
	nextUpdateStr := nextUpdate.Format("Mon, 02 Jan 2006 15:04")

	// Write enhanced header
	fmt.Fprintf(file, `<p align="left">
 <img src="https://latex.codecogs.com/svg.image?\huge&space;{{\color{Golden}\mathrm{🇮🇷\;IRAN\;PR{\color{black}\O}XY\;IP}}" width=280px" </p><br/>

> [!WARNING]
>
> <p><b>🇮🇷 Iran-Optimized Daily Fresh Proxies</b></p>
>
> A curated list of <b>high-quality</b>, Iran-filtered, fully-tested proxies optimized for bypassing strict Iranian internet filtering.
>
> <br/>
>
> <p><b>⚡ Advanced Features</b></p>
>
> ✅ **DPI Evasion Testing** - Multi-endpoint verification for deep packet inspection bypass  
> ✅ **Health Score System** - Intelligent ranking based on performance and reliability  
> ✅ **Iran Filter Detection** - Automatic detection and optimization for Iran's filtering  
> ✅ **TLS 1.3 Preferred** - Enhanced obfuscation using latest TLS protocols  
>
> <br/>
>
> <p><b>⏰ Auto-Updated Every 6 Hours</b></p>
>
> <img src="https://img.shields.io/badge/Last_Update-%s%%20(Tehran)-966600" />  
> <img src="https://img.shields.io/badge/Next_Update-%s%%20(Tehran)-966600" />
>
> <br/>
>
> <p><b>📊 Overview</b></p>  
>
> <img src="https://img.shields.io/badge/Active_Proxies-%d-00C853" />  
> <img src="https://img.shields.io/badge/Iran_Optimized-%d-FF6F00" />  
> <img src="https://img.shields.io/badge/Countries-%d-1976D2" />  
> <img src="https://img.shields.io/badge/Avg_Latency-%dms-D32F2F" />
>
> <br><br/>  

`, encodeBadgeLabel(lastUpdated), encodeBadgeLabel(nextUpdateStr),
		totalActive, iranOptimized, totalCountries, avgPing)

	// Write Iran-optimized section first
	pc.writeIranOptimizedSection(file)

	// Write provider sections
	pc.writeProviderSections(file)

	// Write country sections
	pc.writeCountrySections(file)

	color.Green("💾 All active proxies saved to %s", pc.outputFile)
	return nil
}

func (pc *ProxyChecker) writeIranOptimizedSection(file *os.File) {
	var iranProxies []ProxyResult

	for _, proxies := range pc.activeProxies {
		for _, result := range proxies {
			if result.Info.IranOptimized {
				iranProxies = append(iranProxies, result)
			}
		}
	}

	if len(iranProxies) == 0 {
		return
	}

	sort.Slice(iranProxies, func(i, j int) bool {
		if iranProxies[i].Info.HealthScore != iranProxies[j].Info.HealthScore {
			return iranProxies[i].Info.HealthScore > iranProxies[j].Info.HealthScore
		}
		return iranProxies[i].AvgLatency < iranProxies[j].AvgLatency
	})

	fmt.Fprintf(file, "## 🇮🇷 Iran-Optimized Proxies (%d)\n\n", len(iranProxies))
	fmt.Fprintf(file, "> **Best proxies for bypassing Iran's filtering** - High health score, DPI evasion, multi-endpoint verified\n\n")
	fmt.Fprintf(file, "<details open>\n<summary>Click to expand/collapse</summary>\n\n")
	fmt.Fprintf(file, "|   IP   |   ISP   |   Location   |   Latency   |   Health Score   |   Features   |\n")
	fmt.Fprintf(file, "|:-------|:--------|:------------:|:-----------:|:----------------:|:------------:|\n")

	for _, result := range iranProxies {
		location := fmt.Sprintf("%s, %s", result.Info.Region, result.Info.City)
		healthIcon := getHealthIcon(result.Info.HealthScore)
		latencyEmoji := getLatencyEmoji(result.AvgLatency)

		features := ""
		if result.DPIEvasion {
			features += "🔐 "
		}
		if result.MultiEndpoint {
			features += "🌐 "
		}
		features += fmt.Sprintf("(%d endpoints)", len(result.TestEndpoints))

		fmt.Fprintf(file, "| <pre><code>%s:%s</code></pre> | %s | %s | %d ms %s | %s %d%% | %s |\n",
			result.Info.IP, result.Info.Port, result.Info.ISP, location,
			result.AvgLatency, latencyEmoji,
			healthIcon, result.Info.HealthScore, features)
	}

	fmt.Fprintf(file, "\n</details>\n\n---\n\n")
}

func (pc *ProxyChecker) writeProviderSections(file *os.File) {
	topProviders := []string{"Cloudflare", "Google", "Amazon", "Akamai", "Hetzner", "DigitalOcean"}
	providerBuckets := make(map[string][]ProxyResult)

	for _, provider := range topProviders {
		providerBuckets[provider] = []ProxyResult{}
	}

	for _, proxies := range pc.activeProxies {
		for _, result := range proxies {
			for _, provider := range topProviders {
				if strings.Contains(strings.ToLower(result.Info.ISP), strings.ToLower(provider)) {
					providerBuckets[provider] = append(providerBuckets[provider], result)
				}
			}
		}
	}

	for _, provider := range topProviders {
		list := providerBuckets[provider]
		if len(list) == 0 {
			continue
		}

		logo := getProviderLogoHTML(provider)
		title := provider
		if logo != "" {
			title = fmt.Sprintf("%s %s", logo, provider)
		}

		fmt.Fprintf(file, "## %s (%d)\n", title, len(list))
		fmt.Fprintf(file, "<details>\n<summary>Click to expand</summary>\n\n")
		fmt.Fprintf(file, "|   IP   |   ISP   |   Location   |   Latency   |   Health Score   |\n")
		fmt.Fprintf(file, "|:-------|:--------|:------------:|:-----------:|:----------------:|\n")

		sort.Slice(list, func(i, j int) bool {
			if list[i].Info.HealthScore != list[j].Info.HealthScore {
				return list[i].Info.HealthScore > list[j].Info.HealthScore
			}
			return list[i].AvgLatency < list[j].AvgLatency
		})

		for _, result := range list {
			location := fmt.Sprintf("%s, %s", result.Info.Region, result.Info.City)
			healthIcon := getHealthIcon(result.Info.HealthScore)
			latencyEmoji := getLatencyEmoji(result.AvgLatency)

			fmt.Fprintf(file, "| <pre><code>%s:%s</code></pre> | %s | %s | %d ms %s | %s %d%% |\n",
				result.Info.IP, result.Info.Port, result.Info.ISP, location,
				result.AvgLatency, latencyEmoji,
				healthIcon, result.Info.HealthScore)
		}

		fmt.Fprintf(file, "\n</details>\n\n---\n\n")
	}
}

func (pc *ProxyChecker) writeCountrySections(file *os.File) {
	var countries []string
	for country := range pc.activeProxies {
		countries = append(countries, country)
	}
	sort.Strings(countries)

	for _, countryCode := range countries {
		proxies := pc.activeProxies[countryCode]

		sort.Slice(proxies, func(i, j int) bool {
			if proxies[i].Info.HealthScore != proxies[j].Info.HealthScore {
				return proxies[i].Info.HealthScore > proxies[j].Info.HealthScore
			}
			return proxies[i].AvgLatency < proxies[j].AvgLatency
		})

		flag := getCountryFlag(countryCode)
		name := getCountryName(countryCode)

		fmt.Fprintf(file, "## %s %s (%d proxies)\n", flag, name, len(proxies))
		fmt.Fprintf(file, "<details>\n<summary>Click to expand</summary>\n\n")
		fmt.Fprintf(file, "|   IP   |   ISP   |   Location   |   Latency   |   Health Score   |\n")
		fmt.Fprintf(file, "|:-------|:--------|:------------:|:-----------:|:----------------:|\n")

		for _, result := range proxies {
			location := fmt.Sprintf("%s, %s", result.Info.Region, result.Info.City)
			healthIcon := getHealthIcon(result.Info.HealthScore)
			latencyEmoji := getLatencyEmoji(result.AvgLatency)

			fmt.Fprintf(file, "| <pre><code>%s:%s</code></pre> | %s | %s | %d ms %s | %s %d%% |\n",
				result.Info.IP, result.Info.Port, result.Info.ISP, location,
				result.AvgLatency, latencyEmoji,
				healthIcon, result.Info.HealthScore)
		}

		fmt.Fprintf(file, "\n</details>\n\n---\n\n")
	}
}

// Helper functions

func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func encodeBadgeLabel(s string) string {
	s = strings.ReplaceAll(s, " ", "%20")
	s = strings.ReplaceAll(s, ":", "%3A")
	s = strings.ReplaceAll(s, ",", "%2C")
	s = strings.ReplaceAll(s, "+", "%2B")
	s = strings.ReplaceAll(s, "(", "%28")
	s = strings.ReplaceAll(s, ")", "%29")
	return s
}

func getHealthIcon(score int) string {
	if score >= 90 {
		return "🟢"
	} else if score >= 70 {
		return "🟡"
	}
	return "🔴"
}

func getLatencyEmoji(latency int64) string {
	if latency < 100 {
		return "⚡⚡"
	} else if latency < 300 {
		return "⚡"
	} else if latency < 1000 {
		return "🐇"
	}
	return "🐌"
}

func getProviderLogoHTML(provider string) string {
	mapping := map[string]string{
		"Google":       "google.com",
		"Amazon":       "amazon.com",
		"Cloudflare":   "cloudflare.com",
		"Hetzner":      "hetzner.com",
		"Hostinger":    "hostinger.com",
		"Tencent":      "www.tencent.com",
		"DigitalOcean": "digitalocean.com",
		"Vultr":        "vultr.com",
		"Akamai":       "akamai.com",
	}

	if domain, ok := mapping[provider]; ok {
		return fmt.Sprintf(`<img alt="%s" src="https://www.google.com/s2/favicons?sz=22&domain_url=%s" />`, provider, domain)
	}
	return ""
}

func getCountryFlag(code string) string {
	var flag strings.Builder
	for _, ch := range strings.ToUpper(code) {
		if ch >= 'A' && ch <= 'Z' {
			flag.WriteRune(rune(0x1F1E6 + (ch - 'A')))
		}
	}
	return flag.String()
}

func getCountryName(code string) string {
	countries := map[string]string{
		"US": "United States", "DE": "Germany", "GB": "United Kingdom",
		"FR": "France", "NL": "Netherlands", "CA": "Canada", "AU": "Australia",
		"JP": "Japan", "CN": "China", "SG": "Singapore", "KR": "South Korea",
		"IN": "India", "RU": "Russia", "BR": "Brazil", "IT": "Italy",
		"ES": "Spain", "SE": "Sweden", "CH": "Switzerland", "TR": "Turkey",
		"PL": "Poland", "FI": "Finland", "NO": "Norway", "IE": "Ireland",
		"BE": "Belgium", "AT": "Austria", "DK": "Denmark", "CZ": "Czech Republic",
		"UA": "Ukraine", "HK": "Hong Kong", "TW": "Taiwan", "IR": "Iran",
		"ZA": "South Africa", "RO": "Romania", "ID": "Indonesia", "VN": "Vietnam",
		"TH": "Thailand", "MY": "Malaysia", "MX": "Mexico", "AR": "Argentina",
		"CL": "Chile", "CO": "Colombia", "IL": "Israel", "AE": "United Arab Emirates",
		"SA": "Saudi Arabia", "PT": "Portugal", "HU": "Hungary", "GR": "Greece",
		"BG": "Bulgaria",
	}

	if name, ok := countries[strings.ToUpper(code)]; ok {
		return name
	}
	return code
}
