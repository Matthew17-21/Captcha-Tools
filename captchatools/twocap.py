import requests
from . import exceptions as captchaExceptions
import time

BASEURL =  "http://2captcha.com/in.php"

# There's others way we could've done this, 
# but it's sufficient for right now.
class Twocap:
    '''
    This object will contain the data to interact with 2captcha.com API
    '''
    def __init__(self, parent:object):
        self.user_data = parent

    def get_token(self) -> str:        
        return self.get_answer( self.get_id()  )

    def get_id(self) -> int:
        payload = {
            "key": self.user_data.api_key,
            "method": "userrecaptcha",
            "googlekey": self.user_data.sitekey,
            "pageurl":self.user_data.captcha_url,
            "json": 1
        }

        if self.user_data.captcha_type == "v2":
            if self.user_data.invisible_captcha:
                payload["invisible"] = 1
        
        elif self.user_data.captcha_type == "v3":
            payload["version"] = "v3"
            payload["action"] = self.user_data.action
            payload["min_score"] = self.user_data.min_score

        while True:
            try:
                resp = requests.post(BASEURL, json=payload).json()
                if resp["status"] == 1:
                    return resp["request"] # Return the queue ID

                elif resp["request"] == "ERROR_ZERO_BALANCE":
                    # Throw Exception
                    raise captchaExceptions.NoBalanceException()

                elif resp["request"] == "ERROR_WRONG_GOOGLEKEY":
                    # Throw Exception
                    raise captchaExceptions.WrongSitekeyException()
                
                elif resp["request"] == "ERROR_WRONG_USER_KEY" or resp["request"] == "ERROR_KEY_DOES_NOT_EXIST":
                    # Throw Exception
                    raise captchaExceptions.WrongAPIKeyException()
                break

            except captchaExceptions.NoBalanceException:
                raise captchaExceptions.NoBalanceException()
            except captchaExceptions.WrongSitekeyException:
                raise captchaExceptions.WrongSitekeyException()
            except captchaExceptions.WrongAPIKeyException:
                raise captchaExceptions.WrongAPIKeyException()

            except Exception:
                pass

    def get_answer(self, queue_id) -> str:
        '''
        This method gets the captcha token from the API
        '''
        while True:
            try:

                answer = requests.get(f"http://2captcha.com/res.php?key={self.user_data.api_key}&action=get&id={queue_id}&json=1",timeout=10,).json()
                if answer["status"] == 1: # Solved
                    return answer["request"]
                elif answer["request"] == "ERROR_CAPTCHA_UNSOLVABLE":
                    self.get_token()
                time.sleep(4)
            
            except KeyError:
                self.get_token()
            except Exception:
                pass