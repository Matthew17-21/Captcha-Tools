package twocaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/errors"
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/internal/httputils"
)

const (
	getTaskIDEp    string = "/getTaskID"
	getResultEp    string = "/getTaskResult"
	maxRetries     int    = 30
	pollingTimeout        = 3 * time.Second

	defaultErrorId int64 = -1
)


// getTaskID sends a request to create a new captcha solving task
func (t Twocaptcha) getTaskID(ctx context.Context, baseUrl string, p payload) (int64, error) {
	t.Logger.Info("Creating task for captcha solving...")

	// Create request
	req, err := t.newGetTaskIDReq(ctx, baseUrl, p)
	if err != nil {
		return defaultErrorId, fmt.Errorf("newGetTaskIDReq error: %w", err)
	}

	// Make request
	resp, err := httputils.MakeRequest(ctx, t.Logger, req, maxRetries)
	if err != nil {
		return defaultErrorId, fmt.Errorf("error making getTaskID request: %w", err)
	}

	// Parse response
	return t.parseGetTaskId(resp)
}

func (t Twocaptcha) newGetTaskIDReq(ctx context.Context, baseUrl string, p payload) (*http.Request, error) {
	// Create payload
	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error marshaling task payload: %w", err)
	}
	t.Logger.Debug("Payload for create task: %s", string(payloadBytes))

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseUrl+getTaskIDEp,
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t Twocaptcha) parseGetTaskId(r *http.Response) (int64, error) {
	defer r.Body.Close()

	// Read the body
	resp := struct {
		ErrorID int   `json:"errorId"`
		TaskID  int64 `json:"taskId"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return defaultErrorId, fmt.Errorf("unable to read res body: %s", err)
		}
		return defaultErrorId, fmt.Errorf("unable to decode body (%s) into json: %v", string(b), err)
	}

	// Validate
	if resp.ErrorID != 0 {
		return defaultErrorId, fmt.Errorf("error getting task ID: %w", errors.ErrCodeToError(strconv.Itoa(resp.ErrorID)))
	}
	return resp.TaskID, nil
}

