from . import Harvester
import requests
from . import exceptions as captchaExceptions

class Twocap(Harvester):
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
        return "2captcha_token"
    
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