package errors

import (
	"errors"
	"testing"
)

func TestErrCodeToError(t *testing.T) {
	tests := []struct {
		name     string
		errorID  string
		expected error
	}{
		// Visible/Invisible CAPTCHA errors
		{
			name:     "Visible CAPTCHA error",
			errorID:  "ERROR_VISIBLE_RECAPTCHA",
			expected: ErrVisibleCaptcha,
		},
		{
			name:     "Invisible CAPTCHA error",
			errorID:  "ERROR_INVISIBLE_RECAPTCHA",
			expected: ErrInvisibleCaptcha,
		},

		// Browser compatibility errors
		{
			name:     "Old browser error",
			errorID:  "ERROR_RECAPTCHA_OLD_BROWSER",
			expected: ErrOldUA,
		},

		// Domain errors
		{
			name:     "Invalid domain error 1",
			errorID:  "ERROR_RECAPTCHA_INVALID_DOMAIN",
			expected: ErrDomain,
		},
		{
			name:     "Invalid domain error 2",
			errorID:  "ERROR_DOMAIN_NOT_ALLOWED",
			expected: ErrDomain,
		},
		{
			name:     "Invalid domain error 3",
			errorID:  "ERROR_PAGEURL",
			expected: ErrDomain,
		},

		// Timeout errors
		{
			name:     "CAPTCHA timeout error",
			errorID:  "ERROR_RECAPTCHA_TIMEOUT",
			expected: ErrCaptchaTimeout,
		},

		// Proxy errors
		{
			name:     "Proxy banned error 1",
			errorID:  "ERROR_PROXY_BANNED",
			expected: ErrProxyBanned,
		},
		{
			name:     "Proxy banned error 2",
			errorID:  "ERROR_PROXY_TRANSPARENT",
			expected: ErrProxyBanned,
		},
		{
			name:     "Bad proxy error 1",
			errorID:  "ERROR_BAD_PROXY",
			expected: ErrProxy,
		},
		{
			name:     "Bad proxy error 2",
			errorID:  "ERROR_PROXY_CONNECT_REFUSED",
			expected: ErrProxy,
		},
		{
			name:     "Bad proxy error 3",
			errorID:  "ERROR_PROXY_CONNECT_TIMEOUT",
			expected: ErrProxy,
		},
		{
			name:     "Bad proxy error 4",
			errorID:  "ERROR_PROXY_READ_TIMEOUT",
			expected: ErrProxy,
		},
		{
			name:     "Bad proxy error 5",
			errorID:  "ERROR_PROXY_NOT_AUTHORISED",
			expected: ErrProxy,
		},
		{
			name:     "Bad proxy error 6",
			errorID:  "ERROR_PROXY_FORMAT",
			expected: ErrProxy,
		},

		// CAPTCHA ID errors
		{
			name:     "No CAPTCHA ID error 1",
			errorID:  "ERROR_NO_SUCH_CAPCHA_ID",
			expected: ErrNoCaptchaID,
		},
		{
			name:     "No CAPTCHA ID error 2",
			errorID:  "WRONG_CAPTCHA_ID",
			expected: ErrNoCaptchaID,
		},

		// Unsolvable CAPTCHA errors
		{
			name:     "Unsolvable CAPTCHA error",
			errorID:  "ERROR_CAPTCHA_UNSOLVABLE",
			expected: ErrUnsolvable,
		},

		// Ban-related errors
		{
			name:     "Ban-related error 1",
			errorID:  "MAX_USER_TURN",
			expected: ErrBanned,
		},
		{
			name:     "Ban-related error 2",
			errorID:  "ERROR_IP_NOT_ALLOWED",
			expected: ErrBanned,
		},
		{
			name:     "Ban-related error 3",
			errorID:  "IP_BANNED",
			expected: ErrBanned,
		},
		{
			name:     "Ban-related error 4",
			errorID:  "ERROR_TOO_MUCH_REQUESTS",
			expected: ErrBanned,
		},
		{
			name:     "Ban-related error 5",
			errorID:  "ERROR_IP_BANNED",
			expected: ErrBanned,
		},
		{
			name:     "Ban-related error 6",
			errorID:  "ERROR_IP_BLOCKED",
			expected: ErrBanned,
		},
		{
			name:     "Ban-related error 7",
			errorID:  "ERROR_ACCOUNT_SUSPENDED",
			expected: ErrBanned,
		},

		// Slot availability errors
		{
			name:     "No slot error 1",
			errorID:  "ERROR_NO_SLOT_AVAILABLE",
			expected: ErrNoSlot,
		},
		{
			name:     "No slot error 2",
			errorID:  "ERROR_ALL_WORKERS_FILTERED",
			expected: ErrNoSlot,
		},

		// CAPTCHA data errors
		{
			name:     "CAPTCHA data error 1",
			errorID:  "ERROR_TASK_ABSENT",
			expected: ErrCaptchaData,
		},
		{
			name:     "CAPTCHA data error 2",
			errorID:  "ERROR_TASK_NOT_SUPPORTED",
			expected: ErrCaptchaData,
		},
		{
			name:     "CAPTCHA data error 3",
			errorID:  "ERROR_BAD_TOKEN_OR_PAGEURL",
			expected: ErrCaptchaData,
		},

		// Balance errors
		{
			name:     "No balance error",
			errorID:  "ERROR_ZERO_BALANCE",
			expected: ErrNoBalance,
		},

		// Site key errors
		{
			name:     "Wrong site key error 1",
			errorID:  "ERROR_RECAPTCHA_INVALID_SITEKEY",
			expected: ErrWrongSitekey,
		},
		{
			name:     "Wrong site key error 2",
			errorID:  "ERROR_GOOGLEKEY",
			expected: ErrWrongSitekey,
		},
		{
			name:     "Wrong site key error 3",
			errorID:  "ERROR_SITEKEY",
			expected: ErrWrongSitekey,
		},
		{
			name:     "Wrong site key error 4",
			errorID:  "ERROR_WRONG_GOOGLEKEY",
			expected: ErrWrongSitekey,
		},

		// API key errors
		{
			name:     "Wrong API key error 1",
			errorID:  "ERROR_KEY_DOES_NOT_EXIST",
			expected: ErrWrongAPIKey,
		},
		{
			name:     "Wrong API key error 2",
			errorID:  "ERROR_WRONG_USER_KEY",
			expected: ErrWrongAPIKey,
		},

		// Missing values errors
		{
			name:     "Missing values error",
			errorID:  "ERROR_INCORRECT_SESSION_DATA",
			expected: ErrMissingValues,
		},

		// Unknown error cases
		{
			name:     "Unknown error code",
			errorID:  "SOME_UNKNOWN_ERROR_CODE",
			expected: ErrUnknown,
		},
		{
			name:     "Empty error code",
			errorID:  "",
			expected: ErrUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ErrCodeToError(tc.errorID)

			// Check if the returned error matches the expected error
			if !errors.Is(result, tc.expected) {
				t.Errorf("ErrCodeToError(%q) = %v, want %v", tc.errorID, result, tc.expected)
			}
		})
	}
}

