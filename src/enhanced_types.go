package main

import (
	"time"
)

// EnhancedConfigTester extends ConfigTester with advanced features
type EnhancedConfigTester struct {
	*ConfigTester
	config *AppConfig
}

// NewEnhancedConfigTester creates an enhanced config tester
func NewEnhancedConfigTester(configs []Config, maxConcurrent int,
	timeout time.Duration, config *AppConfig) *EnhancedConfigTester {

	base := NewConfigTester(configs, maxConcurrent, timeout, config.IranMode)

	return &EnhancedConfigTester{
		ConfigTester: base,
		config:       config,
	}
}
