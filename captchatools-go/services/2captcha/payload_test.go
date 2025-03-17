package twocaptcha

import (
	"reflect"
	"testing"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
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
			p := newPayload()
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
			p := newPayload()
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
			p := newPayload()
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
			p := newPayload()
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
			p := newPayload()
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
			p := newPayload()
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
			p := newPayload()
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
	p := newPayload()

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

func TestNewTaskPayload(t *testing.T) {
	tests := []struct {
		name           string
		config         harvester.Config
		options        []harvester.TokenOption
		expectedType   string
		expectedFields map[string]any
	}{
		{
			name: "V2 Captcha without proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options:      []harvester.TokenOption{},
			expectedType: "RecaptchaV2TaskProxyless",
			expectedFields: map[string]any{
				"websiteKey": "test_site_key",
				"websiteURL": "https://example.com",
			},
		},
		{
			name: "V2 Captcha with proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{Ip: "192.168.1.1", Port: "8080"}),
			},
			expectedType: "RecaptchaV2Task",
			expectedFields: map[string]any{
				"websiteKey":   "test_site_key",
				"websiteURL":   "https://example.com",
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
			},
		},
		{
			name: "V2 Captcha with proxy and auth",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{
					Ip:       "192.168.1.1",
					Port:     "8080",
					User:     "user",
					Password: "pass",
				}),
			},
			expectedType: "RecaptchaV2Task",
			expectedFields: map[string]any{
				"websiteKey":    "test_site_key",
				"websiteURL":    "https://example.com",
				"proxyAddress":  "192.168.1.1",
				"proxyPort":     "8080",
				"proxyLogin":    "user",
				"proxyPassword": "pass",
			},
		},
		{
			name: "V3 Captcha without proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V3Captcha,
				MinScore:    0.7,
				Action:      "login",
			},
			options:      []harvester.TokenOption{},
			expectedType: "RecaptchaV3TaskProxyless",
			expectedFields: map[string]any{
				"websiteKey": "test_site_key",
				"websiteURL": "https://example.com",
				"minScore":   float32(0.7),
				"pageAction": "login",
			},
		},
		{
			name: "V3 Captcha with proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V3Captcha,
				MinScore:    0.7,
				Action:      "login",
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{Ip: "192.168.1.1", Port: "8080"}),
				harvester.WithProxyType("HTTP"),
			},
			expectedType: "RecaptchaV3Task",
			expectedFields: map[string]any{
				"websiteKey":   "test_site_key",
				"websiteURL":   "https://example.com",
				"minScore":     float32(0.7),
				"pageAction":   "login",
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
				"proxyType":    "HTTP",
			},
		},
		{
			name: "HCaptcha without proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.HCaptcha,
			},
			options:      []harvester.TokenOption{},
			expectedType: "HCaptchaTaskProxyless",
			expectedFields: map[string]any{
				"sitekey":    "test_site_key",
				"websiteURL": "https://example.com",
			},
		},
		{
			name: "HCaptcha with proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.HCaptcha,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{Ip: "192.168.1.1", Port: "8080"}),
			},
			expectedType: "HCaptchaTask",
			expectedFields: map[string]any{
				"sitekey":      "test_site_key",
				"websiteURL":   "https://example.com",
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
			},
		},
		{
			name: "Cloudflare Turnstile without proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.CFTurnstile,
			},
			options:      []harvester.TokenOption{},
			expectedType: "TurnstileTaskProxyless",
			expectedFields: map[string]any{
				"websiteKey": "test_site_key",
				"websiteURL": "https://example.com",
			},
		},
		{
			name: "Cloudflare Turnstile with proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.CFTurnstile,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{Ip: "192.168.1.1", Port: "8080"}),
			},
			expectedType: "TurnstileTask",
			expectedFields: map[string]any{
				"websiteKey":   "test_site_key",
				"websiteURL":   "https://example.com",
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
			},
		},
		{
			name: "Image Captcha without proxy",
			config: harvester.Config{
				Api_key:     "test_api_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.ImageCaptcha,
			},
			options: []harvester.TokenOption{
				harvester.WithB64Img("base64encodedimage"),
			},
			expectedType: "ImageToTextTask",
			expectedFields: map[string]any{
				"body":       "base64encodedimage",
				"websiteURL": "https://example.com",
			},
		},
		{
			name: "V2 Invisible Captcha without proxy",
			config: harvester.Config{
				Api_key:            "test_api_key",
				Sitekey:            "test_site_key",
				CaptchaURL:         "https://example.com",
				CaptchaType:        harvester.V2Captcha,
				IsInvisibleCaptcha: true,
			},
			options:      []harvester.TokenOption{},
			expectedType: "RecaptchaV2TaskProxyless",
			expectedFields: map[string]any{
				"websiteKey": "test_site_key",
				"websiteURL": "https://example.com",
				"invisible":  1,
			},
		},
		{
			name: "With multiple token options",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.HCaptcha,
				SoftID:      12345,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{Ip: "192.168.1.1", Port: "8080"}),
				harvester.WithUserAgent("Mozilla/5.0"),
				harvester.WithRqData("custom_rqdata"),
			},
			expectedType: "HCaptchaTask",
			expectedFields: map[string]any{
				"sitekey":      "test_site_key",
				"websiteURL":   "https://example.com",
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
				"userAgent":    "Mozilla/5.0",
				"data":         "custom_rqdata",
				"soft_id":      12345,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			twoCaptcha := Twocaptcha{Config: tt.config}

			result := newTaskPayload(twoCaptcha, tt.options...)

			// Check that the client key is set correctly
			if clientKey, ok := result["clientKey"].(string); !ok || clientKey != tt.config.Api_key {
				t.Errorf("newTaskPayload() clientKey = %v, want %v", clientKey, tt.config.Api_key)
			}

			// Get the task data from the result
			taskData, ok := result["task"].(payload)
			if !ok {
				t.Fatalf("newTaskPayload() did not set 'task' field or it's not a payload type")
			}

			// Check the task type
			taskType, ok := taskData["type"].(string)
			if !ok || taskType != tt.expectedType {
				t.Errorf("newTaskPayload() task type = %v, want %v", taskType, tt.expectedType)
			}

			// Check all the expected fields in the task data
			for key, expectedValue := range tt.expectedFields {
				actualValue, exists := taskData[key]
				if !exists {
					t.Errorf("newTaskPayload() task data missing key %v", key)
					continue
				}

				if !reflect.DeepEqual(actualValue, expectedValue) {
					t.Errorf("newTaskPayload() task data[%v] = %v, want %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}

func TestNewTaskPayloadEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		config         harvester.Config
		options        []harvester.TokenOption
		expectedType   string
		expectedFields map[string]any
		checkNegative  bool
		negativeKeys   []string
	}{
		{
			name: "Empty API key",
			config: harvester.Config{
				Api_key:     "",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options:      []harvester.TokenOption{},
			expectedType: "RecaptchaV2TaskProxyless",
			expectedFields: map[string]any{
				"websiteKey": "test_site_key",
				"websiteURL": "https://example.com",
			},
		},
		{
			name: "Nil proxy in options",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(nil),
			},
			expectedType: "RecaptchaV2TaskProxyless", // Should still be proxyless since nil proxy shouldn't change type
			expectedFields: map[string]any{
				"websiteKey": "test_site_key",
				"websiteURL": "https://example.com",
			},
			checkNegative: true,
			negativeKeys:  []string{"proxyAddress", "proxyPort"}, // These shouldn't be set
		},
		{
			name: "Empty proxy fields",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{Ip: "", Port: ""}),
			},
			expectedType: "RecaptchaV2Task", // Even with empty strings, it should switch to proxy version
			expectedFields: map[string]any{
				"websiteKey":   "test_site_key",
				"websiteURL":   "https://example.com",
				"proxyAddress": "",
				"proxyPort":    "",
			},
		},
		{
			name: "Proxy with incomplete auth",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "test_site_key",
				CaptchaURL:  "https://example.com",
				CaptchaType: harvester.V2Captcha,
			},
			options: []harvester.TokenOption{
				harvester.WithProxy(&proxy.Proxy{
					Ip:       "192.168.1.1",
					Port:     "8080",
					User:     "user",
					Password: "", // Missing password
				}),
			},
			expectedType: "RecaptchaV2Task",
			expectedFields: map[string]any{
				"websiteKey":   "test_site_key",
				"websiteURL":   "https://example.com",
				"proxyAddress": "192.168.1.1",
				"proxyPort":    "8080",
			},
			checkNegative: true,
			negativeKeys:  []string{"proxyLogin", "proxyPassword"}, // These shouldn't be set
		},
		{
			name: "Empty sitekey and URL",
			config: harvester.Config{
				Api_key:     "test_api_key",
				Sitekey:     "",
				CaptchaURL:  "",
				CaptchaType: harvester.V2Captcha,
			},
			options:      []harvester.TokenOption{},
			expectedType: "RecaptchaV2TaskProxyless",
			expectedFields: map[string]any{
				"websiteKey": "",
				"websiteURL": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			twoCaptcha := Twocaptcha{Config: tt.config}

			result := newTaskPayload(twoCaptcha, tt.options...)

			// Check that the client key is set correctly
			if clientKey, ok := result["clientKey"].(string); !ok || clientKey != tt.config.Api_key {
				t.Errorf("newTaskPayload() clientKey = %v, want %v", clientKey, tt.config.Api_key)
			}

			// Get the task data from the result
			taskData, ok := result["task"].(payload)
			if !ok {
				t.Fatalf("newTaskPayload() did not set 'task' field or it's not a payload type")
			}

			// Check the task type
			taskType, ok := taskData["type"].(string)
			if !ok || taskType != tt.expectedType {
				t.Errorf("newTaskPayload() task type = %v, want %v", taskType, tt.expectedType)
			}

			// Check all the expected fields in the task data
			for key, expectedValue := range tt.expectedFields {
				actualValue, exists := taskData[key]
				if !exists {
					t.Errorf("newTaskPayload() task data missing key %v", key)
					continue
				}

				if !reflect.DeepEqual(actualValue, expectedValue) {
					t.Errorf("newTaskPayload() task data[%v] = %v, want %v", key, actualValue, expectedValue)
				}
			}

			// Check that certain keys don't exist if specified
			if tt.checkNegative {
				for _, key := range tt.negativeKeys {
					if _, exists := taskData[key]; exists {
						t.Errorf("newTaskPayload() task data should not have key %v but it does", key)
					}
				}
			}
		})
	}
}

