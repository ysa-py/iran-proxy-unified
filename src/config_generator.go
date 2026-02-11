package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ConfigGenerator handles intelligent config generation
type ConfigGenerator struct {
	proxyInfo    ProxyInfo
	iranMode     bool
	configs      []Config
	diversityMap map[string]int // Track config diversity
}

// NewConfigGenerator creates a new config generator
func NewConfigGenerator(proxyInfo ProxyInfo, iranMode bool) *ConfigGenerator {
	rand.Seed(time.Now().UnixNano())
	return &ConfigGenerator{
		proxyInfo:    proxyInfo,
		iranMode:     iranMode,
		configs:      make([]Config, 0),
		diversityMap: make(map[string]int),
	}
}

// GenerateAllConfigs generates configs for all protocols
func (cg *ConfigGenerator) GenerateAllConfigs() []Config {
	color.Cyan("🔧 Generating configs for IP: %s:%s (%s)",
		cg.proxyInfo.IP, cg.proxyInfo.Port, cg.proxyInfo.ISP)

	// Generate VMess configs
	cg.generateVMessConfigs()

	// Generate VLESS configs
	cg.generateVLESSConfigs()

	// Generate Shadowsocks configs
	cg.generateShadowsocksConfigs()

	// Generate Trojan configs
	cg.generateTrojanConfigs()

	color.Green("✅ Generated %d total configs", len(cg.configs))
	return cg.configs
}

// generateVMessConfigs creates VMess configurations
func (cg *ConfigGenerator) generateVMessConfigs() {
	baseUUID := GenerateUUID()

	transports := []string{TransportWebSocket, TransportGRPC, TransportHTTP2}
	if cg.iranMode {
		// Add xhttp for Iran (best DPI evasion)
		transports = append([]string{TransportXHTTP}, transports...)
	}

	for _, transport := range transports {
		// TLS version
		config := cg.createVMessConfig(baseUUID, transport, SecurityTLS)
		if config != nil {
			cg.configs = append(cg.configs, *config)
		}

		// Non-TLS version (some scenarios)
		if !cg.iranMode || transport == TransportTCP {
			config = cg.createVMessConfig(baseUUID, transport, SecurityNone)
			if config != nil {
				cg.configs = append(cg.configs, *config)
			}
		}
	}
}

// createVMessConfig creates a single VMess config
func (cg *ConfigGenerator) createVMessConfig(uuid, transport, security string) *Config {
	config := &Config{
		Protocol: ProtocolVMess,
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		ID:       uuid,
		AlterID:  0, // Modern VMess uses 0
		Security: security,
		Network:  transport,
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,
	}

	// Configure transport-specific settings
	switch transport {
	case TransportWebSocket, TransportXHTTP:
		config.Path = cg.getRandomPath()
		config.Host = cg.getRandomSNI()
	case TransportGRPC:
		config.ServiceName = cg.getRandomServiceName()
		config.Mode = "multi"
	case TransportHTTP2:
		config.Path = cg.getRandomPath()
		config.Host = cg.getRandomSNI()
	}

	// Configure TLS settings
	if security == SecurityTLS {
		config.SNI = cg.getRandomSNI()
		config.Fingerprint = cg.getRandomFingerprint()
		config.ALPN = []string{"h2", "http/1.1"}
	}

	// Calculate Iran optimization score
	config.HealthScore = config.GetIranOptimizedScore()
	config.IranOptimized = config.HealthScore >= 70

	// Generate remark
	config.Remark = cg.generateRemark(config)

	return config
}

// generateVLESSConfigs creates VLESS configurations
func (cg *ConfigGenerator) generateVLESSConfigs() {
	baseUUID := GenerateUUID()

	transports := []string{TransportWebSocket, TransportGRPC, TransportHTTP2, TransportTCP}
	if cg.iranMode {
		transports = append([]string{TransportXHTTP, TransportHTTPUpgrade}, transports...)
	}

	securities := []string{SecurityTLS, SecurityXTLS}
	if cg.iranMode {
		securities = append([]string{SecurityReality}, securities...)
	}

	for _, transport := range transports {
		for _, security := range securities {
			// Skip incompatible combinations
			if security == SecurityReality && transport != TransportTCP && transport != TransportGRPC {
				continue
			}
			if security == SecurityXTLS && transport != TransportTCP {
				continue
			}

			config := cg.createVLESSConfig(baseUUID, transport, security)
			if config != nil {
				cg.configs = append(cg.configs, *config)
			}
		}

		// Also create non-TLS version for some transports
		if transport == TransportWebSocket || transport == TransportHTTP2 {
			config := cg.createVLESSConfig(baseUUID, transport, SecurityNone)
			if config != nil {
				cg.configs = append(cg.configs, *config)
			}
		}
	}
}

