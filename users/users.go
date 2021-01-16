package users

// User is general interface to represent user from different vcs providers (bitbucket, github etc)
type User interface {
	Profile() (*UserProfile, *UserProfileAPIError)
}

// UserProfile describes user profile for bitbucket
type UserProfile struct {
	Username string `json:"username"`
}

// UserProfileAPIError describes error bitbucket will return to us
type UserProfileAPIError struct {
	Error string `json:"error"`
}
