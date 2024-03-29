package dumper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	DefaultBranhcName = "master"
)

type DumpRepositoryOptions struct {
	RepositoryURL          string
	Destination            string
	Creds                  Creds
	OnlyDefaultBranch      *bool
	Output                 *Output
	BranchRestrictions     *BranchRestrictions
	RepositoryRestrictions RepositoryRestrictions
}

type Creds struct {
	Username string
	Password string
}

// Output defines where all the output should be written
type Output struct {
	GitOutput    io.Writer
	DumperOutput io.Writer
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

type RepositoryRestrictions struct {
	IgnoreEmptyRepositories bool
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
		Progress:          os.Stdout,
		Auth: &http.BasicAuth{
			Username: opts.Creds.Username, // in case of access token, username should be empty
			Password: opts.Creds.Password,
		},
	}

	if opts.Output != nil && opts.Output.GitOutput != nil {
		gitCloneOpts.Progress = opts.Output.GitOutput
	}

	switch {
	// default repository cloning with default branch
	case opts.OnlyDefaultBranch != nil && *opts.OnlyDefaultBranch:
		return d.plainClone(opts, gitCloneOpts)
	// single branch cloning (target branch clone)
	case opts.BranchRestrictions != nil && opts.BranchRestrictions.SingleBranch:
		gitCloneOpts.SingleBranch = opts.BranchRestrictions.SingleBranch
		gitCloneOpts.ReferenceName = plumbing.NewBranchReferenceName(opts.BranchRestrictions.BranchName)
		return d.plainClone(opts, gitCloneOpts)
	// all branches clone (mirror clone, bare repository)
	case opts.BranchRestrictions != nil && !opts.BranchRestrictions.SingleBranch:
		gitCloneOpts.Mirror = true
		return d.plainClone(opts, gitCloneOpts)
	default:
		return nil, errors.New("dump repository: unknown error")
	}
}

// plainClone is a helper to clone repository
func (d *Dumper) plainClone(opts *DumpRepositoryOptions, gitCloneOpts *git.CloneOptions) (*git.Repository, error) {
	repository, err := git.PlainClone(
		filepath.Clean(opts.Destination),
		false,
		gitCloneOpts,
	)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrEmptyRemoteRepository):
			if opts.RepositoryRestrictions.IgnoreEmptyRepositories {
				return nil, fmt.Errorf("plain clone repository [%s]: %w", gitCloneOpts.URL, err)
			}

			// taken from workaround from go-git discussions: https://github.com/jmalloc/grit/pull/80/files
			r, err := git.PlainInit(opts.Destination, gitCloneOpts.Mirror)
			if err != nil {
				_ = os.RemoveAll(opts.Destination)
				return nil, fmt.Errorf("plain init repository [%s]: %w", gitCloneOpts.URL, err)
			}

			if _, err := r.CreateRemote(&config.RemoteConfig{Name: git.DefaultRemoteName, URLs: []string{opts.RepositoryURL}}); err != nil {
				_ = os.RemoveAll(opts.Destination)
				return nil, fmt.Errorf("plain create remote [%s]: %w", gitCloneOpts.URL, err)
			}

			if err = r.CreateBranch(&config.Branch{Name: DefaultBranhcName, Remote: git.DefaultRemoteName, Merge: plumbing.Master}); err != nil {
				_ = os.RemoveAll(opts.Destination)
				return nil, fmt.Errorf("plain create branch [%s]: %w", gitCloneOpts.URL, err)
			}
		default: // all other kinds of errors
			return nil, fmt.Errorf("plain clone repository [%s]: %w", gitCloneOpts.URL, err)
		}
	}

	return repository, nil
}
