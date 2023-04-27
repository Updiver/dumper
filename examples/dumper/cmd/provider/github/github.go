package provider

import (
	"context"
	"fmt"
	"path"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
	"github.com/updiver/dumper/pkg/backup"
	"golang.org/x/oauth2"
)

var (
	Username          string
	Token             string
	DestinationFolder string
	GithubCmd         = &cobra.Command{
		Use:   "github",
		Short: "github clones repositories by using user creds passed in",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dumping github repositories")

			ghClient := GetAuthenticatedClient(Token)
			allRepos, err := GetRepositories(ghClient)
			if err != nil {
				fmt.Printf("get repositories: %s\n", err)
				return
			}

			for _, repo := range allRepos {
				fmt.Printf("org [%s] | repo [%s]\n", *repo.Owner.Login, *repo.Name)
				if repo.CloneURL == nil {
					fmt.Printf("skipping repo [%s] as it has no clone url\n", *repo.Name)
					continue
				}

				fullDestFolder := path.Join(DestinationFolder, *repo.Owner.Login, *repo.Name)
				fmt.Printf("=== clone repository to: %s\n", fullDestFolder)
				backup.Clone(*repo.CloneURL, fullDestFolder, struct {
					Username string
					Password string
				}{
					Username: Username,
					Password: Token,
				})
			}
		},
	}
)

func GetAuthenticatedClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetRepositories(ghClient *github.Client) ([]*github.Repository, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := ghClient.Repositories.List(context.Background(), "", opts)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func init() {
	GithubCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	GithubCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	GithubCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	GithubCmd.MarkFlagRequired("username")
	GithubCmd.MarkFlagRequired("token")
	GithubCmd.MarkFlagRequired("destFolder")
	GithubCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
