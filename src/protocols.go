package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Protocol types
const (
	ProtocolVMess       = "vmess"
	ProtocolVLESS       = "vless"
	ProtocolShadowsocks = "shadowsocks"
	ProtocolTrojan      = "trojan"
)

// Transport types (Iran-optimized)
const (
	TransportTCP         = "tcp"
	TransportWebSocket   = "ws"
	TransportGRPC        = "grpc"
	TransportHTTP2       = "h2"
	TransportQUIC        = "quic"
	TransportXHTTP       = "xhttp"       // New xhttp transport for better DPI evasion
	TransportHTTPUpgrade = "httpupgrade" // HTTPUpgrade transport
)

// Security types
const (
	SecurityNone    = "none"
	SecurityTLS     = "tls"
	SecurityXTLS    = "xtls"
	SecurityReality = "reality"
)

// TLS fingerprints (Iran-optimized for better obfuscation)
var TLSFingerprints = []string{
	"chrome",
	"firefox",
	"safari",
	"ios",
	"android",
	"edge",
	"360",
	"qq",
	"random",
	"randomized",
}

// Iran-optimized cipher suites
var ShadowsocksCiphers = []string{
	"chacha20-ietf-poly1305", // Best for Iran
	"aes-256-gcm",
	"aes-128-gcm",
	"2022-blake3-aes-256-gcm", // Modern cipher
	"2022-blake3-aes-128-gcm",
}

// Iran-optimized SNI configurations
var IranOptimizedSNIs = []string{
	"www.speedtest.net",
	"www.cloudflare.com",
	"discord.com",
	"www.google.com",
	"www.microsoft.com",
	"fast.com",
	"www.apple.com",
	"zula.ir",
	"www.visa.com",
	"laravel.com",
}

// Config represents a universal proxy configuration
type Config struct {
	Protocol   string `json:"protocol"`
	Address    string `json:"address"`
	Port       string `json:"port"`
	ID         string `json:"id"`         // UUID for VMess/VLESS
	Password   string `json:"password"`   // For Trojan/SS
	AlterID    int    `json:"alterId"`    // VMess only
	Security   string `json:"security"`   // none, tls, xtls, reality
	Network    string `json:"network"`    // tcp, ws, grpc, h2, quic, xhttp
	Flow       string `json:"flow"`       // VLESS flow control
	Encryption string `json:"encryption"` // VLESS encryption
	Method     string `json:"method"`     // Shadowsocks cipher

	// Transport settings
	Path        string   `json:"path"`
	Host        string   `json:"host"`
	SNI         string   `json:"sni"`
	ServerName  string   `json:"serverName"` // Alias for SNI
	ALPN        []string `json:"alpn"`
	ServiceName string   `json:"serviceName"` // gRPC
	Mode        string   `json:"mode"`        // gRPC mode

	// Advanced settings
	Fingerprint   string `json:"fingerprint"` // TLS fingerprint
	AllowInsecure bool   `json:"allowInsecure"`
	PublicKey     string `json:"publicKey"` // Reality
	ShortID       string `json:"shortId"`   // Reality
	SpiderX       string `json:"spiderX"`   // Reality

	// Iran-specific optimizations
	IranOptimized bool `json:"iranOptimized"`
	HealthScore   int  `json:"healthScore"`
	DPIEvaded     bool `json:"dpiEvaded"`

	// DPI Evasion settings
	Headers             map[string]string `json:"headers"`
	Obfs                string            `json:"obfs"`
	ObfsPassword        string            `json:"obfsPassword"`
	PaddingSize         int               `json:"paddingSize"`
	FragmentationPoints []int             `json:"fragmentationPoints"`
	TimingDelay         int               `json:"timingDelay"`
	Key                 string            `json:"key"`

	// Plugin configuration
	Plugin     string `json:"plugin"`
	PluginOpts string `json:"pluginOpts"`

	// QUIC settings
	CongestionControl string `json:"congestionControl"`
	UDPRelayMode      string `json:"udpRelayMode"`

	// Metadata
	Remark  string `json:"remark"`
	ISP     string `json:"isp"`
	Country string `json:"country"`
}

