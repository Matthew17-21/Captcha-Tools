from .errors import ErrCodeToException
from . import Harvester, exceptions as captchaExceptions
import time
import requests
from typing import Optional

BASEURL = "https://api.solvecha.com"


class Solvecha(Harvester):
    '''
    This object will contain the data to interact with the Solvecha API.
    Solvecha solves image captchas only (https://solvecha.com/docs).
    '''
    def _post(self, path: str, payload: dict) -> dict:
        resp = requests.post(BASEURL + path, json=payload, timeout=20)
        data = resp.json()
        if data.get("errorId", 0) != 0:
            # Errors come back as 4xx with a JSON body; keep the API message
            try:
                ErrCodeToException("Solvecha", data.get("errorCode", ""))
            except captchaExceptions.UnknownError:
                raise captchaExceptions.UnknownError(
                    f"Error returned from Solvecha: {data.get('errorDescription') or data.get('errorCode')}"
                )
        return data

    def get_balance(self) -> float:
        for _ in range(5):
            try:
                return float(self._post("/getBalance", {"clientKey": self.api_key})["balance"])
            except requests.RequestException:
                time.sleep(1)

    def get_token(
            self,
            b64_img: Optional[str] = None,
            user_agent: Optional[str] = None,
            proxy: Optional[str] = None,
            proxy_type: Optional[str] = "HTTP",
            rq_data: Optional[str] = None
        ):
        if self.captcha_type not in ("image", "normal"):
            raise captchaExceptions.NoCaptchaType("Solvecha only solves image captchas")

        payload = {
            "clientKey": self.api_key,
            "task": {"type": "ImageToTextTask", "body": b64_img or ""},
        }
        task_id = None
        for _ in range(50):
            try:
                resp = self._post("/createTask", payload)
                # The answer can already be in the createTask response
                if resp.get("status") == "ready":
                    return resp["solution"]["text"]
                task_id = resp["taskId"]
                break
            except (requests.RequestException, KeyError):
                time.sleep(3)
        if task_id is None:
            return None

        for _ in range(100):
            try:
                resp = self._post("/getTaskResult", {"clientKey": self.api_key, "taskId": task_id})
                if resp.get("status") == "ready":
                    return resp["solution"]["text"]
                time.sleep(3)
            except (requests.RequestException, KeyError):
                time.sleep(3)
