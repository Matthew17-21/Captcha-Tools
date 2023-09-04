from . import exceptions as captchaExceptions

def ErrCodeToException(site: str, error_code:str):
    '''
    Converts a given error code from the solving site and raises an appropriate exception
    '''
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
    elif error_code == "ERROR_IP_NOT_ALLOWED" or error_code == "IP_BANNED":
        return captchaExceptions.Banned(error_code)
    elif error_code == "ERROR_ZERO_CAPTCHA_FILESIZE" or error_code == "ERROR_UPLOAD" or \
        error_code == "ERROR_CAPTCHAIMAGE_BLOCKED" or error_code == "ERROR_IMAGE_TYPE_NOT_SUPPORTED" or \
        error_code == "ERROR_WRONG_FILE_EXTENSION":
        raise captchaExceptions.CaptchaImageError(error_code)
    else: raise captchaExceptions.UnknownError(f"Error returned from {site}: {error_code}")