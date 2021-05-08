# Captcha Tools
Python module to help solve captchas with Capmonster, 2Captcha and Anticaptcha API's!

# Install
```python
pip3 install captchatools
```
##### To update
```python
pip3 install -U captchatools
```

# How to use
```python
import captchatools
solver = captchatools.captcha_harvesters(solving_site="capmonster", api_key="YOUR API KEY", sitekey="SITEKEY", captcha_url="https://www.google.com/recaptcha/api2/demo")
captcha_answer = solver.get_token()
```
or
```python
from captchatools import captcha_harvesters, exceptions
solver = captcha_harvesters(solving_site=1, api_key="YOUR API KEY", sitekey="SITEKEY", captcha_url="https://www.google.com/recaptcha/api2/demo")
captcha_answer = solver.get_token()
```

| Parameter | Required |  Type  | Default | Description|
| :-------------: |:-------------:| :-----:| :-----:| :-----:|
| api_key | true | String| -| The API Key for the captcha solving site|
| solving_site| true| String (name of site) or int (site ID) | "capmonster"| Captcha solving site|
| sitekey| true | String | - | Sitekey from the site where captcha is loaded|
| captcha_url | true| String | - | URL where the captcha is located|
| captcha_type| false| String | "v2" | Type of captcha you are solving. Either captcha `v2`, `v3` or `hcaptcha` (`hcap` works aswell)|
| invisible_captcha| false | bool | false | If the captcha is invisible or not.<br />__This param is only required when solving invisible captchas__|
| min_score | false | double |0.7 | Minimum score for v3 captchas.<br />__This param is only required when solving V3 and it needs a higher / lower score__|
| action | false | String | "verify" | Action that is associated with the V3 captcha.<br />__This param is only required when solving V3 captchas__|


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
2. Handle exceptions appropriately.
    * If a `NoBalanceException` is thrown, tasks should stop. Some sites will temporarily ban IP's if constant requests come in.

# Exceptions
| Exception | Raised |
| :--------:| :-----:|
| `NoBalanceException` | Balance is below 0 for captcha solving site|
| `WrongAPIKeyExceptionException` | Incorrect API Key for captcha solving site|
| `WrongSitekeyException` | Incorrect sitekey |
| `NoHarvesterException` | When the user did not / incorrectly chose a captcha harvester. Refer to the [guide](https://github.com/Matthew17-21/Captcha-Tools#how-to-use) |

```python
from captchatools import captcha_harvesters, exceptions as captchaExceptions
try:
    ...
except captchaExceptions.NoBalanceException:
    print("No balance.")
```
or
```python
import captchatools
try:
    ...
except captchatools.NoBalanceException:
    print("No balance.")
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
6. [] Release in Go
7. [] Allow for refunds