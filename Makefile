# Variables
VERSION := 0.1.0

# .env
ifneq (,$(wildcard ./.env))
	include .env
	export
endif


##
##@ General
##

## Print this message and exit
.PHONY: help
help:
	@awk '																								\
		BEGIN { 																						\
			printf "\nUsage:\n  make \033[36m<target>\033[0m\n"											\
		}																								\
		END {																							\
			printf "\n"																					\
		}																								\
		/^[0-9A-Za-z-]+:/ {																				\
			if (prev ~ /^## /) {																		\
				printf "  \x1b[36m%-23s\x1b[0m %s\n", substr($$1, 0, length($$1)-1), substr(prev, 3)	\
			}																							\
		}																								\
		/^##@/ {																						\
			printf "\n\033[1m%s\033[0m\n", substr($$0, 5)												\
		}																								\
		!/^\.PHONY/ {																					\
			prev = $$0																					\
		}																								\
	' $(MAKEFILE_LIST)


##
##@ Building
##

## Build development version
.PHONY: dev
dev:
	node build.js --dev
	go build

## Build production version
.PHONY: prod
prod:
	node build.js --prod
	go build -ldflags='-s -w -X github.com/tzvetkoff-go/pasteur/pkg/version.Version=$(VERSION)'


##
##@ Running
##

## Run development version
.PHONY: run-dev
run-dev: dev
	./pasteur start

## Run production version
.PHONY: run-prod
run-prod: prod
	./pasteur start


##
##@ Development
##

## Build and watch frontend
.PHONY: watch-frontend
watch-frontend:
	node build.js --dev --watch

## Build and watch backend
.PHONY: watch-backend
watch-backend:
	nwatch -e '.git' -e 'web' -e 'node_modules' -p '*.go' -p '*.html' -p '*.js' -p '*.css' -b 'go build' -s './pasteur start' -w '0.0.0.0:1337'

# vim:ft=make:ts=4:sts=4:noet
