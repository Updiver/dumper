package users

// GithubProfile is a profile from bitbucket vcs provider
type GithubProfile struct{}

// GithubUserProfile describes user profile for bitbucket
type GithubUserProfile struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display"`
}

// GithubUserProfileAPIError describes error bitbucket will return to us
type GithubUserProfileAPIError struct {
	Type  string `json:"type"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Profile pulls user's profile
func (b *GithubProfile) Profile() (*UserProfile, *UserProfileAPIError) {
	return nil, nil
}
