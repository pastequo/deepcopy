# High-level targets

.PHONY: build check run

build: build.local
check: check.imports check.fmt check.lint check.test


## Build targets

TAG=latest

.PHONY: build.vendor build.vendor.full build.prepare build.cmd build.local build.packr2

build.vendor:
	GO111MODULE=on go mod vendor

build.vendor.full:
	@rm -fr $(PWD)/vendor
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

build.local:
	GO111MODULE=on go build -mod=vendor $(BUILD_ARGS) $(PWD)/deepcopy.go


## Check target

LINT_COMMAND=golangci-lint run $(PWD)/deepcopy.go
LINT_RESULT=$(PWD)/lint/result.txt
FILES_LIST=$(PWD)/deepcopy.go

.PHONY: check.fmt check.imports check.lint check.test

check.fmt:
	GO111MODULE=on gofmt -s -w $(FILES_LIST)

check.imports:
	GO111MODULE=on goimports -w $(FILES_LIST)

check.lint:
	@rm -fr $(PWD)/lint
	@mkdir -p $(PWD)/lint
	GO111MODULE=on $(LINT_COMMAND) >> $(LINT_RESULT) 2>&1

check.test:
	GO111MODULE=on go test $(PWD)/