// VMess specific structure
type VMess struct {
	V    string `json:"v"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  int    `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	ALPN string `json:"alpn"`
	FP   string `json:"fp"`
}

// VLESS specific structure
type VLESS struct {
	Address     string
	Port        string
	UUID        string
	Encryption  string
	Security    string
	Type        string
	Host        string
	Path        string
	SNI         string
	FP          string
	ALPN        string
	Flow        string
	PublicKey   string
	ShortID     string
	SpiderX     string
	ServiceName string
	Mode        string
	Remark      string
}

// Shadowsocks specific structure
type Shadowsocks struct {
	Method     string
	Password   string
	Server     string
	Port       string
	Plugin     string
	PluginOpts string
	Remark     string
}

// Trojan specific structure
type Trojan struct {
	Password string
	Address  string
	Port     string
	Security string
	Type     string
	Host     string
	Path     string
	SNI      string
	ALPN     string
	FP       string
	Remark   string
}

// GenerateUUID creates a random UUID v4
func GenerateUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// GeneratePassword creates a random password for Trojan/SS
func GeneratePassword(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)[:length]
}

// GenerateShortID creates random shortId for Reality
func GenerateShortID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// ToVMessLink converts Config to VMess link
func (c *Config) ToVMessLink() (string, error) {
	if c.Protocol != ProtocolVMess {
		return "", fmt.Errorf("not a VMess config")
	}

	vmess := VMess{
		V:    "2",
		PS:   c.Remark,
		Add:  c.Address,
		Port: c.Port,
		ID:   c.ID,
		Aid:  c.AlterID,
		Net:  c.Network,
		Type: "none",
		Host: c.Host,
		Path: c.Path,
		TLS:  c.Security,
		SNI:  c.SNI,
		FP:   c.Fingerprint,
	}

	if len(c.ALPN) > 0 {
		vmess.ALPN = strings.Join(c.ALPN, ",")
	}

	jsonBytes, err := json.Marshal(vmess)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return "vmess://" + encoded, nil
}

// ToVLESSLink converts Config to VLESS link
func (c *Config) ToVLESSLink() (string, error) {
	if c.Protocol != ProtocolVLESS {
		return "", fmt.Errorf("not a VLESS config")
	}

	params := make([]string, 0)
	params = append(params, fmt.Sprintf("encryption=%s", c.Encryption))
	params = append(params, fmt.Sprintf("security=%s", c.Security))
	params = append(params, fmt.Sprintf("type=%s", c.Network))

	if c.Security == SecurityTLS || c.Security == SecurityXTLS || c.Security == SecurityReality {
		params = append(params, fmt.Sprintf("sni=%s", c.SNI))
		if c.Fingerprint != "" {
			params = append(params, fmt.Sprintf("fp=%s", c.Fingerprint))
		}
		if len(c.ALPN) > 0 {
			params = append(params, fmt.Sprintf("alpn=%s", strings.Join(c.ALPN, ",")))
		}
	}

	if c.Security == SecurityReality {
		if c.PublicKey != "" {
			params = append(params, fmt.Sprintf("pbk=%s", c.PublicKey))
		}
		if c.ShortID != "" {
			params = append(params, fmt.Sprintf("sid=%s", c.ShortID))
		}
		if c.SpiderX != "" {
			params = append(params, fmt.Sprintf("spx=%s", c.SpiderX))
		}
	}

	if c.Flow != "" {
		params = append(params, fmt.Sprintf("flow=%s", c.Flow))
	}

	switch c.Network {
	case TransportWebSocket, TransportXHTTP, TransportHTTPUpgrade:
		if c.Host != "" {
			params = append(params, fmt.Sprintf("host=%s", c.Host))
		}
		if c.Path != "" {
			params = append(params, fmt.Sprintf("path=%s", c.Path))
		}
	case TransportGRPC:
		if c.ServiceName != "" {
			params = append(params, fmt.Sprintf("serviceName=%s", c.ServiceName))
		}
		if c.Mode != "" {
			params = append(params, fmt.Sprintf("mode=%s", c.Mode))
		}
	case TransportHTTP2:
		if c.Host != "" {
			params = append(params, fmt.Sprintf("host=%s", c.Host))
		}
		if c.Path != "" {
			params = append(params, fmt.Sprintf("path=%s", c.Path))
		}
	}

	link := fmt.Sprintf("vless://%s@%s:%s?%s#%s",
		c.ID, c.Address, c.Port, strings.Join(params, "&"), c.Remark)

	return link, nil
}

