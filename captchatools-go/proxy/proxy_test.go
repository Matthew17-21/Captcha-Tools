package proxy

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		proxyStr    string
		expected    *Proxy
		expectError bool
	}{
		{
			name:     "Valid IP:Port",
			proxyStr: "192.168.1.1:8080",
			expected: &Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "Valid IP:Port:User:Pass",
			proxyStr: "192.168.1.1:8080:username:password123",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name:        "Empty string",
			proxyStr:    "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Whitespace string",
			proxyStr:    "   ",
			expected:    nil,
			expectError: true,
		},
		{
			name:     "Only IP (missing port)",
			proxyStr: "192.168.1.1",
			expected: &Proxy{
				Ip:   "192.168.1.1",
				Port: "",
			},
			expectError: false,
		},
		{
			name:     "Partial auth (user only)",
			proxyStr: "192.168.1.1:8080:username",
			expected: &Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
				User: "username",
			},
			expectError: false,
		},
		{
			name:     "Extra fields",
			proxyStr: "192.168.1.1:8080:username:password:extra:fields",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "password",
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proxy, err := New(tc.proxyStr)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("didn't expect error but got: %v", err)
			}

			// If we expect error, no need to check the proxy
			if tc.expectError {
				return
			}

			// Check proxy values
			if proxy.Ip != tc.expected.Ip {
				t.Errorf("IP doesn't match: got %q, want %s", proxy.Ip, tc.expected.Ip)
			}
			if proxy.Port != tc.expected.Port {
				t.Errorf("Port doesn't match: got %q, want %s", proxy.Port, tc.expected.Port)
			}
			if proxy.User != tc.expected.User {
				t.Errorf("User doesn't match: got %q, want %s", proxy.User, tc.expected.User)
			}
			if proxy.Password != tc.expected.Password {
				t.Errorf("Password doesn't match: got %q, want %s", proxy.Password, tc.expected.Password)
			}
		})
	}
}

