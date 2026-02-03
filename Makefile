.PHONY: help
help:
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_.-]+:.*##/ {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: serve
serve: ## Serve the site locally with drafts
	hugo serve --buildDrafts

.PHONY: fmt
fmt: fmt-yaml ## Format all files

.PHONY: fmt-yaml
fmt-yaml: install-yamlfmt ## Format YAML files
	yamlfmt

.PHONY: install-yamlfmt
install-yamlfmt: ## Install yamlfmt tool
	pushd tools && go install github.com/google/yamlfmt/cmd/yamlfmt && popd

.PHONY: lint
lint: lint-yaml ## Lint all files

.PHONY: lint-yaml
lint-yaml: ## Lint YAML files
	yamlfmt -lint

.PHONY: spell
spell: ## Check for spelling errors
	codespell content README.md

.PHONY: linkcheck
linkcheck: ## Check for broken links
	linkchecker --no-warnings http://localhost:1313

.PHONY: webp
webp: ## Convert images to WebP format
	./scripts/convert-to-webp.sh
