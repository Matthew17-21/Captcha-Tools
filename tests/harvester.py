import unittest
import sys
import os
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
sys.path.append(os.path.dirname(SCRIPT_DIR))
import captchatools
from captchatools import new_harvester


class TestHarvester(unittest.TestCase):
    '''
    This test ensures that everything initializes correctly
    '''
    def test_2captcha(self):
        h = new_harvester(api_key="key_here", solving_site="2captcha")
        self.assertTrue(isinstance(h, captchatools.Twocap), "Not type captchatools.Twocaptcha")

if __name__ == "__main__":
    unittest.main()