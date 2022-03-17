.DELETE_ON_ERROR:
.DEFAULT_GOAL := help
_YELLOW=\033[0;33m
_NC=\033[0m
SHELL := /bin/bash -o pipefail

.PHONY: help
help: ## prints this help
	@ grep -hE '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## runs Go code formatting and rearranges imports
	@{ \
  		gofmt -s -w $(CURDIR) ;\
  		go run golang.org/x/tools/cmd/goimports -format-only -w -local github.com/ntwrk1/guid $(CURDIR) ;\
  	}
