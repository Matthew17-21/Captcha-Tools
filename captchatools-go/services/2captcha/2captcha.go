package twocaptcha

import (
	"context"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
)

const (
	baseURL      string = "https://api.2captcha.com"
	getBalanceEp string = "/getBalance"
)

type Twocaptcha struct {
	harvester.Config
}

func New(c harvester.Config) *Twocaptcha {
	return &Twocaptcha{c}
}

// GetToken obtains a captcha token from 2Captcha with default context
func (t Twocaptcha) GetToken(opts ...harvester.TokenOption) (harvester.CaptchaAnswer, error) {
	return t.GetTokenWithContext(context.Background(), opts...)
}

// GetTokenWithContext obtains a captcha token from 2Captcha with a custom context
func (t Twocaptcha) GetTokenWithContext(ctx context.Context, opts ...harvester.TokenOption) (harvester.CaptchaAnswer, error) {
	return t.getToken(ctx, baseURL, defaultPollingTimeout, opts...)
}

// Attempt to get the balance from the API
func (t Twocaptcha) GetBalance() (float32, error) {
	return t.getBalance(context.Background(), baseURL)
}
