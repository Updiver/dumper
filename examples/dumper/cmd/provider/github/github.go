package provider

import (
	"fmt"

	"github.com/spf13/cobra"
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
		},
	}
)

func init() {
	GithubCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	GithubCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	GithubCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	GithubCmd.MarkFlagRequired("username")
	GithubCmd.MarkFlagRequired("token")
	GithubCmd.MarkFlagRequired("destFolder")
	GithubCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")
}
