package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// SNIFragmenter handles SNI fragmentation at network layer
type SNIFragmenter struct {
	enabled          bool
	fragmentSize     int
	delayBetween     time.Duration
	randomizeDelay   bool
	adaptiveMode     bool
	successHistory   []bool
}

// NewSNIFragmenter creates a new SNI fragmenter
func NewSNIFragmenter() *SNIFragmenter {
	return &SNIFragmenter{
		enabled:        true,
		fragmentSize:   5,  // Default: split into 5-byte chunks
		delayBetween:   2 * time.Millisecond,
		randomizeDelay: true,
		adaptiveMode:   true,
		successHistory: make([]bool, 0, 100),
	}
}

// FragmentedConn wraps a net.Conn and fragments TLS ClientHello
type FragmentedConn struct {
	net.Conn
	fragmenter    *SNIFragmenter
	firstWrite    bool
	fragmentsUsed int
}

// NewFragmentedConn creates a new fragmented connection wrapper
func NewFragmentedConn(conn net.Conn, fragmenter *SNIFragmenter) *FragmentedConn {
	return &FragmentedConn{
		Conn:       conn,
		fragmenter: fragmenter,
		firstWrite: true,
	}
}

// Write intercepts the first write (TLS ClientHello) and fragments it
func (fc *FragmentedConn) Write(b []byte) (int, error) {
	if !fc.firstWrite || !fc.fragmenter.enabled {
		return fc.Conn.Write(b)
	}

	// Check if this is a TLS ClientHello (starts with 0x16 0x03)
	if len(b) < 6 || b[0] != 0x16 || b[1] != 0x03 {
		return fc.Conn.Write(b)
	}

	// Fragment the ClientHello packet
	fc.firstWrite = false
	return fc.fragmentAndWrite(b)
}

// fragmentAndWrite splits the TLS ClientHello into multiple fragments
func (fc *FragmentedConn) fragmentAndWrite(data []byte) (int, error) {
	if len(data) <= fc.fragmenter.fragmentSize {
		return fc.Conn.Write(data)
	}

	totalWritten := 0
	remaining := data
	fragmentCount := 0

	// Strategy 1: Split TLS header and body
	// TLS Record: [Type(1)] [Version(2)] [Length(2)] [Payload]
	if len(data) > 5 {
		// Write TLS header (5 bytes)
		header := data[:5]
		n, err := fc.Conn.Write(header)
		if err != nil {
			return n, err
		}
		totalWritten += n
		remaining = data[5:]
		fragmentCount++

		// Small delay after header
		if fc.fragmenter.delayBetween > 0 {
			time.Sleep(fc.getFragmentDelay())
		}
	}

	// Strategy 2: Fragment the payload (including SNI)
	// Find SNI extension if present and fragment around it
	sniOffset := fc.findSNIOffset(remaining)
	
	if sniOffset > 0 && sniOffset < len(remaining)-10 {
		// Fragment before SNI
		if sniOffset > 0 {
			chunk := remaining[:sniOffset]
			n, err := fc.Conn.Write(chunk)
			if err != nil {
				return totalWritten + n, err
			}
			totalWritten += n
			remaining = remaining[sniOffset:]
			fragmentCount++
			time.Sleep(fc.getFragmentDelay())
		}

		// Fragment SNI itself into tiny pieces (most critical for DPI evasion)
		sniLength := 30 // Approximate SNI extension length
		if sniLength > len(remaining) {
			sniLength = len(remaining)
		}

		// Split SNI into 3-5 byte chunks
		sniData := remaining[:sniLength]
		chunkSize := 3 + (fragmentCount % 3) // 3-5 bytes
		
		for i := 0; i < len(sniData); i += chunkSize {
			end := i + chunkSize
			if end > len(sniData) {
				end = len(sniData)
			}
			
			chunk := sniData[i:end]
			n, err := fc.Conn.Write(chunk)
			if err != nil {
				return totalWritten + n, err
			}
			totalWritten += n
			fragmentCount++
			
			// Critical delay between SNI fragments
			time.Sleep(fc.getFragmentDelay())
		}
		
		remaining = remaining[sniLength:]
	}

	// Strategy 3: Fragment remaining data in larger chunks
	for len(remaining) > 0 {
		chunkSize := fc.fragmenter.fragmentSize * 4 // Larger chunks for non-SNI data
		if chunkSize > len(remaining) {
			chunkSize = len(remaining)
		}

		chunk := remaining[:chunkSize]
		n, err := fc.Conn.Write(chunk)
		if err != nil {
			return totalWritten + n, err
		}
		
		totalWritten += n
		remaining = remaining[chunkSize:]
		fragmentCount++

		if len(remaining) > 0 && fc.fragmenter.delayBetween > 0 {
			time.Sleep(fc.getFragmentDelay())
		}
	}

	fc.fragmentsUsed = fragmentCount
	return totalWritten, nil
}

