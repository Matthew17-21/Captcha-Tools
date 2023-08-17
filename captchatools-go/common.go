package captchatoolsgo

// General type declarations
const (
	CapmonsterSite  site = iota // The int 1 will represent Capmonter
	AnticaptchaSite             // The int 2 will represent Anticaptcha
	TwoCaptchaSite              // The int 3 will represent 2captcha
	CapsolverSite
	CaptchaAiSite
)

const (
	V2Captcha    captchaType = "v2"
	V3Captcha    captchaType = "v3"
	HCaptcha     captchaType = "hcaptcha"
	ImageCaptcha captchaType = "image"
	CFTurnstile  captchaType = "cfturnstile"
)

type (
	site        int
	captchaType string
)
