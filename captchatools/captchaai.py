from . import Harvester
import requests
from .errors import ErrCodeToException
from typing import Optional
import time

BASEURL = "https://ocr.captchaai.com/in.php"

class CaptchaAI(Harvester):
    def get_balance(self) -> float:
        url = f"https://ocr.captchaai.com/res.php?key={self.api_key}&action=getbalance&json=1"
        for _ in range(5):
            try:
                resp = requests.get(url, timeout=20).json()
                if resp["status"] == 0: # Means there was an error
                    ErrCodeToException("CaptchaAI",resp["request"])
                return float(resp["request"])
            except requests.RequestException:
                pass

    def get_token(self, b64_img: Optional[str] = None, user_agent: Optional[str] = None, proxy: Optional[str] = None, proxy_type: Optional[str] = None):
        # Get ID
        task_id = self.__get_id(
            b64_img=b64_img,
            user_agent=user_agent,
            proxy=proxy,
            proxy_type=proxy_type
        )
        
        # Get Answer
        return self.__get_answer(task_id)
    
    def __create_uri_parms(self, **kwargs):
        '''Creats the URI params to be used for submitting captcha info'''
        params = {
            "key": self.api_key,
            "json": 1,
            "pageurl": self.captcha_url,
        }
        if self.captcha_type == "image" or self.captcha_type == "normal":
            params["method"] = "base64"
            params["body"] = kwargs.get("b64_img", "")
        elif self.captcha_type == "v2":
            params["method"] = "userrecaptcha"
            params["googlekey"] = self.sitekey
            if self.invisible_captcha:
                params["invisible"] = 1
        elif self.captcha_type == "v3":
            params["method"] = "userrecaptcha"
            params["version"] = "v3"
            params["googlekey"] = self.sitekey
            params["pageurl"] = self.captcha_url
            if self.action != "":
                params["action"] = self.action
            if self.min_score is not None:
                params["min_score"] = self.min_score
        elif self.captcha_type == "hcap" or self.captcha_type == "hcaptcha":
            params["method"] = "hcaptcha"
            params["sitekey"] = self.sitekey

        # Add Global Data
        if kwargs.get("proxy", None) is not None:
            params["proxy"] = kwargs.get("proxy")
            pxy_type = kwargs.get("proxy_type", "http")
            params["proxytype"] = pxy_type   
        if kwargs.get("user_agent", None) is not None:
            params["userAgent"] = kwargs.get("user_agent")
        return params
    
    def __get_id(self,**kwargs):
        # Create Payload
        params = self.__create_uri_parms(**kwargs)
        
        # Get token & return it
        for _ in range(50):
            try:
                resp = requests.get(BASEURL, params=params ,timeout=20).json()
                if resp["status"] == 0: # Means there was an error:
                    ErrCodeToException("CaptchaAI",resp["request"])
                return resp["request"]
            except (requests.RequestException, KeyError):
                pass
    
    def __get_answer(self,task_id:int):
        for _ in range(100):
            try:
                response = requests.get(f"https://ocr.captchaai.com/res.php?key={self.api_key}&action=get&id={task_id}&json=1",timeout=20,).json()
                if response["status"] == 0 and response["request"] != "CAPCHA_NOT_READY": # Error checking
                    ErrCodeToException("CaptchaAI",response["request"])
                if response["status"] == 0 and response["request"] == "CAPCHA_NOT_READY":
                    time.sleep(4)
                    continue
                return response["request"] # Return the captcha token
            except (requests.RequestException, KeyError):
                pass
    