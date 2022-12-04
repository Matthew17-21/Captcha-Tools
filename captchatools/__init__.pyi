# https://mypy.readthedocs.io/en/stable/stubs.html
from .harvester import Harvester


def new_harvester(
    api_key : str,
    captcha_type : str,
    solving_site : str | int,
    invisible_captcha : bool,
    captcha_url : int,
    min_score : float,
    sitekey : str,
    action : str,
    soft_id : str | int
    ) -> Harvester: ...