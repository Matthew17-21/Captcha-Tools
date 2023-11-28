#! /bin/bash

echo -e "Running setup.py and creating wheel..."
python setup.py sdist bdist_wheel

echo -e "Uploading to test.pypi ..."
python -m twine upload --repository-url https://test.pypi.org/legacy/ dist/*