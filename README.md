# Captcha Tools
Python module to help solve captchas with Capmonster, 2Captcha and Anticaptcha API's!

#### Go(lang)
To see documentation for the Go implementation, [click here](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go)

# Install
```python
pip3 install captchatools
```
##### To update
```python
pip3 install -U captchatools
```
# How to use
### Basic usage
```python
import captchatools
solver = captchatools.new_harvester(solving_site="capmonster", api_key="YOUR API KEY", sitekey="6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", captcha_url="https://www.google.com/recaptcha/api2/demo")
captcha_answer = solver.get_token()
```
or
```python
from captchatools import new_harvester
solver = new_harvester(solving_site=1, api_key="YOUR API KEY", sitekey="6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", captcha_url="https://www.google.com/recaptcha/api2/demo")
captcha_answer = solver.get_token()
```
### new_harvester() Parameters:
| Parameter | Required |  Type  | Default | Description|
| :-------------: |:-------------:| :-----:| :-----:| :-----:|
| api_key | true | String| -| The API Key for the captcha solving site|
| solving_site| true| String (name of site) or int (site ID) | "capmonster"| The captcha solving site that will be used. Refer to [the site IDs](https://github.com/Matthew17-21/Captcha-Tools/tree/main/captchatools-go#site-specific-support)|
| sitekey| true | String | - | Sitekey from the site where captcha is loaded|
| captcha_url | true| String | - | URL where the captcha is located|
| captcha_type| false| String | "v2" | Type of captcha you are solving. Either captcha `image`, `v2`, `v3` or `hcaptcha` (`hcap` works aswell)|
| invisible_captcha| false | bool | false | If the captcha is invisible or not.<br />__This param is only required when solving invisible captchas__|
| min_score | false | double |0.7 | Minimum score for v3 captchas.<br />__This param is only required when solving V3 and it needs a higher / lower score__|
| action | false | String | "verify" | Action that is associated with the V3 captcha.<br />__This param is only required when solving V3 captchas__|
| soft_id | false | int | - |2captcha Developer ID. <br /> Developers get 10% of spendings of their software users. |
### get_token() Parameters:
| Field | Required |  Type  |  Description|
| :-------------: |:-------------:|  :-----:| :-----:|
| b64_img | false | string |  Base64 encoded captcha image<br />__This param is only required when solving image captchas__|
| proxy| false | string |  Proxy to be used to solve captchas.<br />This will make the captcha be solved from the proxy ip<br /><br />Format: `ip:port:user:pass` |
| proxy_type | false | string |  Type of the proxy being used. Options are:<br /> `HTTP`, `HTTPS`, `SOCKS4`, `SOCKS5`|
| user_agent | false | string |  UserAgent that will be passed to the service and used to solve the captcha |
### Examples
##### Example - V2 Captcha / Basic usage
```python
from captchatools import new_harvester

def main():
    harvester = new_harvester(
        api_key="CHANGE THIS", 
        solving_site="capsolver",
        captcha_type="v2",
        sitekey="6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
        captcha_url="https://www.google.com/recaptcha/api2/demo"
    )
    answer = harvester.get_token()
```
##### Example - V3 Captcha
```python
from captchatools import new_harvester

def main():
    harvester = new_harvester(
    
        api_key="CHANGE THIS", 
        solving_site="capsolver",
        
        captcha_type="v3",
        sitekey="6LcR_okUAAAAAPYrPe-HK_0RULO1aZM15ENyM-Mf",
        captcha_url="https://antcpt.com/score_detector/",
        action="homepage",
        min_score=0.7
    )
    token = harvester.get_token()
```
##### Example - Image captcha
```python
from captchatools import new_harvester

def main():
    harvester = new_harvester(
        api_key="CHANGE THIS", 
        solving_site="capsolver",
        captcha_type="image",
    )
    token = harvester.get_token(b64_img="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")
```
##### Example - Additional captcha data
```python
from captchatools import new_harvester

def main():
    harvester = new_harvester(
        api_key="CHANGE THIS", 
        solving_site="capsolver",
        captcha_type="v2",
        sitekey="6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
        captcha_url="https://www.google.com/recaptcha/api2/demo"
    )
    
    token = harvester.get_token(
        proxy="ip:port:user:pass", 
        proxy_type="http", 
        user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
    )
```

# Supported Sites
- **[Capmonster](https://capmonster.cloud/)**
- **[2Captcha](https://www.2captcha.com/)**
- **[Anticaptcha](https://www.anti-captcha.com/)**
- **[Capsolver](https://www.capsolver.com/)**
- **[CaptchaAI](https://captchaai.com/)**

### Site-Specific Support:
| Site            | Site ID |Captcha Types  Supported    |  Task Types Supported|
| :-------------: |:-------------:| :-----:| :-----:|
| Capmonster      | captchatools.CapmonsterSite| Image captchas,<br/> Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha | ImageToTextTask,<br/>NoCaptchaTask,<br/> NoCaptchaTaskProxyless,<br/> RecaptchaV3TaskProxyless,<br />HCaptchaTaskProxyless |
| Anticaptcha     | captchatools.AnticaptchaSite| Image captchas,<br/> Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha      |    ImageToTextTask,<br/> RecaptchaV2Task<br/>  RecaptchaV2TaskProxyless,<br />RecaptchaV3TaskProxyless,<br />HCaptchaTaskProxyless |
| 2Captcha        | captchatools.TwoCaptchaSite| Image captchas,<br/> Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha      |   - |
| Capsolver        | captchatools.CapsolverSite| Image captchas,<br/> Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha      |   - |
| CaptchaAI        | captchatools.CaptchaAISite| Image captchas,<br/> Recaptcha V2,<br />Recaptcha V3,<br />HCaptcha      |   - |



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
```python
from captchatools import new_harvester, exceptions as captchaExceptions,

def main():
    try:
        harvester = new_harvester()
        token = harvester.get_token()
    except captchaExceptions.NoHarvesterException:
        print("I need to set my captcha harvester!")
```



# Changelog
### 1.4.1
##### What's new
1. Added CaptchaAI
2. Removed internal redundant code
3. Fix creating a new harvester if pass in the ID
### 1.3.0
##### What's new
1. Get Balance Support
2. Proxy Support
3. User Agent Support
4. Text image captcha support
5. Better internal handling
6. Capsolver support 

##### Important Changes
* It is recommend to use the `new_harvester` function rather than the old `captcha_harvesters`