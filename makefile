SHELL=/usr/bin/env bash

# Project specific properties.
application_name = confetti

# Go tools absolute paths.
go_imports	 = $(shell which goimports || echo '')
go_lint		 = $(shell which golint || echo '')
static_check = $(shell which staticcheck || echo '')

# Go tools installation URLs.
go_imports_url   = golang.org/x/tools/cmd/goimports
go_lint_url      = golang.org/x/lint/golint
static_check_url = honnef.co/go/tools/cmd/staticcheck

# Listing all .go files.
go_files = $(shell find . -type f -name '*.go' | grep -v /vendor/)

# Checks the presence of Go tools. If not found, installs them.
tools:
	@echo "+$@"
	$(if $(go_imports), , go install $(go_imports_url))
	$(if $(go_lint), , go install $(go_lint_url))
	$(if $(static_check), , go install $(static_check_url))

# Runs static-check over the project.
check: tools
	@echo "+$@"
	@$(static_check) ./...

# Runs vet over the project.
vet:
	@echo "+$@"
	@go vet ./...

# Runs the go-imports tools over the project.
imports: tools
	@echo "+$@"
	@$(go_imports) -w $(go_files)

# Runs go fmt over the project.
fmt:
	@echo "+$@"
	@go fmt ./...

# Runs golint over the project.
lint: tools
	@echo "+$@"
	@$(go_lint) ./...

# Builds the project and outputs the binary.
build: check vet imports fmt lint
	@echo "+$@"
	@go build .

# Tests the whole project.
test:
	@echo "+$@"
	@go test ./...
