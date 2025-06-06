package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/drycc/workflow-cli/version"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("DRYCC_API_VERSION", "1234")
		w.WriteHeader(200)
	})

	err = cmdr.Version(true)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), fmt.Sprintf(`Workflow CLI Version:            %s
Workflow CLI API Version:        %s
Workflow Controller API Version: 1234
`, version.Version, drycc.APIVersion), "output")

	b.Reset()
	err = cmdr.Version(false)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), version.Version+"\n", "output")
}
