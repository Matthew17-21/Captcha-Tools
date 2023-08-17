package captchatoolsgo

import (
	"context"
	"net/http"
)

/*
NewHarvester returns a captcha harvester based on the info given
by the caller. An error is returned if there is no proper
solving_site argument.

To make the implementation similiar to the Python version,
this function was needed.

For documentation on how to use this, checkout
https://github.com/Matthew17-21/Captcha-Tools
*/

/*
	- type Harvester will be used to represent a captcha harvester.

	- In order to to have the same functionality/implementation as the
	Python version, type Harvester has the `childHarvester` field, which is of type interface.
	This allows us to set the field as a pointer to a real captcha harvester
	and use the `GetToken` method that each captcha harvester struct has.
	(Polymorphism)

	- For documentation, visit https://github.com/Matthew17-21/Captcha-Tools
*/
// Interface that will allow us to interact with the methods from the
// individual structs
type Harvester interface {
	GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error)                                 // Function to get a captcha token
	GetTokenWithContext(ctx context.Context, additional ...*AdditionalData) (*CaptchaAnswer, error) // Function to get a captcha token
	GetBalance() (float32, error)
}

type AdditionalData struct {
	B64Img    string // base64 encoded image
	Proxy     *Proxy // A proxy in correct formatting - such as user:pass@ip:port
	ProxyType string // Type of your proxy: HTTP, HTTPS, SOCKS4, SOCKS5
	UserAgent string // UserAgent that will be passed to the service and used to solve the captcha
}

// Configurations for the captchas you are solving.
// For more a detailed documentation, visit
// https://github.com/Matthew17-21/Captcha-Tools
type Config struct {
	Api_key            string      // The API Key for the captcha solving site.
	Sitekey            string      // Sitekey from the site where captcha is loaded.
	CaptchaURL         string      // URL where the captcha is located.
	CaptchaType        captchaType // Type of captcha you are solving. Visit https://github.com/Matthew17-21/Captcha-Tools for types
	Action             string      // Action that is associated with the V3 captcha.
	IsInvisibleCaptcha bool        // If the captcha is invisible or not.
	MinScore           float32     // Minimum score for v3 captchas.
	SoftID             int         // SoftID for 2captcha. Developers get reward 10% of spendings of their software users.
}

func NewHarvester(solving_site site, config *Config) (Harvester, error) {
	// Check for any errors
	switch config.CaptchaType {
	case ImageCaptcha, V2Captcha, V3Captcha, HCaptcha, CFTurnstile:
	default:
		return nil, ErrIncorrectCapType
	}

	// Get A Harvester
	var h Harvester
	switch solving_site {
	case AnticaptchaSite:
		h = &Anticaptcha{config: config}
	case CapmonsterSite:
		h = &Capmonster{config: config}
	case TwoCaptchaSite:
		h = &Twocaptcha{config: config}
	case CapsolverSite:
		h = &Capsolver{config}
	case CaptchaAiSite:
		h = &CaptchaAi{config}
	default:
		return nil, ErrNoHarvester
	}
	return h, nil
}

// makeRequest is a wrapper function to httpClient.Do
func makeRequest(req *http.Request) (*http.Response, error) {
	c := http.Client{}
	defer c.CloseIdleConnections()
	return c.Do(req)
}
