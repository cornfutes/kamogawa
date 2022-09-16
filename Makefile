cnf ?= config.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

cloc: ## output lines of go code
	find . -name "*.go" -print0 | xargs -0 wc -l

dev: ## start local environment
	docker compose up

build: ## Build and tag for GCR.
	docker build -t gcr.io/linear-cinema-360910/diceduckmonk --platform linux/amd64 .

deploy: build ## Build and Deploy to GCR.
	docker push gcr.io/linear-cinema-360910/diceduckmonk

tidy: ## cleanups code
	go mody tidy

testprod: ## visual regression test on prod
	cd e2e && npx playwright test
