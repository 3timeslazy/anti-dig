//go:build tools
// +build tools

package tools

import (
	// Tools we use during development.
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
