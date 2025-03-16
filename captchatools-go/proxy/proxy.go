package proxy

import (
	"errors"
	"fmt"
	"strings"
)

type Proxy struct {
	Ip       string
	Port     string
	User     string
	Password string
}

// New returns a proxy that can be used to solve captchas
func New(proxy string) (*Proxy, error) {
	if strings.TrimSpace(proxy) == "" {
		return nil, ErrProxyEmpty
	}

	splitted := strings.Split(proxy, ":")
	pLen := len(splitted)
	p := &Proxy{}

	// Handle the case where only IP is provided
	if pLen == 1 {
		p.Ip = splitted[0]
		return p, nil
	}

	// Set IP and port for all cases with at least 2 components
	p.Ip = splitted[0]
	p.Port = splitted[1]

	// Handle authentication (username and password)
	if pLen >= 3 {
		p.User = splitted[2]
	}

	if pLen >= 4 {
		p.Password = splitted[3]
	}

	// We ignore any additional fields beyond the first four

	return p, nil
}

// FromProxyUrl converts a given proxy url string (format: [http(s)://][user:pass@]ip:port) into type Proxy
func FromProxyUrl(proxyUrl string) (*Proxy, error) {
	if strings.TrimSpace(proxyUrl) == "" {
		return nil, errors.New("proxy url is empty") // TODO: Return custom error
	}

	// Remove http:// or https:// prefix if present
	proxyUrl = strings.TrimPrefix(proxyUrl, "http://")
	proxyUrl = strings.TrimPrefix(proxyUrl, "https://")

	// Create a new Proxy instance
	p := &Proxy{}

	// Look for the last @ symbol which separates auth from host:port
	lastAtIndex := strings.LastIndex(proxyUrl, "@")

	if lastAtIndex != -1 {
		// URL contains authentication
		auth := proxyUrl[:lastAtIndex]
		hostPort := proxyUrl[lastAtIndex+1:]

		// Find the first colon in auth part to separate username and password
		// Note that username might contain @ symbols
		colonIndex := strings.Index(auth, ":")
		if colonIndex != -1 {
			p.User = auth[:colonIndex]
			p.Password = auth[colonIndex+1:]
		} else {
			p.User = auth
		}

		// Process the host:port part
		parseHostPort(p, hostPort)
	} else {
		// No authentication
		parseHostPort(p, proxyUrl)
	}

	return p, nil
}

// Helper function to parse the host:port part of a proxy URL
func parseHostPort(p *Proxy, hostPort string) {
	// Handle IPv6 addresses which might contain colons
	if strings.Contains(hostPort, "[") && strings.Contains(hostPort, "]") {
		// IPv6 address format: [IPv6]:port
		ipv6EndBracket := strings.LastIndex(hostPort, "]")

		// Extract the IPv6 address including brackets
		p.Ip = hostPort[:ipv6EndBracket+1]

		// Check if there's a port after the IPv6 address
		if ipv6EndBracket+2 < len(hostPort) && hostPort[ipv6EndBracket+1] == ':' {
			p.Port = hostPort[ipv6EndBracket+2:]
		}
	} else {
		// Regular IPv4 address or hostname
		hostPortParts := strings.Split(hostPort, ":")
		if len(hostPortParts) >= 2 {
			p.Ip = hostPortParts[0]
			p.Port = hostPortParts[1]
		} else {
			// Only IP, no port
			p.Ip = hostPort
		}
	}
}

// IsUserAuth returns if a proxy is user authenticated
func (p Proxy) IsUserAuth() bool {
	return p.User != "" && p.Password != ""
}

// Returns the proxy as a string, unformatted
//
// Example: would return "ip:port" || "ip:port:user:pass"
func (p Proxy) String() string {
	var formatted string = p.Ip + ":" + p.Port
	if p.IsUserAuth() {
		formatted = formatted + ":" + p.User + ":" + p.Password
	}
	return formatted
}

// ProxyUrl returns the proxy as a string in the correct format
//
// Example: returns "user:pass@ip:port" || "ip:port"
func (p Proxy) ProxyUrl() string {
	var formatted string = p.Ip + ":" + p.Port
	if p.IsUserAuth() {
		formatted = fmt.Sprintf("%v:%v@%v", p.User, p.Password, formatted)
	}
	return formatted
}
