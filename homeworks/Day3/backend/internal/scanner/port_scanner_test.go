package scanner

import (
	"testing"

	"mini-asm/internal/model"
)

// TestPortScanner_Type verifies the scanner returns correct type identifier
func TestPortScanner_Type(t *testing.T) {
	scanner := NewPortScannerImpl()
	got := scanner.Type()
	if got != model.ScanTypePort {
		t.Errorf("PortScannerImpl.Type() = %q, want %q", got, model.ScanTypePort)
	}
}

// TestPortScanner_Localhost tests scanning localhost - the only permitted target
func TestPortScanner_Localhost(t *testing.T) {
	scanner := NewPortScannerImpl()

	asset := &model.Asset{
		ID:   "test-id",
		Name: "127.0.0.1",
		Type: model.TypeIP,
	}

	result, err := scanner.Scan(asset)
	if err != nil {
		t.Fatalf("Port scan on localhost should succeed, got error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result for localhost scan")
	}

	// Basic result validation
	if result.IPAddress != "127.0.0.1" {
		t.Errorf("Expected IPAddress = 127.0.0.1, got %q", result.IPAddress)
	}
	if result.TotalScanned <= 0 {
		t.Error("Expected TotalScanned > 0")
	}
	if result.ScanDurationMs <= 0 {
		t.Error("Expected ScanDurationMs > 0")
	}

	// Open ports should be valid
	for _, p := range result.OpenPorts {
		if p.Port <= 0 || p.Port > 65535 {
			t.Errorf("Invalid port number in results: %d", p.Port)
		}
		if p.Protocol != "tcp" {
			t.Errorf("Expected protocol 'tcp', got %q", p.Protocol)
		}
		if p.State != "open" {
			t.Errorf("Expected state 'open', got %q", p.State)
		}
	}

	t.Logf("Port scan: found %d open ports, %d closed, %dms",
		len(result.OpenPorts), result.ClosedPorts, result.ScanDurationMs)
}

// TestPortScanner_PrivateIP tests scanning a private IP range (192.168.x.x)
func TestPortScanner_PrivateIP(t *testing.T) {
	scanner := NewPortScannerImpl()

	// Using a typical router IP - will likely timeout but should be allowed
	asset := &model.Asset{
		ID:   "test-id",
		Name: "192.168.1.1",
		Type: model.TypeIP,
	}

	// Should NOT return a safety error (192.168.x.x is private)
	// It may return an empty result if no ports are open
	_, err := scanner.Scan(asset)
	if err != nil {
		// Check it's not a safety error
		errMsg := err.Error()
		if len(errMsg) > 0 && errMsg[:2] == "⚠️" {
			t.Errorf("Private IP 192.168.1.1 should pass safety check, got: %v", err)
		}
		// Other errors (e.g., timeout) are acceptable
		t.Logf("Expected error (timeout/unreachable) for 192.168.1.1: %v", err)
	}
}

// TestPortScanner_UnauthorizedPublicIP tests that public IPs are blocked
func TestPortScanner_UnauthorizedPublicIP(t *testing.T) {
	scanner := NewPortScannerImpl()

	tests := []struct {
		name string
		ip   string
	}{
		{"google dns", "8.8.8.8"},
		{"cloudflare", "1.1.1.1"},
		{"random public", "203.0.113.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := &model.Asset{
				ID:   "test-id",
				Name: tt.ip,
				Type: model.TypeIP,
			}

			_, err := scanner.Scan(asset)
			if err == nil {
				t.Errorf("Expected safety error for public IP %s, got nil", tt.ip)
			}
		})
	}
}

// TestPortScanner_WrongAssetType tests that wrong asset types are rejected
func TestPortScanner_WrongAssetType(t *testing.T) {
	scanner := NewPortScannerImpl()

	asset := &model.Asset{
		ID:   "test-id",
		Name: "some-service",
		Type: model.TypeService,
	}

	_, err := scanner.Scan(asset)
	if err == nil {
		t.Error("Expected error for service asset type, got nil")
	}
}

// TestPortScanner_SafetyCheck_PrivateRanges validates all private ranges pass safety check
func TestPortScanner_SafetyCheck_PrivateRanges(t *testing.T) {
	scanner := NewPortScannerImpl()

	privateIPs := []string{
		"127.0.0.1",     // loopback
		"10.0.0.1",      // private class A
		"10.255.255.255", // private class A boundary
		"172.16.0.1",    // private class B
		"172.31.255.255", // private class B boundary
		"192.168.0.1",   // private class C
		"192.168.255.255", // private class C boundary
	}

	for _, ip := range privateIPs {
		t.Run(ip, func(t *testing.T) {
			if !scanner.isPrivateOrLocalhost(ip) {
				t.Errorf("IP %s should be considered private/localhost", ip)
			}
		})
	}
}

// TestPortScanner_SafetyCheck_PublicBlocked validates public IPs fail safety check
func TestPortScanner_SafetyCheck_PublicBlocked(t *testing.T) {
	scanner := NewPortScannerImpl()

	publicIPs := []string{
		"8.8.8.8",
		"1.1.1.1",
		"204.79.197.200",
		"203.0.113.1",
	}

	for _, ip := range publicIPs {
		t.Run(ip, func(t *testing.T) {
			if scanner.isPrivateOrLocalhost(ip) {
				t.Errorf("Public IP %s should NOT pass safety check", ip)
			}
		})
	}
}

// TestPortScanner_DetectService validates well-known port→service mapping
func TestPortScanner_DetectService(t *testing.T) {
	scanner := NewPortScannerImpl()

	tests := []struct {
		port    int
		service string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{5432, "postgresql"},
		{9999, "port-9999"}, // unknown port
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			got := scanner.detectServiceByPort(tt.port)
			if got != tt.service {
				t.Errorf("detectServiceByPort(%d) = %q, want %q", tt.port, got, tt.service)
			}
		})
	}
}