func TestSetCustom(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    any
		expected map[string]any
	}{
		{
			name:  "String value",
			key:   "customParam",
			value: "customValue",
			expected: map[string]any{
				"customParam": "customValue",
			},
		},
		{
			name:  "Integer value",
			key:   "customInt",
			value: 42,
			expected: map[string]any{
				"customInt": 42,
			},
		},
		{
			name:  "Float value",
			key:   "customFloat",
			value: 3.14,
			expected: map[string]any{
				"customFloat": 3.14,
			},
		},
		{
			name:  "Boolean value",
			key:   "customBool",
			value: true,
			expected: map[string]any{
				"customBool": true,
			},
		},
		{
			name:  "Nil value",
			key:   "customNil",
			value: nil,
			expected: map[string]any{
				"customNil": nil,
			},
		},
		{
			name:  "Empty key",
			key:   "",
			value: "value",
			expected: map[string]any{
				"": "value",
			},
		},
		{
			name:  "Slice value",
			key:   "customSlice",
			value: []string{"a", "b", "c"},
			expected: map[string]any{
				"customSlice": []string{"a", "b", "c"},
			},
		},
		{
			name:  "Map value",
			key:   "customMap",
			value: map[string]int{"a": 1, "b": 2},
			expected: map[string]any{
				"customMap": map[string]int{"a": 1, "b": 2},
			},
		},
		{
			name:  "Overwrite existing key",
			key:   "existingKey",
			value: "newValue",
			expected: map[string]any{
				"existingKey": "newValue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new payload
			p := newPayload()

			// For the overwrite test, pre-populate the payload
			if tt.name == "Overwrite existing key" {
				p["existingKey"] = "oldValue"
			}

			// Call the method being tested
			p.SetCustom(tt.key, tt.value)

			// Check results
			for expectedKey, expectedValue := range tt.expected {
				// Verify key exists
				actualValue, exists := p[expectedKey]
				if !exists {
					t.Errorf("SetCustom() did not set key %q in payload", expectedKey)
					continue
				}

				// Verify value is correct
				if !reflect.DeepEqual(actualValue, expectedValue) {
					t.Errorf("SetCustom() set %q = %v, want %v", expectedKey, actualValue, expectedValue)
				}
			}

			// For the overwrite test, verify the value changed
			if tt.name == "Overwrite existing key" {
				if p["existingKey"] == "oldValue" {
					t.Errorf("SetCustom() did not overwrite existing value")
				}
			}
		})
	}
}
