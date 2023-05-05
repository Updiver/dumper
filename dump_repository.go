package dumper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type DumpRepositoryOptions struct {
	RepositoryURL string
	Destination   string
	Creds         struct {
		Username string
		Password string
	}
}

func (opts *DumpRepositoryOptions) Validate() error {
	if opts.RepositoryURL == "" {
		return errors.New("repository url required")
	}

	if opts.Destination == "" {
		return errors.New("destination required")
	}

	// opts.Creds.Username is not mandatory as usually people
	// use app passwords (bitbucket) or access tokens (github)
	if opts.Creds.Password == "" {
		return errors.New("username or password required")
	}

	return nil
}

// DumpRepository dumps single repository
func (d *Dumper) DumpRepository(opts *DumpRepositoryOptions) error {
	if opts == nil {
		return errors.New("dump repository: options required")
	}

	if err := opts.Validate(); err != nil {
		return fmt.Errorf("dump repository validate options: %w", err)
	}

	// TODO: implement fetching branches
	// TODO: implement letting user know which repos are > 1Gb
	// and leave links for them to do that manually
	_, err := git.PlainClone(
		filepath.Clean(opts.Destination),
		false,
		&git.CloneOptions{
			URL:               opts.RepositoryURL,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			SingleBranch:      false,
			Progress:          os.Stderr,
			Auth: &http.BasicAuth{
				Username: opts.Creds.Username, // in case of access token, username should be empty
				Password: opts.Creds.Password,
			},
			// ReferenceName:     plumbing.NewBranchReferenceName(branchName),
		})
	if err != nil {
		return fmt.Errorf("pull repository [%s]: %w", opts.RepositoryURL, err)
	}

	return nil
}
