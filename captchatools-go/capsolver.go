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
