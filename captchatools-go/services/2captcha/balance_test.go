package twocaptcha

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalance_Success(t *testing.T) {
	const expectedErrorID int = 0
	const expectedBalance float32 = 123.45
	const expectedAPIKey string = "test-api-key"

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method and path
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, getBalanceEp, r.URL.Path)

		// Validate content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Decode and validate the request payload
		var payload map[string]string
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err, "Failed to decode request body")
		assert.Equal(t, expectedAPIKey, payload["clientKey"], "API key doesn't match")

		// Simulate a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody := map[string]interface{}{
			"errorId": expectedErrorID,
			"balance": expectedBalance,
		}
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	// Run the test with the mock server URL
	twocap := Twocaptcha{
		harvester.Config{
			Api_key: expectedAPIKey,
			Logger:  &mockLogger{},
		},
	}

	// Run the test
	resultBalance, err := twocap.getBalance(server.URL)
	require.NoError(t, err, "GetBalance should not return an error")
	assert.Equal(t, expectedBalance, resultBalance, "Balance doesn't match expected value")
}

func TestGetBalance_Error(t *testing.T) {
	const expectedErrorID int = 10
	const expectedAPIKey string = "test-api-key"

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return an error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 2Captcha uses 200 status even for errors
		responseBody := map[string]interface{}{
			"errorId": expectedErrorID,
			"balance": 0.0,
		}
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	// Run the test with the mock server URL
	twocap := Twocaptcha{
		Config: harvester.Config{
			Api_key: expectedAPIKey,
			Logger:  &mockLogger{},
		},
	}

	// Override the base URL for testing
	_, err := twocap.getBalance(server.URL)
	assert.Error(t, err)
}

// Mock logger to satisfy the Logger interface
type mockLogger struct{}

func (m *mockLogger) Info(format string, args ...interface{})  {}
func (m *mockLogger) Debug(format string, args ...interface{}) {}
func (m *mockLogger) Error(format string, args ...interface{}) {}
func (m *mockLogger) Warn(format string, args ...interface{})  {}
