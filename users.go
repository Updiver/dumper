package dumper

import "bitbucket.org/updiver/dumper/users"

// Profile returns user's profile from vcs provider (bitbucket, github etc)
func (d *Dumper) Profile() (*users.UserProfile, *users.UserProfileAPIError) {

	// TODO: here we use all from ./users folder (package)
	// this is kind of controller, and users package is a model
	// from model you can pull all required data, model will pull
	// all required data from appropriate vcs
	// so you'll have to wait here for response

	// for now we only user bitbucket
	switch d.vcs {
	case VCSBitbucket:
		bitbucket := new(users.BitbucketProfile)
		return bitbucket.Profile(d.credentials.Token)
	case VCSGithub:
		github := new(users.GithubProfile)
		return github.Profile()
	}

	// bb := new(users.BitbucketProfile)
	// return bb.Profile(d.credentials.Token)
	return nil, &users.UserProfileAPIError{
		Error: "cannot get user's profile, seems like no vcs registered",
	}
}
