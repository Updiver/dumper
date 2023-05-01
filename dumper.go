package dumper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type VCSHoster uint

const (
	VCSBitbucket = iota
	VCSGithub
)

type Dumper struct {
	skipLargeRepos     bool
	largeRepoLimitSize uint64 // bytes
}

type Option func(*Dumper)

// New returns new instance of Dumper
func New(options ...Option) *Dumper {
	d := &Dumper{}

	for _, opt := range options {
		opt(d)
	}

	return d
}

// DumpRepository clones single repository
func (d *Dumper) DumpRepository(repoURL, dumpDestination, username, password string) error {
	destination := filepath.Clean(dumpDestination)

	// TODO: implement fetching branches
	// TODO: implement letting user know which repos are > 1Gb
	// and leave links for them to do that manually
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:               repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		SingleBranch:      false,
		Progress:          os.Stderr,
		Auth: &http.BasicAuth{
			Username: username, // in case of access token, username should be empty
			Password: password,
		},
		// ReferenceName:     plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		return fmt.Errorf("pull repository [%s]: %w", repoURL, err)
	}

	return nil
}
