package captchatoolsgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Capsolver struct {
	*Config
}

func (a Capsolver) GetBalance() (float32, error) {
	return a.getBalance()
}

func (c Capsolver) getBalance() (float32, error) {
	// Create payload
	payload := fmt.Sprintf(`{ "clientKey": "%v" }`, c.Api_key)

	// Make POST Request to API
	response := struct {
		ErrorID   int     `json:"errorId"`
		ErrorCode string  `json:"errorCode"`
		Balance   float32 `json:"balance"`
	}{}
	for i := 0; i < 5; i++ {
		resp, err := http.Post(
			"https://api.capsolver.com/getBalance",
			"application/json",
			bytes.NewBufferString(payload),
		)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &response)
		if response.ErrorID != 0 {
			return 0, errCodeToError(response.ErrorCode)
		}
		return response.Balance, nil
	}
	return 0, ErrMaxAttempts
}

func (c Capsolver) GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return c.getCaptchaAnswer(additional...)
}

// Method to get Queue ID from the API.
func (c Capsolver) getID(data *AdditionalData) (string, error) {
	// Get Payload
	payload, err := c.createPayload(data)
	if err != nil {
		return "", err
	}

	// Make request to get answer
	response := &struct {
		ErrorID   int    `json:"errorId"`
		ErrorCode string `json:"errorCode"`
		TaskID    string `json:"taskId"`
	}{}
	for i := 0; i < 50; i++ {
		resp, err := http.Post("https://api.capsolver.com/createTask", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Parse the response
		if response.ErrorID != 0 { // Means there was an error
			return "", errCodeToError(response.ErrorCode)
		}
		return response.TaskID, nil
	}
	return "", ErrMaxAttempts
}

func (c Capsolver) getCaptchaAnswer(additional ...*AdditionalData) (*CaptchaAnswer, error) {
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
	type Payload struct {
		ClientKey string `json:"clientKey"`
		TaskID    string `json:"taskId"`
	}
	payload, _ := json.Marshal(Payload{
		ClientKey: c.Api_key,
		TaskID:    queueID,
	})
	response := &capmonsterTokenResponse{}
	for i := 0; i < 50; i++ {
		resp, err := http.Post("https://api.capsolver.com/getTaskResult", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
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
		switch c.CaptchaType {
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
			c.Api_key,
			c.CaptchaType,
			AnticaptchaSite, // TODO change this
			ua,
		), nil
	}
	return nil, ErrMaxAttempts
}

func (c Capsolver) createPayload(data *AdditionalData) (string, error) {
	type Task struct {
		Type       captchaType `json:"type"`
		WebsiteURL string      `json:"websiteURL"`
		WebsiteKey string      `json:"websiteKey"`

		// Recaptcha V2 Data
		IsInvisible bool   `json:"isInvisible"`
		UserAgent   string `json:"userAgent"`
		Proxy       string `json:"proxy"`
	}
	type Payload struct {
		ClientKey string `json:"clientKey"`
		Task      Task   `json:"task"`
	}

	p := Payload{
		ClientKey: c.Api_key,
		Task: Task{
			Type:       c.CaptchaType,
			WebsiteURL: c.CaptchaURL,
			WebsiteKey: c.Sitekey,
		},
	}

	switch c.CaptchaType {
	case V2Captcha:
		p.Task.Type = "ReCaptchaV2TaskProxyLess"
		p.Task.IsInvisible = c.IsInvisibleCaptcha
		if data != nil && data.Proxy != nil {
			p.Task.Type = "ReCaptchaV2Task"
		}
	}

	// Check for any additional data about the task
	if data != nil && c.CaptchaType != ImageCaptcha {
		if data.UserAgent != "" {
			p.Task.UserAgent = data.UserAgent
		}
		if data.Proxy != nil {
			p.Task.Proxy = data.Proxy.StringFormatted()
		}
	}

	// TODO add softkey id

	// Return
	encoded, _ := json.Marshal(p)
	return string(encoded), nil
}
