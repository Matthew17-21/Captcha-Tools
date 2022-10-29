package captchatoolsgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// This file will contain the code to interact with anticaptcha.com API

func (a Anticaptcha) GetToken() (*CaptchaAnswer, error) {
	return a.getCaptchaAnswer()
}
func (a Anticaptcha) GetBalance() (float32, error) {
	return a.getBalance()
}

// Method to get Queue ID from the API.
func (a Anticaptcha) getID() (int, error) {
	// Get Payload
	payload, err := a.createPayload()
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
func (a Anticaptcha) getCaptchaAnswer() (*CaptchaAnswer, error) {
	// Get Queue ID
	queueID, err := a.getID()
	if err != nil {
		return nil, err
	}

	// Get Captcha Answer
	payload, _ := json.Marshal(capmonsterCapAnswerPayload{
		ClientKey: a.config.Api_key,
		TaskID:    queueID,
	})
	response := &capmonsterTokenResponse{}
	for i := 0; i < 100; i++ {
		resp, err := http.Post("https://api.anti-captcha.com/getTaskResult", "application/json", bytes.NewBuffer([]byte(payload)))
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
			return nil, errCodeToError(response.ErrorCode)
		}

		// Check if the answer is ready or not
		if response.Status == "processing" {
			time.Sleep(3 * time.Second)
			continue
		}

		return newCaptchaAnswer(
			queueID,
			response.Solution.GRecaptchaResponse,
			a.config.Api_key,
			AnticaptchaSite,
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
func (a Anticaptcha) createPayload() (string, error) {
	// Define the payload we are going to send to the API
	payload := capmonsterIDPayload{
		ClientKey: a.config.Api_key,
		Task: struct {
			WebsiteURL  string  "json:\"websiteURL\""
			WebsiteKey  string  "json:\"websiteKey\""
			Type        string  "json:\"type\""
			IsInvisible bool    "json:\"isInvisible,omitempty\""
			MinScore    float32 "json:\"minScore,omitempty\""
			PageAction  string  "json:\"pageAction,omitempty\""
		}{
			WebsiteURL: a.config.CaptchaURL,
			WebsiteKey: a.config.Sitekey,
			Type:       a.config.CaptchaType,
		},
	}

	// Add any other keys to the payload
	switch a.config.CaptchaType {
	case "v2":
		payload.Task.Type = "NoCaptchaTaskProxyless"
		if a.config.IsInvisibleCaptcha {
			payload.Task.IsInvisible = a.config.IsInvisibleCaptcha
		}
	case "v3":
		payload.Task.Type = "RecaptchaV3TaskProxyless"
		payload.Task.MinScore = a.config.MinScore
		payload.Task.PageAction = a.config.Action
	case "hcaptcha", "hcap":
		payload.Task.Type = "HCaptchaTaskProxyless"
	default:
		return "", ErrIncorrectCapType
	}
	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}
