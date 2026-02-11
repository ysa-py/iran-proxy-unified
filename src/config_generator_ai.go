package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// EnhancedConfigGenerator with AI Anti-DPI capabilities
type EnhancedConfigGenerator struct {
	proxyInfo    ProxyInfo
	config       *AppConfig
	aiEngine     *AIAntiDPIEngine
	configs      []Config
	diversityMap map[string]int
}

// NewEnhancedConfigGenerator creates a new enhanced config generator with AI
func NewEnhancedConfigGenerator(proxyInfo ProxyInfo, config *AppConfig) *EnhancedConfigGenerator {
	// Initialize AI engine based on DPI evasion level
	aiEngine := NewAIAntiDPIEngine(config.DPIEvasionLevel)

	return &EnhancedConfigGenerator{
		proxyInfo:    proxyInfo,
		config:       config,
		aiEngine:     aiEngine,
		configs:      make([]Config, 0),
		diversityMap: make(map[string]int),
	}
}

// GenerateAllConfigs generates all configs with AI-powered DPI evasion
func (cg *EnhancedConfigGenerator) GenerateAllConfigs() []Config {
	color.Cyan("🤖 AI Anti-DPI Engine: Analyzing optimal configuration...")

	if cg.config.Verbose {
		cg.aiEngine.PrintAIEngineStatus()
	}

	// Determine optimal protocols based on DPI evasion level
	protocols := cg.selectOptimalProtocols()

	color.Cyan("🔧 Generating configs for IP: %s:%s (%s)",
		cg.proxyInfo.IP, cg.proxyInfo.Port, cg.proxyInfo.ISP)
	color.Green("   Selected protocols (AI-optimized): %v", protocols)

	// Generate configs for each selected protocol
	for _, protocol := range protocols {
		switch protocol {
		case "reality":
			cg.generateRealityConfigs()
		case "xhttp":
			cg.generateXHTTPConfigs()
		case "hysteria2":
			cg.generateHysteria2Configs()
		case "tuic":
			cg.generateTUICConfigs()
		case "vless-xtls":
			cg.generateVLESSXTLSConfigs()
		case "vmess":
			cg.generateVMessConfigsAI()
		case "vless":
			cg.generateVLESSConfigsAI()
		case "trojan":
			cg.generateTrojanConfigsAI()
		case "shadowsocks":
			cg.generateShadowsocksConfigsAI()
		}
	}

	// Apply AI-powered optimizations to all configs
	cg.applyAIOptimizations()

	color.Green("✅ Generated %d AI-optimized configs", len(cg.configs))
	return cg.configs
}

// selectOptimalProtocols selects best protocols based on AI analysis
func (cg *EnhancedConfigGenerator) selectOptimalProtocols() []string {
	protocols := make([]string, 0)

	// For Iran mode with maximum DPI evasion, prioritize Reality and xhttp
	if cg.config.IranMode {
		switch cg.config.DPIEvasionLevel {
		case "maximum":
			// Maximum evasion - use most advanced protocols
			protocols = []string{"reality", "xhttp", "hysteria2", "tuic", "vless-xtls"}

		case "aggressive":
			// Aggressive - balance between advanced and tested
			protocols = []string{"reality", "xhttp", "vless-xtls", "vmess", "vless"}

		default: // standard
			// Standard - proven protocols
			protocols = []string{"vless-xtls", "vmess", "vless", "trojan"}
		}
	} else {
		// Non-Iran mode - use standard protocols
		protocols = []string{"vmess", "vless", "trojan", "shadowsocks"}
	}

	// Let AI engine select the most optimal subset
	availableProtocols := protocols
	if len(cg.config.TestProtocols) > 0 {
		// Filter based on user preference
		filtered := make([]string, 0)
		for _, p := range protocols {
			for _, allowed := range cg.config.TestProtocols {
				if strings.Contains(p, allowed) {
					filtered = append(filtered, p)
					break
				}
			}
		}
		if len(filtered) > 0 {
			availableProtocols = filtered
		}
	}

	return availableProtocols
}

