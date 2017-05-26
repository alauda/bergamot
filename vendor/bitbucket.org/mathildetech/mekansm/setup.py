#!/usr/bin/env python
# -*- coding: utf-8 -*-

from setuptools import find_packages, setup

setup(
    name='mekansm',
    version='0.1.12',
    description='Mathilde common modules.',
    url='https://bitbucket.org/mathildetech/mekansm',
    author='Yongsong You',
    author_email='ysyou@alauda.io',
    license='MIT',
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2.7'
    ],
    packages=find_packages(),
    install_requires=[
        'Django>=1.7',
        'djangorestframework>=2.3.14',
        'requests>=2.4.3',
        'six>=1.10.0'
    ]
)