func TestIsUserAuth(t *testing.T) {
	tests := []struct {
		name     string
		proxy    Proxy
		expected bool
	}{
		{
			name: "With username and password",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "password",
			},
			expected: true,
		},
		{
			name: "Without auth",
			proxy: Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expected: false,
		},
		{
			name: "With username but no password",
			proxy: Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
				User: "username",
			},
			expected: false,
		},
		{
			name: "With password but no username",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				Password: "password",
			},
			expected: false,
		},
		{
			name: "With empty username and password",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "",
				Password: "",
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.proxy.IsUserAuth()
			if result != tc.expected {
				t.Errorf("IsUserAuth() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		proxy    Proxy
		expected string
	}{
		{
			name: "With username and password",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "password",
			},
			expected: "192.168.1.1:8080:username:password",
		},
		{
			name: "Without auth",
			proxy: Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expected: "192.168.1.1:8080",
		},
		{
			name: "IPv6 address",
			proxy: Proxy{
				Ip:   "[2001:db8::1]",
				Port: "8080",
			},
			expected: "[2001:db8::1]:8080",
		},
		{
			name: "With special characters in password",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "user@domain",
				Password: "p@$$w0rd!",
			},
			expected: "192.168.1.1:8080:user@domain:p@$$w0rd!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.proxy.String()
			if result != tc.expected {
				t.Errorf("String() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestProxyUrl(t *testing.T) {
	tests := []struct {
		name     string
		proxy    Proxy
		expected string
	}{
		{
			name: "With username and password",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "password",
			},
			expected: "username:password@192.168.1.1:8080",
		},
		{
			name: "Without auth",
			proxy: Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expected: "192.168.1.1:8080",
		},
		{
			name: "IPv6 address",
			proxy: Proxy{
				Ip:   "[2001:db8::1]",
				Port: "8080",
			},
			expected: "[2001:db8::1]:8080",
		},
		{
			name: "With special characters in auth",
			proxy: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "user@domain",
				Password: "p@$$w0rd!",
			},
			expected: "user@domain:p@$$w0rd!@192.168.1.1:8080",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.proxy.ProxyUrl()
			if result != tc.expected {
				t.Errorf("ProxyUrl() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestFromProxyUrl(t *testing.T) {
	tests := []struct {
		name        string
		proxyUrl    string
		expected    *Proxy
		expectError bool
	}{
		{
			name:     "Standard IP:Port format",
			proxyUrl: "192.168.1.1:8080",
			expected: &Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "With username and password",
			proxyUrl: "username:password@192.168.1.1:8080",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "password",
			},
			expectError: false,
		},
		{
			name:     "Only IP without port",
			proxyUrl: "192.168.1.1",
			expected: &Proxy{
				Ip: "192.168.1.1",
			},
			expectError: false,
		},
		{
			name:     "With username only",
			proxyUrl: "username@192.168.1.1:8080",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "username",
				Password: "",
			},
			expectError: false,
		},
		{
			name:     "IPv6 address with port",
			proxyUrl: "[2001:db8::1]:8080",
			expected: &Proxy{
				Ip:   "[2001:db8::1]",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "IPv6 address without port",
			proxyUrl: "[2001:db8::1]",
			expected: &Proxy{
				Ip: "[2001:db8::1]",
			},
			expectError: false,
		},
		{
			name:     "IPv6 with auth",
			proxyUrl: "user:pass@[2001:db8::1]:8080",
			expected: &Proxy{
				Ip:       "[2001:db8::1]",
				Port:     "8080",
				User:     "user",
				Password: "pass",
			},
			expectError: false,
		},
		{
			name:     "With HTTP prefix",
			proxyUrl: "http://192.168.1.1:8080",
			expected: &Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "With HTTPS prefix",
			proxyUrl: "https://192.168.1.1:8080",
			expected: &Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "With HTTP prefix and auth",
			proxyUrl: "http://user:pass@192.168.1.1:8080",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "user",
				Password: "pass",
			},
			expectError: false,
		},
		{
			name:     "With HTTPS prefix and IPv6",
			proxyUrl: "https://[2001:db8::1]:8080",
			expected: &Proxy{
				Ip:   "[2001:db8::1]",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "User with @ in username",
			proxyUrl: "user@domain:pass@192.168.1.1:8080",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "user@domain",
				Password: "pass",
			},
			expectError: false,
		},
		{
			name:     "Password with special characters",
			proxyUrl: "user:p@$$w0rd!@192.168.1.1:8080",
			expected: &Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "user",
				Password: "p@$$w0rd!",
			},
			expectError: false,
		},
		{
			name:     "Domain name instead of IP",
			proxyUrl: "proxy.example.com:8080",
			expected: &Proxy{
				Ip:   "proxy.example.com",
				Port: "8080",
			},
			expectError: false,
		},
		{
			name:     "Domain with subdomain and auth",
			proxyUrl: "user:pass@sub.example.com:8080",
			expected: &Proxy{
				Ip:       "sub.example.com",
				Port:     "8080",
				User:     "user",
				Password: "pass",
			},
			expectError: false,
		},
		{
			name:        "Empty string",
			proxyUrl:    "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Whitespace only",
			proxyUrl:    "   ",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proxy, err := FromProxyUrl(tc.proxyUrl)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("didn't expect error but got: %v", err)
			}

			// Skip rest of checks if we expected an error
			if tc.expectError {
				return
			}

			// Check all proxy fields
			if proxy.Ip != tc.expected.Ip {
				t.Errorf("IP doesn't match: got %s, want %s", proxy.Ip, tc.expected.Ip)
			}

			if proxy.Port != tc.expected.Port {
				t.Errorf("Port doesn't match: got %s, want %s", proxy.Port, tc.expected.Port)
			}

			if proxy.User != tc.expected.User {
				t.Errorf("User doesn't match: got %s, want %s", proxy.User, tc.expected.User)
			}

			if proxy.Password != tc.expected.Password {
				t.Errorf("Password doesn't match: got %s, want %s", proxy.Password, tc.expected.Password)
			}
		})
	}
}

