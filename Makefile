# Variables
VERSION ?= 0.1.0

# .env
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

## Print this message and exit
.PHONY: help
help:
	@cat $(MAKEFILE_LIST) | awk '														\
		/^([0-9a-z-]+):.*$$/ {															\
			if (description[0] != "") {													\
				printf("\x1b[36mmake %s\x1b[0m\n", substr($$1, 0, length($$1)-1));		\
				for (i in description) {												\
					printf("| %s\n", description[i]);									\
				}																		\
				printf("\n");															\
				split("", description);													\
				descriptionIndex = 0;													\
			}																			\
		}																				\
		/^##/ {																			\
			description[descriptionIndex++] = substr($$0, 4);							\
		}																				\
	'

## Build development version
.PHONY: dev
dev:
	node build.js --dev
	go build

## Run development version
.PHONY: run-dev
run-dev: dev
	./pasteur start

## Build production version
.PHONY: prod
prod:
	node build.js --prod
	go build -ldflags='-s -w -X github.com/tzvetkoff-go/pasteur/pkg/version.Version=$(VERSION)'

## Run production version
.PHONY: run-prod
run-prod: prod
	./pasteur start

## Build and watch frontend
.PHONY: watch-frontend
watch-frontend:
	node build.js --dev --watch

## Build and watch backend
.PHONY: watch-backend
watch-backend:
	nwatch -e '.git' -e 'web' -e 'node_modules' -p '*.go' -p '*.html' -p '*.js' -p '*.css' -b 'go build' -s './pasteur start' -w '0.0.0.0:1337'

# vim:ft=make:ts=4:sts=4:noet
