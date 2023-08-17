package captchatoolsgo

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// Test getting a recap V2 token
// go test -v -run ^TestCaptchaAiV2$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCaptchaAiV2(t *testing.T) {
	// Load ENV
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Create tests
	configs := []Config{
		{Api_key: os.Getenv("CAPTCHAAI_KEY"), Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: V2Captcha},
		{Api_key: os.Getenv("CAPTCHAAI_KEY"), Sitekey: "6LcmDCcUAAAAAL5QmnMvDFnfPTP4iCUYRk2MwC0-", CaptchaURL: "https://recaptcha-demo.appspot.com/recaptcha-v2-invisible.php", CaptchaType: V2Captcha, IsInvisibleCaptcha: true},
	}

	// Run tests
	for testNum, config := range configs {
		t.Run(fmt.Sprintf("Test #%v", testNum+1), func(t *testing.T) {
			h := CaptchaAi{&config}
			answer, err := h.GetToken()
			if err != nil {
				t.Fatalf("Error getting token: %v", err)
			}
			fmt.Println(answer)

		})
	}

}
