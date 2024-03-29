// Code generated by github.com/izumin5210/gex. DO NOT EDIT.

// +build tools

package tools

// tool dependencies
import (
	_ "github.com/frapposelli/wwhrd"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/izumin5210/gex/cmd/gex"
	_ "golang.org/x/tools/cmd/goimports"
	_ "mvdan.cc/gofumpt"
)

// If you want to use tools, please run the following command:
//  go generate ./tools.go
//
//go:generate go build -v -o=./bin/wwhrd github.com/frapposelli/wwhrd
//go:generate go build -v -o=./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go build -v -o=./bin/gex github.com/izumin5210/gex/cmd/gex
//go:generate go build -v -o=./bin/goimports golang.org/x/tools/cmd/goimports
//go:generate go build -v -o=./bin/gofumpt mvdan.cc/gofumpt
