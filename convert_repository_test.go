package dumper

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
)

func TestConver_FromBareToNonBare(t *testing.T) {
	tempDir := os.TempDir()
	fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir)
	fmt.Println("fullDestinationPath: ", fullDestinationPath)

	dumper := New()
	opts := &DumpRepositoryOptions{
		RepositoryURL:     testRepositoryURL,
		Destination:       fullDestinationPath,
		OnlyDefaultBranch: negativeBool(),
		Creds: Creds{
			Password: "blahblah",
		},
		BranchRestrictions: &BranchRestrictions{
			SingleBranch: false, // mirror branch and have it as bare on local fs
		},
	}
	repository, err := dumper.DumpRepository(opts)
	defer os.RemoveAll(fullDestinationPath)
	require.NoError(t, err, "expect to properly dump repository")
	require.IsType(t, &git.Repository{}, repository, "expect to have proper repository instance")

	// convert repository from bare to non-bare
	err = Convert(fullDestinationPath, RepositoryTypeNonBare)
	require.NoError(t, err, "expect to properly convert repository from bare to non-bare")

	// take new repository instance and verify if it is non-bare
	repository, err = git.PlainOpen(fullDestinationPath)
	require.NoError(t, err, "expect to properly open repository")

	worktree, err := repository.Worktree()
	require.NoError(t, err, "expect to properly get worktree")

	files, err := worktree.Filesystem.ReadDir(".")
	require.NoError(t, err, "expect to properly read files from worktree")
	expectedFolderContent := []string{
		".git",
		"LICENSE",
		"test-regular-file.txt",
	}
	actualFolderContent := []string{}
	for _, file := range files {
		actualFolderContent = append(actualFolderContent, file.Name())
	}
	require.ElementsMatch(
		t,
		expectedFolderContent,
		actualFolderContent,
		"expect to have proper folder content",
	)

	// check for branches
	expectedBranches := []string{
		"main",
		"feat/test-regular-file-second-change",
		"feat/test-regular-file-first-change",
	}
	actualBranches := []string{}
	brIter, err := repository.Branches()
	require.NoError(t, err, "expect to properly get branches iterator")
	brIter.ForEach(func(br *plumbing.Reference) error {
		actualBranches = append(actualBranches, br.Name().Short())
		return nil
	})
	require.ElementsMatch(t, expectedBranches, actualBranches, "expect to have proper branches")
}

// TestConver_FromNonBareToBare supposed to test converting
// from non-bare to bare repository but I see no need to have
// this implementation yet
func TestConver_FromNonBareToBare(t *testing.T) {
	tempDir := os.TempDir()
	fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir)
	fmt.Println("fullDestinationPath: ", fullDestinationPath)

	dumper := New()
	opts := &DumpRepositoryOptions{
		RepositoryURL:     testRepositoryURL,
		Destination:       fullDestinationPath,
		OnlyDefaultBranch: negativeBool(),
		Creds: Creds{
			Password: "blahblah",
		},
		BranchRestrictions: &BranchRestrictions{
			SingleBranch: true,
			BranchName:   "main",
		},
	}
	repository, err := dumper.DumpRepository(opts)
	defer os.RemoveAll(fullDestinationPath)
	require.NoError(t, err, "expect to properly dump repository")
	require.IsType(t, &git.Repository{}, repository, "expect to have proper repository instance")

	// convert repository from non-bare to bare
	err = Convert(fullDestinationPath, RepositoryTypeBare)
	require.ErrorAs(t, err, &ErrNotImplemented, "expect to get not implemented error")
}

func TestMoveFolderContent(t *testing.T) {
	// create a temporary source directory
	sourceDir := filepath.Join(os.TempDir(), "sourceDir")
	err := os.Mkdir(sourceDir, 0755)
	require.NoError(t, err, "expect to properly create temporary source directory")
	defer os.RemoveAll(sourceDir)

	// create some files in the source directory
	_, err = os.Create(filepath.Join(sourceDir, "file1.txt"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(sourceDir, "file2.txt"))
	require.NoError(t, err)

	// create a temporary destination directory
	destDir := filepath.Join(os.TempDir(), "destinationDir")
	err = os.Mkdir(destDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(destDir)

	// move the files from the source directory to the destination directory
	err = moveFolderContent(sourceDir, destDir)
	require.NoError(t, err)

	// check that the files were moved
	_, err = os.Stat(filepath.Join(destDir, "file1.txt"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(destDir, "file2.txt"))
	require.NoError(t, err)

	// try to move the files again, which should fail
	unknownDir := filepath.Join(os.TempDir(), "unknownDir")
	err = moveFolderContent(unknownDir, sourceDir)
	require.Error(t, err)

	// check that the error message includes the expected text
	require.Contains(t, err.Error(), "read dir:")
}
