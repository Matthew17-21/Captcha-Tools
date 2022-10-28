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
)
