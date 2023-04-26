package dumper

type VCSHoster uint

const (
	VCSBitbucket = iota
	VCSGithub
)

type Dumper struct {
	vcs VCSHoster
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

func (d *Dumper) VCSHoster() VCSHoster {
	return d.vcs
}
