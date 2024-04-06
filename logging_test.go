package tea

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/rprtr258/assert"
)

func TestLogToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.txt")
	prefix := "logprefix"
	f, err := LogToFile(path, prefix)
	assert.NoError(t, err)

	log.SetFlags(log.Lmsgprefix)
	log.Println("some test log")
	assert.NoError(t, f.Close())

	out := assert.UseFileContent(t, path)

	assert.Equal(t, prefix+" some test log\n", string(out))
}
