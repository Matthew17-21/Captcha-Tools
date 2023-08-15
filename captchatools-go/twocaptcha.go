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

type Twocaptcha struct {
	config *Config
}

// Type that will be used when sending a payload to 2captcha
type twoCapIDPayload struct {
	Key       string  `json:"key"`
	Method    string  `json:"method"`
	Googlekey string  `json:"googlekey"`
	Pageurl   string  `json:"pageurl"`
	JSON      int     `json:"json"`
	Sitekey   string  `json:"sitekey,omitempty"`
	Invisible int     `json:"invisible,omitempty"`
	Version   string  `json:"version,omitempty"`
	Action    string  `json:"action,omitempty"`
	MinScore  float32 `json:"min_score,omitempty"`
	SoftID    int     `json:"soft_id,omitempty"`
	Body      string  `json:"body,omitempty"`      // Base64-encoded captcha image
	UserAgent string  `json:"userAgent,omitempty"` // userAgent that will be used to solve the captcha
	Proxy     string  `json:"proxy,omitempty"`     // Proxy to use to solve captchas from
	ProxyType string  `json:"proxytype,omitempty"` // Type of the proxy
}

// Type that will be used when getting a response from 2captcha
type twocaptchaResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

func (t Twocaptcha) GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return t.getCaptchaAnswer(context.Background(), additional...)
}

func (t Twocaptcha) GetTokenWithContext(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return t.getCaptchaAnswer(ctx, additional...)
}

func (t Twocaptcha) GetBalance() (float32, error) {
	return t.getBalance()
}

// Method to get Queue ID from the API.
func (t Twocaptcha) getID(data *AdditionalData) (string, error) {
	// Get Payload
	payload, err := t.createPayload(data)
	if err != nil {
		return "", err
	}

	// Make request to get answer
	for i := 0; i < 100; i++ {
		resp, err := http.Post("http://2captcha.com/in.php", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		response := &twocaptchaResponse{}
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Parse the response
		if response.Status != 1 { // Means there was an error
			return "", errCodeToError(response.Request)
		}
		return response.Request, nil
	}
	return "", ErrMaxAttempts
}

// This method gets the captcha token from the Capmonster API
func (t Twocaptcha) getCaptchaAnswer(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	var data *AdditionalData = nil
	if len(additional) > 0 {
		data = additional[0]
	}

	// Get Queue ID
	queueID, err := t.getID(data)
	if err != nil {
		return nil, err
	}

	// Get Captcha Answer
	response := &twocaptchaResponse{}
	urlToAnswer := fmt.Sprintf(
		"http://2captcha.com/res.php?key=%v&action=get&id=%v&json=1",
		t.config.Api_key,
		queueID,
	)
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequestWithContext(ctx, "GET", urlToAnswer, nil)
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
		if response.Status == 0 && response.Request != "CAPCHA_NOT_READY" {
			return nil, errCodeToError(response.Request)
		}

		// Check if captcha is ready
		if response.Request == "CAPCHA_NOT_READY" {
			time.Sleep(3 * time.Second)
			continue
		}
		return newCaptchaAnswer(
			queueID,
			response.Request,
			t.config.Api_key,
			t.config.CaptchaType,
			TwoCaptchaSite,
			"",
		), nil
	}
	return nil, ErrMaxAttempts
}

func (t Twocaptcha) getBalance() (float32, error) {
	// Attempt to get the balance from the API
	// Max attempts is 5
	url := fmt.Sprintf("https://2captcha.com/res.php?key=%v&action=getbalance&json=1", t.config.Api_key)
	response := &twocaptchaResponse{}
	for i := 0; i < 5; i++ {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)
		if response.Status == 0 {
			return 0, errCodeToError(response.Request)
		}

		// Convert to float32
		var balance float32
		value, err := strconv.ParseFloat(response.Request, 32)
		if err != nil {
			return 0, errors.New("unable to convert balance")
		}
		balance = float32(value)
		return balance, nil
	}
	return 0, ErrMaxAttempts
}

/*
createPayload returns the payloads required to interact with the API.

Possible errors that can be returned:
1) ErrIncorrectCapType
*/
func (t Twocaptcha) createPayload(data *AdditionalData) (string, error) {
	// Define the payload we are going to send to the API
	payload := twoCapIDPayload{
		Key:     t.config.Api_key,
		Pageurl: t.config.CaptchaURL,
		JSON:    1,
		Method:  "userrecaptcha",
	}

	// Add any other keys to the payload
	if t.config.SoftID != 0 {
		payload.SoftID = t.config.SoftID
	}
	switch t.config.CaptchaType {
	case ImageCaptcha:
		if data == nil {
			return "", ErrAddionalDataMissing
		}
		payload.Method = "base64"
		payload.Body = data.B64Img

	case V2Captcha:
		payload.Googlekey = t.config.Sitekey
		if t.config.IsInvisibleCaptcha {
			payload.Invisible = 1
		}
	case V3Captcha:
		payload.Googlekey = t.config.Sitekey
		payload.Version = "v3"
		payload.Action = t.config.Action
		payload.MinScore = t.config.MinScore
	case HCaptcha:
		payload.Method = "hcaptcha"
		payload.Sitekey = t.config.Sitekey
	case CFTurnstile:
		payload.Method = "turnstile"
		payload.Sitekey = t.config.Sitekey
	default:
		return "", ErrIncorrectCapType
	}

	// Check for any additional data about the task
	if data != nil && t.config.CaptchaType != ImageCaptcha {
		if data.UserAgent != "" {
			payload.UserAgent = data.UserAgent
		}
		if data.Proxy != nil {
			payload.Proxy = data.Proxy.StringFormatted()
		}
		if data.ProxyType != "" {
			payload.ProxyType = data.ProxyType
		}
	}

	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}

func report_2captcha(was_correct bool, c *CaptchaAnswer) error {
	var action string = "reportgood"
	if !was_correct {
		action = "reportbad"
	}

	url := fmt.Sprintf(
		"http://2captcha.com/res.php?key=%v&action=%v&id=%v&json=1",
		c.api_key,
		action,
		c.id,
	)
	response := &twocaptchaResponse{}

	// Make request
	for i := 0; i < 100; i++ {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Check for any errors
		if response.Status == 0 {
			return errCodeToError(response.Request)
		}
		if response.Status == 1 && response.Request != "OK_REPORT_RECORDED" {
			time.Sleep(3 * time.Second)
			continue
		}
		return nil
	}
	return ErrMaxAttempts
}
