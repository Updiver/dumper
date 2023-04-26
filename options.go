package dumper

func SetVCSHoster(vcs VCSHoster) Option {
	return func(d *Dumper) {
		d.vcs = vcs
	}
}
