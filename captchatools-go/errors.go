package captchatoolsgo

import (
	"errors"
)

var (
	// Error type declarations
	ErrNoBalance           = errors.New("no balance on site")
	ErrWrongAPIKey         = errors.New("incorrect API Key for captcha solving site")
	ErrWrongSitekey        = errors.New("incorrect sitekey for captcha")
	ErrNoHarvester         = errors.New("incorrectly chose a captcha harvester. Refer to guide")
	ErrIncorrectCapType    = errors.New("incorrectly chose a captcha type. Refer to guide")
	ErrMaxAttempts         = errors.New("max attempts passed")
	ErrCaptchaData         = errors.New("there was an error with the captcha details")
	ErrNoSlot              = errors.New("no idle captcha workers are available")
	ErrBanned              = errors.New("ip and/or api key banned from solving site")
	ErrUnsolvable          = errors.New("captcha is unsolvable")
	ErrNoCaptchaID         = errors.New("captcha id does not exist")
	ErrUnknown             = errors.New("unknown error")
	ErrProxy               = errors.New("could not connect to provided proxy")
	ErrProxyBanned         = errors.New("proxy IP banned by target service")
	ErrCaptchaTimeout      = errors.New("recaptcha task timeout, probably due to slow proxy server or Google server")
	ErrDomain              = errors.New("captcha provider reported that the domain for this site key is invalid")
	ErrOldUA               = errors.New("captcha provider reported that the browser user-agent is not compatible with their javascript")
	ErrInvisibleCaptcha    = errors.New("an attempt was made to solve an Invisible Recaptcha as if it was a regular one")
	ErrVisibleCaptcha      = errors.New("attempted solution of usual Recaptcha V2 as Recaptcha V2 invisible. Remove flag 'isInvisible' from the API payload")
	ErrMissingValues       = errors.New("some of the required values for successive user emulation are missing")
	ErrAddionalDataMissing = errors.New("additional data is missing. Refer to guide")
	ErrProxyEmpty          = errors.New("proxy is blank")
	ErrNotSupported        = errors.New("captcha type not supported")
)

// errCodeToError converts an error ID returned from the site and
// returns the appropriate error
func errCodeToError(id string) error {
	var err error = ErrUnknown
	switch id {
	case "ERROR_VISIBLE_RECAPTCHA":
		err = ErrVisibleCaptcha
	case "ERROR_INVISIBLE_RECAPTCHA":
		err = ErrInvisibleCaptcha
	case "ERROR_RECAPTCHA_OLD_BROWSER":
		err = ErrOldUA
	case "ERROR_RECAPTCHA_INVALID_DOMAIN", "ERROR_DOMAIN_NOT_ALLOWED", "ERROR_PAGEURL":
		err = ErrDomain
	case "ERROR_RECAPTCHA_TIMEOUT":
		err = ErrCaptchaTimeout
	case "ERROR_PROXY_BANNED", "ERROR_PROXY_TRANSPARENT":
		err = ErrProxyBanned
	case "ERROR_BAD_PROXY", "ERROR_PROXY_CONNECT_REFUSED", "ERROR_PROXY_CONNECT_TIMEOUT", "ERROR_PROXY_READ_TIMEOUT", "ERROR_PROXY_NOT_AUTHORISED", "ERROR_PROXY_FORMAT":
		err = ErrProxy
	case "ERROR_NO_SUCH_CAPCHA_ID", "WRONG_CAPTCHA_ID":
		err = ErrNoCaptchaID
	case "ERROR_CAPTCHA_UNSOLVABLE":
		err = ErrUnsolvable
	case "MAX_USER_TURN", "ERROR_IP_NOT_ALLOWED", "IP_BANNED", "ERROR_TOO_MUCH_REQUESTS", "ERROR_IP_BANNED", "ERROR_IP_BLOCKED", "ERROR_ACCOUNT_SUSPENDED":
		err = ErrBanned
	case "ERROR_NO_SLOT_AVAILABLE", "ERROR_ALL_WORKERS_FILTERED":
		err = ErrNoSlot
	case "ERROR_TASK_ABSENT", "ERROR_TASK_NOT_SUPPORTED", "ERROR_BAD_TOKEN_OR_PAGEURL":
		err = ErrCaptchaData
	case "ERROR_ZERO_BALANCE":
		err = ErrNoBalance
	case "ERROR_RECAPTCHA_INVALID_SITEKEY", "ERROR_GOOGLEKEY", "ERROR_SITEKEY", "ERROR_WRONG_GOOGLEKEY":
		err = ErrWrongSitekey
	case "ERROR_KEY_DOES_NOT_EXIST", "ERROR_WRONG_USER_KEY":
		err = ErrWrongAPIKey
	case "ERROR_INCORRECT_SESSION_DATA":
		err = ErrMissingValues
	}
	return err
}
