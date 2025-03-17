package harvester

// CaptchaType represents the type of CAPTCHA being solved.
// It is used to specify which solving method should be employed.
type CaptchaType string

const (
	// V2Captcha represents Google reCAPTCHA v2 challenges that require
	// clicking the "I'm not a robot" checkbox or solving image puzzles.
	V2Captcha CaptchaType = "v2"

	// V3Captcha represents Google reCAPTCHA v3 challenges that run
	// in the background and assign a score based on user behavior.
	V3Captcha CaptchaType = "v3"

	// HCaptcha represents standard hCaptcha challenges similar to reCAPTCHA v2.
	HCaptcha CaptchaType = "hcaptcha"

	// HcaptchaTurbo represents the accelerated/premium version of hCaptcha
	// that typically offers faster solving times.
	HcaptchaTurbo CaptchaType = "hcaptchaturbo"

	// ImageCaptcha represents traditional image-based captchas where text or
	// numbers must be recognized and entered.
	ImageCaptcha CaptchaType = "image"

	// CFTurnstile represents Cloudflare Turnstile captcha challenges,
	// Cloudflare's alternative to reCAPTCHA.
	CFTurnstile CaptchaType = "cfturnstile"
)

// CaptchaAnswer represents the result of a successfully solved captcha challenge.
// It provides methods to access the solution token, metadata about how the
// captcha was solved, and the ability to report the solution's correctness
// back to the solving service.
type CaptchaAnswer interface {
	// Token returns the actual captcha solution string that can be submitted
	// to the target website to validate the captcha challenge.
	Token() string

	// UserAgent returns the user agent string that was used to solve the captcha.
	// This may be empty if the service doesn't provide this information.
	UserAgent() string

	// ID returns the unique identifier assigned by the solving service.
	// This ID can be used for reporting or tracking purposes.
	ID() any

	// Report submits feedback to the solving service about whether the
	// captcha solution was correct. If the solution was incorrect and the
	// solving service accepts the report, a refund may be credited to the account.
	//
	// Parameters:
	//   wasCorrect - true if the solution worked, false otherwise
	//
	// Returns an error if the reporting attempt fails.
	Report(wasCorrect bool) error
}
