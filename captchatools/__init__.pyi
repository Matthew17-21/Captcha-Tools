from abc import ABC, abstractmethod
from typing import Optional
class Harvester(ABC):
    '''
    Represents a captcha harvester.
    '''
    def __init__(
        self, 
        api_key: Optional[str] = None,
        captcha_type: Optional[str] = "v2",
        solving_site: Optional[str] = None,
        invisible_captcha: Optional[bool] = False,
        captcha_url: Optional[str] = None,
        min_score : Optional[float] = 0.7,
        sitekey : Optional[str] = None,
        action: Optional[str] = "verify",
        soft_id: Optional[str] = None
    ) -> None: ...
    
    @abstractmethod
    def get_balance(self) -> float: ...
    
    @abstractmethod
    def get_token(
        self, b64_img: Optional[str]=None, 
        user_agent: Optional[str]=None, 
        proxy: Optional[str]=None, 
        proxy_type: Optional[str]=None
    ): ... 

def new_harvester(
    api_key: Optional[str] = None,
    captcha_type: Optional[str] = "v2",
    solving_site: Optional[str] = None,
    invisible_captcha: Optional[bool] = False,
    captcha_url: Optional[str] = None,
    min_score : Optional[float] = 0.7,
    sitekey : Optional[str] = None,
    action: Optional[str] = "verify",
    soft_id: Optional[str] = None
) -> Harvester: ...

# Just for backward compatibility
def captcha_harvesters(
    api_key: Optional[str] = None,
    captcha_type: Optional[str] = "v2",
    solving_site: Optional[str] = None,
    invisible_captcha: Optional[bool] = False,
    captcha_url: Optional[str] = None,
    min_score : Optional[float] = 0.7,
    sitekey : Optional[str] = None,
    action: Optional[str] = "verify",
    soft_id: Optional[str] = None
) -> Harvester: ...