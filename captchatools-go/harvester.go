package captchatoolsgo

import "strings"

/*
	NewHarvester returns a captcha harvester based on the info given
	by the caller. An error is returned if there is no proper
	solving_site argument.

	To make the implementation similiar to the Python version,
	this function was needed.

	For documentation on how to use this, checkout
	https://github.com/Matthew17-21/Captcha-Tools
*/
func NewHarvester(solving_site int, config *Config) (*Harvester, error) {
	h := &Harvester{}
	config.CaptchaType = strings.ToLower(config.CaptchaType)
	config.CaptchaURL = strings.ToLower(config.CaptchaURL)

	// Check for any errors
	switch strings.ToLower(config.CaptchaType) {
	case "hcaptcha", "hcap", "v2", "v3":
	default:
		return nil, ErrIncorrectCapType
	}

	// Get A Harvester
	switch solving_site {
	case AnticaptchaSite:
		h.childHarvester = &Anticaptcha{config: config}
	case CapmonsterSite:
		h.childHarvester = &Capmonster{config: config}
	case TwoCaptchaSite:
		h.childHarvester = &Twocaptcha{config: config}
	default:
		return nil, ErrNoHarvester
	}
	return h, nil
}

// GetToken returns a captcha token from the selected solving site
func (h *Harvester) GetToken() (string, error) {
	return h.childHarvester.GetToken()
}
