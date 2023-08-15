package captchatoolsgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

/*
   This file will contain the code to interact with capmonster.cloud API
*/

type Capmonster struct {
	config *Config
}

// This struct will be the payload to get the queue ID from capmonster
type capmonsterIDPayload struct {
	ClientKey string `json:"clientKey"`
	SoftID    int    `json:"softId,omitempty"`
	Task      struct {
		WebsiteURL    string      `json:"websiteURL"`
		WebsiteKey    string      `json:"websiteKey"`
		Type          captchaType `json:"type"`
		IsInvisible   bool        `json:"isInvisible,omitempty"`
		MinScore      float32     `json:"minScore,omitempty"`
		PageAction    string      `json:"pageAction,omitempty"`
		Body          string      `json:"body,omitempty"`
		ProxyType     string      `json:"proxyType,omitempty"`
		ProxyAddress  string      `json:"proxyAddress,omitempty"`
		ProxyPort     int         `json:"proxyPort,omitempty"`
		ProxyLogin    string      `json:"proxyLogin,omitempty"`
		ProxyPassword string      `json:"proxyPassword,omitempty"`
		UserAgent     string      `json:"userAgent,omitempty"`
	} `json:"task"`
}
type capmonsterCapAnswerPayload struct {
	ClientKey string `json:"clientKey"`
	TaskID    int    `json:"taskId"`
}

type capmonsterIDResponse struct {
	ErrorID   int    `json:"errorId"`
	ErrorCode string `json:"errorCode"`
	TaskID    int    `json:"taskId"`
}

type capmonsterTokenResponse struct {
	ErrorID   int    `json:"errorId"`
	ErrorCode string `json:"errorCode"`
	Solution  struct {
		Token              string `json:"token"`
		Text               string `json:"text"`
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
		UserAgent          string `json:"userAgent"`
	} `json:"solution"`
	Status string `json:"status"`
}
type capmonsterBalanceResponse struct {
	Balance   float32 `json:"balance"`
	ErrorCode string  `json:"errorCode"`
	ErrorID   int     `json:"errorId"`
}

func (c Capmonster) GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return c.getCaptchaAnswer(context.Background(), additional...)
}

func (c Capmonster) GetTokenWithContext(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return c.getCaptchaAnswer(ctx, additional...)
}

func (c Capmonster) GetBalance() (float32, error) {
	return c.getBalance()
}

