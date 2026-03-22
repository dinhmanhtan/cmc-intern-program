package scanner

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"mini-asm/internal/model"
)

// IPScanner performs IP geolocation and ASN lookup using free public APIs.
//
// SCAN CATEGORY: 🔵 PASSIVE - No authorization needed
//
// This scanner queries publicly available APIs (ip-api.com) to retrieve
// geolocation and network information for an IP address. No direct connection
// to the target system is made.
type IPScanner struct {
	httpClient *http.Client
}

// NewIPScanner creates a new IP scanner instance
func NewIPScanner() *IPScanner {
	return &IPScanner{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Type returns the scan type identifier
func (s *IPScanner) Type() model.ScanType {
	return model.ScanTypeIP
}

// ipAPIResponse is the response structure from ip-api.com
// Free API: http://ip-api.com/json/{ip}?fields=status,country,countryCode,city,regionName,lat,lon,isp,org,as,query,reverse
type ipAPIResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Region      string  `json:"regionName"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`   // e.g. "AS13335 Cloudflare, Inc."
	Query       string  `json:"query"` // The IP that was queried
	Reverse     string  `json:"reverse"` // Reverse DNS
}

// Scan performs IP geolocation and ASN lookup for the given asset
// Returns IPScanResult populated with geolocation, ASN, and reverse DNS data
func (s *IPScanner) Scan(asset *model.Asset) (*model.IPScanResult, error) {
	if asset.Type != model.TypeIP {
		return nil, fmt.Errorf("IP scan requires an IP asset, got: %s", asset.Type)
	}

	ipAddr := asset.Name

	// Validate it is a valid IP
	if net.ParseIP(ipAddr) == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipAddr)
	}

	// Query ip-api.com for geolocation + ASN info in one call
	// Fields: status, country, countryCode, city, regionName, lat, lon, isp, org, as, query, reverse
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,city,regionName,lat,lon,isp,org,as,query,reverse", ipAddr)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query ip-api.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ip-api.com returned status %d", resp.StatusCode)
	}

	var apiResp ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode ip-api response: %w", err)
	}

	if apiResp.Status != "success" {
		return nil, fmt.Errorf("ip-api.com lookup failed for IP: %s", ipAddr)
	}

	// Parse ASN number and name from "AS13335 Cloudflare, Inc."
	asnNumber, asnName := parseASNField(apiResp.AS)

	result := &model.IPScanResult{
		IPAddress: apiResp.Query,
		GeoLocation: model.IPGeolocation{
			Country:     apiResp.Country,
			CountryCode: apiResp.CountryCode,
			City:        apiResp.City,
			Region:      apiResp.Region,
			Latitude:    apiResp.Lat,
			Longitude:   apiResp.Lon,
			ISP:         apiResp.ISP,
			Org:         apiResp.Org,
		},
		ASN: model.IPASN{
			Number:      asnNumber,
			Name:        asnName,
			Description: apiResp.Org,
		},
		ReverseDNS: apiResp.Reverse,
		CreatedAt:  time.Now(),
	}

	return result, nil
}

// parseASNField parses "AS13335 Cloudflare, Inc." into number=13335 and name="Cloudflare, Inc."
func parseASNField(asField string) (int, string) {
	if asField == "" {
		return 0, ""
	}
	var num int
	var name string
	// Format: "AS<number> <name>"
	_, err := fmt.Sscanf(asField, "AS%d", &num)
	if err != nil {
		return 0, asField
	}
	// Find name after the number
	for i, c := range asField {
		if c == ' ' {
			name = asField[i+1:]
			break
		}
	}
	return num, name
}
