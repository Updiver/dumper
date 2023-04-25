package bitbucket

// This package describes Bitbucket APIs with params passed to them

// BitbucketRepositoriesAPI points to api with list of repositories
const BitbucketRepositoriesAPI = "https://api.bitbucket.org/2.0/repositories/%v?page=%v"

// BitbucketTeamsAPI points to bitbucket teams api
// const BitbucketTeamsAPI = "https://api.bitbucket.org/2.0/teams?role=member&page=%v"
const BitbucketTeamsAPI = "https://api.bitbucket.org/2.0/workspaces?role=member&page=%v"

// BitbucketUsersAPI points to bitbucket users api
const BitbucketUsersAPI = "https://api.bitbucket.org/2.0/user"
