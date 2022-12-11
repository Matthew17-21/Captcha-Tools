class WrongAPIKeyException(Exception):
    '''
    This exception gets thrown when the user provides a wrong API Key
    '''
    def __init__(self, message="[captchatools] Incorrect API key for the captcha solving site"):
        super(WrongAPIKeyException, self).__init__(message)
        

class NoBalanceException(Exception):
    '''
    This exception gets thrown when there is no more funds available for the provider
    '''
    def __init__(self, message="[captchatools] No balance available"):
        super(NoBalanceException, self).__init__(message)


class WrongSitekeyException(Exception):
    '''
    This exception gets thrown when the user provides a wrong google sitekey
    '''
    def __init__(self, message="[captchatools] Incorrect google sitekey"):
        super(WrongSitekeyException, self).__init__(message)

class NoHarvesterException(Exception):
    '''
    This exception gets thrown when a user doesn't properly set a harvester.
    '''

class CaptchaIMGTooBig(Exception):
    '''
    This exception gets thrown when the filesize of the captcha image is too big for the solving site.
    '''
    def __init__(self, message="[captchatools] Size of the captcha image is too big."):
        super(CaptchaIMGTooBig, self).__init__(message)

class FailedToGetCapIMG(Exception):
    '''
    This exception gets thrown when the program fails to get the captcha image after 3 tries
    '''
    def __init__(self, message="[captchatools] Failed to fetch captcha image."):
        super(FailedToGetCapIMG, self).__init__(message)

class Banned(Exception):
    '''
    This exception gets thrown when the user is banned from the solving site
    '''

class TaskDetails(Exception):
    '''
    This exceptions gets thrown when there is missing data
    '''

class NoSlotAvailable(Exception):
    '''
    This exceptions gets thrown when there is no worker available
    '''
class CaptchaImageError(Exception):
    '''
    This exception gets thrown when there is an error with the captcha image
    '''

class UnknownError(Exception):
    '''
    This exceptions gets thrown when there is an unknown error
    '''

class NoCaptchaType(Exception):
    '''
    This exception gets thrown when no captcha type was set
    '''