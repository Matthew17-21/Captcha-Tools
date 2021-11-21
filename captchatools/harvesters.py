from .twocap import Twocap
from .anticaptcha import Anticap
from .capmonster import Capmonster
from .exceptions import NoHarvesterException, FailedToGetCapIMG
import requests
import base64
class captcha_harvesters:
    def __init__(   self, solving_site=1, 
                    api_key="Key the site gave you to solve captchas", 
                    captcha_type="v2", invisible_captcha=False, 
                    sitekey="Sitekey of the page where the captcha is loaded",
                    captcha_url="The site URL of where the captcha is loaded", 
                    action="verify", min_score=0.7, soft_id=4782723):
        self.api_key = api_key
        self.solving_site = solving_site
        self.captcha_type = captcha_type.lower()
        self.invisible_captcha = invisible_captcha
        self.captcha_url = captcha_url
        self.min_score = min_score
        self.sitekey = sitekey
        self.action = action
        self.soft_id = soft_id

        # Get Token from Capmonster API
        if self.solving_site == 1 or str(self.solving_site).lower() == "capmonster":
            self.harvester = Capmonster(self)
        
        # Get Token from Anticaptcha API
        elif self.solving_site == 2 or str(self.solving_site).lower() == "anticaptcha":
            self.harvester = Anticap(self)
        
        #Get Token from 2captcha API
        elif self.solving_site == 3 or str(self.solving_site).lower() == "2captcha":
            self.harvester = Twocap(self)
        else:
            raise NoHarvesterException(
                "No captcha harvester was selected. Double check you entered the site name correctly "
                +
                "|| Double check the site id is type int"
            )
    
    def get_token(self):
        '''
        Returns a recaptcha token for the provided details
        '''
        return self.harvester.get_token()
    
    def get_normal(self, path_to_image):
        '''
        This method will handle returning text from 'Normal Captchas.' 
        
        As per 2captcha, 
        "Normal Captcha is an image that contains distored but human-readable text. To solve the captcha user have to type the text from the image."


        Parameters:
        - path_to_image
            - The path where the captcha is located (you must download the image yourself and then pass it in)
        '''
        return self.harvester.get_normal(path_to_image)
    
    @staticmethod
    def get_cap_img(img_url:str) -> str:
        '''
        This function tries 3 times to get the captcha image.
        
        If successful, it returns the base64 encoded image.
        If not successful, raises a FailedToGetCapIMG exception
        '''
        for _ in range(3):
            try:
                response = requests.get(img_url)
                b64response = base64.b64encode(response.content).decode("utf-8")
                return b64response
            except:
                pass
        raise FailedToGetCapIMG()