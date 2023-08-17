# Captcha Tools (Go)
Go package to help solve captchas with Capmonster, 2Captcha and Anticaptcha API's!

# 2.0.0
### What's new
1. Get Balance Support
2. Proxy Support
3. User Agent Support
4. Text image captcha support
5. Refund support / Report answers support
6. AntiCaptcha SoftID support

### Breaking Changes
* Changed the type of site from type `int` to type `site`
	- This affects the `NewHarvester()` param
* Returns a `CaptchaAnswer` object rather than just the captcha token
* Made captchas their own type

### Upgrading to 2.0.0
* When passing in the config, instead of manually declaring the `CaptchaType` field, use a preset type such as `captchatools.ImageCaptcha`
	- See [captcha types](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#Captchas-Types)
* Instead of a captcha token being returned, a `CaptchaAnswer` is returned. To get the captcha answer, you must call the `.Token` field from the returned answer. [See Examples.](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#How-to-use)


# Install
```go
go get github.com/Matthew17-21/Captcha-Tools/captchatools-go
```
To update:
```go
go get -u github.com/Matthew17-21/Captcha-Tools/captchatools-go
```

# How to use
### Basic usage
```go
package main

import (
	"fmt"

	captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

func main() {
	solver, err := captchatools.NewHarvester(captchatools.CapmonsterSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
		CaptchaType: captchatools.V2Captcha,
	})
	if err != nil {
		...
	}
	answer, err := solver.GetToken()
	if err != nil {
		switch err {
		case captchatools.ErrBanned:
			...
		case captchatools.ErrAddionalDataMissing:
			...
		}
	}

	fmt.Println("Captcha ID:", answer.Id())
	fmt.Println("Captcha token:", answer.Token)

	// ........Use captcha token......
	answer.Report(true) // report good
	// or
	answer.Report(false) // report bad / request refund
}

```
### NewHarvester() Parameters:
| Parameter | Required |  Type  |  Description|
| :-------------: |:-------------:| :-----:| :-----:|
| solving_site | true | int (iota)|  The captcha solving site that will be used. Refer to [the site IDs](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#site-specific-support). Alternatively, you can use shortcuts such as `captchatools.AnticaptchaSite` |
| Config| true | captchatools.Config |  Configurations for the captchas you are solving. |
### Config struct fields:
| Field | Required |  Type  |  Description|
| :-------------: |:-------------:|  :-----:| :-----:|
| Api_key | true | String|  The API Key for the captcha solving site|
| Sitekey| true | String | Sitekey from the site where captcha is loaded|
| CaptchaURL | true| String |  URL where the captcha is located|
| CaptchaType| true| captchaType |  Type of captcha you are solving. See [captcha types](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#Captchas-Types) |
| Action | false | String |  Action that is associated with the V3 captcha.<br />__This param is only required when solving V3 captchas__|
| IsInvisibleCaptcha| false | bool |  If the captcha is invisible or not.<br />__This param is only required when solving invisible captchas__|
| MinScore | false | float32 |  Minimum score for v3 captchas.<br />__This param is only required when solving V3 and it needs a higher / lower score__|
| SoftID | false | int |  2captcha Developer ID. <br /> Developers get 10% of spendings of their software users. |
### AdditionalData struct fields:
| Field | Required |  Type  |  Description|
| :-------------: |:-------------:|  :-----:| :-----:|
| B64Img | false | String |  Base64 encoded captcha image<br />__This param is only required when solving image captchas__|
| Proxy| false | *Proxy |  Proxy to be used to solve captchas.<br />This will make the captcha be solved from the proxy ip|
| ProxyType | false | string |  Type of the proxy being used. Options are:<br /> `HTTP`, `HTTPS`, `SOCKS4`, `SOCKS5`|
| UserAgent | false | string |  UserAgent that will be passed to the service and used to solve the captcha |
### Examples
##### Example - V2 Captcha / Basic usage
```go
package main

import (
	"fmt"

	captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

func main() {
	solver, err := captchatools.NewHarvester(captchatools.CapmonsterSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
		CaptchaType: captchatools.V2Captcha,
	})
	if err != nil {
		...
	}
	answer, err := solver.GetToken()
	if err != nil {
		switch err {
		case captchatools.ErrBanned:
			...
		case captchatools.ErrAddionalDataMissing:
			...
		}
	}

	fmt.Println("Captcha ID:", answer.Id())
	fmt.Println("Captcha token:", answer.Token)

	// ........Use captcha token......
	answer.Report(true) // report good
	// or
	answer.Report(false) // report bad / request refund
}

```
##### Example - V3 Captcha
```go
func v3Example() {
	solver, err := captchatools.NewHarvester(captchatools.AnticaptchaSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "..........",
		CaptchaType: captchatools.V3Captcha,
		Action:      "submit",
		MinScore:    0.9,
	})
	if err != nil {
		panic(err)
	}
	answer, err := solver.GetToken()
	if err != nil {
		switch err {
		case captchatools.ErrBanned:
		case captchatools.ErrAddionalDataMissing:
		}
	}

	fmt.Println("Captcha ID:", answer.Id())
	fmt.Println("Captcha token:", answer.Token)

	// ........Use captcha token......
	answer.Report(true) // report good
	// or
	answer.Report(false) // report bad / request refund

}

```
##### Example - Image captcha
```go
package main

import (
	"fmt"

	captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

func image_captcha() {
	solver, err := captchatools.NewHarvester(captchatools.CapmonsterSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		CaptchaType: captchatools.ImageCaptcha,
	})
	if err != nil {
		panic(err)
	}

	answer, err := solver.GetToken(&captchatools.AdditionalData{
		B64Img: "BASE64_ENCODED_IMAGE",
	})
	if err != nil {
		switch err {
		case captchatools.ErrBanned:
		case captchatools.ErrAddionalDataMissing:
		}
	}

	fmt.Println("Captcha ID:", answer.Id())
	fmt.Println("Captcha token:", answer.Token)

	// ........Use captcha token......
	answer.Report(true) // report good
	// or
	answer.Report(false) // report bad / request refund

}

```
##### Example - Additional captcha data
```go
package main

import (
	"fmt"

	captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

func addtional_data() {
	solver, err := captchatools.NewHarvester(captchatools.CapmonsterSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
		CaptchaType: captchatools.V2Captcha,
	})
	if err != nil {
		panic(err)
	}

	// The following data is OPTIONAL and is not required to solve captchas
	proxy, err := captchatools.NewProxy("IP:PORT:USER:PASS")
	answer, err := solver.GetToken(&captchatools.AdditionalData{
		B64Img:    "BASE64_ENCODED_IMAGE",
		Proxy:     proxy,
		ProxyType: "http",
		UserAgent: "SOME_USER_AGENT",
	})
	if err != nil {
		switch err {
		case captchatools.ErrBanned:
		case captchatools.ErrAddionalDataMissing:
		}
	}

	fmt.Println("Captcha ID:", answer.Id())
	fmt.Println("Captcha token:", answer.Token)

	// ........Use captcha token......
	answer.Report(true) // report good
	// or
	answer.Report(false) // report bad / request refund

}
```

# Supported Sites
- **[Capmonster](https://capmonster.cloud/)**
- **[2Captcha](https://www.2captcha.com/)**
- **[Anticaptcha](https://www.anti-captcha.com/)**
- **[Capsolver](https://capsolver.com/)**
- **[CaptchaAI](https://captchaai.com/)**

### Site-Specific Support:
| Captcha Type            |2Captcha    | Anticaptcha | Capmonster| Capsolver | CaptchaAI|
| :-------------: |:-------------:| :-----:| :-----:| :-----:| :-----:| 
| Recaptcha V2 | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Recaptcha V3 | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Hcaptcha | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Image Captcha | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Cloudflare Turnstile | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x: | :x: |
| Funcaptcha |:x: | :x: | :x: | :x: | :x: |
| GeeTest |:x: | :x: | :x: | :x: | :x: |
| Amazon WAF |:x: | :x: | :x: | :x: | :x: |



# Recommendations
1. For 2Captcha, don't run more than 60 tasks per API key.
2. Handle errors appropriately.
    * If a `ErrNoBalance` is thrown, tasks should stop. Some sites will temporarily ban IP's if constant requests come in.

# Errors
| Errors | Returned When |
| :--------:| :-----:|
| `ErrNoBalance` | Balance is below 0 for captcha solving site|
| `ErrWrongAPIKey` | Incorrect API Key for captcha solving site|
| `ErrWrongSitekey` | Incorrect sitekey |
| `ErrIncorrectCapType` | Incorrectly chose a captcha type. When initializing a new harvester. Refer to [the captcha types](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#config-struct-fields) |
| `ErrNoHarvester` | When the user did not / incorrectly chose a captcha harvester. Refer to the ["how to use" guide](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#how-to-use) |

##### Error Handling
```go
package main

import (
	"fmt"

	captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

func main() {
	solver, err := captchatools.NewHarvester(captchatools.CapmonsterSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
		CaptchaType: captchatools.V2Captcha,
	})
	if err != nil {
		...
	}
	answer, err := solver.GetToken()
	if err != nil {
		switch err {
		case captchatools.ErrBanned:
			...
		case captchatools.ErrAddionalDataMissing:
			...
		}
	}
}

```


# TODOs
- [ ] Add DeadByCaptcha
- [ ] Nocaptchaai (maybe)
- [ ] Context Support
- [ ] FunCaptcha Support
- [ ] Cookie Support