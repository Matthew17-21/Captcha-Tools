import requests
from . import exceptions as captchaExceptions
from . import Harvester
import time

BASEURL =  "http://2captcha.com/in.php"

# There's others way we could've done this, 
# but it's sufficient for right now.
class Twocap(Harvester):
    '''
    This object will contain the data to interact with 2captcha.com API
    '''    
    def get_balance(self) -> float:
        url = f"https://2captcha.com/res.php?key={self.api_key}&action=getbalance&json=1"
        for _ in range(5):
            try:
                resp = requests.get(url).json()
                if resp["status"] == 0: # Means there was an error
                    self.check_error(resp["request"])
                return float(resp["request"])
            except requests.RequestException:
                pass
    
    def get_token(self):
        return "token from 2captcha"
    
    @staticmethod
    def check_error(error_code):
        if error_code == "ERROR_ZERO_BALANCE":
            raise captchaExceptions.NoBalanceException()
        elif error_code == "ERROR_WRONG_GOOGLEKEY":
            raise captchaExceptions.WrongSitekeyException() 
        elif error_code == "ERROR_WRONG_USER_KEY" or error_code == "ERROR_KEY_DOES_NOT_EXIST":
            raise captchaExceptions.WrongAPIKeyException()
        elif error_code == "ERROR_TOO_BIG_CAPTCHA_FILESIZE":
            raise captchaExceptions.CaptchaIMGTooBig()
        else: raise Exception(f"Error returned from 2captcha: {error_code}")

    # def get_token(self) -> str:        
    #     return self.get_answer( self.get_id()  )

    # def get_normal(self, cap_pic_url) -> str:
    #     return self.get_answer( self.get_id(cap_pic_url)  )

    # def get_id(self, cap_pic_url=None) -> int:
    #     payload = {
    #         "key": self.user_data.api_key,
    #         "method": "userrecaptcha", # Because V2 recaptcha is defaulted on init, I'll leave this
    #         "googlekey": self.user_data.sitekey,
    #         "pageurl":self.user_data.captcha_url,
    #         "json": 1
    #     }
    #     if self.user_data.soft_id is not None:
    #         payload["soft_id"] = self.user_data.soft_id
    #     if self.user_data.captcha_type == "v2":
    #         if self.user_data.invisible_captcha:
    #             payload["invisible"] = 1
        
    #     elif self.user_data.captcha_type == "v3":
    #         payload["version"] = "v3"
    #         payload["action"] = self.user_data.action
    #         payload["min_score"] = self.user_data.min_score
        
    #     elif self.user_data.captcha_type == "hcap" or self.user_data.captcha_type == "hcaptcha":
    #         payload["method"] = "hcaptcha"
    #         # We need to remove the "googlekey" ket from the payload
    #         # And replace it with "sitekey"
    #         payload.pop("googlekey")
    #         payload["sitekey"] = self.user_data.sitekey

    #     elif self.user_data.captcha_type == "normal":
    #         payload["method"] = "base64"
    #         payload["body"] = self.user_data.get_cap_img(cap_pic_url)
    #     while True:
    #         try:
    #             resp = requests.post(BASEURL, data=payload).json()
    #             if resp["status"] == 1:
    #                 return resp["request"] # Return the queue ID

    #             elif resp["request"] == "ERROR_ZERO_BALANCE":
    #                 # Throw Exception
    #                 raise captchaExceptions.NoBalanceException()

    #             elif resp["request"] == "ERROR_WRONG_GOOGLEKEY":
    #                 # Throw Exception
    #                 raise captchaExceptions.WrongSitekeyException()
                
    #             elif resp["request"] == "ERROR_WRONG_USER_KEY" or resp["request"] == "ERROR_KEY_DOES_NOT_EXIST":
    #                 # Throw Exception
    #                 raise captchaExceptions.WrongAPIKeyException()
                
    #             elif resp["request"] == "ERROR_TOO_BIG_CAPTCHA_FILESIZE":
    #                 raise captchaExceptions.CaptchaIMGTooBig()
    #             break

    #         except requests.RequestException:
    #             pass

    # def get_answer(self, queue_id) -> str:
    #     '''
    #     This method gets the captcha token from the API
    #     '''
    #     while True:
    #         try:

    #             answer = requests.get(f"http://2captcha.com/res.php?key={self.user_data.api_key}&action=get&id={queue_id}&json=1",timeout=10,).json()
    #             if answer["status"] == 1: # Solved
    #                 return answer["request"]
    #             elif answer["request"] == "ERROR_CAPTCHA_UNSOLVABLE":
    #                 self.get_token()
    #             time.sleep(4)
            
    #         except KeyError:
    #             self.get_token()
    #         except Exception:
    #             pass