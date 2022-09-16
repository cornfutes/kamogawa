cnf ?= config.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help command.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

cloc: ## Total SLOC w/ Tokei. Run `brew install tokei` first
	tokei 

cloc_go: ## Naive SLOC of Go code.
	find . -name "*.go" -print0 | xargs -0 wc -l

dev: ## Run local Postgres server + live-reloading App Server
	docker compose up

dev_clean: ## Start local environment from clean slate 
	docker compose build --no-cache && docker-compose up

build: ## Build and tag for GCR.
	docker build -t gcr.io/linear-cinema-360910/diceduckmonk --platform linux/amd64 ./src

deploy: build ## Build then deploy to GCR.
	docker push gcr.io/linear-cinema-360910/diceduckmonk

go_clean: ## Cleans up go code.
	cd src && go mod tidy && go mode clean

test_prod: ## visual regression against prod
	cd e2e && npx playwright test
