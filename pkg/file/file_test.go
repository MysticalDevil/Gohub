package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPutAndExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	require.False(t, Exists(path))
	require.NoError(t, Put([]byte("hello"), path))
	require.True(t, Exists(path))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "hello", string(data))
}

func TestNameWithoutExtension(t *testing.T) {
	require.Equal(t, "photo", NameWithoutExtension("photo.png"))
	require.Equal(t, "archive.tar", NameWithoutExtension("archive.tar.gz"))
}
