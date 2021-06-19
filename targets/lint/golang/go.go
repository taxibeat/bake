// Package golang contains linting related mage targets.
package golang

import (
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// GolangciFlags are passed directly to the golanci-lint command.
var GolangciFlags = []string{
	"--no-config",
	"--exclude-use-default=false",
	"--deadline=5m",
	"--modules-download-mode=vendor",
	"--build-tags=component,integration",
	"--disable-all",
	"--enable=govet",
	"--enable=govet,revive,gofumpt,gosec,unparam,goconst,prealloc,stylecheck,unconvert",
}

// ConfigFilePath sets the --config flag in the golangci-lint command, instead of GolangciFlags.
var ConfigFilePath string

// Lint groups together lint related tasks.
type Lint mg.Namespace

// Go runs the golangci-lint linter.
// If ConfigFilePath is set a config file is used, otherwise GolangciFlags are used.
func (l Lint) Go() error {
	args := "run "

	if mg.Verbose() {
		args += "-v "
	}

	if ConfigFilePath != "" {
		args += "--config " + ConfigFilePath
	} else {
		args += strings.Join(GolangciFlags, " ")
	}

	return sh.RunV("golangci-lint", strings.Split(args, " ")...)
}
