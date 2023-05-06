package dumper

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/require"
)

// dummy test
func TestNew(t *testing.T) {
	dumper := New()
	require.NotNil(t, dumper, "expect to have proper dumper instance")
}

func TestRepository(t *testing.T) {
	tempDir := os.TempDir()
	fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir)

	dumper := New()
	opts := &DumpRepositoryOptions{
		RepositoryURL:     testRepositoryURL,
		Destination:       fullDestinationPath,
		OnlyDefaultBranch: positiveBool(),
		Creds: Creds{
			Password: "blahblah",
		},
		Output: &Output{
			GitOutput: io.Discard,
		},
	}
	_, err := dumper.DumpRepository(opts)
	defer os.RemoveAll(fullDestinationPath)
	require.NoError(t, err, "expect to properly dump repository")

	// get repository by it's path
	// and verify if it is returned as Repository instance
	repository, err := Repository(fullDestinationPath)
	require.NoError(t, err, "expect to properly get repository with no errors")
	require.IsType(t, &git.Repository{}, repository, "expect to have proper repository instance")
}
