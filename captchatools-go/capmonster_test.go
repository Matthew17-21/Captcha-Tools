package captchatoolsgo

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestCapmonsterGetID tests that it can successfully get a ID from Capmonster
// to run this specific test:
// go test -v -run ^TestCapmonsterGetID$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapmonsterGetID(t *testing.T) {
	// Load Env
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
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
			_, err := a.getID(nil)
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
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
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
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
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
			_, err := a.getCaptchaAnswer(context.Background())
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// TestCapmonsterGetV2Additional tests that it can successfully get a V2 token from 2Captcha
// with additional data
// to run this specific test:
// go test -v -run ^TestCapmonsterGetV2Additional$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func TestCapmonsterGetV2Additional(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
	configs := []testConfigs{
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working V2 Config with custom User Agent",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			AdditionalData: &AdditionalData{
				UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			},
			ExpectError: false,
		},
		{
			SolvingSite: CapmonsterSite,
			Name:        "Working V2 Config with proxy",
			Config:      &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			AdditionalData: &AdditionalData{
				Proxy: &Proxy{
					Ip:       "",
					Port:     "",
					User:     "",
					Password: "",
				},
				UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
			},
			ExpectError: false,
		},
		{
			SolvingSite:    CapmonsterSite,
			Name:           "Bad V2 Config - Sitekey",
			Config:         &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJJ", CaptchaURL: "https://www.google.com/recaptcha/api2/demo", CaptchaType: "v2"},
			AdditionalData: &AdditionalData{},
			ExpectError:    true,
		},
		{
			SolvingSite:    CapmonsterSite,
			Name:           "Bad V2 Config - URL",
			AdditionalData: &AdditionalData{},
			Config:         &Config{Api_key: capmonsterKey, Sitekey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", CaptchaURL: "", CaptchaType: "v2"},
			ExpectError:    true,
		},
	}

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Capmonster{c.Config}
			_, err := a.getCaptchaAnswer(context.Background(), c.AdditionalData)
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
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
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
			_, err := a.getCaptchaAnswer(context.Background())
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
		})
	}
}

