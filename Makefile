ifdef CICD
NO_COLOR ?= true
GOPATH := $(HOME)/go
TIME := true
endif

ifdef DEBUG
SHELL=/bin/bash
.SHELLFLAGS = -vuec
else ifdef TIME
SHELL=./scripts/timetarget.sh
else
SHELL=/bin/bash
.SHELLFLAGS = -uec
endif

.EXPORT_ALL_VARIABLES:
.ONESHELL:

.DEFAULT_GOAL := default

PHONY: default

default: build

include build/Makefile.vars
include build/Makefile.func
include build/Makefile.targets
include build/Makefile.gotargets
