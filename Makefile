
.PHONY: setup
setup: ## Install tools
	go install golang.org/x/tools/cmd/goimports
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0

.PHONY: lint
lint: ## Run the linters
	golangci-lint run

.PHONY: test
test: ## Run all the tests
	go version
	go env
	go list ./... | xargs -n1 -I{} sh -c 'go test -race {}'

.PHONY: fmt
fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

.PHONY: ci
ci: lint test ## Run all the tests and code checks

.PHONY: build
build: ## build binary to .build folder
	go build -o ".build/checkredirects" cmd/main.go

.PHONY: install
install: ## Install to <gopath>/src
	go install ./...

# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help