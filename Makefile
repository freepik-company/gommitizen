.PHONY: help bin install uninstall release test-release git-add-extension new-version clean

CURRENT_VERSION := $(shell git describe --tags --abbrev=0)

SOURCE_PATHS := cmd/gommitizen/main.go $(shell find internal/ -type f)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Common targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bin: ./bin/gommitizen ## Build go application

./bin/gommitizen: $(SOURCE_PATHS)
	go build \
		-ldflags "-X github.com/freepik-company/gommitizen/internal/app/gommitizen/version.version=${CURRENT_VERSION}" \
		-o $@ $<

install: bin ## Install gommitizen
	cp ./bin/gommitizen /usr/local/bin/gommitizen

uninstall:  ## Uninstall gommitizen
	rm /usr/local/bin/gommitizen

bump: ## Bump version using commitizen
	cz bump

bump-alpha: ## Bump version using commitizen (alpha)
	cz bump -pr alpha

bump-beta: ## Bump version using commitizen (beta)
	cz bump -pr beta

bump-rc: ## Bump version using commitizen (release candidate)
	cz bump -pr rc

release: bump ## Release new version
	goreleaser release

test-release: ## Test release new version
	goreleaser release --snapshot

gen-docs: ## Generate docs
	go run cmd/gen-docs/main.go > README.md

clean: ## Clean up
	rm -rf ./bin
	rm -rf ./dist
