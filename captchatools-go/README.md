# Captcha Tools (Go)
Go package to help solve captchas with Capmonster, 2Captcha and Anticaptcha API's!

# Install
```go
go get github.com/Matthew17-21/Captcha-Tools/captchatools-go
```
##### To update
```go
go get -u github.com/Matthew17-21/Captcha-Tools/captchatools-go
```

# How to use
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
		CaptchaType: "V2",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(solver.GetToken())
}

```
V3 Captcha Exmaple:
```go
func v3Example() {
	solver, err := captchatools.NewHarvester(captchatools.AnticaptchaSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "..........",
		CaptchaType: "V3",
		Action:      "submit",
		MinScore:    0.9,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(solver.GetToken())
}

```
### captchatools.NewHarvester() Parameters:
| Parameter | Required |  Type  | Default | Description|
| :-------------: |:-------------:| :-----:| :-----:| :-----:|
| solving_site | true | int| -| The captcha solving site that will be used. Refer to [the site IDs](https://github.com/Matthew17-21/Captcha-Tools/captchatools-go#site-specific-support). Alternatively, you can use shortcuts such as `captchatools.AnticaptchaSite` |
| Config| true | captchatools.Config | - | Configurations for the captchas you are solving. |


### Config struct fields:
| Field | Required |  Type  | Default | Description|
| :-------------: |:-------------:| :-----:| :-----:| :-----:|
| Api_key | true | String| -| The API Key for the captcha solving site|
| Sitekey| true | String | - | Sitekey from the site where captcha is loaded|
| CaptchaURL | true| String | - | URL where the captcha is located|
| CaptchaType| true| String | - | Type of captcha you are solving. Either captcha `v2`, `v3` or `hcaptcha` (`hcap` works aswell)|
| Action | false | String | - | Action that is associated with the V3 captcha.<br />__This param is only required when solving V3 captchas__|
| IsInvisibleCaptcha| false | bool | - | If the captcha is invisible or not.<br />__This param is only required when solving invisible captchas__|
| MinScore | false | float32 | - | Minimum score for v3 captchas.<br />__This param is only required when solving V3 and it needs a higher / lower score__|

<!-- | solving_site| true| String (name of site) or int (site ID) | "capmonster"| Captcha solving site| -->

# Supported Sites
- **[Capmonster](https://capmonster.cloud/)**
- **[2Captcha](https://www.2captcha.com/)**
- **[Anticaptcha](https://www.anti-captcha.com/)**

##### Site-Specific Support:
| Site            |Site ID| Captcha Types  Supported    |  Task Types Supported|
| :-------------: |:-------------:|:-------------:| :-----:|
| Capmonster      |1| Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha | RecaptchaV2TaskProxyless,<br />RecaptchaV3TaskProxyless,<br />HCaptchaTaskProxyless |
| Anticaptcha     |2| Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha      |    RecaptchaV2TaskProxyless,<br />RecaptchaV3TaskProxyless,<br />HCaptchaTaskProxyless |
| 2Captcha        |3| Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha      |   - |


# Recommendations
1. For 2Captcha, don't run more than 60 tasks per API key.
2. Handle errors appropriately.
    * If a `ErrNoBalance` is thrown, tasks should stop. Some sites will temporarily ban IP's if constant requests come in.

# Errors
| Errors | Raised |
| :--------:| :-----:|
| `ErrNoBalance` | Balance is below 0 for captcha solving site|
| `ErrWrongAPIKey` | Incorrect API Key for captcha solving site|
| `ErrWrongSitekey` | Incorrect sitekey |
| `ErrIncorrectCapType` | Incorrectly chose a captcha type. When initializing a new harvester. Refer to [the captcha types](https://github.com/Matthew17-21/Captcha-Tools/captchatools-go#how-to-use) |
| `ErrNoHarvester` | When the user did not / incorrectly chose a captcha harvester. Refer to the [guide](https://github.com/Matthew17-21/Captcha-Tools/captchatools-go#how-to-use) |

```go
package main

import (
	"fmt"

	captchatools "github.com/Matthew17-21/Captcha-Tools/captchatools-go"
)

func main() {
	solver, err := captchatools.NewHarvester(captchatools.AnticaptchaSite, &captchatools.Config{
		Api_key:     "ENTER YOUR API KEY HERE",
		Sitekey:     "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		CaptchaURL:  "https://www.google.com/recaptcha/api2/demo",
		CaptchaType: "V2",
	})
	if err != nil {
		switch err {
		case captchatools.ErrNoBalance:
			fmt.Println("No balance.")
			panic(err)
		}
	}
	fmt.Println(solver.GetToken())
}

```


# TO DO
1. [] Document code better
2. [] 2Captcha
    * [] Clean up code
    * [] Proxy support
    * [] Cookie support
    * [] User Agent Support
    * [] Different type of captchas
3. [] Anticaptcha
    * [] Clean up code
    * [] Proxy support
    * [] Cookie support
    * [] User Agent Support
    * [] Different type of captchas
4. [] Capmonster
    * [] Clean up code
    * [] Proxy support
    * [] Cookie support
    * [] User Agent Support
    * [] Different type of captchas
5. [] Add DeathByCaptcha
6. [] Allow for refunds