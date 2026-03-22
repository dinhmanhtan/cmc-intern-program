package scanner

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"mini-asm/internal/model"
)

/*
⚠️⚠️⚠️ WARNING: ACTIVE SCANNING - READ THIS CAREFULLY ⚠️⚠️⚠️

SCAN CATEGORY: 🔴 ACTIVE / INTRUSIVE

This scanner directly probes target systems by attempting TCP connections.
SAFETY CHECK: Only private/localhost IP ranges are permitted in this training exercise.

Authorized targets:
  - 127.0.0.1 / localhost (loopback)
  - 10.x.x.x (private class A)
  - 172.16-31.x.x (private class B)
  - 192.168.x.x (private class C)
*/

// PortScannerImpl performs port scanning on IP addresses.
// Renamed from PortScanner to PortScannerImpl to avoid conflict with template file.
//
// SCAN CATEGORY: 🔴 ACTIVE - REQUIRES AUTHORIZATION
type PortScannerImpl struct {
	timeout    time.Duration
	maxWorkers int
	commonPorts []int
}

// NewPortScannerImpl creates a new port scanner with safety restrictions
func NewPortScannerImpl() *PortScannerImpl {
	log.Println("⚠️ PORT SCANNING IS ACTIVE RECONNAISSANCE")
	log.Println("⚠️ Only scan systems you own or have explicit permission to scan")
	log.Println("⚠️ For training: ONLY scan private IP ranges (localhost / 192.168.x.x / 10.x.x.x)")

	return &PortScannerImpl{
		timeout:    500 * time.Millisecond, // Short timeout for fast scanning
		maxWorkers: 50,                     // Concurrent goroutines
		commonPorts: []int{
			21,   // FTP
			22,   // SSH
			23,   // Telnet
			25,   // SMTP
			53,   // DNS
			80,   // HTTP
			110,  // POP3
			143,  // IMAP
			443,  // HTTPS
			445,  // SMB
			3306, // MySQL
			3389, // RDP
			5432, // PostgreSQL
			5900, // VNC
			8080, // HTTP Alt
			8443, // HTTPS Alt
		},
	}
}

// Type returns the scan type identifier
func (s *PortScannerImpl) Type() model.ScanType {
	return model.ScanTypePort
}

// Scan performs TCP port scanning on a target IP address.
// CRITICAL: Only allows private/localhost IP ranges to prevent unauthorized scanning.
func (s *PortScannerImpl) Scan(asset *model.Asset) (*model.PortScanResult, error) {
	if asset.Type != model.TypeIP && asset.Type != model.TypeDomain {
		return nil, fmt.Errorf("port scan requires IP or domain asset, got: %s", asset.Type)
	}

	target := asset.Name

	// CRITICAL SAFETY CHECK: Only allow private / loopback IPs
	if !s.isPrivateOrLocalhost(target) {
		return nil, fmt.Errorf(
			"⚠️ UNAUTHORIZED PORT SCAN BLOCKED ⚠️\n"+
				"Target: %s\n"+
				"For this training exercise, port scanning is ONLY permitted on:\n"+
				"  - 127.0.0.1 (localhost)\n"+
				"  - 10.x.x.x (private class A)\n"+
				"  - 172.16-31.x.x (private class B)\n"+
				"  - 192.168.x.x (private class C)\n"+
				"Unauthorized port scanning may violate local cybercrime laws.",
			target,
		)
	}

	log.Printf("🔴 ACTIVE PORT SCAN: target=%s, ports=%d", target, len(s.commonPorts))

	startTime := time.Now()

	// Use worker pool pattern to scan ports concurrently
	type portResult struct {
		port    int
		isOpen  bool
		service string
		version string
	}

	portChan := make(chan int, len(s.commonPorts))
	resultChan := make(chan portResult, len(s.commonPorts))

	// Fill port queue
	for _, port := range s.commonPorts {
		portChan <- port
	}
	close(portChan)

	// Launch workers
	var wg sync.WaitGroup
	for i := 0; i < s.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portChan {
				address := net.JoinHostPort(target, fmt.Sprintf("%d", port))
				conn, err := net.DialTimeout("tcp", address, s.timeout)
				if err != nil {
					// Port closed or filtered
					resultChan <- portResult{port: port, isOpen: false}
					continue
				}
				conn.Close()
				// Port open - detect service by well-known mapping
				svc := s.detectServiceByPort(port)
				resultChan <- portResult{port: port, isOpen: true, service: svc}
			}
		}()
	}

	// Wait for all workers to finish then close result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var openPorts []model.OpenPort
	closedCount := 0

	for res := range resultChan {
		if res.isOpen {
			openPorts = append(openPorts, model.OpenPort{
				Port:     res.port,
				Protocol: "tcp",
				State:    "open",
				Service:  res.service,
				Version:  res.version,
			})
		} else {
			closedCount++
		}
	}

	scanDuration := time.Since(startTime).Milliseconds()

	result := &model.PortScanResult{
		IPAddress:      target,
		OpenPorts:      openPorts,
		ClosedPorts:    closedCount,
		TotalScanned:   len(s.commonPorts),
		ScanDurationMs: scanDuration,
		CreatedAt:      time.Now(),
	}

	log.Printf("✅ Port scan complete: %d open, %d closed, duration=%dms",
		len(openPorts), closedCount, scanDuration)

	return result, nil
}

// isPrivateOrLocalhost checks if IP is in private or loopback range.
// This is the SAFETY CHECK that prevents unauthorized scanning of public systems.
func (s *PortScannerImpl) isPrivateOrLocalhost(target string) bool {
	// Resolve hostname to IP if needed
	ip := net.ParseIP(target)
	if ip == nil {
		// Try resolving as hostname
		ips, err := net.LookupIP(target)
		if err != nil || len(ips) == 0 {
			return false
		}
		ip = ips[0]
	}

	// Check loopback
	if ip.IsLoopback() {
		return true
	}

	// Private ranges (RFC 1918 + RFC 4193)
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7", // IPv6 private
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// detectServiceByPort returns a well-known service name for common ports
func (s *PortScannerImpl) detectServiceByPort(port int) string {
	services := map[int]string{
		21:   "ftp",
		22:   "ssh",
		23:   "telnet",
		25:   "smtp",
		53:   "dns",
		80:   "http",
		110:  "pop3",
		143:  "imap",
		443:  "https",
		445:  "smb",
		3306: "mysql",
		3389: "rdp",
		5432: "postgresql",
		5900: "vnc",
		8080: "http-alt",
		8443: "https-alt",
	}
	if svc, ok := services[port]; ok {
		return svc
	}
	return strings.ToLower(fmt.Sprintf("port-%d", port))
}
