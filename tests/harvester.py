import unittest
import sys
import os
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
sys.path.append(os.path.dirname(SCRIPT_DIR))
import captchatools
from captchatools import new_harvester, exceptions as captchaExceptions
from unittest import mock

def mock_get_response(*args, **kwargs):
    class MockResponse:
        def __init__(self, json_data):
            self.json_data = json_data
        def json(self):
            return self.json_data

    if args[0] == "https://2captcha.com/res.php?key=real_key&action=getbalance&json=1" :
        return MockResponse({"status":1,"request":"17.87942"})
    if args[0] == "https://2captcha.com/res.php?key=fake_key&action=getbalance&json=1" :
        return MockResponse({"status":0,"request":"ERROR_WRONG_USER_KEY","error_text":"You've provided key parameter value in incorrect format], it should contain 32 symbols."})


class Test2Captcha(unittest.TestCase):
    '''
    These tests ensure everything is working correcty with the 2captcha class
    '''
    @mock.patch('requests.get', side_effect=mock_get_response)
    def test_get_balance(self, mock_get):
        real = new_harvester(api_key="real_key", solving_site="2captcha")
        self.assertEqual(real.get_balance(), 17.87942)
        fake = new_harvester(api_key="fake_key", solving_site="2captcha")
        self.assertRaises(captchaExceptions.WrongAPIKeyException, fake.get_balance)

if __name__ == "__main__":
    unittest.main()