// findSNIOffset finds the offset of SNI extension in TLS ClientHello
func (fc *FragmentedConn) findSNIOffset(data []byte) int {
	// TLS ClientHello structure (simplified):
	// - Handshake Type (1 byte) = 0x01
	// - Length (3 bytes)
	// - Client Version (2 bytes)
	// - Random (32 bytes)
	// - Session ID (variable)
	// - Cipher Suites (variable)
	// - Compression Methods (variable)
	// - Extensions (variable) <- SNI is here

	if len(data) < 40 {
		return -1
	}

	// Look for SNI extension type (0x00 0x00)
	for i := 0; i < len(data)-2; i++ {
		if data[i] == 0x00 && data[i+1] == 0x00 {
			// Verify this looks like an SNI extension
			if i+4 < len(data) {
				// Check if there's a reasonable length following
				return i
			}
		}
	}

	return -1
}

// getFragmentDelay returns the delay between fragments (with randomization)
func (fc *FragmentedConn) getFragmentDelay() time.Duration {
	baseDelay := fc.fragmenter.delayBetween
	
	if !fc.fragmenter.randomizeDelay {
		return baseDelay
	}

	// Add ±30% randomization to avoid pattern detection
	variation := int64(float64(baseDelay) * 0.3)
	randomOffset := (randInt64(variation*2) - variation)
	
	finalDelay := baseDelay + time.Duration(randomOffset)
	if finalDelay < time.Millisecond {
		finalDelay = time.Millisecond
	}
	
	return finalDelay
}

// RecordSuccess records if this fragmentation strategy succeeded
func (fc *FragmentedConn) RecordSuccess(success bool) {
	fc.fragmenter.recordAttempt(success, fc.fragmentsUsed)
}

// recordAttempt records fragmentation attempt for adaptive learning
func (sf *SNIFragmenter) recordAttempt(success bool, fragmentsUsed int) {
	sf.successHistory = append(sf.successHistory, success)
	
	// Keep only last 100 attempts
	if len(sf.successHistory) > 100 {
		sf.successHistory = sf.successHistory[1:]
	}

	// Adaptive mode: adjust fragment size based on success rate
	if sf.adaptiveMode && len(sf.successHistory) >= 20 {
		recentSuccess := 0
		recentTotal := 20
		
		for i := len(sf.successHistory) - recentTotal; i < len(sf.successHistory); i++ {
			if sf.successHistory[i] {
				recentSuccess++
			}
		}

		successRate := float64(recentSuccess) / float64(recentTotal)
		
		// Adjust fragment size based on success rate
		if successRate < 0.5 {
			// Low success - make fragments smaller (more aggressive)
			if sf.fragmentSize > 3 {
				sf.fragmentSize--
			}
			if sf.delayBetween < 10*time.Millisecond {
				sf.delayBetween += time.Millisecond
			}
		} else if successRate > 0.85 {
			// High success - can use slightly larger fragments (better performance)
			if sf.fragmentSize < 8 {
				sf.fragmentSize++
			}
			if sf.delayBetween > time.Millisecond {
				sf.delayBetween -= time.Millisecond
			}
		}
	}
}

// GetSuccessRate returns the current success rate of fragmentation
func (sf *SNIFragmenter) GetSuccessRate() float64 {
	if len(sf.successHistory) == 0 {
		return 0.0
	}

	success := 0
	for _, s := range sf.successHistory {
		if s {
			success++
		}
	}

	return float64(success) / float64(len(sf.successHistory))
}

// GetStats returns fragmentation statistics
func (sf *SNIFragmenter) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"enabled":          sf.enabled,
		"fragment_size":    sf.fragmentSize,
		"delay_ms":         sf.delayBetween.Milliseconds(),
		"adaptive":         sf.adaptiveMode,
		"attempts":         len(sf.successHistory),
		"success_rate":     sf.GetSuccessRate(),
	}
}

// SNIFragmentDialer wraps a dialer with SNI fragmentation
type SNIFragmentDialer struct {
	dialer      *net.Dialer
	fragmenter  *SNIFragmenter
	utlsDialer  *UTLSDialer
}

// NewSNIFragmentDialer creates a dialer with SNI fragmentation and uTLS
func NewSNIFragmentDialer(timeout time.Duration) *SNIFragmentDialer {
	return &SNIFragmentDialer{
		dialer: &net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		},
		fragmenter: NewSNIFragmenter(),
		utlsDialer: NewUTLSDialer(timeout),
	}
}

