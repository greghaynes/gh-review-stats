COVER_PROFILE = cover.out

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run unit tests
	go test -coverprofile=$(COVER_PROFILE) $(GO_TEST_FLAGS) ./...
	go tool cover -func=$(COVER_PROFILE)

.PHONY: lint
lint: ## Run go linter
	golangci-lint run
