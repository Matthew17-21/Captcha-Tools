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

from abc import ABC, abstractmethod
from typing import Optional
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
            raise Exception("No API Key!")
        if self.solving_site is None:
            raise Exception("No solving site set")
        if self.captcha_type not in ["v2", "v3", "hcaptcha", "hcap", "image", "normal"]:
            raise Exception("Invalid captcha type")
        if self.soft_id is None: #TODO Set with my own soft_id
            pass
    
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
        proxy_type: Optional[str]=None
    ):
        '''
        Returns a captcha token
        '''

def new_harvester(**kwargs) -> Harvester:
    from .twocap import Twocap
    return Twocap(**kwargs)