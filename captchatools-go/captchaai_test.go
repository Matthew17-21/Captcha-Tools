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

// Test getting a recap V3 token
// go test -v -run ^TestCaptchaAiV3$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCaptchaAiV3(t *testing.T) {
	// Load ENV
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Create tests
	configs := []Config{
		{Api_key: os.Getenv("CAPTCHAAI_KEY"), Sitekey: "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf", CaptchaURL: "https://antcpt.com/score_detector/", CaptchaType: V3Captcha, Action: "homepage", MinScore: 0.7},
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

// Test getting a recap V3 token
// go test -v -run ^TestCaptchaAiHCap$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCaptchaAiHCap(t *testing.T) {
	// Load ENV
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Create tests
	configs := []Config{
		{Api_key: os.Getenv("CAPTCHAAI_KEY"), Sitekey: "a5f74b19-9e45-40e0-b45d-47ff91b7a6c2", CaptchaURL: "https://accounts.hcaptcha.com/demo", CaptchaType: HCaptcha},
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
