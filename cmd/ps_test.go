package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/pkg/time"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestPrintProcesses(t *testing.T) {
	var b bytes.Buffer

	pods := []api.Pods{
		{
			Release: "v3",
			Name:    "benign-quilting-web-4084101150-c871y",
			Type:    "web",
			State:   "up",
			Started: time.Time{},
		},
		{
			Release: "v3",
			Name:    "benign-quilting-worker-4084101150-c871y",
			Type:    "worker",
			State:   "up",
			Started: time.Time{},
		},
	}

	printProcesses("appname", pods, &b)

	assert.Equal(t, b.String(), `=== appname Processes
--- web:
benign-quilting-web-4084101150-c871y up (v3)
--- worker:
benign-quilting-worker-4084101150-c871y up (v3)
`, "output")
}

func TestPsList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/pods/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"release": "v2",
					"type": "web",
					"name": "foo-web-4084101150-c871y",
					"state": "up",
					"started": "2016-02-13T00:47:52"
				}
			]
		}`)
	})

	err = cmdr.PsList("foo", -1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `=== foo Processes
--- web:
foo-web-4084101150-c871y up (v2)
`, "output")
}

func TestPsExec(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	server.Mux.Handle(
		"/v2/apps/foo/pods/foo-web-111/exec/",
		websocket.Handler(func(conn *websocket.Conn) {
			io.Copy(conn, conn)
		}),
	)
	err = cmdr.PsExec("foo", "foo-web-111", true, false, []string{"/bin/sh"})
	assert.NoError(t, err)
}

type psTargetCases struct {
	Targets       []string
	ExpectedError bool
	ExpectedMap   map[string]int
	ExpectedMsg   string
}

func TestParsePsTargets(t *testing.T) {
	t.Parallel()

	cases := []psTargetCases{
		{[]string{"test"}, true, nil, "'test' does not match the pattern 'type=num', ex: web=2"},
		{[]string{"test=a"}, true, nil, "'test=a' does not match the pattern 'type=num', ex: web=2"},
		{[]string{"test="}, true, nil, "'test=' does not match the pattern 'type=num', ex: web=2"},
		{[]string{"Test=2"}, true, nil, "'Test=2' does not match the pattern 'type=num', ex: web=2"},
		{[]string{"test=2"}, false, map[string]int{"test": 2}, ""},
		{[]string{"test-proc=2"}, false, map[string]int{"test-proc": 2}, ""},
		{[]string{"test1=2"}, false, map[string]int{"test1": 2}, ""},
	}

	for _, check := range cases {
		actual, err := parsePsTargets(check.Targets)
		if check.ExpectedError {
			assert.Equal(t, err.Error(), check.ExpectedMsg, "error")
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual, check.ExpectedMap, "error")
		}
	}
}

func TestPsScale(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}
	err = cmdr.PsScale("foo", []string{"test"})
	assert.Equal(t, err.Error(), "'test' does not match the pattern 'type=num', ex: web=2", "error")

	server.Mux.HandleFunc("/v2/apps/foo/pods/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"release": "v2",
					"type": "web",
					"name": "foo-web-4084101150-c871y",
					"state": "up",
					"started": "2016-02-13T00:47:52"
				}
			]
		}`)
	})

	server.Mux.HandleFunc("/v2/apps/foo/scale/", func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertBody(t, map[string]int{"web": 1}, r)
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()
	err = cmdr.PsScale("foo", []string{"web=1"})
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Scaling processes... but first, coffee!
done in 0s
=== foo Processes
--- web:
foo-web-4084101150-c871y up (v2)
`, "output")
}

func TestPsRestart(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/apps/foo/pods/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()
	err = cmdr.PsRestart("foo", "")
	assert.NoError(t, err)

	server.Mux.HandleFunc("/v2/apps/coolapp/pods/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()

	err = cmdr.PsRestart("coolapp", "")
	assert.NoError(t, err)

	server.Mux.HandleFunc("/v2/apps/testapp/pods/web/restart/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	b.Reset()

	err = cmdr.PsRestart("testapp", "web")
	assert.NoError(t, err)
}
