package twocaptcha

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTaskID_Success(t *testing.T) {
	const expectedTaskID int64 = 12345678
	const expectedAPIKey = "test-api-key"

	// Create a mock server for getTaskID endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, getTaskIDEp, r.URL.Path)

		// Verify content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Read and verify request body
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)

		// Verify payload contains expected values
		assert.Equal(t, expectedAPIKey, payload["clientKey"])

		// Send success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Important: This needs to match the exact response structure expected by parseGetTaskId
		responseBody := struct {
			ErrorID int   `json:"errorId"`
			TaskID  int64 `json:"taskId"`
		}{
			ErrorID: 0,
			TaskID:  expectedTaskID,
		}

		err = json.NewEncoder(w).Encode(responseBody)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create Twocaptcha client
	twocap := Twocaptcha{
		Config: harvester.Config{
			Api_key: expectedAPIKey,
			Logger:  &mockLogger{},
		},
	}

	// Create payload with task data as a map
	p := newPayloadWithClientKey(expectedAPIKey)
	p.setToPayload("type", "RecaptchaV2TaskProxyless")
	p.setToPayload("websiteURL", "https://2captcha.com/demo/recaptcha-v2")
	p.setToPayload("websiteKey", "6LfD3PIbAAAAAJs_eEHvoOl75_83eXSqpPSRFJ_u")

	// Call getTaskID method
	taskID, err := twocap.getTaskID(context.Background(), server.URL, p)

	// Verify results
	require.NoError(t, err)
	assert.Equal(t, expectedTaskID, taskID)
}

func TestGetTaskID_WithExtendedPayload(t *testing.T) {
	const expectedTaskID int64 = 12345678
	const expectedAPIKey = "test-api-key"
	const expectedWebsiteURL = "https://extended-example.com"
	const expectedWebsiteKey = "extended-sitekey-value"
	const expectedUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	const expectedRqData = "custom-rqdata-value"

	// Create a mock server for getTaskID endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and parse request body
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)

		// Verify Payload
		assert.Equal(t, expectedAPIKey, payload["clientKey"])
		task, ok := payload["task"].(map[string]interface{})
		require.True(t, ok, "Task should be a map in the payload")
		assert.Equal(t, expectedWebsiteURL, task["websiteURL"], "websiteURL should match")
		assert.Equal(t, expectedWebsiteKey, task["websiteKey"], "websiteKey should match")
		assert.Equal(t, expectedUserAgent, task["userAgent"], "userAgent should match")
		assert.Equal(t, expectedRqData, task["data"], "data should match")

		// Create the expected response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Match the structure expected by parseGetTaskId
		responseBody := struct {
			ErrorID int   `json:"errorId"`
			TaskID  int64 `json:"taskId"`
		}{
			ErrorID: 0,
			TaskID:  expectedTaskID,
		}

		err = json.NewEncoder(w).Encode(responseBody)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create Twocaptcha client
	twocap := Twocaptcha{
		Config: harvester.Config{
			Api_key: expectedAPIKey,
			Logger:  &mockLogger{},
		},
	}

	// Create the task with the proper structure
	// Use a nested map to create the proper structure the API expects
	taskData := map[string]interface{}{
		"type":       "RecaptchaV2TaskProxyless",
		"websiteURL": expectedWebsiteURL,
		"websiteKey": expectedWebsiteKey,
		"userAgent":  expectedUserAgent,
		"data":       expectedRqData,
	}

	// Create the payload with proper structure
	p := newPayloadWithClientKey(expectedAPIKey)
	p["task"] = taskData

	// Call getTaskID method
	taskID, err := twocap.getTaskID(context.Background(), server.URL, p)

	// Verify results
	require.NoError(t, err)
	assert.Equal(t, expectedTaskID, taskID)
}

func TestGetTaskID_WithProxy(t *testing.T) {
	const expectedTaskID int64 = 12345678
	const expectedAPIKey = "test-api-key"
	const expectedProxyAddress = "192.168.1.1"
	const expectedProxyPort = "8080"
	const expectedProxyLogin = "proxyuser"
	const expectedProxyPassword = "proxypass"
	const expectedProxyType = "HTTP"

	// Create a mock server for getTaskID endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and parse request body
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		require.NoError(t, err)

		// Assert the payload
		assert.Equal(t, expectedAPIKey, payload["clientKey"])
		task, ok := payload["task"].(map[string]interface{})
		require.True(t, ok, "Task should be a map in the payload")
		assert.Equal(t, expectedProxyAddress, task["proxyAddress"], "Proxy address should match")
		assert.Equal(t, expectedProxyPort, task["proxyPort"], "Proxy port should match")
		assert.Equal(t, expectedProxyLogin, task["proxyLogin"], "Proxy login should match")
		assert.Equal(t, expectedProxyPassword, task["proxyPassword"], "Proxy password should match")
		assert.Equal(t, expectedProxyType, task["proxyType"], "Proxy type should match")
		assert.Equal(t, "RecaptchaV2Task", task["type"], "Task type should be RecaptchaV2Task")
		assert.Equal(t, "https://example.com", task["websiteURL"], "Website URL should match")
		assert.Equal(t, "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", task["websiteKey"], "Website key should match")

		// Create response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Match the structure expected by parseGetTaskId
		responseBody := struct {
			ErrorID int   `json:"errorId"`
			TaskID  int64 `json:"taskId"`
		}{
			ErrorID: 0,
			TaskID:  expectedTaskID,
		}

		err = json.NewEncoder(w).Encode(responseBody)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create Twocaptcha client
	twocap := Twocaptcha{
		Config: harvester.Config{
			Api_key: expectedAPIKey,
			Logger:  &mockLogger{},
		},
	}

	// Create the task with proper structure
	taskData := map[string]interface{}{
		"type":          "RecaptchaV2Task", // Not proxyless since we're using proxy
		"websiteURL":    "https://example.com",
		"websiteKey":    "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		"proxyAddress":  expectedProxyAddress,
		"proxyPort":     expectedProxyPort,
		"proxyLogin":    expectedProxyLogin,
		"proxyPassword": expectedProxyPassword,
		"proxyType":     expectedProxyType,
	}

	// Create payload with proper structure
	p := newPayloadWithClientKey(expectedAPIKey)
	p["task"] = taskData

	// Call getTaskID method
	taskID, err := twocap.getTaskID(context.Background(), server.URL, p)

	// Verify results
	require.NoError(t, err)
	assert.Equal(t, expectedTaskID, taskID)
}