// createVLESSConfig creates a single VLESS config
func (cg *ConfigGenerator) createVLESSConfig(uuid, transport, security string) *Config {
	config := &Config{
		Protocol:   ProtocolVLESS,
		Address:    cg.proxyInfo.IP,
		Port:       cg.proxyInfo.Port,
		ID:         uuid,
		Encryption: "none", // VLESS uses none encryption
		Security:   security,
		Network:    transport,
		ISP:        cg.proxyInfo.ISP,
		Country:    cg.proxyInfo.CountryCode,
	}

	// Configure transport-specific settings
	switch transport {
	case TransportWebSocket, TransportXHTTP, TransportHTTPUpgrade:
		config.Path = cg.getRandomPath()
		config.Host = cg.getRandomSNI()
	case TransportGRPC:
		config.ServiceName = cg.getRandomServiceName()
		config.Mode = "multi"
	case TransportHTTP2:
		config.Path = cg.getRandomPath()
		config.Host = cg.getRandomSNI()
	}

	// Configure security settings
	switch security {
	case SecurityTLS, SecurityXTLS:
		config.SNI = cg.getRandomSNI()
		config.Fingerprint = cg.getRandomFingerprint()
		config.ALPN = []string{"h2", "http/1.1"}

		if security == SecurityXTLS {
			config.Flow = cg.getRandomFlow()
		}

	case SecurityReality:
		config.SNI = cg.getRandomSNI()
		config.Fingerprint = cg.getRandomFingerprint()
		config.PublicKey = GeneratePassword(43) // Simulated public key
		config.ShortID = GenerateShortID()
		config.SpiderX = "/"
		config.Flow = cg.getRandomFlow()
	}

	// Calculate Iran optimization score
	config.HealthScore = config.GetIranOptimizedScore()
	config.IranOptimized = config.HealthScore >= 70

	// Generate remark
	config.Remark = cg.generateRemark(config)

	return config
}

// generateShadowsocksConfigs creates Shadowsocks configurations
func (cg *ConfigGenerator) generateShadowsocksConfigs() {
	password := GeneratePassword(16)

	ciphers := ShadowsocksCiphers
	if cg.iranMode {
		// Prefer modern ciphers for Iran
		ciphers = []string{
			"chacha20-ietf-poly1305",
			"2022-blake3-aes-256-gcm",
			"aes-256-gcm",
		}
	}

	for _, cipher := range ciphers {
		config := cg.createShadowsocksConfig(password, cipher)
		if config != nil {
			cg.configs = append(cg.configs, *config)
		}
	}
}

// createShadowsocksConfig creates a single Shadowsocks config
func (cg *ConfigGenerator) createShadowsocksConfig(password, cipher string) *Config {
	config := &Config{
		Protocol: ProtocolShadowsocks,
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		Password: password,
		Method:   cipher,
		Security: SecurityNone,
		Network:  TransportTCP,
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,
	}

	// Calculate Iran optimization score
	config.HealthScore = config.GetIranOptimizedScore()
	config.IranOptimized = config.HealthScore >= 70

	// Generate remark
	config.Remark = cg.generateRemark(config)

	return config
}

// generateTrojanConfigs creates Trojan configurations
func (cg *ConfigGenerator) generateTrojanConfigs() {
	password := GeneratePassword(32)

	transports := []string{TransportTCP, TransportWebSocket, TransportGRPC}
	if cg.iranMode {
		transports = append([]string{TransportXHTTP}, transports...)
	}

	for _, transport := range transports {
		config := cg.createTrojanConfig(password, transport)
		if config != nil {
			cg.configs = append(cg.configs, *config)
		}
	}
}

// createTrojanConfig creates a single Trojan config
func (cg *ConfigGenerator) createTrojanConfig(password, transport string) *Config {
	config := &Config{
		Protocol: ProtocolTrojan,
		Address:  cg.proxyInfo.IP,
		Port:     cg.proxyInfo.Port,
		Password: password,
		Security: SecurityTLS, // Trojan always uses TLS
		Network:  transport,
		ISP:      cg.proxyInfo.ISP,
		Country:  cg.proxyInfo.CountryCode,
	}

	// Configure transport-specific settings
	switch transport {
	case TransportWebSocket, TransportXHTTP:
		config.Path = cg.getRandomPath()
		config.Host = cg.getRandomSNI()
	case TransportGRPC:
		config.ServiceName = cg.getRandomServiceName()
		config.Mode = "multi"
	}

	// Configure TLS settings
	config.SNI = cg.getRandomSNI()
	config.Fingerprint = cg.getRandomFingerprint()
	config.ALPN = []string{"h2", "http/1.1"}

	// Calculate Iran optimization score
	config.HealthScore = config.GetIranOptimizedScore()
	config.IranOptimized = config.HealthScore >= 70

	// Generate remark
	config.Remark = cg.generateRemark(config)

	return config
}

// Helper functions for random generation

