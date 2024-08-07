from .errors import ErrCodeToException
from . import Harvester
import time
import requests
from typing import Optional

BASEURL = "https://api.capsolver.com"



class Capsolver(Harvester):
    '''
    This object will contain the data to interact with Capsolver API
    '''
    def get_balance(self) -> float:
        payload = {"clientKey": self.api_key}
        for _ in range(5):
            try:
                resp = requests.post(BASEURL + "/getBalance", json=payload ,timeout=20).json()
                if resp["errorId"] == 1: # Means there was an error
                    ErrCodeToException("Capsolver", resp["errorCode"])
                return float(resp["balance"])
            except requests.RequestException:
                pass

    def get_token(
            self, 
            b64_img: Optional[str] = None, 
            user_agent: Optional[str] = None, 
            proxy: Optional[str] = None, 
            proxy_type: Optional[str] = "HTTP",
            rq_data: Optional[str] = None
        ):
        # Get ID
        task_id, already_solved = self.__get_id(
            b64_img=b64_img,
            user_agent=user_agent,
            proxy=proxy,
            proxy_type=proxy_type,
            rq_data=rq_data
        )

        # Check if token was already retrieved
        if already_solved:
            return task_id
        
        # Get Answer
        return self.__get_answer(task_id)
    
    def __create_payload(self, **kwargs):
        payload = {"clientKey":self.api_key,"task":{}}

        # Add data based on the captcha type
        if self.captcha_type == "image" or self.captcha_type == "normal":
            payload["task"]["type"] = "ImageToTextTask"
            payload["task"]["body"] = kwargs.get("b64_img", "")
        elif self.captcha_type == "v2":
            payload["task"]["type"] = "ReCaptchaV2TaskProxyLess"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
            if self.invisible_captcha:
                payload["task"]["isInvisible"] = self.invisible_captcha
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "ReCaptchaV2Task"
        elif self.captcha_type == "v3":
            payload["task"]["type"] = "RecaptchaV3TaskProxyless"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
            payload["task"]["minScore"] = self.min_score
            payload["task"]["pageAction"] = self.action
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "ReCaptchaV3Task"
        elif self.captcha_type == "hcap" or self.captcha_type == "hcaptcha":
            payload["task"]["type"] = "HCaptchaTaskProxyless"
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "HCaptchaTask"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
        elif self.captcha_type == "hcaptchaturbo":
            payload["task"]["type"] = "HCaptchaTurboTask"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
        elif self.captcha_type == "fun" or self.captcha_type == "funcaptcha":
            payload["task"]["type"] = "FunCaptchaTaskProxyless"
            payload["task"]["websiteURL"] = self.captcha_url
            payload["task"]["websiteKey"] = self.sitekey
            if kwargs.get("proxy", None) is not None:
                payload["task"]["type"] = "FunCaptchaTask"

        # Add Global Data
        if self.soft_id is not None:
            payload["appId"] = self.soft_id
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
        if kwargs.get("rq_data", None) is not None:
            payload["task"]["enterprisePayload"] = {"rqdata":  kwargs.get("rq_data")}
        return payload
    
    def __get_id(self,**kwargs):
        '''
        This method gets a task ID for a captcha token.

        A tuple is returned from this method ( Optional[task_id|captcha_token] | bool)
        If the second index is True, that means a captcha token is returned instead of a task ID
        '''
        # Create Payload
        payload = self.__create_payload(**kwargs)
        
        # Get token & return it
        for _ in range(50):
            try:
                resp = requests.post(BASEURL + "/createTask" , json=payload, timeout=20).json()
                if resp["errorId"] != 0: # Means there was an error:
                    ErrCodeToException("Capsolver", resp["errorCode"])
                
                # Check if there is an answer already available
                if resp["status"] == "ready":
                    if resp["solution"].get("text", None) is not None:
                        answer = resp["solution"].get("text", None)
                    elif resp["solution"].get("gRecaptchaResponse", None) is not None:
                        answer = resp["solution"].get("gRecaptchaResponse", None)
                    return (answer, True)
                return (resp["taskId"], False)
            except (requests.RequestException, KeyError):
                pass
    
    def __get_answer(self,task_id:int):
        payload = {"clientKey":self.api_key,"taskId": task_id}
        for _ in range(100):
            try:
                response = requests.post(BASEURL + "/getTaskResult",json=payload,timeout=20,).json()
                if response["errorId"] != 0: # Means there was an error
                    ErrCodeToException("Capsolver", response["errorDescription"])
                if response["status"] == "processing":
                    time.sleep(3)
                    continue
                if self.captcha_type == "normal" or self.captcha_type == "image":
                    return response["solution"]["text"]
                else:
                    return response["solution"]["gRecaptchaResponse"]
            except (requests.RequestException, KeyError):
                pass    
