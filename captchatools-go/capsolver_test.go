package captchatoolsgo

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// go test -v -timeout 30s -run ^TestCapsolverBalance$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapsolverBalance(t *testing.T) {
	// Load Env
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Run test
	c := Capsolver{
		&Config{Api_key: os.Getenv("CAPSOLVER_KEY")},
	}
	if _, err := c.GetBalance(); err != nil {
		t.Fatalf("Error getting balance: %v", err)
	}
}

// go test -v -run ^TestCapsolverRecapV2$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapsolverRecapV2(t *testing.T) {
	// Load Env
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Run test
	c := Capsolver{
		&Config{
			Api_key:     os.Getenv("CAPSOLVER_KEY"),
			Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
			CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
			CaptchaType: V2Captcha,
		},
	}

	answer, err := c.GetToken()
	if err != nil {
		t.Fatalf("Error getting token: %v", err)
	}
	if answer == nil {
		t.Fatal("Answer is nil")
	}
	fmt.Println(answer)

}

// go test -v -run ^TestCapsolverRecapV3$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapsolverRecapV3(t *testing.T) {
	// Load Env
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Run test
	c := Capsolver{
		&Config{
			Api_key:     os.Getenv("CAPSOLVER_KEY"),
			Sitekey:     "6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf",
			CaptchaURL:  "https://antcpt.com/score_detector/",
			CaptchaType: V3Captcha, Action: "homepage", MinScore: 0.7,
		},
	}

	answer, err := c.GetToken()
	if err != nil {
		t.Fatalf("Error getting token: %v", err)
	}
	if answer == nil {
		t.Fatal("Answer is nil")
	}
	fmt.Println(answer)
}