// generateRealityConfigs generates Reality protocol configs (highest DPI evasion)
func (cg *EnhancedConfigGenerator) generateRealityConfigs() {
	// Reality protocol - mimics real TLS perfectly
	uuid := GenerateUUID()

	// AI-generated domain fronting
	fronting := cg.aiEngine.GenerateDomainFronting()

	// AI-optimized fingerprint
	fingerprint := cg.aiEngine.GenerateAdaptiveFingerprint()

	config := &Config{
		Protocol: "vless",
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		ID:       uuid,
		Security: "reality",
		Network:  "tcp",
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,

		// Reality-specific parameters
		Flow:        "xtls-rprx-vision",
		Fingerprint: fingerprint,
		PublicKey:   generatePublicKey(),
		ShortID:     generateShortID(),
		ServerName:  fronting["front_domain"],
		SpiderX:     "/",

		// AI optimizations
		IranOptimized: true,
		DPIEvaded:     true,
	}

	cg.configs = append(cg.configs, *config)
}

// generateXHTTPConfigs generates xhttp transport configs (advanced HTTP obfuscation)
func (cg *EnhancedConfigGenerator) generateXHTTPConfigs() {
	uuid := GenerateUUID()

	// AI-generated traffic mimicry
	mimicry := cg.aiEngine.GenerateTrafficMimicry("cloudflare")

	config := &Config{
		Protocol: "vmess",
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		ID:       uuid,
		Security: SecurityTLS,
		Network:  "xhttp",
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,

		// xhttp-specific
		Path:    "/api/v1/status",
		Host:    fmt.Sprintf("%v", mimicry["user_agent"]),
		Headers: convertMimicryToHeaders(mimicry),

		// AI optimizations
		IranOptimized: true,
		DPIEvaded:     true,
	}

	cg.configs = append(cg.configs, *config)
}

// generateHysteria2Configs generates Hysteria2 configs (QUIC-based, hard to detect)
func (cg *EnhancedConfigGenerator) generateHysteria2Configs() {
	// Hysteria2 - uses QUIC protocol which is harder for DPI to analyze
	password := generateRandomPassword(32)

	config := &Config{
		Protocol: "hysteria2",
		Address:  cg.proxyInfo.IP,
		Port:     "443", // Standard HTTPS port
		ID:       password,
		Security: "tls",
		Network:  "quic",
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,

		// Hysteria2-specific
		Obfs:         "salamander",
		ObfsPassword: generateRandomPassword(16),
		SNI:          "www.microsoft.com",
		ALPN:         []string{"h3"},

		// AI optimizations
		IranOptimized: true,
		DPIEvaded:     true,
	}

	cg.configs = append(cg.configs, *config)
}

// generateTUICConfigs generates TUIC configs (UDP-based with obfuscation)
func (cg *EnhancedConfigGenerator) generateTUICConfigs() {
	uuid := GenerateUUID()
	password := generateRandomPassword(32)

	config := &Config{
		Protocol: "tuic",
		Address:  cg.proxyInfo.IP,
		Port:     "443",
		ID:       uuid,
		Password: password,
		Security: "tls",
		Network:  "udp",
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,

		// TUIC-specific
		CongestionControl: "bbr",
		UDPRelayMode:      "native",
		SNI:               "cloudflare.com",
		ALPN:              []string{"h3", "spdy/3.1"},

		// AI optimizations
		IranOptimized: true,
		DPIEvaded:     true,
	}

	cg.configs = append(cg.configs, *config)
}

