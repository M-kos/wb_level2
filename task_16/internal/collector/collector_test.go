package collector

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalPath(t *testing.T) {
	outputDir := "result"
	defer func() {
		err := os.RemoveAll(outputDir)
		if err != nil {
			t.Fatalf("failed to remove test output directory: %v", err)
		}
	}()
	c := &Collector{resultDir: outputDir}
	u, _ := url.Parse("https://example.com/path/page")
	p := c.localPath(u, u, "text/html")
	expected := filepath.Join(outputDir, "example.com", "path", "page.html")
	require.Equal(t, expected, p)
}

func TestSaveFile(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "subdir", "file.txt")
	c := &Collector{}
	data := []byte("hello world")

	require.NoError(t, c.saveFile(p, data))

	b, err := os.ReadFile(p)
	require.NoError(t, err)
	require.Equal(t, string(data), string(b))

	fi, err := os.Stat(p)
	require.NoError(t, err)
	require.False(t, fi.Mode().IsDir())
}
