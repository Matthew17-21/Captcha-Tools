package harvester

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