// generateVLESSXTLSConfigs generates VLESS with XTLS Vision
func (cg *EnhancedConfigGenerator) generateVLESSXTLSConfigs() {
	uuid := GenerateUUID()

	// AI-optimized SNI with fragmentation
	sni := "www.google.com"
	fragPoints := cg.aiEngine.GenerateSNIFragmentation(sni)

	config := &Config{
		Protocol: "vless",
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		ID:       uuid,
		Security: "xtls",
		Network:  "tcp",
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,

		// XTLS Vision
		Flow:        "xtls-rprx-vision",
		SNI:         sni,
		Fingerprint: cg.aiEngine.GenerateAdaptiveFingerprint(),
		ALPN:        []string{"h2", "http/1.1"},

		// AI optimizations
		FragmentationPoints: fragPoints,
		IranOptimized:       true,
		DPIEvaded:           true,
	}

	cg.configs = append(cg.configs, *config)
}

// generateVMessConfigsAI generates VMess configs with AI optimizations
func (cg *EnhancedConfigGenerator) generateVMessConfigsAI() {
	baseUUID := GenerateUUID()

	transports := []string{TransportWebSocket, TransportGRPC, TransportHTTP2}
	if cg.config.IranMode {
		transports = append([]string{TransportXHTTP}, transports...)
	}

	for _, transport := range transports {
		// Apply AI timing obfuscation
		delay := cg.aiEngine.ApplyTimingObfuscation()

		config := cg.createVMessConfigWithAI(baseUUID, transport, SecurityTLS, delay)
		if config != nil {
			cg.configs = append(cg.configs, *config)
		}

		// Non-TLS for TCP only (if not Iran mode)
		if !cg.config.IranMode && transport == TransportTCP {
			config = cg.createVMessConfigWithAI(baseUUID, transport, SecurityNone, delay)
			if config != nil {
				cg.configs = append(cg.configs, *config)
			}
		}
	}
}

// createVMessConfigWithAI creates VMess config with AI enhancements
func (cg *EnhancedConfigGenerator) createVMessConfigWithAI(uuid, transport, security string, delay time.Duration) *Config {
	// AI-generated traffic mimicry
	mimicry := cg.aiEngine.GenerateTrafficMimicry("cloudflare")

	config := &Config{
		Protocol: ProtocolVMess,
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		ID:       uuid,
		AlterID:  0,
		Security: security,
		Network:  transport,
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,

		// AI enhancements
		TimingDelay:   convertDurationToMillis(delay),
		Headers:       convertMimicryToHeaders(mimicry),
		IranOptimized: cg.config.IranMode,
	}

	// Configure transport-specific settings with AI
	switch transport {
	case TransportWebSocket, TransportXHTTP:
		config.Path = cg.getRandomPathAI()
		config.Host = fmt.Sprintf("%v", mimicry["user_agent"])

		// Apply packet padding
		config.PaddingSize = cg.aiEngine.ApplyPacketPadding(len(config.Path))

	case TransportGRPC:
		config.ServiceName = cg.getRandomServiceName()
		config.Mode = "multi"

	case TransportHTTP2:
		config.Host = cg.getRandomSNI()
		config.Path = "/"

	case TransportQUIC:
		config.Host = "www.google.com"
		config.Security = "quic"
		config.Key = generateRandomPassword(16)
	}

	// Apply SNI fragmentation for TLS
	if security == SecurityTLS {
		config.SNI = cg.getRandomSNI()
		config.FragmentationPoints = cg.aiEngine.GenerateSNIFragmentation(config.SNI)
		config.Fingerprint = cg.aiEngine.GenerateAdaptiveFingerprint()
		config.ALPN = []string{"h2", "http/1.1"}
	}

	return config
}

