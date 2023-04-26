package bitbucket

import (
	"fmt"
	"log"

	bitbucket "github.com/ktrysmt/go-bitbucket"
	"github.com/spf13/cobra"
)

var (
	Username          string
	Token             string
	DestinationFolder string
	BitbucketCmd      = &cobra.Command{
		Use:   "bitbucket",
		Short: "bitbucket clones repositories by using user creds passed in",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dumping bitbucket repositories")

			client := bitbucket.NewBasicAuth(Username, Token)
			client.Pagelen = 100
			client.DisableAutoPaging = false

			workspaces, err := GetWorkspaces(client)
			if err != nil {
				log.Fatalf("get workspaces: %s\n", err)
			}

			workspaceSlugs := GetWorkspaceSlugs(workspaces)
			for _, workspaceSlug := range workspaceSlugs {
				fmt.Printf("= workspace: %s\n", workspaceSlug)

				workspaceRepos, err := client.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
					Owner: workspaceSlug,
					Role:  "member",
				})
				if err != nil {
					fmt.Printf("get repositories: %s\n", err)
					continue
				}

				for _, repository := range workspaceRepos.Items {
					fmt.Printf("== repository: %s\n", repository.Name)
					// TODO: dump repository to dest folder
				}
			}
		},
	}
)

// Workspaces

func GetWorkspaces(client *bitbucket.Client) (*bitbucket.WorkspaceList, error) {
	workspaces, err := client.Workspaces.List()
	return workspaces, err
}

func GetWorkspaceNames(workspaces *bitbucket.WorkspaceList) []string {
	wList := []string{}
	for _, workspace := range workspaces.Workspaces {
		wList = append(wList, workspace.Name)
	}
	return wList
}

func GetWorkspaceSlugs(workspaces *bitbucket.WorkspaceList) []string {
	wList := []string{}
	for _, workspace := range workspaces.Workspaces {
		wList = append(wList, workspace.Slug)
	}
	return wList
}

// Projects

func GetProjects(client *bitbucket.Client, workspaceSlug string) ([]bitbucket.Project, error) {
	projects, err := client.Workspaces.Projects(workspaceSlug)
	return projects.Items, err
}

func init() {
	BitbucketCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	BitbucketCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	BitbucketCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	BitbucketCmd.MarkFlagRequired("username")
	BitbucketCmd.MarkFlagRequired("token")
	BitbucketCmd.MarkFlagRequired("destFolder")
	BitbucketCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
