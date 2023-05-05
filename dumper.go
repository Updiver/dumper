package dumper

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
