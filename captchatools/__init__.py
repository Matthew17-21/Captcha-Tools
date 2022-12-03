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
__version__ = "2.0.0"
__author__ = "Matthew17-21"
__license__ = "MIT"

from .harvester import Harvester
from .twocap import Twocap
from .exceptions import (
    NoHarvesterException
)

def new_harvester(**kwargs) -> Harvester:
    '''
    Returns a new captcha harvester with the specified config
    '''
    site = kwargs.get("solving_site",None)
    if site is None:
        raise NoHarvesterException(
            "No captcha harvester was selected. Double check you entered the site name correctly "
            +
            "|| Double check the site id is type int"
        )

    if site == 3 or str(site).lower() == "2captcha":
        return Twocap(**kwargs)
