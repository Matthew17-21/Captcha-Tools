package captchatoolsgo

import (
	"testing"
)

const twocapKey = "c9a8a86ed4e59e331e2ca6a304155d6b"

// Test2CaptchaGetID tests that it can successfully get a ID from 2Captcha
// to run this specific test:
// go test -v -run ^Test2CaptchaGetID$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func Test2CaptchaGetID(t *testing.T) {
	var tests = []testConfigs{
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Working V2 Config",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Working V3 Config",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: false,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Empty Config",
			Config:      &Config{},
			ExpectError: true,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Incorrect V2 Config - bad sitekey",
			Config:      &Config{Api_key: twocapKey, Sitekey: "", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: true,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Incorrect V2 Config - bad captcha url",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			a := &Twocaptcha{c.Config}
			_, err := a.getID()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// Test2CaptchaGetBalance tests that it can successfully get the balance from 2Captcha
// to run this specific test:
// go test -v -run ^Test2CaptchaGetBalance$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func Test2CaptchaGetBalance(t *testing.T) {
	var tests = []testConfigs{
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Working API Key",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Incorrect API Key",
			Config:      &Config{Api_key: "9f47074b59d3d4cf5c07961f90deb7d9", Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: true,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Blank API Key",
			Config:      &Config{Api_key: "", Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			a := &Twocaptcha{c.Config}
			_, err := a.GetBalance()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// Test2CaptchaGetV2 tests that it can successfully get a V2 token from 2Captcha
// to run this specific test:
// go test -v -run ^Test2CaptchaGetV2$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func Test2CaptchaGetV2(t *testing.T) {
	configs := []testConfigs{
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Working V2 Config",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Bad V2 Config - Sitekey",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJJ", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: true,
		},
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Bad V2 Config - URL",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError: true,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Twocaptcha{c.Config}
			_, err := a.getCaptchaAnswer()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// Test2CaptchaGetV3 tests that it can successfully get a V2 token from 2Captcha
// to run this specific test:
// go test -v -run ^Test2CaptchaGetV3$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func Test2CaptchaGetV3(t *testing.T) {
	configs := []testConfigs{
		{
			SolvingSite: TwoCaptchaSite,
			Name:        "Working V3 Config",
			Config:      &Config{Api_key: twocapKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: false,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Twocaptcha{c.Config}
			_, err := a.getCaptchaAnswer()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}
