package httputils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseError(t *testing.T) {
	tests := []struct {
		name          string
		inputErr      error
		expectedType  interface{}
		expectedField string
	}{
		{
			name:          "Invalid Port Error",
			inputErr:      &net.AddrError{Err: "invalid port", Addr: "192.168.1.1:70000"},
			expectedType:  &ProxyErr{},
			expectedField: "InvalidPort",
		},
		{
			name:          "Invalid Host Error",
			inputErr:      &net.DNSError{Err: "no such host", Name: "nonexistenthost.com", IsNotFound: true},
			expectedType:  &ProxyErr{},
			expectedField: "InvalidHost",
		},
		{
			name:          "Connection Refused Error",
			inputErr:      &url.Error{Op: "Get", URL: "http://localhost:9999", Err: syscall.ECONNREFUSED},
			expectedType:  &ProxyErr{},
			expectedField: "",
		},
		{
			name:          "EOF Error",
			inputErr:      &url.Error{Op: "Get", URL: "http://example.com", Err: fmt.Errorf("EOF")},
			expectedType:  &RequestErr{},
			expectedField: "Message",
		},
		{
			name:          "Proxy Auth Error",
			inputErr:      &url.Error{Op: "Get", URL: "http://example.com", Err: fmt.Errorf("Proxy responded with non 200 code: 407 Proxy Authentication Required")},
			expectedType:  &ProxyErr{},
			expectedField: "InvalidAuth",
		},
		{
			name:          "Proxy Relay Error",
			inputErr:      &url.Error{Op: "Get", URL: "http://example.com", Err: fmt.Errorf("Proxy responded with non 200 code: 502 Proxy Error (The selected relay is offline or busy processing other threads)")},
			expectedType:  &ProxyErr{},
			expectedField: "",
		},
		{
			name:          "Invalid URL Scheme Error",
			inputErr:      &url.Error{Op: "Get", URL: "http://example.com", Err: fmt.Errorf("invalid URL scheme: []")},
			expectedType:  &ClientErr{},
			expectedField: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseError(tt.inputErr)

			// Check the type
			assert.IsType(t, tt.expectedType, result, "Expected type %T but got %T", tt.expectedType, result)

			// Check specific fields based on the error type
			switch tt.expectedType.(type) {
			case *ProxyErr:
				if proxyErr, ok := result.(*ProxyErr); ok {
					switch tt.expectedField {
					case "InvalidAuth":
						assert.True(t, proxyErr.InvalidAuth, "Expected InvalidAuth to be true")
					case "InvalidHost":
						assert.True(t, proxyErr.InvalidHost, "Expected InvalidHost to be true")
					case "InvalidPort":
						assert.True(t, proxyErr.InvalidPort, "Expected InvalidPort to be true")
					}
				}
			case *ClientErr:
				if clientErr, ok := result.(*ClientErr); ok && tt.expectedField == "TimedOut" {
					assert.True(t, clientErr.TimedOut, "Expected TimedOut to be true")
				}
			case *RequestErr:
				if requestErr, ok := result.(*RequestErr); ok && tt.expectedField == "Message" {
					assert.Equal(t, "EOF", requestErr.Message, "Expected Message to be 'EOF'")
				}
			}
		})
	}
}

func TestIsProxyError(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		expected bool
	}{
		{
			name:     "Proxy Error - Invalid Auth",
			inputErr: &ProxyErr{InvalidAuth: true},
			expected: true,
		},
		{
			name:     "Proxy Error - Invalid Host",
			inputErr: &ProxyErr{InvalidHost: true},
			expected: true,
		},
		{
			name:     "Proxy Error - Invalid Port",
			inputErr: &ProxyErr{InvalidPort: true},
			expected: true,
		},
		{
			name:     "Wrapped Proxy Error",
			inputErr: fmt.Errorf("wrapped: %w", &ProxyErr{}),
			expected: true,
		},
		{
			name:     "Client Error",
			inputErr: &ClientErr{},
			expected: false,
		},
		{
			name:     "Generic Error",
			inputErr: errors.New("generic error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsProxyError(tt.inputErr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsClientError(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		expected bool
	}{
		{
			name:     "Client Error - Timeout",
			inputErr: &ClientErr{TimedOut: true},
			expected: true,
		},
		{
			name:     "Client Error - Generic",
			inputErr: &ClientErr{},
			expected: true,
		},
		{
			name:     "Wrapped Client Error",
			inputErr: fmt.Errorf("wrapped: %w", &ClientErr{}),
			expected: true,
		},
		{
			name:     "Proxy Error",
			inputErr: &ProxyErr{},
			expected: false,
		},
		{
			name:     "Generic Error",
			inputErr: errors.New("generic error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsClientError(tt.inputErr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
