package backup

import (
	"encoding/json"
	"fmt"
	"log"

	// "bitbucket-backup/pkg/backup"

	"net/http"

	"github.com/spf13/cobra"
)

// BitbucketRepositoriesAPI points to api with list of repositories
const BitbucketRepositoriesAPI = "https://api.bitbucket.org/2.0/repositories/%s"

// RepositoryWrapper Contains meta information about repository object
type RepositoryWrapper struct {
	// Total number of objects in the existing page
	PageLength int `json:"pagelen"`
	// Repositories of specified user
	Repositories Repositories `json:"values"`
	// Current page
	Page int `json:"page"`
	// Pointer to next page if any
	Next string `json:"next"`
}

// Repository contains repository description
type Repository struct {
	// Only reposttory name without username
	Name string `json:"name"`
	// Repository name with username
	FullName string `json:"full_name"`
}

// Repositories is a list of Repository struct
type Repositories []Repository

// NewCommand creates new command and returns new cmd pointer
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Backup command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is a backup command")
			// flag := cmd.Flag("username")
			// fmt.Printf("%#+v", cmd)

			var repositories RepositoryWrapper
			url := fmt.Sprintf(BitbucketRepositoriesAPI, cmd.Flag("username").DefValue)

			client := http.DefaultClient
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf("Failed sending GET request to Bitbucket: %s", err)
			}

			request.Header.Add("Content-Type", "application/json")
			response, err := client.Do(request)
			if err != nil {
				log.Fatalf("Failed making request: %s", err)
				return
			}
			defer response.Body.Close()

			err = json.NewDecoder(response.Body).Decode(&repositories)

			fmt.Printf("%v %v", len(repositories.Repositories), repositories.PageLength)

			// backup.Clone()
		},
	}
}
