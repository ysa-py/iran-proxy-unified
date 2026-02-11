package main

// This file is reference documentation only
// All executable code has been moved to main.go

// Original implementation reference (for documentation only):
/*
import (
	"flag"
	"fmt"
	"os"
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
║      🇮🇷 Advanced Iran Proxy Checker & Config Generator 🇮🇷              ║
║                                                                          ║
║     Proxy Check • Config Generation • Smart Testing • DPI Evasion       ║
║                                                                          ║
╚══════════════════════════════════════════════════════════════════════════╝
`

const version = "v2.0.0-iran-optimized"

func main() {
	// Command line flags
	proxyFile := flag.String("proxy-file", DefaultProxyFile, "Path to proxy file (CSV format)")
	outputFile := flag.String("output-file", DefaultOutputFile, "Path to output markdown file")
	configOutput := flag.String("config-output", "configs/iran-configs.txt", "Path to config output file")
	maxConcurrent := flag.Int("max-concurrent", DefaultMaxConcurrent, "Maximum concurrent connections")
	timeoutSecs := flag.Int("timeout", DefaultTimeoutSecs, "Timeout in seconds for each proxy check")
	iranMode := flag.Bool("iran-mode", true, "Enable Iran-specific optimizations")
	generateConfigs := flag.Bool("generate-configs", true, "Generate configs from proxies")
	testConfigs := flag.Bool("test-configs", true, "Test generated configs")
	configsOnly := flag.Bool("configs-only", false, "Only generate configs, skip proxy check")
	help := flag.Bool("help", false, "Show help message")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *showVersion {
		fmt.Printf("Proxy Checker & Config Generator %s\n", version)
		return
	}

	// Show banner
	showBanner()

	// Configuration summary
	printConfigSummary(*proxyFile, *outputFile, *configOutput, *maxConcurrent, *timeoutSecs, *iranMode, *generateConfigs, *testConfigs)

	// Create proxy checker
	timeout := time.Duration(*timeoutSecs) * time.Second
	checker := NewProxyChecker(*proxyFile, *outputFile, *maxConcurrent, timeout)
	checker.iranMode = *iranMode

	// Overall start time
	overallStartTime := time.Now()

	// Step 1: Run proxy checker (unless configs-only mode)
	if !*configsOnly {
		startTime := time.Now()

		color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
		color.Cyan("║                  STEP 1: PROXY CHECKING                       ║")
		color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

		if err := checker.Run(); err != nil {
			color.Red("❌ Error: %v", err)
			os.Exit(1)
		}

		// Write proxy results
		if err := checker.WriteMarkdownFile(); err != nil {
			color.Red("❌ Failed to write markdown file: %v", err)
			os.Exit(1)
		}

		elapsed := time.Since(startTime)
		color.Green("\n✅ Proxy checking completed in %s", elapsed.String())
		printProxyStats(checker, elapsed)
	}

	// Step 2: Generate configs from active proxies
	var allGeneratedConfigs []Config
	var allPassedConfigs []TestedConfig

	if *generateConfigs && (checker.stats.TotalActive > 0 || *configsOnly) {
		startTime := time.Now()

		color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
		color.Cyan("║                  STEP 2: CONFIG GENERATION                    ║")
		color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

		configCount := 0

		// Generate configs for each active proxy
		for _, proxyResults := range checker.activeProxies {
			for _, result := range proxyResults {
				color.Yellow("Generating configs for: %s:%s (%s)",
					result.Info.IP, result.Info.Port, result.Info.ISP)

				generator := NewConfigGenerator(result.Info, *iranMode)
				configs := generator.GenerateAllConfigs()
				generator.PrintConfigSummary()

				allGeneratedConfigs = append(allGeneratedConfigs, configs...)
				configCount += len(configs)
			}
		}

		elapsed := time.Since(startTime)
		color.Green("\n✅ Config generation completed in %s", elapsed.String())
		color.Green("📊 Total configs generated: %d", configCount)
	}

	// Step 3: Test generated configs
	if *testConfigs && len(allGeneratedConfigs) > 0 {
		startTime := time.Now()

		color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
		color.Cyan("║                  STEP 3: CONFIG TESTING                       ║")
		color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

		tester := NewConfigTester(allGeneratedConfigs, *maxConcurrent, timeout, *iranMode)
		allPassedConfigs = tester.TestAllConfigs()
		tester.PrintConfigDetails()

		elapsed := time.Since(startTime)
		color.Green("\n✅ Config testing completed in %s", elapsed.String())
	}

	// Step 4: Write config output files
	if len(allPassedConfigs) > 0 {
		color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
		color.Cyan("║                  STEP 4: WRITING CONFIG FILES                 ║")
		color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

		if err := writeConfigFiles(*configOutput, allPassedConfigs); err != nil {
			color.Red("❌ Failed to write config files: %v", err)
			os.Exit(1)
		}

		color.Green("\n✅ Config files written successfully")
	}

	// Final summary
	overallElapsed := time.Since(overallStartTime)
	printFinalSummary(checker, allGeneratedConfigs, allPassedConfigs, overallElapsed)
}

func showBanner() {
	color.Cyan(banner)
	color.White("Version: %s", version)
	fmt.Println()
}

func printConfigSummary(proxyFile, outputFile, configOutput string, maxConcurrent, timeout int, iranMode, generateConfigs, testConfigs bool) {
	color.Cyan("╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                    CONFIGURATION SUMMARY                      ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Proxy File:          %-40s║", truncate(proxyFile, 40))
	color.Yellow("║  Output File:         %-40s║", truncate(outputFile, 40))
	color.Yellow("║  Config Output:       %-40s║", truncate(configOutput, 40))
	color.Yellow("║  Max Concurrent:      %-40d║", maxConcurrent)
	color.Yellow("║  Timeout:             %-37ds ║", timeout)

	if iranMode {
		color.Magenta("║  Iran Mode:           %-40s║", "🇮🇷 ENABLED")
		color.Magenta("║  DPI Evasion:         %-40s║", "✅ ACTIVE")
		color.Magenta("║  Multi-Endpoint:      %-40s║", "✅ ACTIVE")
	} else {
		color.Yellow("║  Iran Mode:           %-40s║", "DISABLED")
	}

	color.Yellow("║  Generate Configs:    %-40s║", boolToStatus(generateConfigs))
	color.Yellow("║  Test Configs:        %-40s║", boolToStatus(testConfigs))

	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printProxyStats(checker *ProxyChecker, elapsed time.Duration) {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                   PROXY CHECK RESULTS                         ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Total Tested:        %-40d║", checker.stats.TotalTested)
	color.Green("║  Active Proxies:      %-40d║", checker.stats.TotalActive)
	color.Red("║  Failed Proxies:      %-40d║", checker.stats.TotalFailed)

	if checker.iranMode {
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

func printFinalSummary(checker *ProxyChecker, generatedConfigs []Config, passedConfigs []TestedConfig, elapsed time.Duration) {
	fmt.Println()
	color.Cyan("╔═══════════════════════════════════════════════════════════════╗")
	color.Green("║                  ✅ FINAL SUMMARY                             ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Status:              %-40s║", "SUCCESS")
	color.Green("║  Total Time:          %-40s║", elapsed.String())
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Yellow("║  Active Proxies:      %-40d║", checker.stats.TotalActive)
	color.Yellow("║  Generated Configs:   %-40d║", len(generatedConfigs))
	color.Green("║  Passed Configs:      %-40d║", len(passedConfigs))

	if checker.iranMode {
		iranCount := 0
		for _, tc := range passedConfigs {
			if tc.Config.IranOptimized {
				iranCount++
			}
		}
		color.Magenta("║  Iran-Optimized:      %-40d║", iranCount)
	}

	// Count by protocol
	protocolCounts := make(map[string]int)
	for _, tc := range passedConfigs {
		protocolCounts[tc.Config.Protocol]++
	}

	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Cyan("║                  PROTOCOL BREAKDOWN                           ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	for protocol, count := range protocolCounts {
		color.Yellow("║  %-20s: %-37d║", protocol, count)
	}

	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func showHelp() {
	color.Cyan("╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║   🇮🇷 Advanced Iran Proxy Checker & Config Generator - HELP   ║")
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	color.Yellow("DESCRIPTION:")
	fmt.Println("  This tool checks proxy IPs, generates VPN configs for multiple protocols,")
	fmt.Println("  and tests them with Iran-specific optimizations for bypassing DPI filtering.")
	fmt.Println()

	color.Yellow("FEATURES:")
	fmt.Println("  🔐 Multi-Protocol Support - VMess, VLESS, Shadowsocks, Trojan")
	fmt.Println("  🌐 Advanced Transports    - WebSocket, gRPC, xhttp, HTTP/2, QUIC")
	fmt.Println("  🔒 Security Options       - TLS, XTLS, Reality")
	fmt.Println("  🧪 Smart Config Testing   - Automatic validation before output")
	fmt.Println("  🇮🇷 Iran Filter Bypass    - DPI evasion, multi-endpoint testing")
	fmt.Println("  📊 Health Scoring         - Intelligent ranking (0-100%)")
	fmt.Println("  ⚡ Concurrent Processing  - Fast parallel operations")
	fmt.Println()

	color.Yellow("USAGE:")
	fmt.Println("  go run *.go [options]")
	fmt.Println()

	color.Yellow("OPTIONS:")
	flag.PrintDefaults()
	fmt.Println()

	color.Yellow("EXAMPLES:")
	fmt.Println()
	fmt.Println("  1. Full pipeline (proxy check + config generation + testing):")
	color.Green("     go run *.go")
	fmt.Println()

	fmt.Println("  2. Only generate and test configs (skip proxy check):")
	color.Green("     go run *.go -configs-only=true")
	fmt.Println()

	fmt.Println("  3. High-performance mode:")
	color.Green("     go run *.go -max-concurrent=200 -timeout=5")
	fmt.Println()

	fmt.Println("  4. Generate configs without testing:")
	color.Green("     go run *.go -test-configs=false")
	fmt.Println()

	fmt.Println("  5. Disable Iran optimizations:")
	color.Green("     go run *.go -iran-mode=false")
	fmt.Println()

	color.Yellow("SUPPORTED PROTOCOLS:")
	fmt.Println("  • VMess      - with TCP, WebSocket, gRPC, HTTP/2, xhttp")
	fmt.Println("  • VLESS      - with Reality, XTLS, TLS, multiple transports")
	fmt.Println("  • Shadowsocks - with modern ciphers (chacha20, aes-gcm, 2022)")
	fmt.Println("  • Trojan     - with TLS, WebSocket, gRPC, xhttp")
	fmt.Println()

	color.Yellow("IRAN-SPECIFIC OPTIMIZATIONS:")
	fmt.Println("  ✅ xhttp transport (best DPI evasion)")
	fmt.Println("  ✅ Reality protocol (undetectable)")
	fmt.Println("  ✅ Port 443 filtering")
	fmt.Println("  ✅ Whitelisted ISPs (CDN & cloud providers)")
	fmt.Println("  ✅ Multi-endpoint validation")
	fmt.Println("  ✅ TLS fingerprinting (chrome, firefox, random)")
	fmt.Println("  ✅ SNI obfuscation")
	fmt.Println()

	color.Yellow("OUTPUT FILES:")
	fmt.Println("  • Proxy list: sub/ProxyIP-Daily.md (markdown format)")
	fmt.Println("  • Configs: configs/iran-configs.txt (subscription format)")
	fmt.Println("  • By protocol: configs/vmess.txt, configs/vless.txt, etc.")
	fmt.Println("  • Iran-optimized: configs/iran-optimized.txt")
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
}*/
// Reference documentation only - actual implementation is in main.go