func TestParseHostPort(t *testing.T) {
	tests := []struct {
		name     string
		hostPort string
		expected Proxy
	}{
		{
			name:     "Standard IPv4 with port",
			hostPort: "192.168.1.1:8080",
			expected: Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
		},
		{
			name:     "IPv4 without port",
			hostPort: "192.168.1.1",
			expected: Proxy{
				Ip: "192.168.1.1",
			},
		},
		{
			name:     "Hostname with port",
			hostPort: "example.com:8080",
			expected: Proxy{
				Ip:   "example.com",
				Port: "8080",
			},
		},
		{
			name:     "Hostname without port",
			hostPort: "example.com",
			expected: Proxy{
				Ip: "example.com",
			},
		},
		{
			name:     "IPv6 with port",
			hostPort: "[2001:db8::1]:8080",
			expected: Proxy{
				Ip:   "[2001:db8::1]",
				Port: "8080",
			},
		},
		{
			name:     "IPv6 without port",
			hostPort: "[2001:db8::1]",
			expected: Proxy{
				Ip: "[2001:db8::1]",
			},
		},
		{
			name:     "IPv6 with multiple colons",
			hostPort: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080",
			expected: Proxy{
				Ip:   "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
				Port: "8080",
			},
		},
		{
			name:     "IP with multiple ports (should take first port)",
			hostPort: "192.168.1.1:8080:8081",
			expected: Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
		},
		{
			name:     "Empty string",
			hostPort: "",
			expected: Proxy{
				Ip: "",
			},
		},
		{
			name:     "Port only",
			hostPort: ":8080",
			expected: Proxy{
				Ip:   "",
				Port: "8080",
			},
		},
		{
			name:     "IPv6 localhost with port",
			hostPort: "[::1]:8080",
			expected: Proxy{
				Ip:   "[::1]",
				Port: "8080",
			},
		},
		{
			name:     "IPv6 with zone index and port",
			hostPort: "[fe80::1%eth0]:8080",
			expected: Proxy{
				Ip:   "[fe80::1%eth0]",
				Port: "8080",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create an empty proxy to populate
			p := &Proxy{}

			// Call the function to test
			parseHostPort(p, tc.hostPort)

			// Check IP
			if p.Ip != tc.expected.Ip {
				t.Errorf("IP doesn't match: got %s, want %s", p.Ip, tc.expected.Ip)
			}

			// Check Port
			if p.Port != tc.expected.Port {
				t.Errorf("Port doesn't match: got %s, want %s", p.Port, tc.expected.Port)
			}
		})
	}
}

// TestParseHostPortWithExistingFields tests that parseHostPort doesn't affect
// existing fields in the Proxy struct (like User and Password)
func TestParseHostPortWithExistingFields(t *testing.T) {
	tests := []struct {
		name     string
		proxy    Proxy
		hostPort string
		expected Proxy
	}{
		{
			name: "Existing auth fields with IPv4",
			proxy: Proxy{
				User:     "testuser",
				Password: "testpass",
			},
			hostPort: "192.168.1.1:8080",
			expected: Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "testuser",
				Password: "testpass",
			},
		},
		{
			name: "Existing auth fields with IPv6",
			proxy: Proxy{
				User:     "testuser",
				Password: "testpass",
			},
			hostPort: "[2001:db8::1]:8080",
			expected: Proxy{
				Ip:       "[2001:db8::1]",
				Port:     "8080",
				User:     "testuser",
				Password: "testpass",
			},
		},
		{
			name: "Existing IP with new hostPort",
			proxy: Proxy{
				Ip:       "10.0.0.1",
				Port:     "9090",
				User:     "testuser",
				Password: "testpass",
			},
			hostPort: "192.168.1.1:8080",
			expected: Proxy{
				Ip:       "192.168.1.1", // Should be overwritten
				Port:     "8080",        // Should be overwritten
				User:     "testuser",    // Should remain unchanged
				Password: "testpass",    // Should remain unchanged
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a copy of the initial proxy
			p := &Proxy{
				Ip:       tc.proxy.Ip,
				Port:     tc.proxy.Port,
				User:     tc.proxy.User,
				Password: tc.proxy.Password,
			}

			// Call the function to test
			parseHostPort(p, tc.hostPort)

			// Check all fields
			if p.Ip != tc.expected.Ip {
				t.Errorf("IP doesn't match: got %s, want %s", p.Ip, tc.expected.Ip)
			}
			if p.Port != tc.expected.Port {
				t.Errorf("Port doesn't match: got %s, want %s", p.Port, tc.expected.Port)
			}
			if p.User != tc.expected.User {
				t.Errorf("User doesn't match: got %s, want %s", p.User, tc.expected.User)
			}
			if p.Password != tc.expected.Password {
				t.Errorf("Password doesn't match: got %s, want %s", p.Password, tc.expected.Password)
			}
		})
	}
}
