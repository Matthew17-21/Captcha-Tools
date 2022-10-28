package captchatoolsgo

import (
	"testing"
)

const antiCapKey = "9f47074b59d3d4cf5c07961f90deb7d8"

var antiCapConfigs = []testConfigs{
	{
		SolvingSite: AnticaptchaSite,
		Name:        "Working Config",
		Config:      &Config{Api_key: antiCapKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
		ExpectError: false,
	},
	{
		SolvingSite: AnticaptchaSite,
		Name:        "Empty Config",
		Config:      &Config{},
		ExpectError: true,
	},
}

// TestGetAnticapID tests that it can successfully get a ID from anticaptcha
// to run this specific test:
// go test -v -run ^TestGetAnticapID$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestGetAnticapID(t *testing.T) {
	for _, c := range antiCapConfigs {
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
	for _, c := range antiCapConfigs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Anticaptcha{c.Config}
			_, err := a.GetBalance()
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}