// TestErrorMessages ensures that all predefined errors have the expected error messages
func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{name: "ErrNoBalance", err: ErrNoBalance, expected: "no balance on site"},
		{name: "ErrWrongAPIKey", err: ErrWrongAPIKey, expected: "incorrect API Key for captcha solving site"},
		{name: "ErrWrongSitekey", err: ErrWrongSitekey, expected: "incorrect sitekey for captcha"},
		{name: "ErrNoHarvester", err: ErrNoHarvester, expected: "incorrectly chose a captcha harvester. Refer to guide"},
		{name: "ErrIncorrectCapType", err: ErrIncorrectCapType, expected: "incorrectly chose a captcha type. Refer to guide"},
		{name: "ErrMaxAttempts", err: ErrMaxAttempts, expected: "max attempts passed"},
		{name: "ErrCaptchaData", err: ErrCaptchaData, expected: "there was an error with the captcha details"},
		{name: "ErrNoSlot", err: ErrNoSlot, expected: "no idle captcha workers are available"},
		{name: "ErrBanned", err: ErrBanned, expected: "ip and/or api key banned from solving site"},
		{name: "ErrUnsolvable", err: ErrUnsolvable, expected: "captcha is unsolvable"},
		{name: "ErrNoCaptchaID", err: ErrNoCaptchaID, expected: "captcha id does not exist"},
		{name: "ErrUnknown", err: ErrUnknown, expected: "unknown error"},
		{name: "ErrProxy", err: ErrProxy, expected: "could not connect to provided proxy"},
		{name: "ErrProxyBanned", err: ErrProxyBanned, expected: "proxy IP banned by target service"},
		{name: "ErrCaptchaTimeout", err: ErrCaptchaTimeout, expected: "recaptcha task timeout, probably due to slow proxy server or Google server"},
		{name: "ErrDomain", err: ErrDomain, expected: "captcha provider reported that the domain for this site key is invalid"},
		{name: "ErrOldUA", err: ErrOldUA, expected: "captcha provider reported that the browser user-agent is not compatible with their javascript"},
		{name: "ErrInvisibleCaptcha", err: ErrInvisibleCaptcha, expected: "an attempt was made to solve an Invisible Recaptcha as if it was a regular one"},
		{name: "ErrVisibleCaptcha", err: ErrVisibleCaptcha, expected: "attempted solution of usual Recaptcha V2 as Recaptcha V2 invisible. Remove flag 'isInvisible' from the API payload"},
		{name: "ErrMissingValues", err: ErrMissingValues, expected: "some of the required values for successive user emulation are missing"},
		{name: "ErrAddionalDataMissing", err: ErrAddionalDataMissing, expected: "additional data is missing. Refer to guide"},
		{name: "ErrNotSupported", err: ErrNotSupported, expected: "captcha type not supported"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.expected {
				t.Errorf("%s error message = %q, want %q", tc.name, tc.err.Error(), tc.expected)
			}
		})
	}
}

// TestErrorComparisons ensures that we can properly compare errors
func TestErrorComparisons(t *testing.T) {
	// Test direct comparison
	if ErrCodeToError("ERROR_ZERO_BALANCE") != ErrNoBalance {
		t.Errorf("Error comparison failed: ErrCodeToError(\"ERROR_ZERO_BALANCE\") != ErrNoBalance")
	}

	// Test with errors.Is()
	if !errors.Is(ErrCodeToError("ERROR_ZERO_BALANCE"), ErrNoBalance) {
		t.Errorf("errors.Is() comparison failed: !errors.Is(ErrCodeToError(\"ERROR_ZERO_BALANCE\"), ErrNoBalance)")
	}
}
