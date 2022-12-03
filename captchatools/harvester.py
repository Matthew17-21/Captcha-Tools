from abc import ABC, abstractmethod

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
        if self.captcha_type not in ["v2", "v3", "hcaptcha", "hcap"]:
            raise Exception("Invalid captcha type")
        if self.soft_id is None: # Set with my own soft_id
            pass
    
    @abstractmethod
    def get_balance(self) -> float:
        '''
        Returns the balance for the current captcha harvster
        '''
    
    @abstractmethod
    def get_token(self):
        '''
        Returns a captcha token
        '''