package captchatoolsgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
   This file will contain the code to interact with the Solvecha API (https://solvecha.com/docs).
   Solvecha solves image captchas only. Token captchas return ErrNotSupported.
*/

// Base URL of the Solvecha API (a variable so tests can point it at a local server).
var solvechaBaseURL = "https://api.solvecha.com"

type Solvecha struct {
	*Config
}

type solvechaResponse struct {
	ErrorID          int     `json:"errorId"`
	ErrorCode        string  `json:"errorCode"`
	ErrorDescription string  `json:"errorDescription"`
	TaskID           string  `json:"taskId"`
	Status           string  `json:"status"`
	Balance          float32 `json:"balance"`
	Solution         struct {
		Text string `json:"text"`
	} `json:"solution"`
}

// solvechaError converts an API error to an error of this package.
// Unknown codes keep the message of the API.
func solvechaError(r *solvechaResponse) error {
	err := errCodeToError(r.ErrorCode)
	if err == ErrUnknown {
		if r.ErrorDescription != "" {
			return fmt.Errorf("%w: %s", ErrUnknown, r.ErrorDescription)
		}
		return fmt.Errorf("%w: %s", ErrUnknown, r.ErrorCode)
	}
	return err
}

// post sends a JSON request to the Solvecha API. Errors come back as
// 4xx with a JSON body, so the body is parsed whatever the status is.
func (c Solvecha) post(ctx context.Context, path string, payload interface{}) (*solvechaResponse, error) {
	encoded, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", solvechaBaseURL+path, bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := makeRequest(req)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	response := &solvechaResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("solvecha: unexpected response (HTTP %d)", resp.StatusCode)
	}
	if response.ErrorID != 0 {
		return nil, solvechaError(response)
	}
	return response, nil
}

// GetBalance returns the prepaid balance only.
func (c Solvecha) GetBalance() (float32, error) {
	var lastErr error = ErrMaxAttempts
	for i := 0; i < 5; i++ {
		r, err := c.post(context.Background(), "/getBalance", map[string]string{"clientKey": c.Api_key})
		if err != nil {
			if isSolvechaAPIError(err) {
				return 0, err
			}
			lastErr = err
			time.Sleep(1 * time.Second)
			continue
		}
		return r.Balance, nil
	}
	return 0, lastErr
}

func (c Solvecha) GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return c.getCaptchaAnswer(context.Background(), additional...)
}

func (c Solvecha) GetTokenWithContext(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	return c.getCaptchaAnswer(ctx, additional...)
}

// isSolvechaAPIError reports whether err came from an API error body (no retry needed).
func isSolvechaAPIError(err error) bool {
	return errors.Is(err, ErrUnknown) || errors.Is(err, ErrWrongAPIKey) || errors.Is(err, ErrNoBalance) ||
		errors.Is(err, ErrBanned) || errors.Is(err, ErrNoSlot) || errors.Is(err, ErrUnsolvable) ||
		errors.Is(err, ErrNoCaptchaID) || errors.Is(err, ErrCaptchaData)
}

func (c Solvecha) getCaptchaAnswer(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) {
	// Image captchas only
	if c.CaptchaType != ImageCaptcha {
		return nil, ErrNotSupported
	}
	var data *AdditionalData
	if len(additional) > 0 {
		data = additional[0]
	}
	body := ""
	if data != nil {
		body = data.B64Img
	}

	created, err := c.post(ctx, "/createTask", map[string]interface{}{
		"clientKey": c.Api_key,
		"task":      map[string]string{"type": "ImageToTextTask", "body": body},
	})
	if err != nil {
		return nil, err
	}

	// The answer can already be in the createTask response.
	if created.Status == "ready" {
		return newCaptchaAnswer(created.TaskID, created.Solution.Text, c.Api_key, c.CaptchaType, SolvechaSite, ""), nil
	}

	for i := 0; i < 50; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("getCaptchaAnswer error: %w", ctx.Err())
		case <-time.After(solvechaPollInterval):
		}
		r, err := c.post(ctx, "/getTaskResult", map[string]string{"clientKey": c.Api_key, "taskId": created.TaskID})
		if err != nil {
			if isSolvechaAPIError(err) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			continue
		}
		if r.Status == "ready" {
			return newCaptchaAnswer(created.TaskID, r.Solution.Text, c.Api_key, c.CaptchaType, SolvechaSite, ""), nil
		}
	}
	return nil, ErrMaxAttempts
}

var solvechaPollInterval = 3 * time.Second
