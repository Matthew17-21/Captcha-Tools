#! /bin/bash

echo -e "Running setup.py and creating wheel..."
python setup.py sdist bdist_wheel

echo -e "Uploading to pypi..."
python -m twine upload dist/*