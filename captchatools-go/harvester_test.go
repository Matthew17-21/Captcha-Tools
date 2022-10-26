package captchatoolsgo

// Test ran with the command:
// go test -v -run ^TestHarvester$ github.com/Matthew17-21/Captcha-Tools/captchatools-go

import (
	"testing"
)

func TestHarvester(t *testing.T) {
	solver, err := NewHarvester(CapmonsterSite, &Config{
		Api_key:     "2cffb45a7f3b15b7f7bfd5c53c08d0cd",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
		CaptchaType: "V2",
	})
	if err != nil {
		t.Fatalf(`NewHarvester() Error: %v, wanted: %v`, err, "nil")
	}
	_, err = solver.GetToken()
	if err != nil {
		t.Fatalf(`GetToken() error: %v, wanted: %v`, err, "nil")
	}
}
