package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const banner = `
╔══════════════════════════════════════════════════════════════════════════╗
║                                                                          ║
║   ██████╗ ██████╗  ██████╗ ██╗  ██╗██╗   ██╗     ██████╗██╗  ██╗███████╗ ║
║   ██╔══██╗██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗ ██╔╝    ██╔════╝██║  ██║██╔════╝ ║
║   ██████╔╝██████╔╝██║   ██║ ╚███╔╝  ╚████╔╝     ██║     ███████║█████╗   ║
║   ██╔═══╝ ██╔══██╗██║   ██║ ██╔██╗   ╚██╔╝      ██║     ██╔══██║██╔══╝   ║
║   ██║     ██║  ██║╚██████╔╝██╔╝ ██╗   ██║       ╚██████╗██║  ██║███████╗ ║
║   ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝        ╚═════╝╚═╝  ╚═╝╚══════╝ ║
║                                                                          ║
║      🇮🇷 Advanced Iran Proxy & Config System - Enterprise Grade 🇮🇷      ║
║                                                                          ║
║   Proxy Check • Config Generation • DPI Evasion • Self-Healing System   ║
║                                                                          ║
╚══════════════════════════════════════════════════════════════════════════╝
`

const version = "v3.2.0-ultimate-iran-enterprise"

// Config holds application configuration
type AppConfig struct {
	// File paths
	ProxyFile    string
	OutputFile   string
	ConfigOutput string

	// Performance settings
	MaxConcurrent int
	TimeoutSecs   int

	// Mode flags
	IranMode        bool
	GenerateConfigs bool
	TestConfigs     bool
	ConfigsOnly     bool

	// Advanced features - NEW
	PerformanceMode   string // speed, balanced, quality
	DPIEvasionLevel   string // standard, aggressive, maximum
	EnableFallback    bool
	EnableSelfHealing bool
	EnableMonitoring  bool

	// Protocol testing
	TestProtocols []string

	// Emergency & Recovery
	EmergencyMode bool
	DeepAnalysis  bool

	// Display
	Help        bool
	ShowVersion bool
	Verbose     bool
}

