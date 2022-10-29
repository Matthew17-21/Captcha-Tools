package captchatoolsgo

import (
	"testing"
)

const capmonsterKey = "2cffb45a7f3b15b7f7bfd5c53c08d0cd"

// TestCapmonsterGetID tests that it can successfully get a ID from Capmonster
// to run this specific test:
// go test -v -run ^TestCapmonsterGetID$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapmonsterGetID(t *testing.T) {
	var tests = []testConfigs{
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working V2 Config",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working V3 Config",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: false,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Empty Config",
			Config:      &Config{},
			ExpectError: true,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Incorrect V2 Config - bad sitekey",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: true,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Incorrect V2 Config - bad captcha url",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			a := &Capmonster{c.Config}
			_, err := a.getID()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestCapmonsterGetBalance tests that it can successfully get the balance from Capmonster
// to run this specific test:
// go test -v -run ^TestCapmonsterGetBalance$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapmonsterGetBalance(t *testing.T) {
	var tests = []testConfigs{
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working API Key",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Incorrect API Key",
			Config:      &Config{Api_key: "9f47074b59d3d4cf5c07961f90deb7d9", Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: true,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Blank API Key",
			Config:      &Config{Api_key: "", Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.Name, func(t *testing.T) {
			a := &Capmonster{c.Config}
			_, err := a.GetBalance()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestCapmonsterGetV2 tests that it can successfully get a V2 token from Capmonster
// to run this specific test:
// go test -v -run ^TestCapmonsterGetV2$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapmonsterGetV2(t *testing.T) {
	configs := []testConfigs{
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working V2 Config",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: false,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Bad V2 Config - Sitekey",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJJ", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			ExpectError: true,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Bad V2 Config - URL",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError: true,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Capmonster{c.Config}
			_, err := a.getCaptchaAnswer()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestCapmonsterGetV3 tests that it can successfully get a V2 token from Capmonster
// to run this specific test:
// go test -v -run ^TestCapmonsterGetV3$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapmonsterGetV3(t *testing.T) {
	configs := []testConfigs{
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working V3 Config",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: "v3", Action: "homepage", MinScore: 0.7},
			ExpectError: false,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Capmonster{c.Config}
			_, err := a.getCaptchaAnswer()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}
