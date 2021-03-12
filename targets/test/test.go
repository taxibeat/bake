// Package test contains test related mage targets.
package test

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/taxibeat/bake/docker"
)

var (
	// GoBuildTags used when running all tests.
	GoBuildTags = []string{
		"component",
		"integration",
	}
	// TestArgs used in test targets.
	TestArgs = []string{
		"test",
		"-mod=vendor",
		"-cover",
		"-race",
	}
	// CoverArgs used in coverage targets.
	CoverArgs = []string{
		"test",
		"-mod=vendor",
		"-coverpkg=./...",
		"-covermode=atomic",
		"-coverprofile=coverage.txt",
		"-race",
	}
	// Pkgs is the pkg pattern to target.
	Pkgs = "./..."
)

const (
	goCmd = "go"
)

// Test groups together test related tasks.
type Test mg.Namespace

// Unit runs unit tests.
func (Test) Unit() error {
	args := append(TestArgs, Pkgs)
	return run(args)
}

// All runs all tests.
func (Test) All() error {
	args := append(TestArgs, getBuildTagFlag(GoBuildTags))
	args = append(args, Pkgs)
	return run(args)
}

// CoverUnit runs unit tests and produces a coverage report.
func (Test) CoverUnit() error {
	args := append(CoverArgs, Pkgs)
	return run(args)
}

// CoverAll runs all tests and produces a coverage report.
func (Test) CoverAll() error {
	args := append(CoverArgs, getBuildTagFlag(GoBuildTags))
	args = append(args, Pkgs)
	return run(args)
}

// Cleanup removes any local resources created by `mage test:all`.
func (Test) Cleanup() error {
	return docker.CleanupResources()
}

func run(args []string) error {
	fmt.Printf("test: running tests with args: %v\n", args)
	fmt.Printf("Executing cmd: %s %s\n", goCmd, strings.Join(args, " "))
	return sh.RunV(goCmd, args...)
}

func getBuildTagFlag(buildTags []string) string {
	return "-tags=" + strings.Join(buildTags, ",")
}
