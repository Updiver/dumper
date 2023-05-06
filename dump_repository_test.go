package dumper

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
)

var (
	testRepositoryURL        = "https://github.com/Updiver/test-repository.git"
	destinationRepositoryDir = "repository-clone-example"
)

func positiveBool() *bool {
	var posBool = true
	return &posBool
}

func negativeBool() *bool {
	var negBool = false
	return &negBool
}

func TestDumpRepository_PositiveNegativeCases(t *testing.T) {
	tempDir := os.TempDir()
	fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir)

	tests := []struct {
		name          string
		opts          *DumpRepositoryOptions
		shouldFail    bool
		expectedError string
	}{
		{
			name: "empty repository url",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     "",
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: positiveBool(),
				Creds: Creds{
					Username: "",
					Password: "blahblah",
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: repository url required",
		},
		{
			name: "empty destination",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       "",
				OnlyDefaultBranch: positiveBool(),
				Creds: Creds{
					Username: "",
					Password: "blahblah",
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: destination required",
		},
		{
			name: "empty credentials",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: positiveBool(),
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: username or password required",
		},
		{
			name: "incorrect credentials (password is missing)",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: positiveBool(),
				Creds: Creds{
					Username: "",
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: username or password required",
		},
		{
			name: "incorrect credentials (username is missing)",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: positiveBool(),
				Creds: Creds{
					Password: "blahblah",
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    false,
			expectedError: "dump repository validate options: username or password required",
		},
		{
			name: "target branch in restrictions without SingleBranch flag set",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: negativeBool(),
				Creds: Creds{
					Password: "blahblah",
				},
				BranchRestrictions: &BranchRestrictions{
					SingleBranch: false,
					BranchName:   "feat/test-regular-file-second-change",
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: branch restrictions validate: branch name is not required when single branch is not set",
		},
		{
			name: "target branch in restrictions with BranchName empty",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: negativeBool(),
				Creds: Creds{
					Password: "blahblah",
				},
				BranchRestrictions: &BranchRestrictions{
					SingleBranch: true,
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: branch restrictions validate: branch name required",
		},
		{
			name: "neither target branch nor SingleBranch flag in restrictions",
			opts: &DumpRepositoryOptions{
				RepositoryURL: testRepositoryURL,
				Destination:   fullDestinationPath,
				Creds: Creds{
					Password: "blahblah",
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail:    true,
			expectedError: "dump repository validate options: either OnlyDefaultBranch or BranchRestrictions required",
		},
		{
			name: "branch restrictions to clone all branches",
			opts: &DumpRepositoryOptions{
				RepositoryURL:     testRepositoryURL,
				Destination:       fullDestinationPath,
				OnlyDefaultBranch: negativeBool(),
				Creds: Creds{
					Password: "blahblah",
				},
				BranchRestrictions: &BranchRestrictions{
					SingleBranch: false,
				},
				Output: &Output{
					GitOutput: io.Discard,
				},
			},
			shouldFail: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dumper := &Dumper{}
			repository, err := dumper.DumpRepository(tc.opts)
			defer os.RemoveAll(fullDestinationPath)

			if tc.shouldFail {
				require.Error(t, err, "should fail with error")
				require.Equal(t, tc.expectedError, err.Error(), "error message should match")
				require.Nil(t, repository, "repository should be nil")
				return
			}

			require.NoError(t, err, "dump repository")
			require.NotNil(t, repository, "repository should not be nil")
		})
	}
}

func TestDumpRepository_DefaultBranch(t *testing.T) {
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
	repository, err := dumper.DumpRepository(opts)
	defer os.RemoveAll(fullDestinationPath)
	require.NoError(t, err, "dump repository")
	require.NotNil(t, repository, "repository should not be nil")

	fileContent, err := os.Open(path.Join(fullDestinationPath, "test-regular-file.txt"))
	require.NoError(t, err, "open file")

	txt, err := io.ReadAll(fileContent)
	require.NoError(t, err, "read file content")

	require.Equal(t, "Test regular file content", string(txt), "expect to have proper file content")

	refIter, err := repository.Branches()
	require.NoError(t, err, "get branches iterator")

	branches := make([]string, 0)
	refIter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})
	require.Len(t, branches, 1, "expect to have only one branch")
	require.Equal(t, "main", branches[0], "expect to have proper branch name")
}

func TestDumpRepository_NonDefaultBranch(t *testing.T) {
	tempDir := os.TempDir()
	fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir)

	dumper := New()
	opts := &DumpRepositoryOptions{
		RepositoryURL: testRepositoryURL,
		Destination:   fullDestinationPath,
		Creds: Creds{
			Password: "blahblah",
		},
		BranchRestrictions: &BranchRestrictions{
			SingleBranch: true,
			BranchName:   "feat/test-regular-file-first-change",
		},
		Output: &Output{
			GitOutput: io.Discard,
		},
	}
	repository, err := dumper.DumpRepository(opts)
	defer os.RemoveAll(fullDestinationPath)
	require.NoError(t, err, "dump repository")
	require.NotNil(t, repository, "repository should not be nil")

	fileContent, err := os.Open(path.Join(fullDestinationPath, "test-regular-file.txt"))
	require.NoError(t, err, "open file")

	txt, err := io.ReadAll(fileContent)
	require.NoError(t, err, "read file content")

	expectedFileContext := "Test regular file content\nThis is first change to this file\n"
	require.Equal(t, expectedFileContext, string(txt), "expect to have proper file content")

	refIter, err := repository.Branches()
	require.NoError(t, err, "get branches iterator")

	branches := make([]string, 0)
	refIter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})
	require.Len(t, branches, 1, "expect to have only one branch")
	require.Equal(t, "feat/test-regular-file-first-change", branches[0], "expect to have proper branch name")
}

func TestDumpRepository_DumpAllBranches(t *testing.T) {
	tempDir := os.TempDir()
	fullDestinationPath := path.Join(filepath.Clean(tempDir), destinationRepositoryDir)

	dumper := New()
	opts := &DumpRepositoryOptions{
		RepositoryURL: testRepositoryURL,
		Destination:   fullDestinationPath,
		Creds: Creds{
			Password: "blahblah",
		},
		BranchRestrictions: &BranchRestrictions{
			SingleBranch: false,
		},
		Output: &Output{
			GitOutput: io.Discard,
		},
	}
	repository, err := dumper.DumpRepository(opts)
	defer os.RemoveAll(fullDestinationPath)
	require.NoError(t, err, "dump repository")
	require.NotNil(t, repository, "repository should not be nil")

	// as this is a bare repository, we have different file structure
	refIter, err := repository.Branches()
	require.NoError(t, err, "get branches iterator")

	branches := make([]string, 0)
	refIter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})
	require.Len(t, branches, 3, "expect to have three branches")

	expectedBranches := []string{
		"feat/test-regular-file-first-change",
		"feat/test-regular-file-second-change",
		"main",
	}
	require.ElementsMatch(t, expectedBranches, branches, "expect to have proper branch names")
}
