package anti_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnti(t *testing.T) {
	const root = "testcases"
	dirs, err := os.ReadDir(root)
	require.NoError(t, err)

	for _, dir := range dirs {
		t.Run(dir.Name(), func(t *testing.T) {
			if !dir.IsDir() {
				return
			}

			dirpath := filepath.Join(root, dir.Name())

			maingo := filepath.Join(dirpath, "main.go")
			actual, err := exec.Command("go", "run", maingo).Output()
			require.NoError(t, err)

			expected, err := os.ReadFile(filepath.Join(dirpath, "expected.txt"))
			require.NoError(t, err)

			assert.Equal(t, string(expected), string(actual))
		})
	}
}
