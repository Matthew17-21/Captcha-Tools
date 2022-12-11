from . import Harvester
from . import exceptions as captchaExceptions
import requests
from typing import Optional
import time

BASEURL =  "https://api.anti-captcha.com"

class Anticaptcha(Harvester):
    def get_balance(self) -> float:
        payload = {"clientKey": self.api_key}
        for _ in range(5):
            try:
                resp = requests.post("https://api.anti-captcha.com/getBalance", json=payload ,timeout=20).json()
                if resp["errorId"] == 1: # Means there was an error
                    self.check_error(resp["errorCode"])
                return float(resp["balance"])
            except requests.RequestException:
                pass

    def get_token(self, b64_img: Optional[str] = None, user_agent: Optional[str] = None, proxy: Optional[str] = None, proxy_type: Optional[str] = "HTTP"):
        # Get ID
        task_id = self.__get_id(
            b64_img=b64_img,
            user_agent=user_agent,
            proxy=proxy,
            proxy_type=proxy_type
        )
        
        # Get Answer
        return self.__get_answer(task_id)
    
    def __create_payload(self, **kwargs):
        payload = {"clientKey":self.api_key,"task":{}}

        # Add data based on the captcha type
        if self.captcha_type == "image" or self.captcha_type == "normal":
            payload["task"]["type"] = "ImageToTextTask"
            payload["task"]["body"] = kwargs.get("b64_img", "")
        elif self.captcha_type == "v2":
            payload["task"]["type"] = "RecaptchaV2TaskProxyless"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
            if self.invisible_captcha:
                payload["task"]["isInvisible"] = self.invisible_captcha
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "RecaptchaV2Task"
        elif self.captcha_type == "v3":
            payload["task"]["type"] = "RecaptchaV3TaskProxyless"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
            payload["task"]["minScore"] = self.min_score
            payload["task"]["pageAction"] = self.action
        elif self.captcha_type == "hcap" or self.captcha_type == "hcaptcha":
            payload["task"]["type"] = "HCaptchaTaskProxyless"
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "HCaptchaTask"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey


        # Add Global Data
        if self.soft_id is not None:
            payload["softId"] = self.soft_id
        if kwargs.get("proxy", None) is not None:
            splitted = kwargs.get("proxy").split(":")
            payload["task"]["proxyType"] = kwargs.get('proxy_type',"HTTP")
            payload["task"]["proxyAddress"] = splitted[0]
            try:
                payload["task"]["proxyPort"] = int(splitted[1])
            except Exception:
                payload["task"]["proxyPort"] = splitted[1]
            if len(splitted) >=4:
                payload["task"]["proxyLogin"] = splitted[2]
                payload["task"]["proxyPassword"] = splitted[3]
        if kwargs.get("user_agent", None) is not None:
            payload["task"]["userAgent"] = kwargs.get("user_agent")
        return payload
    
    def __get_id(self,**kwargs):
        # Create Payload
        payload = self.__create_payload(**kwargs)
        
        # Get token & return it
        for _ in range(50):
            try:
                resp = requests.post(BASEURL + "/createTask " , json=payload, timeout=20).json()
                if resp["errorId"] != 0: # Means there was an error:
                    self.check_error(resp["errorCode"])
                return resp["taskId"]
            except (requests.RequestException, KeyError):
                pass
    
    def __get_answer(self,task_id:int):
        payload = {"clientKey":self.api_key,"taskId": task_id}
        for _ in range(100):
            try:
                response = requests.post(BASEURL + "/getTaskResult",json=payload,timeout=20,).json()
                if response["errorId"] != 0: # Means there was an error
                    self.check_error(response["errorId"])
                if response["status"] == "processing":
                    time.sleep(3)
                    continue
                if self.captcha_type == "normal" or self.captcha_type == "image":
                    return response["solution"]["text"]
                else:
                    return response["solution"]["gRecaptchaResponse"]
            except (requests.RequestException, KeyError):
                pass
    
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
        elif error_code == "ERROR_PAGEURL":
            raise captchaExceptions.TaskDetails(f"Error: {error_code}")
        elif error_code == "MAX_USER_TURN" or error_code == "ERROR_NO_SLOT_AVAILABLE":
            raise captchaExceptions.NoSlotAvailable("No slot available")
        elif error_code == "ERROR_IP_NOT_ALLOWED" or error_code == "IP_BANNED" or error_code == "ERROR_IP_BLOCKED":
            return captchaExceptions.Banned(error_code)
        elif error_code == "ERROR_ZERO_CAPTCHA_FILESIZE" or error_code == "ERROR_UPLOAD" or \
            error_code == "ERROR_CAPTCHAIMAGE_BLOCKED" or error_code == "ERROR_IMAGE_TYPE_NOT_SUPPORTED" or \
            error_code == "ERROR_WRONG_FILE_EXTENSION":
            raise captchaExceptions.CaptchaImageError(error_code)
        else: raise captchaExceptions.UnknownError(f"Error returned from 2captcha: {error_code}")
