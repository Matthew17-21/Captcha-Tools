package captchatoolsgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
   This file will contain the code to interact with capmonster.cloud API
*/

func (c Capmonster) GetToken() (string, error) {
	return c.getCaptchaAnswer()
}
func (c Capmonster) GetBalance() (float32, error) {
	return c.getBalance()
}

// Method to get Queue ID from the API.
func (c Capmonster) getID() (int, error) {
	// Get Payload
	payload, _ := c.createPayload()

	// Make request to get answer
	for {
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
}

// This method gets the captcha token from the Capmonster API
func (c Capmonster) getCaptchaAnswer() (string, error) {
	// Get Queue ID
	queueID, err := c.getID()
	if err != nil {
		return "", err
	}

	// Get Captcha Answer
	payload, _ := json.Marshal(capmonsterCapAnswerPayload{
		ClientKey: c.config.Api_key,
		TaskID:    queueID,
	})
	response := &capmonsterTokenResponse{}
	for {
		resp, err := http.Post("https://api.capmonster.cloud/getTaskResult", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)
		if response.Status == "ready" {
			return response.Solution.GRecaptchaResponse, nil
		} else if response.ErrorID == 12 || response.ErrorID == 16 { // Captcha unsolvable || TaskID doesn't exist
			c.GetToken()
		}
		time.Sleep(3 * time.Second)
	}
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
func (c Capmonster) createPayload() (string, error) {
	// Define the payload we are going to send to the API
	payload := capmonsterIDPayload{
		ClientKey: c.config.Api_key,
		Task: struct {
			WebsiteURL  string  "json:\"websiteURL\""
			WebsiteKey  string  "json:\"websiteKey\""
			Type        string  "json:\"type\""
			IsInvisible bool    "json:\"isInvisible,omitempty\""
			MinScore    float32 "json:\"minScore,omitempty\""
			PageAction  string  "json:\"pageAction,omitempty\""
		}{
			WebsiteURL: c.config.CaptchaURL,
			WebsiteKey: c.config.Sitekey,
			Type:       c.config.CaptchaType,
		},
	}

	// Add any other keys to the payload
	switch c.config.CaptchaType {
	case "v2":
		payload.Task.Type = "NoCaptchaTaskProxyless"
		if c.config.IsInvisibleCaptcha {
			payload.Task.IsInvisible = c.config.IsInvisibleCaptcha
		}
	case "v3":
		payload.Task.Type = "RecaptchaV3TaskProxyless"
		payload.Task.MinScore = c.config.MinScore
		payload.Task.PageAction = c.config.Action
	case "hcaptcha", "hcap":
		payload.Task.Type = "HCaptchaTaskProxyless"
	default:
		return "", ErrIncorrectCapType
	}
	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}
