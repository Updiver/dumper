package dumper

var (
	instance *Dumper
)

// VCSHoster represents vcs hoster constant
type VCSHoster uint

const (
	// VCSBitbucket bitbucket vcs hoster
	VCSBitbucket = iota
	// VCSGithub github vcs hoster
	VCSGithub
)

// Dumper is a main struct for Dumper library
// although `er` ending indicates in golang that it's an interface,
// in this case Dumper is a name of core library as below
type Dumper struct {
	// credentials is a token / secret required for vcs hosters
	credentials Credentials
	vcs         VCSHoster
}

// New returns new instance of Dumper
func New() *Dumper {
	instance = new(Dumper)

	return instance
}

// Instance returns existing Dumper instance or creates news one
func Instance() *Dumper {
	if instance != nil {
		return instance
	}

	return New()
}
