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
	./node_modules/.bin/grunt dev
	go build

## Run development version
.PHONY: run-dev
run-dev: dev
	./pasteur start

## Build production version
.PHONY: prod
prod:
	./node_modules/.bin/grunt default
	go build -ldflags='-s -w'

## Run production version
.PHONY: run-prod
run-prod: prod
	./pasteur start

# vim:ft=make:ts=4:sts=4:noet
