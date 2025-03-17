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

	errs "errors"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/errors"
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/internal/httputils"
)

const (
	createTaskEp string = "/createTask"
	getResultEp  string = "/getTaskResult"
	maxRetries   int    = 30

	defaultPollingTimeout       = 3 * time.Second
	defaultErrorId        int64 = -1
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
		baseUrl+createTaskEp,
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

// Response type for 2Captcha getTaskID and getTaskResult endpoints
type taskResponse struct {
	ErrorID          int                    `json:"errorId"`
	TaskId           string                 `json:"taskId,omitempty"`
	Status           string                 `json:"status,omitempty"`
	ErrorCode        string                 `json:"errorCode,omitempty"`
	Solution         map[string]interface{} `json:"solution,omitempty"`
	ErrorDescription string                 `json:"errorDescription,omitempty"`
}

// getToken is the internal implementation that handles the API communication
func (t Twocaptcha) getToken(ctx context.Context, baseUrl string, timeout time.Duration, opts ...harvester.TokenOption) (harvester.CaptchaAnswer, error) {
	// Create default payload
	payload := newTaskPayload(t)

	// Apply user options
	for _, opt := range opts {
		opt(&payload)
	}

	// Get task ID
	taskID, err := t.getTaskID(ctx, baseUrl, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Get task result
	return t.getTaskResult(ctx, baseUrl, taskID, timeout)
}

// getTaskResult polls 2Captcha for the result of a captcha solving task
func (t Twocaptcha) getTaskResult(ctx context.Context, baseUrl string, taskID int64, pollingTimeout time.Duration) (harvester.CaptchaAnswer, error) {
	t.Logger.Info("Getting result for task %d...", taskID)

	// Create payload
	payload := newPayloadWithClientKey(t.Config.Api_key)
	payload.setToPayload("taskId", taskID)
	t.Logger.Debug("Payload to get task result: %v", payload)

	// Poll for result with retries
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default: // So it does not block
		}

		// Create request
		req, err := t.newTaskResultReq(ctx, baseUrl, payload)
		if err != nil {
			t.Logger.Error("Error creating task result request: %v", err)
			continue
		}

		// Make request
		resp, err := httputils.MakeRequest(ctx, t.Logger, req, 1) // We use 1 here because we are already retrying with `maxRetries`
		if err != nil {
			t.Logger.Error("Error making task result request: %v", err)
			continue
		}

		// Parse response
		result, err := t.parseGetTaskResult(resp, taskID)
		if err != nil {
			// Check if it's a "not ready" error, and continue polling if so
			if err.Error() == "TASK_NOT_READY" {
				t.Logger.Debug("Task not ready, retrying in %s seconds...", pollingTimeout)
				time.Sleep(pollingTimeout)
				continue
			}

			// For any other error, return it
			return nil, fmt.Errorf("parseGetTaskResult error: %w", err)
		}

		// If we got here, we have a valid result
		return result, nil
	}

	return nil, errors.ErrMaxAttempts
}

func (t Twocaptcha) newTaskResultReq(ctx context.Context, baseUrl string, p payload) (*http.Request, error) {
	// Create payload
	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error marshaling task payload: %w", err)
	}
	t.Logger.Debug("Payload for task result: %s", string(payloadBytes))

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseUrl+getResultEp,
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (t Twocaptcha) parseGetTaskResult(r *http.Response, taskID int64) (harvester.CaptchaAnswer, error) {
	defer r.Body.Close()

	// Parse the JSON response
	result := &taskResponse{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read res body: %s", err)
		}
		return nil, fmt.Errorf("unable to decode body (%s) into json: %v", string(b), err)
	}

	// Check for errors
	if result.ErrorID != 0 {
		// Return the error code as the error
		return nil, errors.ErrCodeToError(result.ErrorCode)
	}

	// Check if the task is ready
	if result.Status != "ready" {
		// If the task is not ready, return a specific error
		return nil, errs.New("TASK_NOT_READY")
	}

	// If we got a successful result and the task is ready, create and return the CaptchaAnswer
	return t.createCaptchaAnswer(taskID, result), nil
}
