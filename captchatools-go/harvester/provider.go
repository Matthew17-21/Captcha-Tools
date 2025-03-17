package harvester

type Provider string

const (
	CapmonsterSite  Provider = "capmonster"
	AnticaptchaSite Provider = "anticaptcha"
	TwoCaptchaSite  Provider = "2captcha"
	CapsolverSite   Provider = "capsolver"
	CaptchaAiSite   Provider = "captchaai"
)
