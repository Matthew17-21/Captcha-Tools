package harvester

// CaptchaAnswer represents the result of a solved captcha
type CaptchaAnswer interface {
	// Token returns the actual captcha solution
	Token() string

	// UserAgent returns the user agent that was used to solve the captcha
	UserAgent() string

	// ID returns the identifier from the solving site
	ID() interface{}

	// Report submits feedback about whether the captcha solution was correct
	Report(wasCorrect bool) error
}
