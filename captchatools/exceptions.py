
class WrongAPIKeyException(Exception):
    '''
    This exception gets thrown when the user provides a wrong API Key
    '''
    def __init__(self):
        super().__init__()
        

class NoBalanceException(Exception):
    '''
    This exception gets thrown when there is no more funds available for the provider
    '''
    def __init__(self):
        super().__init__()


class WrongSitekeyException(Exception):
    '''
    This exception gets thrown when the user provides a wrong google sitekey
    '''
    def __init__(self):
        super().__init__()