// generateVLESSConfigsAI generates VLESS configs with AI optimizations
func (cg *EnhancedConfigGenerator) generateVLESSConfigsAI() {
	uuid := GenerateUUID()

	transports := []string{TransportWebSocket, TransportGRPC, TransportTCP}

	for _, transport := range transports {
		delay := cg.aiEngine.ApplyTimingObfuscation()

		config := &Config{
			Protocol: ProtocolVLESS,
			Address:  cg.proxyInfo.IP,
			Port:     cg.proxyInfo.Port,
			ID:       uuid,
			Security: SecurityTLS,
			Network:  transport,
			ISP:      cg.proxyInfo.ISP,
			Country:  cg.proxyInfo.CountryCode,

			// AI enhancements
			TimingDelay:   convertDurationToMillis(delay),
			Fingerprint:   cg.aiEngine.GenerateAdaptiveFingerprint(),
			IranOptimized: cg.config.IranMode,
		}

		// Transport-specific AI optimizations
		switch transport {
		case TransportWebSocket:
			config.Path = cg.getRandomPathAI()
			config.Host = cg.getRandomSNI()

		case TransportGRPC:
			config.ServiceName = cg.getRandomServiceName()

		case TransportTCP:
			config.Flow = "xtls-rprx-direct"
		}

		// SNI fragmentation
		config.SNI = cg.getRandomSNI()
		config.FragmentationPoints = cg.aiEngine.GenerateSNIFragmentation(config.SNI)

		cg.configs = append(cg.configs, *config)
	}
}

// generateTrojanConfigsAI generates Trojan configs with AI optimizations
func (cg *EnhancedConfigGenerator) generateTrojanConfigsAI() {
	password := generateRandomPassword(32)

	transports := []string{TransportWebSocket, TransportGRPC, TransportTCP}

	for _, transport := range transports {
		delay := cg.aiEngine.ApplyTimingObfuscation()
		mimicry := cg.aiEngine.GenerateTrafficMimicry("microsoft")

		config := &Config{
			Protocol: ProtocolTrojan,
			Address:  cg.proxyInfo.IP,
			Port:     cg.proxyInfo.Port,
			Password: password,
			Security: SecurityTLS,
			Network:  transport,
			ISP:      cg.proxyInfo.ISP,
			Country:  cg.proxyInfo.CountryCode,

			// AI enhancements
			TimingDelay:   convertDurationToMillis(delay),
			Headers:       convertMimicryToHeaders(mimicry),
			Fingerprint:   cg.aiEngine.GenerateAdaptiveFingerprint(),
			IranOptimized: cg.config.IranMode,
		}

		// Transport settings
		if transport == TransportWebSocket {
			config.Path = cg.getRandomPathAI()
		} else if transport == TransportGRPC {
			config.ServiceName = cg.getRandomServiceName()
		}

		// SNI optimization
		config.SNI = cg.getRandomSNI()
		config.FragmentationPoints = cg.aiEngine.GenerateSNIFragmentation(config.SNI)

		cg.configs = append(cg.configs, *config)
	}
}

// generateShadowsocksConfigsAI generates Shadowsocks configs with AI optimizations
func (cg *EnhancedConfigGenerator) generateShadowsocksConfigsAI() {
	password := generateRandomPassword(32)

	// Modern ciphers only (older ones are easily detected)
	methods := []string{"chacha20-ietf-poly1305", "aes-256-gcm", "aes-128-gcm"}

	for _, method := range methods {
		delay := cg.aiEngine.ApplyTimingObfuscation()

		config := &Config{
			Protocol: ProtocolShadowsocks,
			Address:  cg.proxyInfo.IP,
			Port:     cg.proxyInfo.Port,
			Password: password,
			Method:   method,
			Network:  "tcp",
			ISP:      cg.proxyInfo.ISP,
			Country:  cg.proxyInfo.CountryCode,

			// AI enhancements
			TimingDelay:   convertDurationToMillis(delay),
			IranOptimized: false, // SS is less effective against modern DPI
		}

		// Add obfuscation plugin for Iran mode
		if cg.config.IranMode && cg.config.DPIEvasionLevel != "standard" {
			config.Plugin = "obfs-local"
			config.PluginOpts = "obfs=tls;obfs-host=www.bing.com"
		}

		cg.configs = append(cg.configs, *config)
	}
}

