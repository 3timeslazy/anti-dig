package antidig

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var generate = flag.Bool("generate", false, "generates output to testdata/ if set")

func VerifyVisualization(t *testing.T, testname string, c *Container, opts ...VisualizeOption) {
	var b bytes.Buffer
	require.NoError(t, Visualize(c, &b, opts...))

	dotFile := filepath.Join("testdata", testname+".dot")

	if *generate {
		err := os.WriteFile(dotFile, b.Bytes(), 0644)
		require.NoError(t, err)
		return
	}

	wantBytes, err := os.ReadFile(dotFile)
	require.NoError(t, err)

	got := b.String()
	want := string(wantBytes)
	assert.Equal(t, want, got,
		"Output did not match. Make sure you updated the testdata by running 'go test -generate'")
}
