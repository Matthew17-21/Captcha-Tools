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
4 = Capsolver
5 = CaptchaAI
'''
__version__ = "1.4.1"
__author__ = "Matthew17-21"
__license__ = "MIT"

from abc import ABC, abstractmethod
from typing import Optional
from . import exceptions as captchaExceptions
from warnings import warn
class Harvester(ABC):
    '''
    Represents a captcha harvester.
    '''
    def __init__(self, **kwargs) -> None:
        self.api_key = kwargs.get("api_key",None)
        self.captcha_type = kwargs.get("captcha_type","v2").lower()
        self.solving_site = kwargs.get("solving_site",None)
        self.invisible_captcha = kwargs.get("invisible_captcha",False)
        self.captcha_url = kwargs.get("captcha_url",None)
        self.min_score = kwargs.get("min_score",0.7)
        self.sitekey = kwargs.get("sitekey",None)
        self.action = kwargs.get("action","verify")
        self.soft_id = kwargs.get("soft_id",None)

        # Validate Data
        if self.api_key is None:
            raise captchaExceptions.WrongAPIKeyException()
        if self.solving_site is None:
            raise captchaExceptions.NoHarvesterException("No solving site selected")
        if self.captcha_type not in ["v2", "v3", "hcaptcha", "hcap", "image", "normal", "hcaptchaturbo"]:
            raise captchaExceptions.NoCaptchaType("Invalid captcha type")
        if self.soft_id is None:
            if self.solving_site == 3 or self.solving_site == "2captcha":
                self.soft_id = 4782723
    
    @abstractmethod
    def get_balance(self) -> float:
        '''
        Returns the balance for the current captcha harvster
        '''
    
    @abstractmethod
    def get_token(
        self, b64_img: Optional[str]=None, 
        user_agent: Optional[str]=None, 
        proxy: Optional[str]=None, 
        proxy_type: Optional[str]=None,
        rq_data: Optional[str] = None
    ):
        '''
        Returns a captcha token
        '''

def new_harvester(**kwargs) -> Harvester:
    # Need to import here to prevent circular imports
    from .twocap import Twocap
    from .anticaptcha import Anticaptcha
    from .capmonster import Capmonster
    from .capsolver import Capsolver
    from .captchaai import CaptchaAI
    
    site = kwargs.get("solving_site","")
    if isinstance(site, str):
        site = site.lower()
    if site == 1 or site == "capmonster":
        return Capmonster(**kwargs)
    elif site == 2 or site == "anticaptcha":
        return Anticaptcha(**kwargs)
    elif site == 3 or site == "2captcha":
        return Twocap(**kwargs)
    elif site == 4 or site == "capsolver":
        return Capsolver(**kwargs)
    elif site == 5 or site == "captchaai":
        return CaptchaAI(**kwargs)
    raise captchaExceptions.NoHarvesterException("No solving site selected")


# Just for backward compatibility
def captcha_harvesters(**kwargs) -> Harvester:
    warn('This function is deprecated. Use the `new_harvester()` function', DeprecationWarning, stacklevel=2)
    return new_harvester(**kwargs)