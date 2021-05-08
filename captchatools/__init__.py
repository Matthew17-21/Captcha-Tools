'''
All-in-one library for solving captchas!
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

How to use:
>>> import captchatools
>>> solver = captchatools.captcha_harvesters(solving_site="capmonster", api_key="YOUR API KEY", sitekey="SITEKEY", captcha_url="https://www.google.com/recaptcha/api2/demo")
>>> captcha_answer = solver.get_token()
or
>>> from captchatools import captcha_harvesters, exceptions
>>> solver = captcha_harvesters(solving_site=1, api_key="YOUR API KEY", sitekey="SITEKEY", captcha_url="https://www.google.com/recaptcha/api2/demo")
>>> captcha_answer = solver.get_token()

Sites:
1 = Capmonster
2 = Anticaptcha
3 = 2captcha
'''
from .harvesters import captcha_harvesters
from .exceptions import(
    WrongAPIKeyException, NoBalanceException, WrongSitekeyException, NoHarvesterException
)