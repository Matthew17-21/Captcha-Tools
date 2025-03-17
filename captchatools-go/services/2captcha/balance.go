package twocaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/errors"
)

// Attempt to get the balance from the API
func (t Twocaptcha) GetBalance() (float32, error) {
	return t.getBalance(baseURL)
}

func (t Twocaptcha) getBalance(baseUrl string) (float32, error) {
	t.Logger.Info("Attempting to get balance for key %q...", t.Config.Api_key)

	// Create payload
	p := map[string]string{"clientKey": t.Config.Api_key}
	payload, err := json.Marshal(p)
	if err != nil {
		return 0, fmt.Errorf("error creating payload: %w", err)
	}
	t.Logger.Debug("Payload to get balance: %q", string(payload))

	// Create & make request
	resp, err := http.Post(
		baseUrl+getBalanceEp,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}

	// Parse response
	return parseGetBalance(resp)
}

// newGetBalanceReq creates a new HTTP request to get the account balance
func newGetBalanceReq(baseUrl, apiKey string) (*http.Request, error) {
	// Create payload
	p := map[string]string{"clientKey": apiKey}
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error creating payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest(
		http.MethodPost,
		baseUrl+getBalanceEp,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func parseGetBalance(r *http.Response) (float32, error) {
	defer r.Body.Close()

	// Read the body
	resp := &struct {
		ErrorID int     `json:"errorId"`
		Balance float32 `json:"balance"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return 0, fmt.Errorf("unable to read res body: %s", err)
		}
		return 0, fmt.Errorf("unable to decode body (%s) into json: %v", string(b), err)
	}

	// Validate & return
	if resp.ErrorID != 0 {
		return 0, errors.ErrCodeToError(strconv.Itoa(resp.ErrorID))
	}
	return resp.Balance, nil
}
