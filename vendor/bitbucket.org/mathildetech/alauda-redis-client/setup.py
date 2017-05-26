#!/usr/bin/env python
# -*- coding: utf-8 -*-

from setuptools import setup, find_packages


setup(
    name='alauda_redis_client',
    version='0.1.0',
    description='redis client, used to switch normal redis and redis cluster',
    url='https://bitbucket.org/mathildetech/alauda_redis_client',
    author='Jian Liao',
    author_email='jliao@alauda.io',
    license='MIT',
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2.7'
    ],
    packages=find_packages(),
    install_requires=[
        'redis>=2.10.3',
        'redis-py-cluster>=1.3.4'
    ]
)
