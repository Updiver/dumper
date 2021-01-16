package dumper

// SetVCS sets vcs hoster to current dumper instance
func (d *Dumper) SetVCS(vcs VCSHoster) {
	d.vcs = vcs
}
