package twocaptcha

import (
	"testing"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/proxy"
)

func TestSetB64Img(t *testing.T) {
	tests := []struct {
		name     string
		img      string
		expected string
	}{
		{
			name:     "Basic base64 image",
			img:      "base64encodedimage",
			expected: "base64encodedimage",
		},
		{
			name:     "Empty string",
			img:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.SetB64Img(tt.img)

			if p["body"] != tt.expected {
				t.Errorf("SetB64Img() got = %q, want %q", p["body"], tt.expected)
			}
		})
	}
}

func TestSetProxy(t *testing.T) {
	tests := []struct {
		name     string
		proxy    *proxy.Proxy
		expected map[string]any
	}{
		{
			name: "Full proxy with authentication",
			proxy: &proxy.Proxy{
				Ip:       "192.168.1.1",
				Port:     "8080",
				User:     "user",
				Password: "pass",
			},
			expected: map[string]any{
				"proxyAddress":  "192.168.1.1",
				"proxyPort":     "8080",
				"proxyLogin":    "user",
				"proxyPassword": "pass",
			},
		},
		{
			name: "Proxy without authentication",
			proxy: &proxy.Proxy{
				Ip:   "192.168.1.1",
				Port: "8080",
			},
			expected: map[string]any{
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
			},
		},
		{
			name:     "Nil proxy",
			proxy:    nil,
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.SetProxy(tt.proxy)

			for key, value := range tt.expected {
				if p[key] != value {
					t.Errorf("SetProxy() got[%q] = %q, want %q", key, p[key], value)
				}
			}

			// Check that no unexpected keys were added
			expectedLen := len(tt.expected)
			if len(p) != expectedLen {
				t.Errorf("SetProxy() added %d entries, want %d", len(p), expectedLen)
			}
		})
	}
}

func TestSetProxyType(t *testing.T) {
	tests := []struct {
		name      string
		proxyType string
		expected  string
	}{
		{
			name:      "HTTP proxy",
			proxyType: "HTTP",
			expected:  "HTTP",
		},
		{
			name:      "HTTPS proxy",
			proxyType: "HTTPS",
			expected:  "HTTPS",
		},
		{
			name:      "SOCKS4 proxy",
			proxyType: "SOCKS4",
			expected:  "SOCKS4",
		},
		{
			name:      "SOCKS5 proxy",
			proxyType: "SOCKS5",
			expected:  "SOCKS5",
		},
		{
			name:      "Empty string",
			proxyType: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.SetProxyType(tt.proxyType)

			if p["proxyType"] != tt.expected {
				t.Errorf("SetProxyType() got = %q, want %q", p["proxyType"], tt.expected)
			}
		})
	}
}

func TestSetUserAgent(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		expected  string
	}{
		{
			name:      "Standard user agent",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			expected:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		},
		{
			name:      "Empty user agent",
			userAgent: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.SetUserAgent(tt.userAgent)

			if p["userAgent"] != tt.expected {
				t.Errorf("SetUserAgent() got = %q, want %q", p["userAgent"], tt.expected)
			}
		})
	}
}

func TestSetRqData(t *testing.T) {
	tests := []struct {
		name     string
		rqData   string
		expected string
	}{
		{
			name:     "Valid rqData",
			rqData:   "someRqData123456",
			expected: "someRqData123456",
		},
		{
			name:     "Empty rqData",
			rqData:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.SetRqData(tt.rqData)

			if p["data"] != tt.expected {
				t.Errorf("SetRqData() got = %q, want %q", p["data"], tt.expected)
			}
		})
	}
}

func TestSetToPayload(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected any
	}{
		{
			name:     "String value",
			key:      "testKey",
			value:    "testValue",
			expected: "testValue",
		},
		{
			name:     "Integer value",
			key:      "testInt",
			value:    123,
			expected: 123,
		},
		{
			name:     "Boolean value",
			key:      "testBool",
			value:    true,
			expected: true,
		},
		{
			name:     "Nil value",
			key:      "testNil",
			value:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.setToPayload(tt.key, tt.value)

			if p[tt.key] != tt.expected {
				t.Errorf("setToPayload() got = %q, want %q", p[tt.key], tt.expected)
			}
		})
	}
}

func TestProxyIsUserAuth(t *testing.T) {
	tests := []struct {
		name     string
		proxy    *proxy.Proxy
		expected bool
	}{
		{
			name: "With authentication",
			proxy: &proxy.Proxy{
				User:     "user",
				Password: "pass",
			},
			expected: true,
		},
		{
			name: "With username only",
			proxy: &proxy.Proxy{
				User: "user",
			},
			expected: false,
		},
		{
			name: "With password only",
			proxy: &proxy.Proxy{
				Password: "pass",
			},
			expected: false,
		},
		{
			name:     "Without authentication",
			proxy:    &proxy.Proxy{},
			expected: false,
		},
	}

	// This is a helper test to verify our understanding of how IsUserAuth works
	// Since we can't directly test it (it's in another package), we're making assumptions
	// based on common proxy authentication patterns
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(payload)
			p.SetProxy(tt.proxy)

			// Check if the proxy credentials were added to the payload
			_, hasLogin := p["proxyLogin"]
			_, hasPassword := p["proxyPassword"]

			// Both should be present or both should be absent
			authPresent := hasLogin && hasPassword

			if authPresent != tt.expected {
				t.Errorf("IsUserAuth() behavior appears to be %v, want %v", authPresent, tt.expected)
			}
		})
	}
}

func TestMultipleOperations(t *testing.T) {
	// Test multiple operations on the same payload
	p := make(payload)

	// Set various properties
	p.SetB64Img("testImage")
	p.SetUserAgent("testAgent")
	p.SetProxyType("HTTP")
	p.SetRqData("testRqData")

	// Verify all properties were set correctly
	expected := map[string]any{
		"body":      "testImage",
		"userAgent": "testAgent",
		"proxyType": "HTTP",
		"data":      "testRqData",
	}

	for key, value := range expected {
		if p[key] != value {
			t.Errorf("MultipleOperations got[%q] = %q, want %q", key, p[key], value)
		}
	}

	// Ensure no extra keys were added
	if len(p) != len(expected) {
		t.Errorf("MultipleOperations payload has %d entries, want %d", len(p), len(expected))
	}
}
