## High-level targets

.DEFAULT_GOAL := help

.PHONY: help tools build check

help: help.all
tools: tools.clean tools.get
format: format.imports format.code
check: check.imports check.fmt check.lint check.test

# Colors used in this Makefile
escape=$(shell printf '\033')
RESET_COLOR=$(escape)[0m
COLOR_YELLOW=$(escape)[38;5;220m
COLOR_BLUE=$(escape)[94m


#####################
# Help targets      #
#####################

.PHONY: help.highlevel help.all

#help help.highlevel: show help for high level targets
help.highlevel:
	@grep -hE '^[a-z_-]+:' $(MAKEFILE_LIST) | LANG=C sort -d | \
	awk 'BEGIN {FS = ":"}; {printf("$(COLOR_YELLOW)%-25s$(RESET_COLOR) %s\n", $$1, $$2)}'

#help help.all: display all targets' help messages
help.all:
	@grep -hE '^#help|^[a-z_-]+:' $(MAKEFILE_LIST) | sed "s/#help //g" | LANG=C sort -d | \
	awk 'BEGIN {FS = ":"}; {if ($$1 ~ /\./) printf("    $(COLOR_BLUE)%-21s$(RESET_COLOR) %s\n", $$1, $$2); else printf("$(COLOR_YELLOW)%-25s$(RESET_COLOR) %s\n", $$1, $$2)}'


#####################
# Tools targets     #
#####################

TOOLS_DIR=$(CURDIR)/tools/bin

.PHONY: tools.clean tools.get

#help tools.clean: remove every tools installed in tools/bin directory
tools.clean:
	rm -fr $(TOOLS_DIR)/*

#help tools.get: retrieve all tools specified in gex
tools.get:
	cd $(CURDIR)/tools && go generate tools.go


##################
# Format targets #
##################

GO_MODULE := $(shell head -n 1 go.mod | cut -d ' ' -f 2)
FILE_LIST := $(shell ls -d *.go)

.PHONY: format.imports format.code

#help format.imports: fix and format go imports
format.imports:
	@$(TOOLS_DIR)/goimports -w -local $(GO_MODULE) $(FILE_LIST)

#help format.code: format go code
format.code:
	@$(TOOLS_DIR)/gofumpt -w $(FILE_LIST)


#####################
# Check targets     #
#####################

LINT_COMMAND=$(TOOLS_DIR)/golangci-lint run -c $(CURDIR)/.golangci.yml

.PHONY: check.fmt check.imports check.lint check.test check.licenses

#help check.fmt: check if code is formated according gofumpt rules
check.fmt:
	@$(TOOLS_DIR)/gofumpt -l $(FILE_LIST) | wc -l | grep 0

#help check.imports: check if imports are well formated
check.imports:
	@$(TOOLS_DIR)/goimports -l -local $(GO_MODULE) $(FILE_LIST) | wc -l | grep 0

#help check.lint: check if the go code is properly written, rules are in .golangci.yml
check.lint:
	$(TOOLS_DIR)/golangci-lint run -c $(CURDIR)/.golangci.yml

#help check.test: execute go unit tests
check.test:
	go test ./...

#help check.licenses: check if the thirdparties' licences are whitelisted (in .wwhrd.yml)
check.licenses:
	$(TOOLS_DIR)/wwhrd check

