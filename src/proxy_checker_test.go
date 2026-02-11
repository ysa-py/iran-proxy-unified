package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewProxyChecker(t *testing.T) {
	checker := NewProxyChecker("test.txt", "output.md", 50, 5*time.Second)

	if checker.proxyFile != "test.txt" {
		t.Errorf("Expected proxyFile to be 'test.txt', got '%s'", checker.proxyFile)
	}

	if checker.outputFile != "output.md" {
		t.Errorf("Expected outputFile to be 'output.md', got '%s'", checker.outputFile)
	}

	if checker.maxConcurrent != 50 {
		t.Errorf("Expected maxConcurrent to be 50, got %d", checker.maxConcurrent)
	}

	if checker.timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5s, got %v", checker.timeout)
	}

	if checker.activeProxies == nil {
		t.Error("Expected activeProxies map to be initialized")
	}
}

func TestReadProxyFile(t *testing.T) {
	// Create a temporary test file
	testFile := "test_proxies.txt"
	content := `1.1.1.1,443,US,Cloudflare
8.8.8.8,443,US,Google
192.168.1.1,80,US,Local
9.9.9.9,443,US,Unknown ISP
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	checker := NewProxyChecker(testFile, "output.md", 50, 5*time.Second)
	proxies, err := checker.ReadProxyFile()

	if err != nil {
		t.Fatalf("Failed to read proxy file: %v", err)
	}

	// Should filter out the 192.168.1.1 (port 80) and 9.9.9.9 (unknown ISP)
	// Only Cloudflare and Google should pass
	if len(proxies) != 2 {
		t.Errorf("Expected 2 proxies, got %d", len(proxies))
	}

	// Check if filtered proxies contain only valid ones
	for _, proxy := range proxies {
		parts := strings.Split(proxy, ",")
		if len(parts) < 4 {
			t.Errorf("Invalid proxy format: %s", proxy)
			continue
		}

		port := strings.TrimSpace(parts[1])
		if port != "443" {
			t.Errorf("Expected port 443, got %s", port)
		}

		isp := parts[3]
		isGoodISP := false
		for _, goodISP := range GoodISPs {
			if strings.Contains(isp, goodISP) {
				isGoodISP = true
				break
			}
		}

		if !isGoodISP {
			t.Errorf("ISP %s should not pass filter", isp)
		}
	}
}

func TestEncodeBadgeLabel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello%20World"},
		{"Time: 12:30", "Time%3A%2012%3A30"},
		{"Value (UTC+3:30)", "Value%20%28UTC%2B3%3A30%29"},
		{"A, B, C", "A%2C%20B%2C%20C"},
	}

	for _, test := range tests {
		result := encodeBadgeLabel(test.input)
		if result != test.expected {
			t.Errorf("encodeBadgeLabel(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestGetCountryFlag(t *testing.T) {
	tests := []struct {
		code     string
		expected int // Number of flag characters (should be 2 for valid country codes)
	}{
		{"US", 2},
		{"IR", 2},
		{"GB", 2},
		{"", 0},
	}

	for _, test := range tests {
		result := getCountryFlag(test.code)
		// Count runes in result
		count := len([]rune(result))
		if count != test.expected {
			t.Errorf("getCountryFlag(%q) returned %d characters, want %d", test.code, count, test.expected)
		}
	}
}

func TestGetCountryName(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"US", "United States"},
		{"IR", "Iran"},
		{"GB", "United Kingdom"},
		{"XX", "XX"}, // Unknown code should return the code itself
	}

	for _, test := range tests {
		result := getCountryName(test.code)
		if result != test.expected {
			t.Errorf("getCountryName(%q) = %q, want %q", test.code, result, test.expected)
		}
	}
}

func TestGetLatencyEmoji(t *testing.T) {
	tests := []struct {
		latency  int64
		expected string
	}{
		{50, "⚡"},
		{500, "⚡"},
		{1098, "⚡"},
		{1100, "🐇"},
		{1500, "🐇"},
		{1598, "🐇"},
		{1600, "🐌"},
		{2000, "🐌"},
	}

	for _, test := range tests {
		result := getLatencyEmoji(test.latency)
		if result != test.expected {
			t.Errorf("getLatencyEmoji(%d) = %q, want %q", test.latency, result, test.expected)
		}
	}
}

func TestGetProviderLogoHTML(t *testing.T) {
	tests := []struct {
		provider string
		hasLogo  bool
	}{
		{"Google", true},
		{"Amazon", true},
		{"Cloudflare", true},
		{"UnknownProvider", false},
	}

	for _, test := range tests {
		result := getProviderLogoHTML(test.provider)
		if test.hasLogo && result == "" {
			t.Errorf("Expected logo HTML for %s, got empty string", test.provider)
		}
		if !test.hasLogo && result != "" {
			t.Errorf("Expected no logo HTML for %s, got %q", test.provider, result)
		}
		if test.hasLogo && !strings.Contains(result, "img") {
			t.Errorf("Logo HTML should contain img tag, got %q", result)
		}
	}
}

func TestGetStringValue(t *testing.T) {
	m := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"city":  "Tehran",
		"valid": true,
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"name", "John"},
		{"city", "Tehran"},
		{"age", ""},      // Not a string
		{"valid", ""},    // Not a string
		{"missing", ""},  // Key doesn't exist
	}

	for _, test := range tests {
		result := getStringValue(m, test.key)
		if result != test.expected {
			t.Errorf("getStringValue(m, %q) = %q, want %q", test.key, result, test.expected)
		}
	}
}

func TestProcessProxyTimeout(t *testing.T) {
	// This test ensures that proxy checking respects timeout
	checker := NewProxyChecker("test.txt", "output.md", 1, 1*time.Second)
	checker.selfIP = "1.2.3.4"

	ctx := context.Background()

	// Use an IP that will timeout (non-routable IP)
	start := time.Now()
	checker.ProcessProxy(ctx, "192.0.2.1,443,US,TestISP")
	elapsed := time.Since(start)

	// Should timeout within reasonable time (allow some overhead)
	if elapsed > 3*time.Second {
		t.Errorf("Proxy check took too long: %v (expected < 3s)", elapsed)
	}
}

func TestGoodISPsNotEmpty(t *testing.T) {
	if len(GoodISPs) == 0 {
		t.Error("GoodISPs list should not be empty")
	}

	// Verify some expected ISPs are in the list
	expectedISPs := []string{"Google", "Amazon", "Cloudflare", "Hetzner"}
	for _, expected := range expectedISPs {
		found := false
		for _, isp := range GoodISPs {
			if isp == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected ISP %s not found in GoodISPs list", expected)
		}
	}
}

// Benchmark tests
func BenchmarkEncodeBadgeLabel(b *testing.B) {
	input := "Mon, 02 Jan 2006 15:04 (UTC+3:30)"
	for i := 0; i < b.N; i++ {
		encodeBadgeLabel(input)
	}
}

func BenchmarkGetCountryFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryFlag("US")
	}
}

func BenchmarkGetCountryName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCountryName("US")
	}
}
