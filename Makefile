SHELL := /bin/sh
.DEFAULT_GOAL := default
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MAKEFILE_PATH)))

.PHONY: default
default: help

.PHONY: deps-install
# @HELP Intsall project dependencies
deps-install:
	@go get -v github.com/golang/mock/mockgen

.PHONY: mocks-generate
# @HELP Generate mocks for Go interfaces
mocks-generate:
	@mockgen -source=hub.go -destination=mock/hub.go -package=mock
	@mockgen -source=client_factory.go -destination=mock/client_factory.go -package=mock

.PHONY: test
# Examples:
# 	1. Run all tests:
#		$ make test
# 	2. Run all tests in a package:
#		$ make test p=./path/to/pkg
# 	3. Run a single test in a package:
#		$ make test p=./path/to/pkg t=TestFunc/Name
# @HELP Run tests
test:
    ifneq ($(p),)
        ifneq ($(t),)
			@go test -cover -coverprofile=coverage.out -race ${p} -v -run ${t} && go tool cover -html=coverage.out -o coverage.html
        else
			@go test -cover -coverprofile=coverage.out -race ${p} -v && go tool cover -html=coverage.out -o coverage.html
        endif
    else
		@go test -cover -coverprofile=coverage.out -race -v ./... && go tool cover -html=coverage.out -o coverage.html
    endif

.PHONY: help
# @HELP Print help
help:
	@echo "\033[33mPlease use \`make <target>\` where <target> is one of:\033[0m"
	@sed -ne"/^# @HELP /{h;s/.*//;:d" -e"H;n;s/^# @HELP //;td" -e"s/:.*//;G;s/\\n# @HELP /---/;s/\\n/ /g;p;}" ${MAKEFILE_LIST}|LC_ALL='C' sort -f|awk -F --- -v n=$$(tput cols) -v i=32 -v a="$$(tput setaf 6)" -v z="$$(tput sgr0)" '{printf"%s%*s%s ",a,-i,$$1,z;m=split($$2,w," ");l=n-i;for(j=1;j<=m;j++){l-=length(w[j])+1;if(l<= 0){l=n-i-length(w[j])-1;printf"\n%*s ",-i," ";}printf"%s ",w[j];}printf"\n";}'
