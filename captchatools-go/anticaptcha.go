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

// This file will contain the code to interact with anticaptcha.com API

type Anticaptcha struct {
	config *Config
}

type anticaptchaBalanceResponse struct {
	ErrorID          int     `json:"errorId"`
	ErrorCode        string  `json:"errorCode"`
	ErrorDescription string  `json:"errorDescription"`
	Balance          float32 `json:"balance"`
}

type anticaptchaResponse struct {
	ErrorID   int    `json:"errorId"`
	ErrorCode string `json:"errorCode"`
	Status    string `json:"status"`
	Solution  struct {
		Token              string `json:"token"`
		Text               string `json:"text"`
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
		UserAgent          string `json:"userAgent"`
	} `json:"solution"`
}

func (a Anticaptcha) GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return a.getCaptchaAnswer(context.Background(), additional...)
}

func (a Anticaptcha) GetTokenWithContext(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return a.getCaptchaAnswer(ctx, additional...)
}

func (a Anticaptcha) GetBalance() (float32, error) {
	return a.getBalance()
}

// Method to get Queue ID from the API.
func (a Anticaptcha) getID(data *AdditionalData) (int, error) {
	// Get Payload
	payload, err := a.createPayload(data)
	if err != nil {
		return 0, err
	}

	// Make request to get answer
	response := &capmonsterIDResponse{}
	for i := 0; i < 100; i++ {
		resp, err := http.Post("https://api.anti-captcha.com/createTask", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
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
func (a Anticaptcha) getCaptchaAnswer(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	var data *AdditionalData = nil
	if len(additional) > 0 {
		data = additional[0]
	}

	// Get Queue ID
	queueID, err := a.getID(data)
	if err != nil {
		return nil, err
	}

	// Get Captcha Answer
	payload, _ := json.Marshal(capmonsterCapAnswerPayload{
		ClientKey: a.config.Api_key,
		TaskID:    queueID,
	})
	response := &anticaptchaResponse{}
	for i := 0; i < 100; i++ {
		reqToMake, _ := http.NewRequestWithContext(ctx, "POST", "https://api.anti-captcha.com/getTaskResult", bytes.NewBufferString(string(payload)))
		reqToMake.Header.Add("Content-Type", "application/json")
		resp, err := makeRequest(reqToMake)
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
		if response.ErrorID > 0 { // means there was an error
			return nil, errCodeToError(response.ErrorCode)
		}

		// Check if the answer is ready or not
		if response.Status == "processing" {
			time.Sleep(3 * time.Second)
			continue
		}

		var solution string
		var ua string = response.Solution.UserAgent
		switch a.config.CaptchaType {
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
			a.config.Api_key,
			a.config.CaptchaType,
			AnticaptchaSite,
			ua,
		), nil
	}
	return nil, ErrMaxAttempts
}

func (a Anticaptcha) getBalance() (float32, error) {
	// Attempt to get the balance from the API
	// Max attempts is 5
	payload := fmt.Sprintf(`{"clientKey": "%v"}`, a.config.Api_key)
	response := &anticaptchaBalanceResponse{}
	for i := 0; i < 5; i++ {
		resp, err := http.Post("https://api.anti-captcha.com/getBalance", "application/json", bytes.NewBuffer([]byte(payload)))
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
func (a Anticaptcha) createPayload(data *AdditionalData) (string, error) {
	// Define the payload we are going to send to the API
	payload := capmonsterIDPayload{
		ClientKey: a.config.Api_key,
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
			WebsiteURL: a.config.CaptchaURL,
			WebsiteKey: a.config.Sitekey,
			Type:       a.config.CaptchaType,
		},
	}

	// Add any other keys to the payload
	if a.config.SoftID != 0 {
		payload.SoftID = a.config.SoftID
	}
	switch a.config.CaptchaType {
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
			payload.Task.Type = "RecaptchaV2Task"
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

			payload.Task.UserAgent = data.UserAgent // REQUIRED WITH PROXIES

		}

		if a.config.IsInvisibleCaptcha {
			payload.Task.IsInvisible = a.config.IsInvisibleCaptcha
		}
	case V3Captcha:
		payload.Task.Type = "RecaptchaV3TaskProxyless"
		payload.Task.MinScore = a.config.MinScore
		payload.Task.PageAction = a.config.Action
	case HCaptcha:
		payload.Task.Type = "HCaptchaTaskProxyless"
	case CFTurnstile:
		payload.Task.Type = "TurnstileTaskProxyless"
	default:
		return "", ErrIncorrectCapType
	}
	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}

func report_anticaptcha(was_correct bool, c *CaptchaAnswer) error {
	var endpoint string
	switch c.capType {
	case V2Captcha, V3Captcha:
		endpoint = "reportIncorrectRecaptcha"
	case HCaptcha:
		endpoint = "reportIncorrectHcaptcha"
	default:
		return ErrIncorrectCapType
	}
	if was_correct && (c.capType == V2Captcha || c.capType == V3Captcha) {
		endpoint = "reportCorrectRecaptcha"
	}

	// Make request
	response := &capmonsterTokenResponse{}
	payload, _ := json.Marshal(capmonsterCapAnswerPayload{
		ClientKey: c.api_key,
		TaskID:    c.id.(int),
	})
	for i := 0; i < 100; i++ {
		resp, err := http.Post("https://api.anti-captcha.com/"+endpoint, "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Check for any errors
		if response.ErrorID > 0 { // means there was an error
			return errCodeToError(response.ErrorCode)
		}
		return nil
	}
	return ErrMaxAttempts
}
