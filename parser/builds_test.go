package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func (d FakeDryccCmd) BuildsInfo(string, int) error {
	return errors.New("builds:info")
}

func (d FakeDryccCmd) BuildsCreate(string, string, string, string, string, string) error {
	return errors.New("builds:create")
}

func (d FakeDryccCmd) BuildsFetch(string, int, string, string, string, bool) error {
	return errors.New("builds:fetch")
}

func TestBuilds(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := FakeDryccCmd{WOut: &b, ConfigFile: cf}

	// cases defines the arguments and expected return of the call.
	// if expected is "", it defaults to args[0].
	cases := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"builds:info"},
			expected: "",
		},
		{
			args:     []string{"builds:create", "drycc/example-go:latest"},
			expected: "",
		},
		{
			args:     []string{"builds:fetch"},
			expected: "",
		},
		{
			args:     []string{"builds"},
			expected: "builds:info",
		},
	}

	// For each case, check that calling the route with the arguments
	// returns the expected error, which is args[0] if not provided.
	for _, c := range cases {
		var expected string
		if c.expected == "" {
			expected = c.args[0]
		} else {
			expected = c.expected
		}
		err = Builds(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