func main() {
	// Parse configuration from flags and environment
	config := parseConfig()

	if config.Help {
		showHelp()
		return
	}

	if config.ShowVersion {
		fmt.Printf("🇮🇷 Iran Proxy & Config System %s\n", version)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	// Show banner
	showBanner()

	// Initialize monitoring if enabled
	var monitor *SystemMonitor
	if config.EnableMonitoring {
		monitor = NewSystemMonitor()
		monitor.Start()
		defer monitor.Stop()
	}

	// Configuration summary
	printConfigSummary(config)

	// Optimize system resources based on performance mode
	optimizeSystemResources(config.PerformanceMode)

	// Create enhanced proxy checker with advanced anti-DPI features
	timeout := time.Duration(config.TimeoutSecs) * time.Second
	checker := NewEnhancedProxyChecker(
		config.ProxyFile,
		config.OutputFile,
		config.MaxConcurrent,
		timeout,
		config,
	)

	// Configure advanced features based on DPI evasion level
	enableUTLS := config.DPIEvasionLevel != "standard"
	enableSNI := config.DPIEvasionLevel == "aggressive" || config.DPIEvasionLevel == "maximum"
	enableAI := config.DPIEvasionLevel == "maximum"
	enableAdaptive := config.EnableSelfHealing

	checker.SetAdvancedFeatures(enableUTLS, enableSNI, enableAI, enableAdaptive)

	// Show enabled features
	if enableUTLS || enableSNI || enableAI {
		color.Cyan("\n🛡️  Advanced Anti-DPI Features Enabled:")
		if enableUTLS {
			color.Green("   ✓ uTLS Fingerprint Spoofing")
		}
		if enableSNI {
			color.Green("   ✓ SNI Fragmentation")
		}
		if enableAI {
			color.Green("   ✓ AI-Powered DPI Evasion")
		}
		if enableAdaptive {
			color.Green("   ✓ Adaptive Learning System")
		}
		fmt.Println()
	}

	// Overall start time
	overallStartTime := time.Now()

	// Step 1: Proxy Checking (unless configs-only mode)
	if !config.ConfigsOnly {
		if err := runProxyChecking(checker, config, monitor); err != nil {
			color.Red("❌ Proxy checking failed: %v", err)
			if config.EmergencyMode {
				color.Yellow("🔄 Emergency mode: Attempting recovery...")
				if err := runEmergencyRecovery(checker, config); err != nil {
					os.Exit(1)
				}
			} else {
				os.Exit(1)
			}
		}
	}

	// Step 2: Config Generation
	var allGeneratedConfigs []Config
	var allPassedConfigs []TestedConfig

	if config.GenerateConfigs && (checker.stats.TotalActive > 0 || config.ConfigsOnly) {
		var err error
		allGeneratedConfigs, err = runConfigGeneration(checker, config, monitor)
		if err != nil {
			color.Red("❌ Config generation failed: %v", err)
			if !config.ConfigsOnly {
				// Continue with partial results
				color.Yellow("⚠️  Continuing with partial results...")
			} else {
				os.Exit(1)
			}
		}
	}

	// Step 3: Config Testing
	if config.TestConfigs && len(allGeneratedConfigs) > 0 {
		var err error
		allPassedConfigs, err = runConfigTesting(allGeneratedConfigs, config, monitor)
		if err != nil {
			color.Red("❌ Config testing failed: %v", err)
			// Continue with generated configs
			color.Yellow("⚠️  Continuing with untested configs...")
		}
	}

	// Step 4: Write Config Files
	if len(allPassedConfigs) > 0 || (len(allGeneratedConfigs) > 0 && !config.TestConfigs) {
		configsToWrite := allPassedConfigs
		if len(allPassedConfigs) == 0 {
			// Convert generated configs to tested configs
			for _, cfg := range allGeneratedConfigs {
				configsToWrite = append(configsToWrite, TestedConfig{
					Config:      cfg,
					Passed:      true,
					Latency:     0,
					AvgLatency:  0,
					SuccessRate: 0,
				})
			}
		}

		if err := runConfigWriting(config.ConfigOutput, configsToWrite, config); err != nil {
			color.Red("❌ Failed to write config files: %v", err)
			os.Exit(1)
		}
	}

	// Final Summary
	overallElapsed := time.Since(overallStartTime)
	printFinalSummary(checker, allGeneratedConfigs, allPassedConfigs, overallElapsed, config)

	// Generate metrics report if monitoring enabled
	if monitor != nil {
		monitor.GenerateReport()
	}

	// Self-healing check
	if config.EnableSelfHealing {
		performSelfHealingCheck(checker, allPassedConfigs, config)
	}
}

func parseConfig() *AppConfig {
	config := &AppConfig{}

	// File paths
	flag.StringVar(&config.ProxyFile, "proxy-file", getEnv("PROXY_FILE", DefaultProxyFile), "Path to proxy file")
	flag.StringVar(&config.OutputFile, "output-file", getEnv("OUTPUT_FILE", DefaultOutputFile), "Path to output file")
	flag.StringVar(&config.ConfigOutput, "config-output", getEnv("CONFIG_OUTPUT", "configs/iran-configs.txt"), "Config output path")

	// Performance settings
	flag.IntVar(&config.MaxConcurrent, "max-concurrent", getEnvInt("MAX_CONCURRENT", DefaultMaxConcurrent), "Max concurrent connections (50-500)")
	flag.IntVar(&config.TimeoutSecs, "timeout", getEnvInt("TIMEOUT", DefaultTimeoutSecs), "Timeout in seconds (5-30)")

	// Mode flags - ALIGNED WITH GITHUB ACTIONS
	flag.BoolVar(&config.IranMode, "iran-mode", getEnvBool("IRAN_MODE", true), "Enable Iran-specific optimizations")
	flag.BoolVar(&config.GenerateConfigs, "generate-configs", getEnvBool("GENERATE_CONFIGS", true), "Generate configs from proxies")
	flag.BoolVar(&config.TestConfigs, "test-configs", getEnvBool("TEST_CONFIGS", true), "Test generated configs")
	flag.BoolVar(&config.ConfigsOnly, "configs-only", getEnvBool("CONFIGS_ONLY", false), "Only generate configs")

	// Advanced features - NEW & ALIGNED WITH WORKFLOW
	perfMode := flag.String("performance-mode", getEnv("PERFORMANCE_MODE", "balanced"), "Performance mode: speed, balanced, quality")
	dpiLevel := flag.String("dpi-evasion-level", getEnv("DPI_EVASION_LEVEL", "aggressive"), "DPI evasion: standard, aggressive, maximum")
	flag.BoolVar(&config.EnableFallback, "enable-fallback", getEnvBool("ENABLE_FALLBACK", true), "Enable multi-tier fallback")
	flag.BoolVar(&config.EnableSelfHealing, "enable-self-healing", getEnvBool("ENABLE_SELF_HEALING", true), "Enable self-healing")
	flag.BoolVar(&config.EnableMonitoring, "enable-monitoring", getEnvBool("ENABLE_MONITORING", true), "Enable monitoring & metrics")

	// Protocol testing
	protocolsStr := flag.String("test-protocols", getEnv("TEST_PROTOCOLS", "vmess,vless,trojan,shadowsocks"), "Protocols to test")

	// Emergency & Deep Analysis
	flag.BoolVar(&config.EmergencyMode, "emergency-mode", getEnvBool("EMERGENCY_MODE", false), "Emergency recovery mode")
	flag.BoolVar(&config.DeepAnalysis, "deep-analysis", getEnvBool("DEEP_ANALYSIS", false), "Deep analysis mode")

	// Display
	flag.BoolVar(&config.Help, "help", false, "Show help message")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version")
	flag.BoolVar(&config.Verbose, "verbose", getEnvBool("VERBOSE", false), "Verbose output")

	flag.Parse()

	// Validate and set advanced configs
	config.PerformanceMode = validatePerformanceMode(*perfMode)
	config.DPIEvasionLevel = validateDPILevel(*dpiLevel)
	config.TestProtocols = parseProtocols(*protocolsStr)

	// Validate ranges
	if config.MaxConcurrent < 50 {
		config.MaxConcurrent = 50
	} else if config.MaxConcurrent > 500 {
		config.MaxConcurrent = 500
	}

	if config.TimeoutSecs < 5 {
		config.TimeoutSecs = 5
	} else if config.TimeoutSecs > 30 {
		config.TimeoutSecs = 30
	}

	return config
}

func validatePerformanceMode(mode string) string {
	mode = strings.ToLower(mode)
	switch mode {
	case "speed", "balanced", "quality":
		return mode
	default:
		return "balanced"
	}
}

func validateDPILevel(level string) string {
	level = strings.ToLower(level)
	switch level {
	case "standard", "aggressive", "maximum":
		return level
	default:
		return "aggressive"
	}
}

func parseProtocols(protocolsStr string) []string {
	protocols := strings.Split(protocolsStr, ",")
	result := make([]string, 0, len(protocols))
	for _, p := range protocols {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func showBanner() {
	color.Cyan(banner)
	color.White("Version: %s", version)
	color.White("Go: %s | Platform: %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}

func printConfigSummary(config *AppConfig) {
	color.Cyan("╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                    CONFIGURATION SUMMARY                      ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Proxy File:          %-40s║", truncate(config.ProxyFile, 40))
	color.Yellow("║  Output File:         %-40s║", truncate(config.OutputFile, 40))
	color.Yellow("║  Config Output:       %-40s║", truncate(config.ConfigOutput, 40))
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Max Concurrent:      %-40d║", config.MaxConcurrent)
	color.Yellow("║  Timeout:             %-37ds ║", config.TimeoutSecs)
	color.Yellow("║  Performance Mode:    %-40s║", strings.ToUpper(config.PerformanceMode))
	color.Yellow("║  DPI Evasion Level:   %-40s║", strings.ToUpper(config.DPIEvasionLevel))
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	if config.IranMode {
		color.Magenta("║  Iran Mode:           %-40s║", "🇮🇷 ENABLED")
		color.Magenta("║  DPI Evasion:         %-40s║", "✅ ACTIVE")
		color.Magenta("║  Multi-Endpoint:      %-40s║", "✅ ACTIVE")
	} else {
		color.Yellow("║  Iran Mode:           %-40s║", "DISABLED")
	}

	color.Yellow("║  Generate Configs:    %-40s║", boolToStatus(config.GenerateConfigs))
	color.Yellow("║  Test Configs:        %-40s║", boolToStatus(config.TestConfigs))
	color.Yellow("║  Self-Healing:        %-40s║", boolToStatus(config.EnableSelfHealing))
	color.Yellow("║  Monitoring:          %-40s║", boolToStatus(config.EnableMonitoring))
	color.Yellow("║  Fallback System:     %-40s║", boolToStatus(config.EnableFallback))

	if config.EmergencyMode {
		color.Red("║  Emergency Mode:      %-40s║", "🚨 ACTIVE")
	}
	if config.DeepAnalysis {
		color.Magenta("║  Deep Analysis:       %-40s║", "🔬 ACTIVE")
	}

	if len(config.TestProtocols) > 0 {
		protocols := strings.Join(config.TestProtocols, ", ")
		color.Yellow("║  Test Protocols:      %-40s║", truncate(protocols, 40))
	}

	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func optimizeSystemResources(mode string) {
	switch mode {
	case "speed":
		runtime.GOMAXPROCS(runtime.NumCPU())
		color.Green("⚡ Speed mode: Using all CPU cores (%d)", runtime.NumCPU())
	case "balanced":
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
		color.Cyan("⚖️  Balanced mode: Using %d CPU cores", runtime.NumCPU()/2)
	case "quality":
		runtime.GOMAXPROCS(runtime.NumCPU() / 4)
		color.Magenta("💎 Quality mode: Using %d CPU cores for precision", runtime.NumCPU()/4)
	}
}

func runProxyChecking(checker *EnhancedProxyChecker, config *AppConfig, monitor *SystemMonitor) error {
	startTime := time.Now()

	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                  STEP 1: PROXY CHECKING                       ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	if monitor != nil {
		monitor.RecordEvent("proxy_check_start")
	}

	if err := checker.Run(); err != nil {
		if monitor != nil {
			monitor.RecordEvent("proxy_check_failed")
		}
		return err
	}

	if err := checker.WriteMarkdownFile(); err != nil {
		return fmt.Errorf("failed to write markdown: %w", err)
	}

	elapsed := time.Since(startTime)
	if monitor != nil {
		monitor.RecordEvent("proxy_check_complete")
		monitor.RecordMetric("proxy_check_duration", elapsed.Seconds())
	}

	color.Green("\n✅ Proxy checking completed in %s", elapsed.String())
	printProxyStats(checker, elapsed, config)

	return nil
}

func runConfigGeneration(checker *EnhancedProxyChecker, config *AppConfig, monitor *SystemMonitor) ([]Config, error) {
	startTime := time.Now()

	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                  STEP 2: CONFIG GENERATION                    ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	if monitor != nil {
		monitor.RecordEvent("config_gen_start")
	}

	var allConfigs []Config
	configCount := 0

	for _, proxyResults := range checker.activeProxies {
		for _, result := range proxyResults {
			if config.Verbose {
				color.Yellow("Generating configs for: %s:%s (%s)",
					result.Info.IP, result.Info.Port, result.Info.ISP)
			}

			generator := NewEnhancedConfigGenerator(result.Info, config)
			configs := generator.GenerateAllConfigs()

			if config.Verbose {
				generator.PrintConfigSummary()
			}

			allConfigs = append(allConfigs, configs...)
			configCount += len(configs)
		}
	}

	elapsed := time.Since(startTime)
	if monitor != nil {
		monitor.RecordEvent("config_gen_complete")
		monitor.RecordMetric("config_gen_duration", elapsed.Seconds())
		monitor.RecordMetric("configs_generated", float64(configCount))
	}

	color.Green("\n✅ Config generation completed in %s", elapsed.String())
	color.Green("📊 Total configs generated: %d", configCount)

	return allConfigs, nil
}

func runConfigTesting(configs []Config, config *AppConfig, monitor *SystemMonitor) ([]TestedConfig, error) {
	startTime := time.Now()

	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                  STEP 3: CONFIG TESTING                       ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	if monitor != nil {
		monitor.RecordEvent("config_test_start")
	}

	timeout := time.Duration(config.TimeoutSecs) * time.Second
	tester := NewEnhancedConfigTester(configs, config.MaxConcurrent, timeout, config)
	passedConfigs := tester.TestAllConfigs()

	if config.Verbose {
		tester.PrintConfigDetails()
	}

	elapsed := time.Since(startTime)
	if monitor != nil {
		monitor.RecordEvent("config_test_complete")
		monitor.RecordMetric("config_test_duration", elapsed.Seconds())
		monitor.RecordMetric("configs_passed", float64(len(passedConfigs)))
	}

	color.Green("\n✅ Config testing completed in %s", elapsed.String())
	color.Green("📊 Passed configs: %d / %d", len(passedConfigs), len(configs))

	return passedConfigs, nil
}

func runConfigWriting(outputPath string, configs []TestedConfig, config *AppConfig) error {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                  STEP 4: WRITING CONFIG FILES                 ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	if err := writeEnhancedConfigFiles(outputPath, configs, config); err != nil {
		return err
	}

	color.Green("\n✅ Config files written successfully")
	color.Green("📁 Output directory: %s", outputPath)

	return nil
}

func runEmergencyRecovery(checker *EnhancedProxyChecker, config *AppConfig) error {
	color.Yellow("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Yellow("║               🚨 EMERGENCY RECOVERY MODE                      ║")
	color.Yellow("╚═══════════════════════════════════════════════════════════════╝\n")

	// Try with relaxed settings
	checker.maxConcurrent = config.MaxConcurrent / 2
	checker.timeout = checker.timeout * 2

	color.Yellow("🔄 Retrying with relaxed settings...")
	color.Yellow("   Concurrent: %d -> %d", config.MaxConcurrent, checker.maxConcurrent)
	color.Yellow("   Timeout: %s -> %s", time.Duration(config.TimeoutSecs)*time.Second, checker.timeout)

	if err := checker.Run(); err != nil {
		return fmt.Errorf("emergency recovery failed: %w", err)
	}

	color.Green("✅ Emergency recovery successful")
	return nil
}

func performSelfHealingCheck(checker *EnhancedProxyChecker, configs []TestedConfig, config *AppConfig) {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                  🔧 SELF-HEALING CHECK                        ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	// Check success rate
	totalActive := checker.stats.TotalActive
	totalTested := checker.stats.TotalTested

	if totalTested == 0 {
		color.Yellow("⚠️  No proxies tested, skipping self-healing")
		return
	}

	successRate := float64(totalActive) / float64(totalTested) * 100

	if successRate < 60 {
		color.Red("⚠️  Low success rate (%.1f%%) - System health degraded", successRate)
		color.Yellow("💡 Recommendations:")
		color.Yellow("   • Increase timeout value")
		color.Yellow("   • Reduce concurrent connections")
		color.Yellow("   • Check network connectivity")
		color.Yellow("   • Verify proxy source list")
	} else if successRate < 80 {
		color.Yellow("⚠️  Moderate success rate (%.1f%%) - Room for improvement", successRate)
		color.Yellow("💡 Consider enabling aggressive DPI evasion")
	} else {
		color.Green("✅ Excellent success rate (%.1f%%) - System healthy", successRate)
	}

	// Check config quality
	if len(configs) > 0 {
		iranOptimized := 0
		for _, cfg := range configs {
			if cfg.Config.IranOptimized {
				iranOptimized++
			}
		}
		iranRatio := float64(iranOptimized) / float64(len(configs)) * 100

		if iranRatio < 50 {
			color.Yellow("⚠️  Low Iran-optimized config ratio (%.1f%%)", iranRatio)
			color.Yellow("💡 Enable iran-mode and aggressive DPI evasion")
		} else {
			color.Green("✅ Good Iran-optimized config ratio (%.1f%%)", iranRatio)
		}
	}
}

func printProxyStats(checker *EnhancedProxyChecker, elapsed time.Duration, config *AppConfig) {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                   PROXY CHECK RESULTS                         ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Total Tested:        %-40d║", checker.stats.TotalTested)
	color.Green("║  Active Proxies:      %-40d║", checker.stats.TotalActive)
	color.Red("║  Failed Proxies:      %-40d║", checker.stats.TotalFailed)

	if config.IranMode {
		color.Magenta("║  Iran-Optimized:      %-40d║", checker.stats.IranOptimized)
		color.Magenta("║  DPI Evaded:          %-40d║", checker.stats.DPIEvaded)
		color.Magenta("║  Multi-Endpoint OK:   %-40d║", checker.stats.MultiEndpointOK)
	}

	color.Green("║  Countries Found:     %-40d║", len(checker.activeProxies))

	if checker.stats.TotalActive > 0 {
		successRate := float64(checker.stats.TotalActive) / float64(checker.stats.TotalTested) * 100
		color.Green("║  Success Rate:        %-37.1f%% ║", successRate)
	}

	color.Yellow("║  Time Elapsed:        %-40s║", elapsed.String())
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
}

func printFinalSummary(checker *EnhancedProxyChecker, generated []Config, passed []TestedConfig,
	elapsed time.Duration, config *AppConfig) {
	fmt.Println()
	color.Cyan("╔═══════════════════════════════════════════════════════════════╗")
	color.Green("║                  ✅ FINAL SUMMARY                             ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Status:              %-40s║", "SUCCESS")
	color.Green("║  Total Time:          %-40s║", elapsed.String())
	color.Green("║  Performance Mode:    %-40s║", strings.ToUpper(config.PerformanceMode))
	color.Green("║  DPI Evasion Level:   %-40s║", strings.ToUpper(config.DPIEvasionLevel))
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Active Proxies:      %-40d║", checker.stats.TotalActive)
	color.Yellow("║  Generated Configs:   %-40d║", len(generated))
	color.Green("║  Passed Configs:      %-40d║", len(passed))

	if config.IranMode && len(passed) > 0 {
		iranCount := 0
		for _, tc := range passed {
			if tc.Config.IranOptimized {
				iranCount++
			}
		}
		color.Magenta("║  Iran-Optimized:      %-40d║", iranCount)

		if checker.stats.TotalActive > 0 {
			iranRatio := float64(iranCount) / float64(len(passed)) * 100
			color.Magenta("║  Iran Optimization:   %-37.1f%% ║", iranRatio)
		}
	}

	// Protocol breakdown
	if len(passed) > 0 {
		protocolCounts := make(map[string]int)
		for _, tc := range passed {
			protocolCounts[tc.Config.Protocol]++
		}

		color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
		color.Cyan("║                  PROTOCOL BREAKDOWN                           ║")
		color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

		for protocol, count := range protocolCounts {
			percentage := float64(count) / float64(len(passed)) * 100
			color.Yellow("║  %-20s: %4d (%.1f%%)                    ║",
				strings.ToUpper(protocol), count, percentage)
		}
	}

	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Performance insights
	if checker.stats.TotalActive > 0 {
		avgPerProxy := elapsed.Seconds() / float64(checker.stats.TotalTested)
		color.Green("⚡ Performance: %.2f seconds per proxy", avgPerProxy)
	}

	// Save summary to file
	summaryFile := "results/run-summary.txt"
	os.MkdirAll("results", 0755)
	if f, err := os.Create(summaryFile); err == nil {
		defer f.Close()
		fmt.Fprintf(f, "Run Summary - %s\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(f, "Version: %s\n", version)
		fmt.Fprintf(f, "Total Time: %s\n", elapsed.String())
		fmt.Fprintf(f, "Active Proxies: %d\n", checker.stats.TotalActive)
		fmt.Fprintf(f, "Generated Configs: %d\n", len(generated))
		fmt.Fprintf(f, "Passed Configs: %d\n", len(passed))
		color.Green("📄 Summary saved to: %s", summaryFile)
	}
}

func showHelp() {
	color.Cyan("╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║   🇮🇷 Iran Proxy & Config System - Enterprise Grade - HELP   ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	color.Yellow("DESCRIPTION:")
	fmt.Println("  Advanced proxy checking and VPN config generation system optimized")
	fmt.Println("  for bypassing Iran's strict DPI filtering with enterprise-grade features.")
	fmt.Println()

	color.Yellow("FEATURES:")
	fmt.Println("  🔐 Multi-Protocol       - VMess, VLESS, Shadowsocks, Trojan")
	fmt.Println("  🌐 Advanced Transports  - WebSocket, gRPC, xhttp, HTTP/2, QUIC")
	fmt.Println("  🔒 Security             - TLS, XTLS, Reality")
	fmt.Println("  🇮🇷 Iran DPI Bypass     - Multi-layer evasion techniques")
	fmt.Println("  🔧 Self-Healing         - Automatic recovery and optimization")
	fmt.Println("  📊 Monitoring           - Real-time metrics and analytics")
	fmt.Println("  ⚡ Performance Modes    - Speed, Balanced, Quality")
	fmt.Println("  🚨 Emergency Recovery   - Fallback systems for reliability")
	fmt.Println()

	color.Yellow("USAGE:")
	fmt.Println("  go run *.go [options]")
	fmt.Println()

	color.Yellow("CORE OPTIONS:")
	fmt.Println("  -proxy-file string")
	fmt.Println("        Path to proxy file (default: edge/assets/p-list-february.txt)")
	fmt.Println("  -output-file string")
	fmt.Println("        Path to output file (default: sub/ProxyIP-Daily.md)")
	fmt.Println("  -config-output string")
	fmt.Println("        Config output path (default: configs/iran-configs.txt)")
	fmt.Println()

	color.Yellow("PERFORMANCE OPTIONS:")
	fmt.Println("  -max-concurrent int")
	fmt.Println("        Max concurrent connections: 50-500 (default: 100)")
	fmt.Println("  -timeout int")
	fmt.Println("        Timeout in seconds: 5-30 (default: 10)")
	fmt.Println("  -performance-mode string")
	fmt.Println("        Performance mode: speed, balanced, quality (default: balanced)")
	fmt.Println()

	color.Yellow("IRAN OPTIMIZATION:")
	fmt.Println("  -iran-mode")
	fmt.Println("        Enable Iran-specific optimizations (default: true)")
	fmt.Println("  -dpi-evasion-level string")
	fmt.Println("        DPI evasion level: standard, aggressive, maximum (default: aggressive)")
	fmt.Println()

	color.Yellow("ADVANCED FEATURES:")
	fmt.Println("  -enable-fallback")
	fmt.Println("        Enable multi-tier fallback system (default: true)")
	fmt.Println("  -enable-self-healing")
	fmt.Println("        Enable automatic recovery (default: true)")
	fmt.Println("  -enable-monitoring")
	fmt.Println("        Enable metrics and analytics (default: true)")
	fmt.Println("  -test-protocols string")
	fmt.Println("        Protocols to test, comma-separated (default: vmess,vless,trojan,shadowsocks)")
	fmt.Println()

	color.Yellow("MODE FLAGS:")
	fmt.Println("  -generate-configs")
	fmt.Println("        Generate configs from proxies (default: true)")
	fmt.Println("  -test-configs")
	fmt.Println("        Test generated configs (default: true)")
	fmt.Println("  -configs-only")
	fmt.Println("        Only generate configs, skip proxy check (default: false)")
	fmt.Println("  -emergency-mode")
	fmt.Println("        Emergency recovery mode (default: false)")
	fmt.Println("  -deep-analysis")
	fmt.Println("        Deep analysis mode (default: false)")
	fmt.Println()

	color.Yellow("EXAMPLES:")
	fmt.Println()
	fmt.Println("  1. Full pipeline with Iran optimizations:")
	color.Green("     go run *.go -iran-mode=true -dpi-evasion-level=aggressive")
	fmt.Println()

	fmt.Println("  2. High-speed mode with maximum concurrency:")
	color.Green("     go run *.go -performance-mode=speed -max-concurrent=500")
	fmt.Println()

	fmt.Println("  3. Quality mode with deep analysis:")
	color.Green("     go run *.go -performance-mode=quality -deep-analysis=true")
	fmt.Println()

	fmt.Println("  4. Emergency recovery mode:")
	color.Green("     go run *.go -emergency-mode=true -enable-fallback=true")
	fmt.Println()

	fmt.Println("  5. Test specific protocols only:")
	color.Green("     go run *.go -test-protocols=vless,vmess -dpi-evasion-level=maximum")
	fmt.Println()

	color.Cyan("═══════════════════════════════════════════════════════════════")
	fmt.Println()
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func boolToStatus(b bool) string {
	if b {
		return "✅ ENABLED"
	}
	return "❌ DISABLED"
}
