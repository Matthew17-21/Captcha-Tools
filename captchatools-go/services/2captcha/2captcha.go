package twocaptcha

import (
	"context"
	"time"

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

// GetBalance returns the current account balance
func (t Twocaptcha) GetBalance() (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return t.getBalance(ctx, baseURL)
}

// GetBalanceWithContext returns the account balance with a custom context
func (t Twocaptcha) GetBalanceWithContext(ctx context.Context) (float32, error) {
	return t.getBalance(ctx, baseURL)
}