// Method to get Queue ID from the API.
func (c Capmonster) getID(data *AdditionalData) (int, error) {
	// Get Payload
	payload, err := c.createPayload(data)
	if err != nil {
		return 0, err
	}

	// Make request to get answer
	for i := 0; i < 100; i++ {
		resp, err := http.Post("https://api.capmonster.cloud/createTask", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		response := &capmonsterIDResponse{}
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Parse the response
		if response.ErrorID != 0 { // Means there was an error
			return 0, errCodeToError(response.ErrorCode)
		}
		return response.TaskID, nil
	}
	return 0, ErrMaxAttempts
}

// This method gets the captcha token from the Capmonster API
func (c Capmonster) getCaptchaAnswer(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	var data *AdditionalData = nil
	if len(additional) > 0 {
		data = additional[0]
	}

	// Get Queue ID
	queueID, err := c.getID(data)
	if err != nil {
		return nil, err
	}

	// Get Captcha Answer
	payload, _ := json.Marshal(capmonsterCapAnswerPayload{
		ClientKey: c.config.Api_key,
		TaskID:    queueID,
	})
	response := &capmonsterTokenResponse{}
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.capmonster.cloud/getTaskResult", bytes.NewBufferString(string(payload)))
		req.Header.Add("Content-Type", "application/json")
		resp, err := makeRequest(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, fmt.Errorf("getCaptchaAnswer error: %w", err)
			}
			time.Sleep(3 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Check for any errors
		if response.ErrorID != 0 { // means there was an error
			return nil, errCodeToError(response.ErrorCode)
		}

		// Check if captcha is ready
		if response.Status == "processing" {
			time.Sleep(3 * time.Second)
			continue
		}

		var solution string
		var ua string = response.Solution.UserAgent
		switch c.config.CaptchaType {
		case V2Captcha, V3Captcha:
			solution = response.Solution.GRecaptchaResponse
		case ImageCaptcha:
			solution = response.Solution.Text
		case CFTurnstile:
			solution = response.Solution.Token
		}

		return newCaptchaAnswer(
			queueID,
			solution,
			c.config.Api_key,
			c.config.CaptchaType,
			AnticaptchaSite,
			ua,
		), nil
	}
	return nil, ErrMaxAttempts
}

// getBalance() returns the balance on the API key
func (c Capmonster) getBalance() (float32, error) {
	// Attempt to get the balance from the API
	// Max attempts is 5
	payload := fmt.Sprintf(`{"clientKey": "%v"}`, c.config.Api_key)
	response := &capmonsterBalanceResponse{}
	for i := 0; i < 5; i++ {
		resp, err := http.Post("https://api.capmonster.cloud/getBalance", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)
		if response.ErrorID != 0 {
			return 0, errCodeToError(response.ErrorCode)
		}
		return response.Balance, nil
	}
	return 0, ErrMaxAttempts
}

/*
createPayload returns the payloads required to interact with the API.

Possible errors that can be returned:
1) ErrIncorrectCapType
*/
func (c Capmonster) createPayload(data *AdditionalData) (string, error) {
	// Define the payload we are going to send to the API
	payload := capmonsterIDPayload{
		ClientKey: c.config.Api_key,
		Task: struct {
			WebsiteURL    string      "json:\"websiteURL\""
			WebsiteKey    string      "json:\"websiteKey\""
			Type          captchaType "json:\"type\""
			IsInvisible   bool        "json:\"isInvisible,omitempty\""
			MinScore      float32     "json:\"minScore,omitempty\""
			PageAction    string      "json:\"pageAction,omitempty\""
			Body          string      "json:\"body,omitempty\""
			ProxyType     string      "json:\"proxyType,omitempty\""
			ProxyAddress  string      "json:\"proxyAddress,omitempty\""
			ProxyPort     int         "json:\"proxyPort,omitempty\""
			ProxyLogin    string      "json:\"proxyLogin,omitempty\""
			ProxyPassword string      "json:\"proxyPassword,omitempty\""
			UserAgent     string      "json:\"userAgent,omitempty\""
		}{
			WebsiteURL: c.config.CaptchaURL,
			WebsiteKey: c.config.Sitekey,
			Type:       c.config.CaptchaType,
		},
	}

	// Add any other keys to the payload
	switch c.config.CaptchaType {
	case ImageCaptcha:
		if data == nil {
			return "", ErrAddionalDataMissing
		}
		payload.Task.Type = "ImageToTextTask"
		payload.Task.Body = data.B64Img
	case V2Captcha:
		payload.Task.Type = "NoCaptchaTaskProxyless"

		// Check for proxy data
		if data != nil && data.Proxy != nil {
			payload.Task.Type = "NoCaptchaTask"
			if data.ProxyType == "" {
				data.ProxyType = "http"
			}
			payload.Task.ProxyType = data.ProxyType
			payload.Task.ProxyAddress = data.Proxy.Ip
			portInt, err := strconv.Atoi(data.Proxy.Port)
			if err != nil {
				return "", errors.New("error converting proxy port to int")
			}
			payload.Task.ProxyPort = portInt
			if data.Proxy.IsUserAuth() {
				payload.Task.ProxyLogin = data.Proxy.User
				payload.Task.ProxyPassword = data.Proxy.Password
			}
		}

		if c.config.IsInvisibleCaptcha {
			payload.Task.IsInvisible = c.config.IsInvisibleCaptcha
		}
	case V3Captcha:
		payload.Task.Type = "RecaptchaV3TaskProxyless"
		payload.Task.MinScore = c.config.MinScore
		payload.Task.PageAction = c.config.Action
	case HCaptcha:
		payload.Task.Type = "HCaptchaTaskProxyless"
	case CFTurnstile:
		payload.Task.Type = "TurnstileTaskProxyless"
	default:
		return "", ErrIncorrectCapType
	}

	// Check for addtional data
	if data != nil && data.UserAgent != "" {
		payload.Task.UserAgent = data.UserAgent
	}

	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}
