from setuptools import setup, find_packages
from io import open

PACKAGE_NAME = "captchatools"
VERSION = "1.3.0"
SHORT_DESCRIPTION = "Python module to help solve captchas with Capmonster, 2captcha and Anticaptcha API's!"
GITHUB_URL = "https://github.com/Matthew17-21/Captcha-Tools"

with open('README.md', 'r', encoding='utf-8') as fp:
    README = fp.read()

setup(
    name = PACKAGE_NAME,
    author = 'Matthew17-21',
    packages=find_packages(),
    author_email = 'admin@monumentalshopping.com',
    version=VERSION,
    description = SHORT_DESCRIPTION,
    long_description=README,
    long_description_content_type='text/markdown',
    url = GITHUB_URL,
    license="MIT",
    keywords = [
        'captcha',
        '2captcha',
        'capmonster',
        'anticaptcha',
        'scraping',
        'scrape',
        'challenge',
        "sneakers"
    ],
    include_package_data = True,
    install_requires = [
        'requests',
    ],
    classifiers=[
        'Intended Audience :: Developers',
        'Natural Language :: English',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Programming Language :: Python :: 3.8',
        'Programming Language :: Python :: 3.9',
        'Topic :: Internet :: WWW/HTTP',
        'Topic :: Software Development :: Libraries :: Python Modules'
    ],
)