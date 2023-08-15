package captchatoolsgo

// Test ran with the command:
// go test -v -run ^TestHarvester$ github.com/Matthew17-21/Captcha-Tools/captchatools-go

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type testConfigs struct {
	Name           string
	Config         *Config
	SolvingSite    site
	ExpectError    bool
	Image          string
	AdditionalData *AdditionalData
}

var configs = []testConfigs{
	{
		Name:        "Working 2captcha Key",
		Config:      &Config{Api_key: "c9a8a86ed4e59e331e2ca6a304155d6b"},
		SolvingSite: TwoCaptchaSite,
		ExpectError: false,
	},
	{
		Name:        "Working capmonster Key",
		Config:      &Config{Api_key: "2cffb45a7f3b15b7f7bfd5c53c08d0cd"},
		SolvingSite: CapmonsterSite,
		ExpectError: false,
	},
	{
		Name:        "Working anticaptcha Key",
		Config:      &Config{Api_key: "9f47074b59d3d4cf5c07961f90deb7d8"},
		SolvingSite: AnticaptchaSite,
		ExpectError: false,
	},
	{
		Name:        "Bad 2cap Key",
		SolvingSite: TwoCaptchaSite,
		Config:      &Config{Api_key: "abc"},
		ExpectError: true,
	},
	{
		Name:        "Bad anticap Key",
		SolvingSite: AnticaptchaSite,
		Config:      &Config{Api_key: "abc"},
		ExpectError: true,
	},
	{
		Name:        "Bad capmonster Key",
		SolvingSite: CapmonsterSite,
		Config:      &Config{Api_key: "abc"},
		ExpectError: true,
	},
}

func TestHarvester(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
	solver, err := NewHarvester(CapmonsterSite, &Config{
		Api_key:     capmonsterKey,
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

// TestGetBalance tests getting a balance from sites
// to run this test:
// go test -v -run ^TestGetBalance$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestGetBalance(t *testing.T) {
	for _, c := range configs {
		c.Config.CaptchaType = "v2"
		t.Run(c.Name, func(t *testing.T) {
			h, err := NewHarvester(c.SolvingSite, c.Config)
			if err != nil {
				t.Fatalf(`NewHarvester() Error: %v, wanted: %v`, err, "nil")
			}
			_, err = h.GetBalance()
			if err != nil && !c.ExpectError {
				t.Fatalf(`GetBalance() Error: %v, wanted: %v`, err, nil)
			}
		})
	}
}

// TestGetV2 tests getting a V2 recaptcha token
func TestGetV2(t *testing.T) {

}

// TestGetV3 tests getting a V3 recaptcha token
func TestGetV3(t *testing.T) {
}

// go test -v -run ^TestCapsolverHarvester$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapsolverHarvester(t *testing.T) {
	// Load Env
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	h, err := NewHarvester(CapsolverSite, &Config{Api_key: os.Getenv("CAPSOLVER_KEY"), CaptchaType: ImageCaptcha})
	if err != nil {
		t.Fatalf("Error getting new harvester: %v", err)
	}
	if _, err := h.GetBalance(); err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
}
