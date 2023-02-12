GO ?= go
TOOLS = $(CURDIR)/.tools

.PHONY: tools
tools:
	@mkdir -p $(TOOLS)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS) v1.51.1

.PHONY: lint lint-fix
lint-fix: ARGS=--fix
lint-fix: lint
lint:
	@$(TOOLS)/golangci-lint run $(ARGS)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: test test-race
test-race: ARGS=-race
test-race: test
test:
	$(GO) test $(ARGS) ./...

.PHONY: check
check: tidy lint test-race
