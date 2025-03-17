package twocaptcha

import (
	"encoding/json"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
)

// CaptchaAnswer implementation
type captchaAnswer struct {
	id          int64
	token       string
	apiKey      string
	captchaType string
	provider    string
}

// createCaptchaAnswer creates a CaptchaAnswer from the 2Captcha response
func (t Twocaptcha) createCaptchaAnswer(taskID int64, result *taskResponse) harvester.CaptchaAnswer {
	// Extract token from solution (depends on captcha type)
	var token string

	if result.Solution != nil {
		// Check common response formats based on different captcha types
		if gRecaptchaResponse, ok := result.Solution["gRecaptchaResponse"].(string); ok {
			token = gRecaptchaResponse
		} else if text, ok := result.Solution["text"].(string); ok {
			token = text
		} else if token, ok = result.Solution["token"].(string); !ok {
			// If none of the above, try to marshal the entire solution
			solutionJson, _ := json.Marshal(result.Solution)
			token = string(solutionJson)
		}
	}

	// Create captcha answer
	answer := &captchaAnswer{
		id:       taskID,
		token:    token,
		apiKey:   t.Config.Api_key,
		provider: "2captcha",
	}

	return answer
}

func (c *captchaAnswer) ID() any {
	return c.id
}

func (c *captchaAnswer) Token() string {
	return c.token
}

// TODO: Add UserAgent implementation
func (c *captchaAnswer) UserAgent() string {
	return ""
}

// TODO: Add report implementation
func (c *captchaAnswer) Report(wasCorrect bool) error {
	// Implementation of captcha reporting
	// For now, return nil as implementation detail
	return nil
}
