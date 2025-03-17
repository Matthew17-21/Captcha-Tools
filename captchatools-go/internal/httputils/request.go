package httputils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/logger"
)

func MakeRequest(ctx context.Context, logger logger.Logger, req *http.Request, retries int) (*http.Response, error) {
	client := http.Client{}
	for i := 0; i < retries; i++ {
		// Attempt to make request
		resp, err := client.Do(req)
		if err != nil {
			// Check if context is alive
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context is done: %w", ctx.Err())
			default: // So it doesn't block
			}

			// handle error, if any
			networkErr := ClassifyNetworkError(err)
			switch networkErr {
			case BadURL:
				return nil, fmt.Errorf("URL %q is malformed", req.URL.String())
			default:
			}

			// After parsing the error, retry the request
			logger.Debug("error while making request: %q. Retrying...\n", networkErr)
			continue
		}
		return resp, nil
	}

	// If unable to make the request in x tries, return an error
	return nil, fmt.Errorf("failed to make request in %d tries", retries)
}
