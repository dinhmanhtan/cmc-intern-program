package scanner

import (
	"net"
	"testing"

	"mini-asm/internal/model"
)

// TestIPScanner_ValidIP tests scanning a valid public IP address
// Note: This test makes a real HTTP call to ip-api.com
// In production, use a mock HTTP client
func TestIPScanner_ValidIP(t *testing.T) {
	scanner := NewIPScanner()

	asset := &model.Asset{
		ID:   "test-id",
		Name: "8.8.8.8",
		Type: model.TypeIP,
	}

	result, err := scanner.Scan(asset)
	if err != nil {
		t.Skipf("Skipping: IP API call failed (network not available): %v", err)
		return
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify required fields are populated
	if result.IPAddress == "" {
		t.Error("Expected IPAddress to be populated")
	}
	if result.GeoLocation.Country == "" {
		t.Error("Expected Country to be populated")
	}
}

// TestIPScanner_WrongAssetType tests that IP scanner rejects non-IP assets
func TestIPScanner_WrongAssetType(t *testing.T) {
	scanner := NewIPScanner()

	asset := &model.Asset{
		ID:   "test-id",
		Name: "example.com",
		Type: model.TypeDomain, // Wrong type!
	}

	_, err := scanner.Scan(asset)
	if err == nil {
		t.Error("Expected error for domain asset, got nil")
	}
}

// TestIPScanner_InvalidIPAddress tests that scanner rejects invalid IP strings
func TestIPScanner_InvalidIPAddress(t *testing.T) {
	scanner := NewIPScanner()

	asset := &model.Asset{
		ID:   "test-id",
		Name: "not-an-ip-address!!",
		Type: model.TypeIP,
	}

	_, err := scanner.Scan(asset)
	if err == nil {
		t.Error("Expected error for invalid IP address, got nil")
	}
}

// TestIPScanner_LocalhostIP tests scanning localhost IP
func TestIPScanner_LocalhostIP(t *testing.T) {
	scanner := NewIPScanner()

	asset := &model.Asset{
		ID:   "test-id",
		Name: "127.0.0.1",
		Type: model.TypeIP,
	}

	// 127.0.0.1 is a valid IP, but ip-api.com may return "fail" for private IPs
	// so we just verify no panic and handle both success/fail gracefully
	_, err := scanner.Scan(asset)
	// Either succeeds or fails - both are acceptable for this test
	// The important thing is no panic occurs
	t.Logf("Localhost scan result: err=%v", err)
}

// TestParseASNField tests the ASN parsing utility function
func TestParseASNField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantNum  int
		wantName string
	}{
		{
			name:     "valid cloudflare ASN",
			input:    "AS13335 Cloudflare, Inc.",
			wantNum:  13335,
			wantName: "Cloudflare, Inc.",
		},
		{
			name:     "valid google ASN",
			input:    "AS15169 Google LLC",
			wantNum:  15169,
			wantName: "Google LLC",
		},
		{
			name:     "empty string",
			input:    "",
			wantNum:  0,
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotName := parseASNField(tt.input)
			if gotNum != tt.wantNum {
				t.Errorf("parseASNField(%q) number = %d, want %d", tt.input, gotNum, tt.wantNum)
			}
			if tt.wantNum > 0 && gotName != tt.wantName {
				t.Errorf("parseASNField(%q) name = %q, want %q", tt.input, gotName, tt.wantName)
			}
		})
	}
}

// TestIPScanner_Type verifies the scanner returns correct type identifier
func TestIPScanner_Type(t *testing.T) {
	scanner := NewIPScanner()
	got := scanner.Type()
	if got != model.ScanTypeIP {
		t.Errorf("IPScanner.Type() = %q, want %q", got, model.ScanTypeIP)
	}
}

// TestIPScanner_ParseIP verifies that only valid IPs are accepted
func TestIPScanner_ParseIP(t *testing.T) {
	tests := []struct {
		ip    string
		valid bool
	}{
		{"8.8.8.8", true},
		{"192.168.1.1", true},
		{"127.0.0.1", true},
		{"2001:4860:4860::8888", true}, // IPv6
		{"not-an-ip", false},
		{"256.0.0.1", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			parsed := net.ParseIP(tt.ip)
			isValid := parsed != nil
			if isValid != tt.valid {
				t.Errorf("net.ParseIP(%q) valid = %v, want %v", tt.ip, isValid, tt.valid)
			}
		})
	}
}
