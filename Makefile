GO ?= go
TOOLS = $(CURDIR)/.tools

$(TOOLS):
	@mkdir -p $@
$(TOOLS)/%: | $(TOOLS)
	@GOBIN=$(TOOLS) go install $(PACKAGE)

GOLANGCI_LINT = $(TOOLS)/golangci-lint
$(GOLANGCI_LINT): PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest

GOVULNCHECK = $(TOOLS)/govulncheck
$(GOVULNCHECK): PACKAGE=golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: tools
tools: $(GOLANGCI_LINT) $(GOVULNCHECK)

.PHONY: govulncheck
govulncheck: | $(GOVULNCHECK)
	@$(TOOLS)/govulncheck ./...

.PHONY: lint lint-fix
lint-fix: ARGS=--fix
lint-fix: lint
lint: | $(GOLANGCI_LINT)
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
check: tidy lint govulncheck test-race
