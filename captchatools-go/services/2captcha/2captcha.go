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

