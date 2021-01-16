package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/updiver/dumper/variables/bitbucket"
)

// BitbucketProfile is a profile from bitbucket vcs provider
type BitbucketProfile struct{}

// BitbucketUserProfile describes user profile for bitbucket
type BitbucketUserProfile struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display"`
}

// BitbucketUserProfileAPIError describes error bitbucket will return to us
type BitbucketUserProfileAPIError struct {
	Type  string `json:"type"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Profile pulls user's profile
func (b *BitbucketProfile) Profile(bearerToken string) (*UserProfile, *UserProfileAPIError) {
	var (
		// TODO: right now we return UserProfile but that's was used from bitbucket
		// we have to use bitbucket response and convert that into UserProfile response
		user         *BitbucketUserProfile
		errorMessage *BitbucketUserProfileAPIError
	)

	// TODO: need http interaction library to be shared, lots of trash code
	client := &http.Client{}
	request, err := http.NewRequest("GET", bitbucket.BitbucketUsersAPI, nil)
	if err != nil {
		return &UserProfile{}, &UserProfileAPIError{err.Error()}
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+bearerToken)

	response, err := client.Do(request)
	if err != nil {
		return &UserProfile{}, &UserProfileAPIError{fmt.Sprintf("failed making request to user bitbucket api: %s", err)}
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		err := json.NewDecoder(response.Body).Decode(&errorMessage)
		if err != nil {
			return nil, &UserProfileAPIError{fmt.Sprintf("failed decoding response: %s", err)}
		}

		// INFO: error always will be empty because that depends on vcs provider response
		return nil, &UserProfileAPIError{fmt.Sprintf("failed pulling user profile: %s", errorMessage.Error.Message)}
	} else if response.StatusCode == http.StatusOK {
		err := json.NewDecoder(response.Body).Decode(&user)
		if err != nil {
			return &UserProfile{}, &UserProfileAPIError{fmt.Sprintf("failed decoding response: %s", err)}
		}
	}

	return &UserProfile{
		Username: user.Username,
	}, nil
}