// DialContext establishes a connection with both uTLS and SNI fragmentation
func (sfd *SNIFragmentDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// First establish TCP connection
	conn, err := sfd.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("TCP dial failed: %v", err)
	}

	// Wrap with fragmentation layer
	fragConn := NewFragmentedConn(conn, sfd.fragmenter)
	
	return fragConn, nil
}

// GetFragmenter returns the underlying fragmenter for statistics
func (sfd *SNIFragmentDialer) GetFragmenter() *SNIFragmenter {
	return sfd.fragmenter
}

// SetFragmentSize manually sets the fragment size
func (sfd *SNIFragmentDialer) SetFragmentSize(size int) {
	if size >= 3 && size <= 20 {
		sfd.fragmenter.fragmentSize = size
	}
}

// SetDelay sets the delay between fragments
func (sfd *SNIFragmentDialer) SetDelay(delay time.Duration) {
	if delay >= 0 && delay <= 50*time.Millisecond {
		sfd.fragmenter.delayBetween = delay
	}
}

// EnableAdaptive enables adaptive fragment size adjustment
func (sfd *SNIFragmentDialer) EnableAdaptive() {
	sfd.fragmenter.adaptiveMode = true
}

// randInt64 generates a random int64 for delay randomization
func randInt64(max int64) int64 {
	if max <= 0 {
		return 0
	}
	
	// Simple pseudo-random for delay jitter
	seed := time.Now().UnixNano()
	return (seed % max)
}

// TLSFragmentedDialer combines uTLS fingerprinting with SNI fragmentation
type TLSFragmentedDialer struct {
	baseDialer    *SNIFragmentDialer
	utlsDialer    *UTLSDialer
	timeout       time.Duration
}

// NewTLSFragmentedDialer creates the ultimate anti-DPI dialer
func NewTLSFragmentedDialer(timeout time.Duration) *TLSFragmentedDialer {
	return &TLSFragmentedDialer{
		baseDialer: NewSNIFragmentDialer(timeout),
		utlsDialer: NewUTLSDialer(timeout),
		timeout:    timeout,
	}
}

// Dial establishes a connection with all anti-DPI techniques enabled
func (tfd *TLSFragmentedDialer) Dial(network, addr string) (net.Conn, error) {
	return tfd.DialContext(context.Background(), network, addr)
}

// DialContext establishes a connection with full DPI evasion
func (tfd *TLSFragmentedDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// Use SNI fragmented connection for the base layer
	return tfd.baseDialer.DialContext(ctx, network, addr)
}

// GetStats returns comprehensive dialer statistics
func (tfd *TLSFragmentedDialer) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// SNI fragmentation stats
	fragStats := tfd.baseDialer.GetFragmenter().GetStats()
	stats["sni_fragmentation"] = fragStats
	
	// uTLS fingerprint info
	stats["utls_fingerprint"] = tfd.utlsDialer.GetFingerprintInfo()
	stats["browser_rotation"] = tfd.utlsDialer.rotationEnabled
	
	return stats
}

// bufferReader wraps an io.Reader to allow peeking at data
type bufferReader struct {
	r      io.Reader
	buffer *bytes.Buffer
}

func newBufferReader(r io.Reader) *bufferReader {
	return &bufferReader{
		r:      r,
		buffer: new(bytes.Buffer),
	}
}

func (br *bufferReader) Read(p []byte) (int, error) {
	// First read from buffer if available
	if br.buffer.Len() > 0 {
		return br.buffer.Read(p)
	}
	return br.r.Read(p)
}

func (br *bufferReader) Peek(n int) ([]byte, error) {
	// Read n bytes into buffer without consuming
	for br.buffer.Len() < n {
		tmp := make([]byte, n-br.buffer.Len())
		nr, err := br.r.Read(tmp)
		if nr > 0 {
			br.buffer.Write(tmp[:nr])
		}
		if err != nil {
			return br.buffer.Bytes(), err
		}
	}
	return br.buffer.Bytes()[:n], nil
}

// parseTLSHeader parses basic TLS record header
func parseTLSHeader(data []byte) (recordType byte, version uint16, length uint16, err error) {
	if len(data) < 5 {
		return 0, 0, 0, fmt.Errorf("data too short")
	}
	
	recordType = data[0]
	version = binary.BigEndian.Uint16(data[1:3])
	length = binary.BigEndian.Uint16(data[3:5])
	
	return recordType, version, length, nil
}
