default: lint

lint: ## run linter
	@echo "Running golangci-lint with auto-fix..."
	golangci-lint run --fix --config .golangci.yml

.PHONY: lint