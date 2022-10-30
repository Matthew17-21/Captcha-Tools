package captchatoolsgo

/*
NewHarvester returns a captcha harvester based on the info given
by the caller. An error is returned if there is no proper
solving_site argument.

To make the implementation similiar to the Python version,
this function was needed.

For documentation on how to use this, checkout
https://github.com/Matthew17-21/Captcha-Tools
*/
func NewHarvester(solving_site site, config *Config) (Harvester, error) {
	// Check for any errors
	switch config.CaptchaType {
	case "hcaptcha", "hcap", "v2", "v3":
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
	default:
		return nil, ErrNoHarvester
	}
	return h, nil
}
