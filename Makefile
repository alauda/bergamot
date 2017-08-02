OS = Linux

VERSION = 0.0.1

CURDIR = $(shell pwd)
SOURCEDIR = $(CURDIR)
COVER = $($3)

ECHO = echo
RM = rm -rf
MKDIR = mkdir

# If the first argument is "cover"...
ifeq (cover,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: test

setup:
	@$(ECHO) "Installing gomock..."
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen
	@$(ECHO) "Installing govendor..."
	go get -u github.com/kardianos/govendor
	@$(ECHO) "Installing cobra..."
	go get -v github.com/spf13/cobra/cobra

add:
	govendor add +external

all: build

build:
	go build -ldflags "-w -s" -v -o bergamot github.com/alauda/bergamot

help:
	@$(ECHO) "Targets:"
	@$(ECHO) "all				- test"
	@$(ECHO) "setup				- install necessary libraries"
	@$(ECHO) "add				- runs govendor add +external command"
	@$(ECHO) "build				- build and exports using ALAUDACI_DEST_DIR"
	
