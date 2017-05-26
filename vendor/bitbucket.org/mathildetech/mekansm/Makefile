OS = Linux

VERSION = 0.0.1

PACKAGE_NAME = mekansm
CURDIR = $(shell pwd)
SOURCEDIR = $(CURDIR)

ECHO = echo
RM = rm -rf
MKDIR = mkdir
FLAKE8 = flake8
PEP8-NAMING = pep8-naming
PIP_INSTALL = pip install

.PHONY: setup build test help

all: build

setup:
	$(PIP_INSTALL) $(FLAKE8) $(PEP8-NAMING)

build:
	$(FLAKE8) $(SOURCEDIR)/$(PACKAGE_NAME) --show-source --show-pep8 --statistics --count

help:
	@$(ECHO) "Targets:"
	@$(ECHO) "all     - setup, build and test"
	@$(ECHO) "setup   - set up prerequisites for build"
	@$(ECHO) "build   - perform static analysis"
