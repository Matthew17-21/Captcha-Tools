package captchatoolsgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CaptchaAi struct {
	*Config
}

func (t CaptchaAi) GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return t.getCaptchaAnswer(context.Background(), additional...)
}

func (t CaptchaAi) GetTokenWithContext(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return t.getCaptchaAnswer(ctx, additional...)
}

func (t CaptchaAi) GetBalance() (float32, error) {
	return t.getBalance()
}

// Method to get Queue ID from the API.
func (t CaptchaAi) getID(data *AdditionalData) (string, error) {
	// Get Payload
	uri, err := t.createUrl(data)
	if err != nil {
		return "", err
	}

	// Make request to get answer
	response := &struct {
		Status  int `json:"status"`
		Request int `json:"request"`
	}{}
	for i := 0; i < 100; i++ {
		resp, err := http.Get(uri)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Parse the response
		if response.Status != 1 { // Means there was an error
			// Have to read the error into an interface
			temp := make(map[string]string)
			json.Unmarshal(body, &temp)
			return "", errCodeToError(temp["request"])
		}
		return strconv.Itoa(response.Request), nil
	}
	return "", ErrMaxAttempts
}

// This method gets the captcha token from the Capmonster API
func (t CaptchaAi) getCaptchaAnswer(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
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
		"https://ocr.captchaai.com/res.php?key=%v&action=get&id=%v&json=1",
		t.Api_key,
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
			t.Api_key,
			t.CaptchaType,
			CaptchaAiSite,
			"",
		), nil
	}
	return nil, ErrMaxAttempts
}

func (t CaptchaAi) getBalance() (float32, error) {
	// Attempt to get the balance from the API
	// Max attempts is 5
	url := fmt.Sprintf("https://ocr.captchaai.com/res.php?key=%v&action=getbalance&json=1", t.Api_key)
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
createUrl creates the Uri needed to submit data to CaptchaAi

Possible errors that can be returned:
1) ErrIncorrectCapType
*/
func (t CaptchaAi) createUrl(data *AdditionalData) (string, error) {

	// Create base uri
	u, err := url.Parse("https://ocr.captchaai.com/in.php")
	if err != nil {
		return "", fmt.Errorf("createUrl error: %w", err)
	}

	// Dynamically add queries
	query := u.Query()
	query.Add("key", t.Api_key)
	query.Add("json", "1")
	query.Add("pageurl", t.CaptchaURL)
	switch t.CaptchaType {
	case ImageCaptcha:
		query.Add("method", "base64")
		if data != nil && data.B64Img != "" {
			query.Add("body", data.B64Img)
		}
	case V2Captcha:
		query.Add("method", "userrecaptcha")
		query.Add("googlekey", t.Sitekey)
		if t.IsInvisibleCaptcha {
			query.Add("invisible", "1")
		}
	case V3Captcha:
		query.Add("method", "userrecaptcha")
		query.Add("version", "v3")
		query.Add("googlekey", t.Sitekey)
		if t.Action != "" {
			query.Add("action", t.Action)
		}
		if t.MinScore > 0 {
			query.Add("min_score", fmt.Sprintf("%v", t.MinScore))
		}
	case HCaptcha:
		query.Add("method", "hcaptcha")
		query.Add("sitekey", t.Sitekey)

	case CFTurnstile:
		return "", ErrNotSupported
	default:
		return "", ErrIncorrectCapType
	}
	if data != nil && t.CaptchaType != ImageCaptcha {
		if data.UserAgent != "" {
			query.Add("userAgent", data.UserAgent)
		}
		if data.Proxy != nil {
			query.Add("proxy", data.Proxy.StringFormatted())
		}
		if data.ProxyType != "" {
			query.Add("proxytype", data.ProxyType)
		}
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}
