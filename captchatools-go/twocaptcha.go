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

func (t *Twocaptcha) GetToken() (string, error) {
	return t.getCaptchaAnswer()
}

// Method to get Queue ID from the API.
func (t *Twocaptcha) getID() (string, error) {
	// Get Payload
	payload, _ := t.createPayload()

	// Make request to get answer
	for {
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
		if response.Status == 1 { // Means there was no error
			return response.Request, nil
		}
		switch response.Request {
		case "ERROR_ZERO_BALANCE":
			return "", ErrNoBalance
		case "ERROR_WRONG_GOOGLEKEY":
			return "", ErrWrongSitekey
		case "ERROR_WRONG_USER_KEY", "ERROR_KEY_DOES_NOT_EXIST":
			return "", ErrWrongAPIKey
		}

	}
}

// This method gets the captcha token from the Capmonster API
func (t *Twocaptcha) getCaptchaAnswer() (string, error) {
	// Get Queue ID
	queueID, err := t.getID()
	if err != nil {
		return "", err
	}

	// Get Captcha Answer
	response := &twocaptchaResponse{}
	urlToAnswer := fmt.Sprintf(
		"http://2captcha.com/res.php?key=%v&action=get&id=%v&json=1",
		t.config.Api_key,
		queueID,
	)
	for {
		resp, err := http.Get(urlToAnswer)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// Parse Response
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, response)
		if response.Status == 1 {
			return response.Request, nil
		} else if response.Request == "ERROR_CAPTCHA_UNSOLVABLE" {
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
func (t *Twocaptcha) createPayload() (string, error) {
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
	case "v2":
		payload.Googlekey = t.config.Sitekey
		if t.config.IsInvisibleCaptcha {
			payload.Invisible = 1
		}
	case "v3":
		payload.Googlekey = t.config.Sitekey
		payload.Version = "v3"
		payload.Action = t.config.Action
		payload.MinScore = t.config.MinScore
	case "hcaptcha", "hcap":
		payload.Method = "hcaptcha"
		payload.Sitekey = t.config.Sitekey
	default:
		return "", ErrIncorrectCapType
	}
	encoded, _ := json.Marshal(payload)
	return string(encoded), nil
}
