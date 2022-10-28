package captchatoolsgo

import "errors"

var (
	// Error type declarations
	ErrNoBalance        = errors.New("no balance on site")
	ErrWrongAPIKey      = errors.New("incorrect API Key for captcha solving site")
	ErrWrongSitekey     = errors.New("incorrect sitekey for captcha")
	ErrNoHarvester      = errors.New("incorrectly chose a captcha harvester. Refer to guide")
	ErrIncorrectCapType = errors.New("incorrectly chose a captcha type. Refer to guide")
	ErrMaxAttempts      = errors.New("max attempts passed")
	ErrCaptchaData      = errors.New("there was an error with the captcha details")
)

// errCodeToError converts an error ID returned from the site and
// returns the appropriate error
func errCodeToError(id string) error {
	var err error = nil
	switch id {
	case "ERROR_TASK_ABSENT":
		err = ErrCaptchaData
	case "ERROR_ZERO_BALANCE":
		err = ErrNoBalance
	case "ERROR_RECAPTCHA_INVALID_SITEKEY":
		err = ErrWrongSitekey
	case "ERROR_KEY_DOES_NOT_EXIST":
		err = ErrWrongAPIKey
	}
	return err
}
