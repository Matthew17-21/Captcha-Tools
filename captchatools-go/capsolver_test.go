package captchatoolsgo

import (
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
