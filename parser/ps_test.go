package parser

import (
	"bytes"
	"errors"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

// Create fake implementations of each method that return the argument
// we expect to have called the function (as an error to satisfy the interface).

func (d FakeDryccCmd) PsList(string, int) error {
	return errors.New("ps:list")
}

func (d FakeDryccCmd) PsLogs(string, string, int, bool, string, bool) error {
	return errors.New("ps:logs")
}

func (d FakeDryccCmd) PsExec(string, string, bool, bool, []string) error {
	return errors.New("ps:exec")
}

func (d FakeDryccCmd) PsDescribe(string, string) error {
	return errors.New("ps:describe")
}

func (d FakeDryccCmd) PsDelete(string, []string) error {
	return errors.New("ps:delete")
}
func TestPs(t *testing.T) {
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
			args:     []string{"ps:list"},
			expected: "",
		},
		{
			args:     []string{"ps:logs", "myapp-web-111"},
			expected: "",
		},
		{
			args:     []string{"ps:exec", "myapp-web-111", "-it", "--", "bash"},
			expected: "",
		},
		{
			args:     []string{"ps:scale", "web", "5"},
			expected: "",
		},
		{
			args:     []string{"ps:describe", "myapp-web-111"},
			expected: "",
		},
		{
			args:     []string{"ps:delete", "myapp-web-111"},
			expected: "",
		},
		{
			args:     []string{"ps"},
			expected: "ps:list",
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
		err = Ps(c.args, cmdr)
		assert.Error(t, errors.New(expected), err)
	}
}
