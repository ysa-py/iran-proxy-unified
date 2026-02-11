package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

// writeConfigFiles writes configs to various output files
func writeConfigFiles(baseOutput string, configs []TestedConfig) error {
	// Create output directory
	outputDir := filepath.Dir(baseOutput)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create subdirectories
	protocolDir := filepath.Join(outputDir, "by-protocol")
	regionDir := filepath.Join(outputDir, "by-region")
	base64Dir := filepath.Join(outputDir, "base64")

	for _, dir := range []string{protocolDir, regionDir, base64Dir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	color.Cyan("📝 Writing config files...")

	// Write all configs to main file
	if err := writeAllConfigs(baseOutput, configs); err != nil {
		return err
	}

	// Write all configs in base64 format
	base64File := filepath.Join(base64Dir, "all-configs-base64.txt")
	if err := writeAllConfigsBase64(base64File, configs); err != nil {
		return err
	}

	// Write configs by protocol
	if err := writeByProtocol(protocolDir, configs); err != nil {
		return err
	}

	// Write Iran-optimized configs
	iranFile := filepath.Join(regionDir, "iran-optimized.txt")
	if err := writeIranOptimized(iranFile, configs); err != nil {
		return err
	}

	// Write configs by country
	if err := writeByCountry(regionDir, configs); err != nil {
		return err
	}

	// Write subscription files
	if err := writeSubscriptionFiles(outputDir, configs); err != nil {
		return err
	}

	// Print summary
	printWriteSummary(outputDir, configs)

	return nil
}

// writeAllConfigs writes all configs to a single file
func writeAllConfigs(filename string, configs []TestedConfig) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# 🇮🇷 Iran-Optimized VPN Configs\n")
	fmt.Fprintf(file, "# Generated: %s\n", getCurrentTime())
	fmt.Fprintf(file, "# Total Configs: %d\n", len(configs))
	fmt.Fprintf(file, "#\n")
	fmt.Fprintf(file, "# Protocols: VMess, VLESS, Shadowsocks, Trojan\n")
	fmt.Fprintf(file, "# Transports: TCP, WebSocket, gRPC, HTTP/2, xhttp, QUIC\n")
	fmt.Fprintf(file, "# Security: TLS, XTLS, Reality\n")
	fmt.Fprintf(file, "#\n\n")

	// Write configs
	for i, tested := range configs {
		link, err := tested.Config.ToLink()
		if err != nil {
			color.Yellow("⚠️  Skipping config %d: %v", i+1, err)
			continue
		}
		fmt.Fprintf(file, "%s\n", link)
	}

	color.Green("✅ Written: %s (%d configs)", filename, len(configs))
	return nil
}

// writeAllConfigsBase64 writes all configs in base64 subscription format
func writeAllConfigsBase64(filename string, configs []TestedConfig) error {
	var allLinks []string

	for _, tested := range configs {
		link, err := tested.Config.ToLink()
		if err != nil {
			continue
		}
		allLinks = append(allLinks, link)
	}

	// Join with newlines and encode
	content := strings.Join(allLinks, "\n")
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	// Write to file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", encoded)

	color.Green("✅ Written: %s (base64 encoded)", filename)
	return nil
}

// writeByProtocol writes configs grouped by protocol
func writeByProtocol(dir string, configs []TestedConfig) error {
	byProtocol := make(map[string][]TestedConfig)

	for _, tested := range configs {
		protocol := tested.Config.Protocol
		byProtocol[protocol] = append(byProtocol[protocol], tested)
	}

	for protocol, protocolConfigs := range byProtocol {
		filename := filepath.Join(dir, fmt.Sprintf("%s.txt", protocol))
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		// Write header
		fmt.Fprintf(file, "# %s Configs\n", strings.ToUpper(protocol))
		fmt.Fprintf(file, "# Total: %d\n", len(protocolConfigs))
		fmt.Fprintf(file, "# Generated: %s\n\n", getCurrentTime())

		// Write configs
		for _, tested := range protocolConfigs {
			link, err := tested.Config.ToLink()
			if err != nil {
				continue
			}
			fmt.Fprintf(file, "%s\n", link)
		}

		file.Close()
		color.Green("✅ Written: %s (%d configs)", filename, len(protocolConfigs))
	}

	return nil
}

