package dumper

import (
	"fmt"
	"path"

	"github.com/go-git/go-git/v5"
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

// Repository returns git repository
func Repository(destination string) (*git.Repository, error) {
	destination = path.Clean(destination)
	repository, err := git.PlainOpen(destination)
	if err != nil {
		return nil, fmt.Errorf("open repository: %w", err)
	}

	return repository, nil
}
