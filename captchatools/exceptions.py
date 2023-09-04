# ------------------------------------------------------------------------------- #

class HarvesterException(Exception):
    '''
    This exception is the base exception for all things related to the harvester
    '''

class NoHarvesterException(HarvesterException):
    '''
    This exception gets thrown when a user doesn't properly set a harvester.
    '''
    def __init__(self, message):
        super().__init__(message)

class NoCaptchaType(HarvesterException):
    '''
    This exception gets thrown when no captcha type was set
    '''
    def __init__(self, message):
        super().__init__(message)

class WrongAPIKeyException(HarvesterException):
    '''
    This exception gets thrown when the user provides a wrong API Key
    '''
    def __init__(self, message="[captchatools] Incorrect API key for the captcha solving site"):
        super().__init__(message)

class TaskDetails(HarvesterException):
    '''
    This exceptions gets thrown when there is missing data
    '''

class Banned(HarvesterException):
    '''
    This exception gets thrown when the user is banned from the solving site
    '''


class NoBalanceException(HarvesterException):
    '''
    This exception gets thrown when there is no more funds available for the provider
    '''
    def __init__(self, message="[captchatools] No balance available"):
        super().__init__(message)

class WrongSitekeyException(HarvesterException):
    '''
    This exception gets thrown when the user provides a wrong google sitekey
    '''
    def __init__(self, message="[captchatools] Incorrect google sitekey"):
        super().__init__(message)
# ------------------------------------------------------------------------------- #

class CaptchaAnswerException(Exception):
    '''
    Base exception for errors related to captcha answers
    '''

class CaptchaIMGTooBig(CaptchaAnswerException):
    '''
    This exception gets thrown when the filesize of the captcha image is too big for the solving site.
    '''
    def __init__(self, message="[captchatools] Size of the captcha image is too big."):
        super().__init__(message)

class FailedToGetCapIMG(CaptchaAnswerException):
    '''
    This exception gets thrown when the program fails to get the captcha image after 3 tries
    '''
    def __init__(self, message="[captchatools] Failed to fetch captcha image."):
        super().__init__(message)

class NoSlotAvailable(CaptchaAnswerException):
    '''
    This exceptions gets thrown when there is no worker available
    '''
class CaptchaImageError(CaptchaAnswerException):
    '''
    This exception gets thrown when there is an error with the captcha image
    '''

# ------------------------------------------------------------------------------- #

class UnknownError(Exception):
    '''
    This exceptions gets thrown when there is an unknown error
    '''
