package dumper

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testRepository = "https://github.com/Updiver/test-repository.git"

func TestDumpRepository(t *testing.T) {
	var (
		destinationDir = "repository-clone-example"
	)

	dir := os.TempDir()
	fmt.Printf("temp dir: %s\n", dir)
	fmt.Printf("desitnation dir: %s\n", destinationDir)
	tempDir := path.Join(filepath.Clean(dir), destinationDir)
	defer os.RemoveAll(tempDir)

	dumper := New()
	dumpOpts := &DumpRepositoryOptions{
		RepositoryURL: testRepository,
		Destination:   tempDir,
		Creds: struct {
			Username string
			Password string
		}{
			Username: "",
			Password: "somerandomnonexistentpassword",
		},
	}
	err := dumper.DumpRepository(dumpOpts)
	require.NoError(t, err, "dump repository")

	fileContent, err := os.Open(path.Join(tempDir, "test-regular-file.txt"))
	require.NoError(t, err, "open file")

	txt, err := io.ReadAll(fileContent)
	require.NoError(t, err, "read file content")

	require.Equal(t, "Test regular file content", string(txt), "expect to have proper file content")
}
