SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=-race -covermode=atomic -coverprofile=coverage.txt

setup: ## Install tools
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.27.0
	mkdir .bin || true; mv bin/golangci-lint .bin/golangci-lint && rm -rf bin

lint: ## Run the linters
	golangci-lint run

test: ## Run all the tests
	gotestcover $(TEST_OPTIONS)  $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=1m

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

ci: lint test ## Run all the tests and code checks

build: ## build binary to .build folder
	go build -o ".build/checkredirects" main.go

install: ## Install to <gopath>/src
	go install ./...

build-release: ## builds the checked out version into the .release/${tag} folder
	.release/build.sh


# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help