// Test2CapmonsterGetImage tests that it can successfully solve a image captcha
// to run this specific test:
// go test -v -run ^Test2CapmonsterGetImage$ github.com/Matthew17-21/Captcha-Tools/captchatools-go
func Test2CapmonsterGetImage(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}
	capmonsterKey := os.Getenv("CAPMONSTER_KEY")
	configs := []testConfigs{
		{
			SolvingSite: CapmonsterSite,
			Name:        "Image Captcha with correct config",
			Config:      &Config{Api_key: capmonsterKey, CaptchaType: ImageCaptcha},
			Image:       "iVBORw0KGgoAAAANSUhEUgAAAUQAAAAxCAYAAACictAAAAAACXBIWXMAAAsTAAALEwEAmpwYAAAGt0lEQVR4nO1di23bMBC9bKBsIG+QbqBOUHcCa4N4A2uDeIIqGzgTRCNkg2qDaIO6EEABgmDxezz+7gEEglbm8emOj38KgMFgMBgMBoPBYDAYDAaDwWAwGAwGg8FgMBgMhjleAKAFgA4ArgDQi79fAaABgCohDq8rDv2KSyueyQE5+IuSV70TG1eR/1E8w7DEOwDcd9KYiK1GBMQkyX+dBgA4IVY2V15L+QcDDkvevUEFGAzyNklzBQ3pr1x5gfi/WfxuhrHxJfK2RSX4X0VePjWC0pbSgaoKF7Ot+UW+OQR8HQkvjMp70bATWjh8+StXXq8IHEaDOG9FAzsa5m8DSlvaGAkLgG2rNnyZJoFIzQurAs8BJkNI4fDpr1x5dYhcThpcbPK11QhKW1q4EBYA2xZGEMoCkZoXZiWeeyoQmXD49leuvDAF8a4x71ysINaEBfBhC6sC1JHwwq7IsQmHb3/lyqvTiK0l6dj5UvApVhBHwgJg22o18ulEa7hMVldirq/blKeOhNf6+UlUtKvgehRlb8TfsgWcJc2/pRSOOrC/cuW1FcRB/NucxxZzHmeNhZdGwqlIQbwQFsCHrVEhBDorx40IrjoSXpMou+52DNVQba8nUInfmiZZT+UjAn/lymsp36AQsjVeFKIom2d+FLu9EP/RsyD6tOU8zMMogA9brcOCAhZ88KqQyzFXCEzIemC/IvdXyrxs9xR2Fo0liFhdRibbmMQWKUpbu7jtGBk8FMCHrZvkt1SbUSnfoQqyPVtYkAnvmIC/SuOl2grmImBU8U1ia69Vqz1UZl+2psC9Dcp3qIMbQQ9RNmfZR+6vEnm5in0Rgrg35zR3rQG5MvuyJWv1KI6xUb5D18AZCGwsDUGs/iqVF4hy2K40FyGI74rMMSuzL1tnw97QMtGe4jvUgewUg8uRLd25siFyf5XIy+cc55iLILYarRlWZfZp66oRwMed88DLVpZTAu9QhUqxyo05jyVbdDhF7K9SeWHwy1oQq53Mt60ERmX2besm2R5RKxYYtnaOEfHaw3LZwzrp7H004eZzHiqUv0rlteDoyC9rQXzT7EFgVGbftvZ+pysU23SOhJcpX1mw6+5Vc926oTPsCuWvUnnJGu+7SKeSBbE2eCmulZnClk2wqVIT0TvUze9RRew83Btou+gQ0l8l8/J9ld+YuiDqDPOwKjOFLR+BOEb0DnXzW6ev1QkXqmGX7gp2CH+VygsU88oTwrxy0oJ4MZxsd6nMVLZGwx5TJ5nvWacmgndokp+MN9Ziyg1h2EXtr5J5XQiG5WOqgniwcLhtZaa0pQrEb8mh99HwzColL1l+/zZJVakmhIl6rE29lP4qmZfqkokOcJCsID5a5fpU/Ma2MlPaGh1aQdlm2ikwL1mZ201arpJXieJLBHvYKP1VKq/fRGKYrCDKjpZhV2ZKW6pAnByHoHVAXjZQbe8YAi46UPurVF4/FLfadICLJAVx7zRFp0h7BZg2z50C2VIFks5wY++EwX0zxKHmZYtKIYo2c1IN4pEvKn+VyOtALIZZCSJmGgLZUtnbuxBVd4WxDcjLBQ1ypZBxNxVyKn+Vxuug6KV24AcsiJEJ4tUxCJoMBREkPYUP5PsVTfc5UvmrJF6hxHAGC2Jkgnj2GIhNwoI4Ig0FsS8EoPJXKbxCiuEMFsTIBLFxDGzZUIUFUV7ZsOcjMf1VAq/nwGI4gwUxMkGsDJ41HepUifYQayQ7Pm5YpvJX7ryeFYtnFGKYrCAeNVZDsVZIKW3prPA9I71kal4/AeAJ8E8ofARadKD2V868niIRw2QF0RaUe+hcbJ0Vw5UnC+HoA/O6i6BvxTzRk6ZAqo5rnZA+muVyHDCkv3Lg1SueP6ySydcGH6Ha5LdNMpGS2aoC28paECuNTyweNs+rhKOJQBDX6UucRunESYSfq3QUQy7Vt3dNKrzPL8iF9FfqvGTDfddkei2aS+oD28paEHVf5qh5CB+TH5YgUgcG9ubuWPyVOi8WRGBB1MHc2v5FEo46Q0GU3chj8jEirMYihL9y4MWCCCyIYBDw345BiD0hHYMgmoihz0WH0P7KgRcLIrAgmgajTQs9IV/b7iqIuld7Yd+HqJoHw/5SHJW/cuHFgggsiKaoRQv7V1NUMC9RxRLEWlzr9Sl6G6YCOVjOibUP7ltc0h/wAwp/5cKrkfB4lEzE+BE6T7b6wLaKRSNe9J/VV+p6sSrbevjmiC+8rPZCvgk+S1r4vAq+qXDK2V+l8GIwGAwGg8FgMBgMBoPBYDAYDAaDwWAwGAwGA/LEf2oS4NVP9R70AAAAAElFTkSuQmCC",
			ExpectError: false,
		},
	}
	EXPECTED_ANSWER := "446437676211"

	for _, c := range configs {
		t.Run(c.Name, func(t *testing.T) {
			a := &Capmonster{c.Config}
			answer, err := a.getCaptchaAnswer(context.Background(), &AdditionalData{B64Img: c.Image})
			if err != nil && !c.ExpectError {
				t.Fatalf(`getID() Error: %v , wanted: %v`, err, nil)
			}
			if answer.Token != EXPECTED_ANSWER {
				t.Fatalf(`getCaptchaAnswer() Got: %v || Wanted: %v`, answer.Token, EXPECTED_ANSWER)
			}
		})
	}
}
