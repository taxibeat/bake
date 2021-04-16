package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPruneFileNoLines(t *testing.T) {
	err := pruneCoverageFile("", []string{})
	assert.NoError(t, err)
}

func TestPruneLines(t *testing.T) {
	patterns := []string{"mongo", "jaeger", `^te`}
	input := strings.NewReader(`mode: atomic
github.com/taxibeat/bake/targets/test/test.go:47.26,50.2 2 0
github.com/taxibeat/bake/docker/component/jaeger/component.go:18.85,29.27 2 1
github.com/taxibeat/bake/docker/component/mongodb/component.go:55.17,57.4 1 0
github.com/taxibeat/bake/targets/code/golang/golang.go:20.27,23.54 2 0
test
`)
	exp := []string{
		"mode: atomic",
		"github.com/taxibeat/bake/targets/test/test.go:47.26,50.2 2 0",
		"github.com/taxibeat/bake/targets/code/golang/golang.go:20.27,23.54 2 0",
	}

	got, err := pruneCoverageLines(input, patterns)
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
}