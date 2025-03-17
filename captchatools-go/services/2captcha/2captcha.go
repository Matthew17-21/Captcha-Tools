package twocaptcha

import (
	"context"
	"fmt"

	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
)

const (
	baseURL      string = "https://api.2captcha.com"
	getBalanceEp string = "/getBalance"
)

type Twocaptcha struct {
	harvester.Config
}

// Attempt to get the balance from the API
func (t Twocaptcha) GetBalance() (float32, error) {
	return t.getBalance(context.Background(), baseURL)
}
