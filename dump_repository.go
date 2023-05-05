package dumper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type DumpRepositoryOptions struct {
	RepositoryURL      string
	Destination        string
	Creds              Creds
	OnlyDefaultBranch  *bool
	BranchRestrictions *BranchRestrictions
}

type Creds struct {
	Username string
	Password string
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

	if opts.OnlyDefaultBranch == nil && opts.BranchRestrictions == nil {
		return errors.New("either OnlyDefaultBranch or BranchRestrictions required")
	}

	if (opts.OnlyDefaultBranch != nil && *opts.OnlyDefaultBranch) && opts.BranchRestrictions != nil {
		return errors.New("only one branch restriction can be set, use either OnlyDefaultBranch or BranchRestrictions")
	}

	if (opts.OnlyDefaultBranch == nil || !(*opts.OnlyDefaultBranch)) && opts.BranchRestrictions != nil {
		if err := opts.BranchRestrictions.Validate(); err != nil {
			return fmt.Errorf("branch restrictions validate: %w", err)
		}
	}

	return nil
}

// BranchRestrictions has power in case OnlyDefaultBranch is not set
type BranchRestrictions struct {
	// SingleBranch indicates that only one branch should be cloned
	// if not set, all branches will be cloned
	// Depends on a BranchName restriction
	SingleBranch bool
	BranchName   string
}

func (br *BranchRestrictions) Validate() error {
	if br.SingleBranch && br.BranchName == "" {
		return errors.New("branch name required")
	}

	if !br.SingleBranch && br.BranchName != "" {
		return errors.New("branch name is not required when single branch is not set")
	}

	return nil
}

// DumpRepository dumps single repository
func (d *Dumper) DumpRepository(opts *DumpRepositoryOptions) (*git.Repository, error) {
	if opts == nil {
		return nil, errors.New("dump repository: options required")
	}

	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("dump repository validate options: %w", err)
	}

	// TODO: implement fetching branches
	// TODO: implement letting user know which repos are > 1Gb
	// and leave links for them to do that manually
	gitCloneOpts := &git.CloneOptions{
		URL:               opts.RepositoryURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stderr,
		Auth: &http.BasicAuth{
			Username: opts.Creds.Username, // in case of access token, username should be empty
			Password: opts.Creds.Password,
		},
	}

	// default repostiory cloning
	if opts.OnlyDefaultBranch != nil && *opts.OnlyDefaultBranch {
		repository, err := git.PlainClone(
			filepath.Clean(opts.Destination),
			false,
			gitCloneOpts,
		)
		if err != nil {
			return nil, fmt.Errorf("pull repository [%s]: %w", opts.RepositoryURL, err)
		}

		return repository, nil
	}

	// single branch cloning (target branch clone)
	if opts.BranchRestrictions != nil && opts.BranchRestrictions.SingleBranch {
		gitCloneOpts.SingleBranch = opts.BranchRestrictions.SingleBranch
		gitCloneOpts.ReferenceName = plumbing.NewBranchReferenceName(opts.BranchRestrictions.BranchName)

		// TODO: duplicate code, move to separate func
		repository, err := git.PlainClone(
			filepath.Clean(opts.Destination),
			false,
			gitCloneOpts,
		)
		if err != nil {
			return nil, fmt.Errorf("pull repository [%s]: %w", opts.RepositoryURL, err)
		}

		return repository, nil
	}

	// all branches clone (mirror clone, bare repository)
	if opts.BranchRestrictions != nil && !opts.BranchRestrictions.SingleBranch {
		gitCloneOpts.Mirror = true

		// TODO: duplicate code, move to separate func
		repository, err := git.PlainClone(
			filepath.Clean(opts.Destination),
			false,
			gitCloneOpts,
		)
		if err != nil {
			return nil, fmt.Errorf("pull repository [%s]: %w", opts.RepositoryURL, err)
		}

		return repository, nil
	}

	return nil, nil
}
