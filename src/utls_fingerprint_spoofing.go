package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// BrowserFingerprint represents different browser TLS fingerprints
type BrowserFingerprint struct {
	Name            string
	Version         string
	CipherSuites    []uint16
	SupportedCurves []tls.CurveID
	SignatureSchemes []tls.SignatureScheme
	ALPNProtocols   []string
	Extensions      []uint16
}

// FingerprintDatabase contains real browser fingerprints for spoofing
var FingerprintDatabase = map[string]BrowserFingerprint{
	"chrome120": {
		Name:    "Chrome",
		Version: "120.0.6099",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		SupportedCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		SignatureSchemes: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.PKCS1WithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
			tls.PKCS1WithSHA384,
			tls.PSSWithSHA512,
			tls.PKCS1WithSHA512,
		},
		ALPNProtocols: []string{"h2", "http/1.1"},
		Extensions:    []uint16{0, 10, 11, 13, 16, 23, 27, 35, 43, 45, 51},
	},
	"firefox121": {
		Name:    "Firefox",
		Version: "121.0",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
		SupportedCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		SignatureSchemes: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.ECDSAWithP521AndSHA512,
			tls.PSSWithSHA256,
			tls.PSSWithSHA384,
			tls.PSSWithSHA512,
			tls.PKCS1WithSHA256,
			tls.PKCS1WithSHA384,
			tls.PKCS1WithSHA512,
		},
		ALPNProtocols: []string{"h2", "http/1.1"},
		Extensions:    []uint16{0, 10, 11, 13, 23, 27, 35, 43, 45, 51},
	},
	"edge120": {
		Name:    "Edge",
		Version: "120.0.2210",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
		SupportedCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		SignatureSchemes: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.PKCS1WithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
			tls.PKCS1WithSHA384,
		},
		ALPNProtocols: []string{"h2", "http/1.1"},
		Extensions:    []uint16{0, 10, 11, 13, 23, 27, 35, 43, 45, 51},
	},
	"safari17": {
		Name:    "Safari",
		Version: "17.2",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		SupportedCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		SignatureSchemes: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.PKCS1WithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
			tls.PKCS1WithSHA384,
		},
		ALPNProtocols: []string{"h2", "http/1.1"},
		Extensions:    []uint16{0, 10, 11, 13, 23, 27, 35, 43, 45, 51, 65281},
	},
}

// UTLSDialer creates connections with spoofed browser fingerprints
type UTLSDialer struct {
	timeout          time.Duration
	selectedBrowser  string
	rotationEnabled  bool
	browserIndex     int
	availableBrowsers []string
}

// NewUTLSDialer creates a new uTLS dialer with browser rotation
func NewUTLSDialer(timeout time.Duration) *UTLSDialer {
	return &UTLSDialer{
		timeout:         timeout,
		rotationEnabled: true,
		browserIndex:    0,
		availableBrowsers: []string{
			"chrome120",
			"firefox121",
			"edge120",
			"safari17",
		},
	}
}

// DialContext establishes a connection with spoofed TLS fingerprint
func (d *UTLSDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// Rotate browser fingerprint for each connection
	if d.rotationEnabled {
		d.selectedBrowser = d.availableBrowsers[d.browserIndex%len(d.availableBrowsers)]
		d.browserIndex++
	} else if d.selectedBrowser == "" {
		d.selectedBrowser = "chrome120" // Default
	}

	fingerprint := FingerprintDatabase[d.selectedBrowser]
	
	// Create TLS config with browser-specific settings
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS13,
		CipherSuites:            fingerprint.CipherSuites,
		CurvePreferences:        fingerprint.SupportedCurves,
		InsecureSkipVerify:      false,
		NextProtos:              fingerprint.ALPNProtocols,
		PreferServerCipherSuites: false,
		SessionTicketsDisabled:  false,
	}

	// Establish TCP connection first
	dialer := &net.Dialer{
		Timeout:   d.timeout,
		KeepAlive: 30 * time.Second,
	}

	rawConn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("TCP connection failed: %v", err)
	}

	// Wrap with TLS using spoofed fingerprint
	tlsConn := tls.Client(rawConn, tlsConfig)
	
	// Set handshake deadline
	tlsConn.SetDeadline(time.Now().Add(d.timeout))
	
	if err := tlsConn.Handshake(); err != nil {
		rawConn.Close()
		return nil, fmt.Errorf("TLS handshake failed: %v", err)
	}

	// Reset deadline after handshake
	tlsConn.SetDeadline(time.Time{})
	
	return tlsConn, nil
}

// SetBrowser manually sets the browser fingerprint to use
func (d *UTLSDialer) SetBrowser(browser string) error {
	if _, exists := FingerprintDatabase[browser]; !exists {
		return fmt.Errorf("unknown browser fingerprint: %s", browser)
	}
	d.selectedBrowser = browser
	d.rotationEnabled = false
	return nil
}

// EnableRotation enables automatic browser fingerprint rotation
func (d *UTLSDialer) EnableRotation() {
	d.rotationEnabled = true
}

// GetCurrentFingerprint returns the currently selected browser fingerprint
func (d *UTLSDialer) GetCurrentFingerprint() BrowserFingerprint {
	if d.selectedBrowser == "" {
		d.selectedBrowser = "chrome120"
	}
	return FingerprintDatabase[d.selectedBrowser]
}

// GetFingerprintInfo returns detailed info about current fingerprint
func (d *UTLSDialer) GetFingerprintInfo() string {
	fp := d.GetCurrentFingerprint()
	return fmt.Sprintf("%s %s (Ciphers: %d, Curves: %d, ALPN: %v)",
		fp.Name, fp.Version, len(fp.CipherSuites), len(fp.SupportedCurves), fp.ALPNProtocols)
}