// ToShadowsocksLink converts Config to Shadowsocks link
func (c *Config) ToShadowsocksLink() (string, error) {
	if c.Protocol != ProtocolShadowsocks {
		return "", fmt.Errorf("not a Shadowsocks config")
	}

	userInfo := fmt.Sprintf("%s:%s", c.Method, c.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(userInfo))

	link := fmt.Sprintf("ss://%s@%s:%s#%s",
		encoded, c.Address, c.Port, c.Remark)

	return link, nil
}

// ToTrojanLink converts Config to Trojan link
func (c *Config) ToTrojanLink() (string, error) {
	if c.Protocol != ProtocolTrojan {
		return "", fmt.Errorf("not a Trojan config")
	}

	params := make([]string, 0)
	params = append(params, fmt.Sprintf("security=%s", c.Security))
	params = append(params, fmt.Sprintf("type=%s", c.Network))

	if c.Security == SecurityTLS {
		params = append(params, fmt.Sprintf("sni=%s", c.SNI))
		if c.Fingerprint != "" {
			params = append(params, fmt.Sprintf("fp=%s", c.Fingerprint))
		}
		if len(c.ALPN) > 0 {
			params = append(params, fmt.Sprintf("alpn=%s", strings.Join(c.ALPN, ",")))
		}
	}

	switch c.Network {
	case TransportWebSocket, TransportXHTTP:
		if c.Host != "" {
			params = append(params, fmt.Sprintf("host=%s", c.Host))
		}
		if c.Path != "" {
			params = append(params, fmt.Sprintf("path=%s", c.Path))
		}
	case TransportGRPC:
		if c.ServiceName != "" {
			params = append(params, fmt.Sprintf("serviceName=%s", c.ServiceName))
		}
	}

	link := fmt.Sprintf("trojan://%s@%s:%s?%s#%s",
		c.Password, c.Address, c.Port, strings.Join(params, "&"), c.Remark)

	return link, nil
}

// ToLink converts Config to appropriate protocol link
func (c *Config) ToLink() (string, error) {
	switch c.Protocol {
	case ProtocolVMess:
		return c.ToVMessLink()
	case ProtocolVLESS:
		return c.ToVLESSLink()
	case ProtocolShadowsocks:
		return c.ToShadowsocksLink()
	case ProtocolTrojan:
		return c.ToTrojanLink()
	default:
		return "", fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

// GetIranOptimizedScore calculates Iran optimization score
func (c *Config) GetIranOptimizedScore() int {
	score := 0

	// Protocol scoring
	switch c.Protocol {
	case ProtocolVLESS:
		score += 30 // VLESS is best for Iran
	case ProtocolVMess:
		score += 25
	case ProtocolTrojan:
		score += 20
	case ProtocolShadowsocks:
		score += 15
	}

	// Transport scoring
	switch c.Network {
	case TransportXHTTP:
		score += 25 // xhttp is best for DPI evasion
	case TransportGRPC:
		score += 20
	case TransportWebSocket:
		score += 18
	case TransportHTTPUpgrade:
		score += 15
	case TransportHTTP2:
		score += 12
	case TransportQUIC:
		score += 10
	case TransportTCP:
		score += 5
	}

	// Security scoring
	switch c.Security {
	case SecurityReality:
		score += 25 // Reality is best for Iran
	case SecurityXTLS:
		score += 20
	case SecurityTLS:
		score += 15
	case SecurityNone:
		score += 0
	}

	// TLS fingerprint
	if c.Fingerprint == "chrome" || c.Fingerprint == "firefox" || c.Fingerprint == "random" {
		score += 10
	}

	// ALPN support
	if len(c.ALPN) > 0 {
		score += 5
	}

	// Flow control (VLESS)
	if c.Flow != "" {
		score += 5
	}

	return score
}
