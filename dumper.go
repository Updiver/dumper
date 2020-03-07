package dumper

// Dumper is a main struct for Dumper library
// although `er` ending indicates in golang that it's an interface,
// in this case Dumper is a name of core library as below
type Dumper struct {
	credentials Credentials
}

// New returns new instance of Dumper
func New() *Dumper {
	return new(Dumper)
}
