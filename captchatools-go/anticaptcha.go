package captchatoolsgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// This file will contain the code to interact with anticaptcha.com API

func (t *Anticaptcha) GetToken() (string, error) {
	return t.getCaptchaAnswer()
}

// Method to get Queue ID from the API.
func (t *Anticaptcha) getID() (int, error) {
	// Get Payload
	payload, _ := t.createPayload()

	// Make request to get answer
	response := &capmonsterIDResponse{}
	for {
		resp, err := http.Post("https://api.anti-captcha.com/createTask", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)

		// Parse the response
		if response.ErrorID == 0 { // Means there was no error
			return response.TaskID, nil
		}
		switch response.ErrorCode {
		case "ERROR_ZERO_BALANCE":
			return 0, ErrNoBalance
		case "ERROR_RECAPTCHA_INVALID_SITEKEY":
			return 0, ErrWrongSitekey
		case "ERROR_KEY_DOES_NOT_EXIST":
			return 0, ErrWrongAPIKey
		}

	}
}

// This method gets the captcha token from the Capmonster API
func (t *Anticaptcha) getCaptchaAnswer() (string, error) {
	// Get Queue ID
	queueID, err := t.getID()
	if err != nil {
		return "", err
	}

	// Get Captcha Answer
	payload, _ := json.Marshal(capmonsterCapAnswerPayload{
		ClientKey: t.config.Api_key,
		TaskID:    queueID,
	})
	response := &capmonsterTokenResponse{}
	for {
		resp, err := http.Post("https://api.anti-captcha.com/getTaskResult", "application/json", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// Parse Response
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)
		if response.Status == "ready" {
			return response.Solution.GRecaptchaResponse, nil
		} else if response.ErrorID == 12 || response.ErrorID == 16 { // Captcha unsolvable || TaskID doesn't exist
			t.GetToken()
		}
		time.Sleep(3 * time.Second)
	}
}

/*
	createPayload returns the payloads required to interact with the API.

	Possible errors that can be returned:
	1) ErrIncorrectCapType
*/
func (t *Anticaptcha) createPayload() (string, error) {
	// Define the payload we are going to send to the API
	payload := capmonsterIDPayload{
		ClientKey: t.config.Api_key,
		Task: struct {
			WebsiteURL  string  "json:\"websiteURL\""
			WebsiteKey  string  "json:\"websiteKey\""
			Type        string  "json:\"type\""
			IsInvisible bool    "json:\"isInvisible,omitempty\""
			MinScore    float32 "json:\"minScore,omitempty\""
			PageAction  string  "json:\"pageAction,omitempty\""
		}{
			WebsiteURL: t.config.CaptchaURL,
			WebsiteKey: t.config.Sitekey,
			Type:       t.config.CaptchaType,
		},
	}

	// Add any other keys to the payload
	switch t.config.CaptchaType {
	case "v2":
		payload.Task.Type = "NoCaptchaTaskProxyless"
		if t.config.IsInvisibleCaptcha {
			payload.Task.IsInvisible = t.config.IsInvisibleCaptcha
		}
	case "v3":
		payload.Task.Type = "RecaptchaV3TaskProxyless"
		payload.Task.MinScore = t.config.MinScore
		payload.Task.PageAction = t.config.Action
	case "hcaptcha", "hcap":
		payload.Task.Type = "HCaptchaTaskProxyless"
	default:
		return "", ErrIncorrectCapType
	}
	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}
