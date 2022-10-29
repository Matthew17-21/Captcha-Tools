package captchatoolsgo

import (
	"testing"
)

const antiCapKey = "9f47074b59d3d4cf5c07961f90deb7d8"

// TestGetAnticapID tests that it can successfully get a ID from anticaptcha
// to run this specific test:
// go test -v -run ^TestGetAnticapID$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestGetAnticapID(t *testing.T) {
	var tests = []testConfigs{
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Working V2 Config",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Working V3 Config",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: false,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Empty Config",
			Config:      &Config{},
			ExpectError: true,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Incorrect V2 Config - bad sitekey",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: true,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Incorrect V2 Config - bad captcha url",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			a := &Anticaptcha{c.Config}
			_, err := a.getID()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestGetAnticapBalance tests that it can successfully get the balance from anticaptcha
// to run this specific test:
// go test -v -run ^TestGetAnticapBalance$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestGetAnticapBalance(t *testing.T) {
	var tests = []testConfigs{
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Working API Key",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Incorrect API Key",
			Config:      &Config{Api_key: "9f47074b59d3d4cf5c07961f90deb7d9", Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: true,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Blank API Key",
			Config:      &Config{Api_key: "", Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			a := &Anticaptcha{c.Config}
			_, err := a.GetBalance()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestGetAnticapV2 tests that it can successfully get a V2 token from anticaptcha
// to run this specific test:
// go test -v -run ^TestGetAnticapV2$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestGetAnticapV2(t *testing.T) {
	configs := []testConfigs{
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Working V2 Config",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Bad V2 Config - Sitekey",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJJ", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: true,
		},
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Bad V2 Config - URL",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError: true,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Anticaptcha{c.Config}
			_, err := a.getCaptchaAnswer()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestGetAnticapV3 tests that it can successfully get a V2 token from anticaptcha
// to run this specific test:
// go test -v -run ^TestGetAnticapV3$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestGetAnticapV3(t *testing.T) {
	configs := []testConfigs{
		{
			SolvingSite: AnticaptchaSite,
			Name:        "Working V3 Config",
			Config:      &Config{Api_key: antiCapKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: false,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Anticaptcha{c.Config}
			_, err := a.getCaptchaAnswer()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}