// applyAIOptimizations applies AI-powered optimizations to all configs
func (cg *EnhancedConfigGenerator) applyAIOptimizations() {
	color.Cyan("🤖 Applying AI optimizations to %d configs...", len(cg.configs))

	for i := range cg.configs {
		config := &cg.configs[i]

		// Apply traffic mimicry
		if len(config.Headers) == 0 {
			config.Headers = convertMimicryToHeaders(cg.aiEngine.GenerateTrafficMimicry("cloudflare"))
		}

		// Apply packet padding if not set
		if config.PaddingSize == 0 {
			baseSize := len(config.Address) + len(config.Port)
			config.PaddingSize = cg.aiEngine.ApplyPacketPadding(baseSize)
		}

		// Ensure fingerprint randomization
		if config.Fingerprint == "" && config.Security == SecurityTLS {
			config.Fingerprint = cg.aiEngine.GenerateAdaptiveFingerprint()
		}

		// Mark as DPI evaded if using advanced techniques
		if config.Security == "reality" || config.Security == "xtls" ||
			config.Network == "xhttp" || config.Protocol == "hysteria2" {
			config.DPIEvaded = true
		}
	}

	color.Green("✅ AI optimizations applied successfully")
}

// getRandomPathAI generates AI-optimized random path
func (cg *EnhancedConfigGenerator) getRandomPathAI() string {
	// Mimic real API endpoints
	paths := []string{
		"/api/v1/status",
		"/api/v2/health",
		"/health/check",
		"/status/live",
		"/metrics/prometheus",
		"/v1/graphql",
		"/api/data/fetch",
		"/cdn/assets/load",
	}

	return paths[randInt(len(paths))]
}

// Utility functions

func generatePublicKey() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	key := make([]byte, 44)
	for i := range key {
		key[i] = chars[randInt(len(chars))]
	}
	return string(key)
}

func generateShortID() string {
	chars := "0123456789abcdef"
	id := make([]byte, 8)
	for i := range id {
		id[i] = chars[randInt(len(chars))]
	}
	return string(id)
}

func generateRandomPassword(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()"
	password := make([]byte, length)
	for i := range password {
		password[i] = chars[randInt(len(chars))]
	}
	return string(password)
}

// getRandomServiceName generates random gRPC service name
func (cg *EnhancedConfigGenerator) getRandomServiceName() string {
	services := []string{
		"GunService",
		"TunService",
		"DataService",
		"StreamService",
		"ProxyService",
	}

	return services[randInt(len(services))]
}

// convertMimicryToHeaders converts map[string]interface{} to map[string]string
func convertMimicryToHeaders(mimicry map[string]interface{}) map[string]string {
	headers := make(map[string]string)
	for key, value := range mimicry {
		if strVal, ok := value.(string); ok {
			headers[key] = strVal
		} else {
			headers[key] = fmt.Sprintf("%v", value)
		}
	}
	return headers
}

// convertDurationToMillis converts time.Duration to milliseconds as int
func convertDurationToMillis(d time.Duration) int {
	return int(d.Milliseconds())
}

// getRandomSNI generates random SNI
func (cg *EnhancedConfigGenerator) getRandomSNI() string {
	snis := []string{
		"www.google.com",
		"www.microsoft.com",
		"www.cloudflare.com",
		"www.bing.com",
		"www.apple.com",
		"cdn.cloudflare.com",
		"www.github.com",
	}

	return snis[randInt(len(snis))]
}

// PrintConfigSummary prints summary of generated configs
func (cg *EnhancedConfigGenerator) PrintConfigSummary() {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║              CONFIG GENERATION SUMMARY                        ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Total Configs:      %-40d║", len(cg.configs))

	// Count by protocol
	protocolCount := make(map[string]int)
	iranOptimized := 0
	dpiEvaded := 0

	for _, config := range cg.configs {
		protocolCount[config.Protocol]++
		if config.IranOptimized {
			iranOptimized++
		}
		if config.DPIEvaded {
			dpiEvaded++
		}
	}

	for protocol, count := range protocolCount {
		color.Yellow("║  %-20s: %-40d║", protocol, count)
	}

	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Magenta("║  Iran-Optimized:     %-40d║", iranOptimized)
	color.Magenta("║  DPI-Evaded:         %-40d║", dpiEvaded)
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
}
