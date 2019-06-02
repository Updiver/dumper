package backup

import (
	"encoding/json"
	"fmt"
	"log"

	"bitbucket-backup/pkg/backup"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// BitbucketRepositoriesAPI points to api with list of repositories
const BitbucketRepositoriesAPI = "https://api.bitbucket.org/2.0/repositories/%v?page=%v"

// BitbucketTeamsAPI points to bitbucket teams api
const BitbucketTeamsAPI = "https://api.bitbucket.org/2.0/teams?role=member&page=%v"

// RepositoryWrapper Contains meta information about repository object
type RepositoryWrapper struct {
	// Total number of objects in the existing page
	PageLength int `json:"pagelen"`
	// Repositories of specified user
	Repositories Repositories `json:"values"`
	// Current page
	Page int `json:"page"`
}

// Repository contains repository description
type Repository struct {
	// Only reposttory name without username
	Name string `json:"name"`
	// Repository name with username
	FullName string `json:"full_name"`
	// Size in bytes
	Size int `json:"size"`
}

// TeamWrapper describes teams linked to current bitbucket user
type TeamWrapper struct {
	// Total number of objects in the existing page
	PageLength int `json:"pagelen"`
	// Teams of specified user
	Teams []struct {
		TeamName string `json:"username"`
	} `json:"values"`
	// Current page
	Page int `json:"page"`
}

// Creds is bitbucket user credentials
type Creds struct {
	Username string
	Password string
}

// Repositories is a list of Repository struct
type Repositories []Repository

// NewCommand creates new command and returns new cmd pointer
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Backup command",
		Run: func(cmd *cobra.Command, args []string) {
			username := cmd.Flag("username").Value.String()
			password := cmd.Flag("password").Value.String()
			saveFolder := cmd.Flag("saveFolder").Value.String()

			fmt.Println("Finished pulling list of repositories")

			repos, tooLargeRepos := repos(Creds{username, password})
			for _, repoName := range repos {
				urlWithCreds := "https://bitbucket.org/" + repoName
				// FIXME: make sure we have no identical folder names
				// directory := strings.Split(repoName, "/")[1] // take repository name
				backup.Clone(urlWithCreds, saveFolder+"/"+repoName, Creds{username, password})
			}

			fmt.Println("=== COULDN'T CLONE THESE REPOS BECAUSE OF THEIR LARGE SIZE ===")
			fmt.Println("=== Please clone those repos manually or clean them to have size under 500MB ===")
			for _, rName := range tooLargeRepos {
				fmt.Println(rName)
			}
		},
	}
}

// teams pulling list of teams current user is member of
func teams(creds Creds) []string {
	client := &http.Client{}
	var teams TeamWrapper
	var teamNames []string

	for page := 1; ; page++ {
		url := fmt.Sprintf(BitbucketTeamsAPI, page)
		fmt.Printf("Sending request to: %s\n", url)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed sending GET request to Bitbucket: %s\n", err)
		}

		request.Header.Add("Content-Type", "application/json")
		request.SetBasicAuth(creds.Username, creds.Password)

		response, err := client.Do(request)
		if err != nil {
			log.Fatalf("Failed making request: %s\n", err)
			return []string{}
		}
		defer response.Body.Close()

		err = json.NewDecoder(response.Body).Decode(&teams)

		for _, team := range teams.Teams {
			teamNames = append(teamNames, team.TeamName)
		}

		time.Sleep(1 * time.Second)

		if len(teams.Teams) == 0 {
			break
		}

	}

	return teamNames
}

// repos returns list of repositories of current bitbucket user
func repos(creds Creds) ([]string, []string) {
	var repositories RepositoryWrapper
	var respositoryNames []string
	var tooLargeRepositories []string
	client := &http.Client{}

	fullRepoList := []string{creds.Username}
	fullRepoList = append(fullRepoList, teams(creds)...)

	for _, teamName := range fullRepoList {
		for page := 1; ; page++ {
			fmt.Printf("[ %s ] Doing %v page\n", teamName, page)

			url := fmt.Sprintf(BitbucketRepositoriesAPI, teamName, page)
			fmt.Printf("Sending request to: %s\n", url)

			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf("Failed sending GET request to Bitbucket: %s\n", err)
			}

			request.Header.Add("Content-Type", "application/json")
			request.SetBasicAuth(creds.Username, creds.Password)

			response, err := client.Do(request)
			if err != nil {
				log.Fatalf("Failed making request: %s\n", err)
				return []string{}, []string{}
			}
			defer response.Body.Close()

			err = json.NewDecoder(response.Body).Decode(&repositories)

			for _, repo := range repositories.Repositories {

				// checking if size is not too big
				// becuase system won't be able to hanle 500MB+
				size := repo.Size
				if (size / (1024 * 1024)) < 500 {
					// FullName is a teamname or username + repo name
					respositoryNames = append(respositoryNames, repo.FullName)
				} else {
					fmt.Println("== ADDING REPO TO TOO LARGE ==")
					tooLargeRepositories = append(tooLargeRepositories, repo.FullName)
				}
			}

			// TODO: allow to specify interval from command line
			time.Sleep(1 * time.Second)

			if len(repositories.Repositories) == 0 {
				break
			}
		}
	}

	return respositoryNames, tooLargeRepositories
}
