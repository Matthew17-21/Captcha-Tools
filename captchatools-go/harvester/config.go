package harvester

import "github.com/Matthew17-21/Captcha-Tools/captchatools-go/logger"

// Config defines the configuration settings required for solving various types of captchas.
// It contains all the parameters needed to initialize a captcha-solving request,
// including authentication, target information, and solving preferences.
type Config struct {
	Api_key            string      // The API Key for the captcha solving site.
	Sitekey            string      // Sitekey from the site where captcha is loaded.
	CaptchaURL         string      // URL where the captcha is located.
	CaptchaType        CaptchaType // Type of captcha you are solving. Visit https://github.com/Matthew17-21/Captcha-Tools for types
	Action             string      // Action that is associated with the V3 captcha.
	IsInvisibleCaptcha bool        // If the captcha is invisible or not.
	MinScore           float32     // Minimum score for v3 captchas.
	SoftID             int         // SoftID for 2captcha. Developers get reward 10% of spendings of their software users.
	Logger             logger.Logger
}
