'''
Library for solving captchas!
~~~~~~~~~~~~~~~~~~~~~

How to use:
>>> import captcha_harvesters
>>> solver = captcha_harvesters.captcha_harvesters()
>>> captcha_answer = solver.get_token()
or

>>> from captcha_harvesters import captcha_harvesters, exceptions
>>> solver = captcha_harvesters()
>>> captcha_answer = solver.get_token()

Sites:
1 = Capmonster
2 = Anticaptcha
3 = 2captcha
'''
from .harvesters import captcha_harvesters
from . import exceptions