func (cg *ConfigGenerator) getRandomPath() string {
	paths := []string{
		"/",
		"/api",
		"/v2ray",
		"/vless",
		"/vmess",
		"/speedtest",
		"/graphql",
		"/ws",
		"/download",
		"/path",
		"/socket.io/",
		"/cdn-cgi/trace",
	}

	if cg.iranMode {
		// Add Iran-friendly paths
		paths = append(paths, []string{
			"/cloudflare",
			"/google",
			"/microsoft",
			"/apple",
			"/amazon",
		}...)
	}

	return paths[rand.Intn(len(paths))]
}

func (cg *ConfigGenerator) getRandomSNI() string {
	if cg.iranMode {
		return IranOptimizedSNIs[rand.Intn(len(IranOptimizedSNIs))]
	}

	// Fallback SNIs
	snis := append(IranOptimizedSNIs, []string{
		"www.yahoo.com",
		"www.github.com",
		"www.digitalocean.com",
	}...)

	return snis[rand.Intn(len(snis))]
}

func (cg *ConfigGenerator) getRandomFingerprint() string {
	if cg.iranMode {
		// Prefer chrome, firefox, random for Iran
		fps := []string{"chrome", "firefox", "random", "randomized"}
		return fps[rand.Intn(len(fps))]
	}

	return TLSFingerprints[rand.Intn(len(TLSFingerprints))]
}

func (cg *ConfigGenerator) getRandomServiceName() string {
	services := []string{
		"grpc",
		"GunService",
		"VMessGRPC",
		"VLESSGRPC",
		"TrojanGRPC",
		"CloudflareGRPC",
		"GoogleGRPC",
	}

	return services[rand.Intn(len(services))]
}

func (cg *ConfigGenerator) getRandomFlow() string {
	flows := []string{
		"xtls-rprx-vision",
		"xtls-rprx-direct",
	}

	return flows[rand.Intn(len(flows))]
}

func (cg *ConfigGenerator) generateRemark(config *Config) string {
	var parts []string

	// Add flag
	flag := getCountryFlag(config.Country)
	if flag != "" {
		parts = append(parts, flag)
	}

	// Add ISP (shortened)
	ispShort := config.ISP
	if len(ispShort) > 20 {
		ispShort = ispShort[:20]
	}
	parts = append(parts, ispShort)

	// Add protocol
	parts = append(parts, strings.ToUpper(config.Protocol))

	// Add transport
	transportName := config.Network
	switch config.Network {
	case TransportWebSocket:
		transportName = "WS"
	case TransportGRPC:
		transportName = "gRPC"
	case TransportXHTTP:
		transportName = "xHTTP"
	case TransportHTTPUpgrade:
		transportName = "HTTPUp"
	case TransportHTTP2:
		transportName = "H2"
	case TransportQUIC:
		transportName = "QUIC"
	}
	parts = append(parts, transportName)

	// Add security
	if config.Security != SecurityNone {
		parts = append(parts, strings.ToUpper(config.Security))
	}

	// Add Iran optimization marker
	if config.IranOptimized {
		parts = append(parts, "🇮🇷")
	}

	return strings.Join(parts, " | ")
}

// GetConfigsByScore returns configs sorted by health score
func (cg *ConfigGenerator) GetConfigsByScore() []Config {
	configs := make([]Config, len(cg.configs))
	copy(configs, cg.configs)

	// Sort by health score (descending)
	for i := 0; i < len(configs); i++ {
		for j := i + 1; j < len(configs); j++ {
			if configs[i].HealthScore < configs[j].HealthScore {
				configs[i], configs[j] = configs[j], configs[i]
			}
		}
	}

	return configs
}

// GetIranOptimizedConfigs returns only Iran-optimized configs
func (cg *ConfigGenerator) GetIranOptimizedConfigs() []Config {
	var iranConfigs []Config
	for _, config := range cg.configs {
		if config.IranOptimized {
			iranConfigs = append(iranConfigs, config)
		}
	}
	return iranConfigs
}

// GetConfigsByProtocol returns configs grouped by protocol
func (cg *ConfigGenerator) GetConfigsByProtocol() map[string][]Config {
	grouped := make(map[string][]Config)

	for _, config := range cg.configs {
		grouped[config.Protocol] = append(grouped[config.Protocol], config)
	}

	return grouped
}

// PrintConfigSummary prints a summary of generated configs
func (cg *ConfigGenerator) PrintConfigSummary() {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║              CONFIG GENERATION SUMMARY                        ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	byProtocol := cg.GetConfigsByProtocol()

	for protocol, configs := range byProtocol {
		iranCount := 0
		for _, c := range configs {
			if c.IranOptimized {
				iranCount++
			}
		}
		color.Yellow("║  %-15s: %3d total, %3d Iran-optimized          ║",
			strings.ToUpper(protocol), len(configs), iranCount)
	}

	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  TOTAL CONFIGS: %-3d                                         ║", len(cg.configs))

	iranTotal := len(cg.GetIranOptimizedConfigs())
	color.Magenta("║  IRAN-OPTIMIZED: %-3d                                        ║", iranTotal)
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")
}
