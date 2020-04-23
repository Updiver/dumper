package dumper

import (
	"dumper/variables/bitbucket"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

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

// Repositories is a list of Repository struct
type Repositories []Repository

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

// Teams returns list of teams user belongs to
func (d *Dumper) Teams() []string {
	var (
		teams     TeamWrapper
		teamNames []string
	)

	client := new(http.Client)

	for page := 1; ; page++ {
		url := fmt.Sprintf(bitbucket.BitbucketTeamsAPI, page)
		fmt.Printf("Sending request to: %s\n", url)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed sending GET request to Bitbucket: %s\n", err)
		}

		request.Header.Add("Content-Type", "application/json")
		// INFO: before you should initialize creds via SetCreds()
		request.SetBasicAuth(d.credentials.Username, d.credentials.Password)

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

// Repositories returns list of repositories and repositories which are too large
func (d *Dumper) Repositories() ([]string, []string) {
	var (
		repositories         RepositoryWrapper
		respositoryNames     []string
		tooLargeRepositories []string
	)

	// TODO: perhaps we need a better naming for this variable
	// I guess teamsList or something similar
	teamRepoList := []string{d.credentials.Username}
	// TODO: bring teams function here
	teamRepoList = append(teamRepoList, d.Teams()...)

	for _, teamName := range teamRepoList {
		for page := 1; ; page++ {
			fmt.Printf("[ %s ] Doing %v page\n", teamName, page)

			url := fmt.Sprintf(bitbucket.BitbucketRepositoriesAPI, teamName, page)
			fmt.Printf("Sending request to: %s\n", url)

			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				// TODO: for now leave Fatalf, but further we need
				// to return errors instead of failing
				log.Fatalf("Failed sending GET request to Bitbucket: %s\n", err)
			}

			request.Header.Add("Content-Type", "application/json")
			request.SetBasicAuth(d.credentials.Username, d.credentials.Password)

			client := &http.Client{}
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
