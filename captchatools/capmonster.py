import requests
from . import exceptions as captchaExceptions
import time

BASEURL = " https://api.capmonster.cloud"


# There's others way we could've done this, 
# but it's sufficient for right now.
class Capmonster:
    '''
    This object will contain the data to interact with capmonster.cloud API
    '''
    def __init__(self, parent:object):
        self.user_data = parent

    def get_token(self) -> str:        
        return self.get_answer( self.get_id()  )

    def get_id(self) -> int:
        '''
        Method to get Queue ID from the API.
        '''

        # Define the payload we are going to send to the API
        payload = {
            "clientKey":self.user_data.api_key,
            "task":{
                "websiteURL":self.user_data.captcha_url,
                "websiteKey": self.user_data.sitekey
            }
        }

        # Add anything else to the API, based off the user's input
        if self.user_data.captcha_type == "v2":
            payload["task"]["type"] = "NoCaptchaTaskProxyless"
            if self.user_data.invisible_captcha:
                payload["task"]["isInvisible"] = True
        
        elif self.user_data.captcha_type == "v3":
            payload["task"]["type"] = "RecaptchaV3TaskProxyless"
            payload["task"]["minScore"] = self.user_data.min_score
            payload["task"]["pageAction"] = self.user_data.action

        # Get the Queue ID b sending a POST request to their API
        while True:
            try:
                resp = requests.post(BASEURL + "/createTask", json=payload).json()
                if resp["errorId"] == 0: # Means there was no error
                    return resp["taskId"] # Return the queue ID

                elif resp["errorCode"] == "ERROR_ZERO_BALANCE":
                    # Throw Exception
                    raise captchaExceptions.NoBalanceException()

                elif resp["errorCode"] == "ERROR_RECAPTCHA_INVALID_SITEKEY":
                    # Throw Exception
                    raise captchaExceptions.WrongSitekeyException()
                
                elif resp["errorCode"] == "ERROR_KEY_DOES_NOT_EXIST":
                    # Throw Exception
                    raise captchaExceptions.WrongAPIKeyException()

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
                payload_getAnswer = {
                    "clientKey":self.user_data.api_key,
                    "taskId": queue_id
                }
            
                answer = requests.post(BASEURL + "/getTaskResult", json=payload_getAnswer).json()
                if answer["status"] == "ready":
                    return answer["solution"]["gRecaptchaResponse"]
                
                elif answer["errorId"] == 12 or answer["errorId"] == 16:
                    # Captcha unsolvable || TaskID doesn't exist
                    self.get_token()


                time.sleep(4)
            except KeyError:
                self.get_token()
            except Exception:
                pass