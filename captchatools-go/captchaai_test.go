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

// Test getting a Image captcha
// go test -v -run ^TestCaptchaAiNormalCap$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCaptchaAiNormalCap(t *testing.T) {
	// Load ENV
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Create tests
	configs := []Config{
		{Api_key: os.Getenv("CAPTCHAAI_KEY"), CaptchaType: ImageCaptcha},
	}

	// Run tests
	for testNum, config := range configs {
		t.Run(fmt.Sprintf("Test #%v", testNum+1), func(t *testing.T) {
			h := CaptchaAi{&config}
			answer, err := h.GetToken(&AdditionalData{
				B64Img: "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",
			})
			if err != nil {
				t.Fatalf("Error getting token: %v", err)
			}
			fmt.Println(answer)
		})
	}

}

// Test getting balance info
// go test -v -run ^TestCaptchaAiGetBalance$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCaptchaAiGetBalance(t *testing.T) {
	// Load ENV
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Create tests
	configs := []Config{
		{Api_key: os.Getenv("CAPTCHAAI_KEY"), CaptchaType: ImageCaptcha},
	}

	// Run tests
	for testNum, config := range configs {
		t.Run(fmt.Sprintf("Test #%v", testNum+1), func(t *testing.T) {
			h := CaptchaAi{&config}
			balance, err := h.GetBalance()
			if err != nil {
				t.Fatalf("Error getting balance: %v", err)
			}
			fmt.Println(balance)
		})
	}

}
