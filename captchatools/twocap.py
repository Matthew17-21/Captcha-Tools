from . import Harvester
import requests
from .errors import ErrCodeToException
from typing import Optional
import time

BASEURL_IN = "http://2captcha.com/in.php"

class Twocap(Harvester):
    def get_balance(self) -> float:
        url = f"https://2captcha.com/res.php?key={self.api_key}&action=getbalance&json=1"
        for _ in range(5):
            try:
                resp = requests.get(url, timeout=20).json()
                if resp["status"] == 0: # Means there was an error
                    ErrCodeToException("2Captcha",resp["request"])
                return float(resp["request"])
            except requests.RequestException:
                pass

    def get_token(
            self, 
            b64_img: Optional[str] = None, 
            user_agent: Optional[str] = None, 
            proxy: Optional[str] = None, 
            proxy_type: Optional[str] = None,
            rq_data: Optional[str] = None
        ):
        # Get ID
        task_id = self.__get_id(
            b64_img=b64_img,
            user_agent=user_agent,
            proxy=proxy,
            proxy_type=proxy_type,
            rq_data=rq_data
        )
        
        # Get Answer
        return self.__get_answer(task_id)
    
    def __create_payload(self, **kwargs):
        payload = {"key": self.api_key,"json": 1}

        # Add data based on the captcha type
        if self.captcha_type == "image" or self.captcha_type == "normal":
            payload["method"] = "base64"
            payload["body"] = kwargs.get("b64_img", "")
        elif self.captcha_type == "v2":
            payload["method"] = "userrecaptcha"
            payload["googlekey"] = self.sitekey
            payload["pageurl"] = self.captcha_url
            if self.invisible_captcha:
                payload["invisible"] = 1
        elif self.captcha_type == "v3":
            payload["method"] = "userrecaptcha"
            payload["version"] = "v3"
            payload["action"] = self.action
            payload["googlekey"] = self.sitekey
            payload["pageurl"] = self.captcha_url
            if self.min_score is not None:
                payload["min_score"] = self.min_score
        elif self.captcha_type == "hcap" or self.captcha_type == "hcaptcha":
            payload["method"] = "hcaptcha"
            payload["sitekey"] = self.sitekey
            payload["pageurl"] = self.captcha_url

        # Add Global Data
        if self.soft_id is not None:
            payload["soft_id"] = self.soft_id
        if kwargs.get("proxy", None) is not None:
            payload["proxy"] = kwargs.get("proxy")
            pxy_type = kwargs.get("proxy_type", "http")
            payload["proxytype"] = pxy_type   
        if kwargs.get("user_agent", None) is not None:
            payload["userAgent"] = kwargs.get("user_agent")
        if kwargs.get("rq_data", None) is not None:
            payload["data"] = kwargs.get("rq_data") 
        return payload
    
    def __get_id(self,**kwargs):
        # Create Payload
        payload = self.__create_payload(**kwargs)
        
        # Get token & return it
        for _ in range(50):
            try:
                resp = requests.post(BASEURL_IN, json=payload, timeout=20).json()
                if resp["status"] == 0: # Means there was an error:
                    ErrCodeToException("2Captcha",resp["request"])
                return resp["request"]
            except (requests.RequestException, KeyError):
                pass
    
    def __get_answer(self,task_id:int):
        for _ in range(100):
            try:
                response = requests.get(f"http://2captcha.com/res.php?key={self.api_key}&action=get&id={task_id}&json=1",timeout=20,).json()
                if response["status"] == 0 and response["request"] != "CAPCHA_NOT_READY": # Error checking
                    ErrCodeToException("2Captcha",response["request"])
                if response["status"] == 0 and response["request"] == "CAPCHA_NOT_READY":
                    time.sleep(4)
                    continue
                return response["request"] # Return the captcha token
            except (requests.RequestException, KeyError):
                pass