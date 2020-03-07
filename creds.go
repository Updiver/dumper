package dumper

// Credentials is for creds for user git-service account
// TODO: in future, we can have other creds like tokens, some keys etc
type Credentials struct {
	Username string
	Password string
}

// SetCreds sets credential for user to access repositories/github accounts
func (d *Dumper) SetCreds(creds Credentials) {
	d.credentials = creds
}
