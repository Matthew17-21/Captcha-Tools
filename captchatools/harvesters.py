from .twocap import Twocap
from .anticaptcha import Anticap
from .capmonster import Capmonster

class captcha_harvesters:
    def __init__(   self, solving_site=1, 
                    api_key="Key the site gave you to solve captchas", 
                    captcha_type="v2", invisible_captcha=False, 
                    sitekey="Sitekey of the page where the captcha is loaded",
                    captcha_url="The site URL of where the captcha is loaded", 
                    action="verify", min_score=0.7):
        self.api_key = api_key
        self.solving_site = solving_site
        self.captcha_type = captcha_type.lower()
        self.invisible_captcha = invisible_captcha
        self.captcha_url = captcha_url
        self.min_score = min_score
        self.sitekey = sitekey
        self.action = action

        # Get Token from Capmonster API
        if self.solving_site == 1 or str(self.solving_site).lower() == "capmonster":
            self.harvester = Capmonster(self)
        
        # Get Token from Anticaptcha API
        elif self.solving_site == 2 or str(self.solving_site).lower() == "anticaptcha":
            self.harvester = Anticap(self)
        
        #Get Token from 2captcha API
        elif self.solving_site == 3 or str(self.solving_site).lower() == "2captcha":
            self.harvester = Twocap(self)
    
    def get_token(self):
        '''
        Returns a captcha token for the provided details
        '''
        return self.harvester.get_token()
