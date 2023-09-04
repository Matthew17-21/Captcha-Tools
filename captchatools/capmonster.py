from .errors import ErrCodeToException
from . import Harvester
import time
import requests
from typing import Optional

BASEURL = "https://api.capmonster.cloud"



class Capmonster(Harvester):
    '''
    This object will contain the data to interact with capmonster.cloud API
    '''
    def get_balance(self) -> float:
        payload = {"clientKey": self.api_key}
        for _ in range(5):
            try:
                resp = requests.post(BASEURL + "/getBalance", json=payload ,timeout=20).json()
                if resp["errorId"] == 1: # Means there was an error
                    ErrCodeToException("Capmonster", resp["errorCode"])
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
            payload["task"]["type"] = "NoCaptchaTaskProxyless"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
            if self.invisible_captcha:
                payload["task"]["isInvisible"] = self.invisible_captcha
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "NoCaptchaTask"
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
            payload["task"]["proxyAddress"] = splitted[0]
            try:
                payload["task"]["proxyPort"] = int(splitted[1])
            except Exception:
                payload["task"]["proxyPort"] = splitted[1]
            if len(splitted) >=4:
                payload["task"]["proxyLogin"] = splitted[2]
                payload["task"]["proxyPassword"] = splitted[3]
            payload["task"]["proxyType"] = kwargs.get("proxy_type", "http")
        if kwargs.get("user_agent", None) is not None:
            payload["task"]["userAgent"] = kwargs.get("user_agent")
        return payload
    
    def __get_id(self,**kwargs):
        # Create Payload
        payload = self.__create_payload(**kwargs)
        
        # Get token & return it
        for _ in range(50):
            try:
                resp = requests.post(BASEURL + "/createTask" , json=payload, timeout=20).json()
                if resp["errorId"] != 0: # Means there was an error:
                    ErrCodeToException("Capmonster", resp["errorCode"])
                return resp["taskId"]
            except (requests.RequestException, KeyError):
                pass
    
    def __get_answer(self,task_id:int):
        payload = {"clientKey":self.api_key,"taskId": task_id}
        for _ in range(100):
            try:
                response = requests.post(BASEURL + "/getTaskResult",json=payload,timeout=20,).json()
                if response["errorId"] != 0: # Means there was an error
                    ErrCodeToException("Capmonster", response["errorCode"])
                if response["status"] == "processing":
                    time.sleep(3)
                    continue
                if self.captcha_type == "normal" or self.captcha_type == "image":
                    return response["solution"]["text"]
                else:
                    return response["solution"]["gRecaptchaResponse"]
            except (requests.RequestException, KeyError):
                pass