// writeIranOptimized writes Iran-optimized configs
func writeIranOptimized(filename string, configs []TestedConfig) error {
	var iranConfigs []TestedConfig

	for _, tested := range configs {
		if tested.Config.IranOptimized {
			iranConfigs = append(iranConfigs, tested)
		}
	}

	if len(iranConfigs) == 0 {
		color.Yellow("⚠️  No Iran-optimized configs to write")
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# 🇮🇷 Iran-Optimized Configs\n")
	fmt.Fprintf(file, "# Best configs for bypassing Iran filtering\n")
	fmt.Fprintf(file, "# Total: %d\n", len(iranConfigs))
	fmt.Fprintf(file, "# Generated: %s\n\n", getCurrentTime())

	// Sort by health score
	for i := 0; i < len(iranConfigs); i++ {
		for j := i + 1; j < len(iranConfigs); j++ {
			if iranConfigs[i].Config.HealthScore < iranConfigs[j].Config.HealthScore {
				iranConfigs[i], iranConfigs[j] = iranConfigs[j], iranConfigs[i]
			}
		}
	}

	// Write configs with health score comments
	for _, tested := range iranConfigs {
		link, err := tested.Config.ToLink()
		if err != nil {
			continue
		}

		// Add comment with details
		fmt.Fprintf(file, "# Health: %d%% | Latency: %dms | %s\n",
			tested.Config.HealthScore, tested.AvgLatency, tested.Config.Remark)
		fmt.Fprintf(file, "%s\n\n", link)
	}

	color.Green("✅ Written: %s (%d Iran-optimized configs)", filename, len(iranConfigs))
	return nil
}

// writeByCountry writes configs grouped by country
func writeByCountry(dir string, configs []TestedConfig) error {
	byCountry := make(map[string][]TestedConfig)

	for _, tested := range configs {
		country := tested.Config.Country
		if country == "" {
			country = "Unknown"
		}
		byCountry[country] = append(byCountry[country], tested)
	}

	for country, countryConfigs := range byCountry {
		filename := filepath.Join(dir, fmt.Sprintf("%s.txt", strings.ToLower(country)))
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		flag := getCountryFlag(country)
		countryName := getCountryName(country)

		// Write header
		fmt.Fprintf(file, "# %s %s Configs\n", flag, countryName)
		fmt.Fprintf(file, "# Total: %d\n", len(countryConfigs))
		fmt.Fprintf(file, "# Generated: %s\n\n", getCurrentTime())

		// Write configs
		for _, tested := range countryConfigs {
			link, err := tested.Config.ToLink()
			if err != nil {
				continue
			}
			fmt.Fprintf(file, "%s\n", link)
		}

		file.Close()
		color.Green("✅ Written: %s (%d configs)", filename, len(countryConfigs))
	}

	return nil
}

// writeSubscriptionFiles writes subscription format files
func writeSubscriptionFiles(dir string, configs []TestedConfig) error {
	// Write by transport type for Iran
	iranTransports := map[string][]TestedConfig{
		"xhttp":     {},
		"websocket": {},
		"grpc":      {},
		"reality":   {},
	}

	for _, tested := range configs {
		if !tested.Config.IranOptimized {
			continue
		}

		switch tested.Config.Network {
		case TransportXHTTP:
			iranTransports["xhttp"] = append(iranTransports["xhttp"], tested)
		case TransportWebSocket, TransportHTTPUpgrade:
			iranTransports["websocket"] = append(iranTransports["websocket"], tested)
		case TransportGRPC:
			iranTransports["grpc"] = append(iranTransports["grpc"], tested)
		}

		if tested.Config.Security == SecurityReality {
			iranTransports["reality"] = append(iranTransports["reality"], tested)
		}
	}

	// Write each transport type
	for transport, transportConfigs := range iranTransports {
		if len(transportConfigs) == 0 {
			continue
		}

		filename := filepath.Join(dir, fmt.Sprintf("iran-%s.txt", transport))
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		fmt.Fprintf(file, "# 🇮🇷 Iran %s Configs\n", strings.ToUpper(transport))
		fmt.Fprintf(file, "# Total: %d\n\n", len(transportConfigs))

		for _, tested := range transportConfigs {
			link, err := tested.Config.ToLink()
			if err != nil {
				continue
			}
			fmt.Fprintf(file, "%s\n", link)
		}

		file.Close()
		color.Green("✅ Written: %s (%d configs)", filename, len(transportConfigs))
	}

	return nil
}

// printWriteSummary prints a summary of written files
func printWriteSummary(dir string, configs []TestedConfig) {
	color.Cyan("\n╔═══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                   FILE WRITE SUMMARY                          ║")
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	// Count by protocol
	protocolCounts := make(map[string]int)
	iranCount := 0

	for _, tested := range configs {
		protocolCounts[tested.Config.Protocol]++
		if tested.Config.IranOptimized {
			iranCount++
		}
	}

	color.Green("║  Total Configs Written:     %-34d║", len(configs))
	color.Magenta("║  Iran-Optimized Configs:    %-34d║", iranCount)
	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")

	for protocol, count := range protocolCounts {
		color.Yellow("║  %-28s: %-31d║", strings.ToUpper(protocol), count)
	}

	color.Cyan("╠═══════════════════════════════════════════════════════════════╣")
	color.Green("║  Output Directory:          %-34s║", truncate(dir, 34))
	color.Cyan("╚═══════════════════════════════════════════════════════════════╝\n")

	// Print file locations
	color.Yellow("\n📂 Generated Files:")
	color.White("   Main configs:        %s", filepath.Join(dir, "iran-configs.txt"))
	color.White("   By protocol:         %s", filepath.Join(dir, "by-protocol/"))
	color.White("   By region:           %s", filepath.Join(dir, "by-region/"))
	color.White("   Base64 subscription: %s", filepath.Join(dir, "base64/"))
	color.White("   Iran-optimized:      %s", filepath.Join(dir, "by-region/iran-optimized.txt"))
	fmt.Println()
}

// getCurrentTime returns formatted current time in Tehran timezone
func getCurrentTime() string {
	// In production, use proper timezone handling
	// For now, return simple format
	return time.Now().Format("2006-01-02 15:04:05")
}
