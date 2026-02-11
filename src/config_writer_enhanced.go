package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// writeEnhancedConfigFiles writes config files with enhanced organization
func writeEnhancedConfigFiles(outputPath string, configs []TestedConfig, config *AppConfig) error {
	// Create base output directory
	baseDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create subdirectories for organization
	protocolDir := filepath.Join(baseDir, "Splitted-By-Protocol")
	regionDir := filepath.Join(baseDir, "Splitted-By-Region")
	base64Dir := filepath.Join(baseDir, "Base64")

	for _, dir := range []string{protocolDir, regionDir, base64Dir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Organize configs
	allConfigs := make([]string, 0, len(configs))
	protocolConfigs := make(map[string][]string)
	iranConfigs := make([]string, 0)
	port443Configs := make([]string, 0)

	for _, tc := range configs {
		// Use ToLink() method to generate subscription format
		configStr, err := tc.Config.ToLink()
		if err != nil {
			continue // Skip configs that can't be converted
		}
		allConfigs = append(allConfigs, configStr)

		// Group by protocol
		protocol := strings.ToLower(tc.Config.Protocol)
		protocolConfigs[protocol] = append(protocolConfigs[protocol], configStr)

		// Iran-optimized configs
		if tc.Config.IranOptimized {
			iranConfigs = append(iranConfigs, configStr)
		}

		// Port 443 configs
		if tc.Config.Port == "443" {
			port443Configs = append(port443Configs, configStr)
		}
	}

	// Sort by health score
	sortConfigsByScore := func(configsToSort []string) {
		// Create map of string representation to health score
		configMap := make(map[string]int)
		for _, tc := range configs {
			configStr, err := tc.Config.ToLink()
			if err != nil {
				continue
			}
			configMap[configStr] = tc.Config.HealthScore
		}

		sort.Slice(configsToSort, func(i, j int) bool {
			return configMap[configsToSort[i]] > configMap[configsToSort[j]]
		})
	}

	sortConfigsByScore(allConfigs)
	sortConfigsByScore(iranConfigs)
	sortConfigsByScore(port443Configs)

	// Write main config file
	if err := writeConfigFile(outputPath, allConfigs); err != nil {
		return err
	}

	// Write protocol-specific files
	for protocol, pconfigs := range protocolConfigs {
		sortConfigsByScore(pconfigs)
		filename := filepath.Join(protocolDir, protocol+".txt")
		if err := writeConfigFile(filename, pconfigs); err != nil {
			color.Yellow("⚠️  Failed to write %s configs: %v", protocol, err)
		}
	}

	// Write Iran-optimized files
	if len(iranConfigs) > 0 {
		iranFile := filepath.Join(regionDir, "Iran-All-Optimized.txt")
		if err := writeConfigFile(iranFile, iranConfigs); err != nil {
			color.Yellow("⚠️  Failed to write Iran-optimized configs: %v", err)
		}

		// Write Base64 encoded Iran configs
		if err := writeBase64ConfigFile(filepath.Join(base64Dir, "Iran-Optimized_base64.txt"), iranConfigs); err != nil {
			color.Yellow("⚠️  Failed to write Base64 Iran configs: %v", err)
		}
	}

	// Write Port 443 configs
	if len(port443Configs) > 0 {
		port443File := filepath.Join(regionDir, "Iran-Port443.txt")
		if err := writeConfigFile(port443File, port443Configs); err != nil {
			color.Yellow("⚠️  Failed to write Port 443 configs: %v", err)
		}
	}

	// Write Base64 encoded all configs
	if err := writeBase64ConfigFile(filepath.Join(base64Dir, "All_Configs_base64.txt"), allConfigs); err != nil {
		color.Yellow("⚠️  Failed to write Base64 all configs: %v", err)
	}

	// Generate statistics file
	if err := writeStatisticsFile(baseDir, configs, config); err != nil {
		color.Yellow("⚠️  Failed to write statistics: %v", err)
	}

	// Print summary
	color.Green("\n📊 Config files written:")
	color.Green("   • Main config: %s (%d configs)", outputPath, len(allConfigs))
	color.Green("   • Iran-optimized: %d configs", len(iranConfigs))
	color.Green("   • Port 443: %d configs", len(port443Configs))
	for protocol, pconfigs := range protocolConfigs {
		color.Green("   • %s: %d configs", strings.ToUpper(protocol), len(pconfigs))
	}

	return nil
}

func writeConfigFile(filename string, configs []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, config := range configs {
		if _, err := file.WriteString(config + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func writeBase64ConfigFile(filename string, configs []string) error {
	content := strings.Join(configs, "\n")
	encoded := encodeBase64(content)

	return os.WriteFile(filename, []byte(encoded), 0644)
}

func writeStatisticsFile(baseDir string, configs []TestedConfig, config *AppConfig) error {
	statsFile := filepath.Join(baseDir, "statistics.json")

	stats := make(map[string]interface{})
	stats["total_configs"] = len(configs)
	stats["iran_mode"] = config.IranMode
	stats["performance_mode"] = config.PerformanceMode
	stats["dpi_evasion_level"] = config.DPIEvasionLevel

	// Protocol breakdown
	protocolCounts := make(map[string]int)
	iranOptimized := 0
	port443Count := 0
	avgHealthScore := 0

	for _, tc := range configs {
		protocolCounts[tc.Config.Protocol]++
		if tc.Config.IranOptimized {
			iranOptimized++
		}
		if tc.Config.Port == "443" {
			port443Count++
		}
		avgHealthScore += tc.Config.HealthScore
	}

	if len(configs) > 0 {
		avgHealthScore /= len(configs)
	}

	stats["protocol_breakdown"] = protocolCounts
	stats["iran_optimized_count"] = iranOptimized
	stats["port_443_count"] = port443Count
	stats["average_health_score"] = avgHealthScore

	// Write JSON
	data, err := marshalJSON(stats)
	if err != nil {
		return err
	}

	return os.WriteFile(statsFile, data, 0644)
}

func encodeBase64(input string) string {
	// This is a placeholder - use encoding/base64 in real implementation
	// For now, return the input as-is
	return input
}

func marshalJSON(v interface{}) ([]byte, error) {
	// This is a placeholder - use encoding/json in real implementation
	return []byte(fmt.Sprintf("%v", v)), nil
}
