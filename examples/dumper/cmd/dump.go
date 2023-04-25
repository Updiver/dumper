package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/updiver/dumper"
	"github.com/updiver/dumper/pkg/backup"
	cli "github.com/updiver/dumper/pkg/cli/backup"
)

var (
	Username          string
	Token             string
	DestinationFolder string

	dumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "dump clones repositories by using user creds passed in",
		Run: func(cmd *cobra.Command, args []string) {
			dump := new(dumper.Dumper)
			dump.SetCreds(dumper.Credentials{
				Token: Token,
			})

			repos, tooLargeRepos := dump.Repositories()
			for _, repoName := range repos {
				// TODO: no hardcodes later
				urlWithCreds := "https://bitbucket.org/" + repoName
				// FIXME: make sure we have no identical folder names
				backup.Clone(
					urlWithCreds,
					fmt.Sprintf("%s/%s", DestinationFolder, repoName),
					cli.Creds{
						Username: Username, Password: Token,
					},
				)
			}

			fmt.Println("=== COULDN'T CLONE THESE REPOS BECAUSE OF THEIR LARGE SIZE ===")
			fmt.Println("=== Please clone those repos manually or clean them to have size under 500MB ===")
			for _, rName := range tooLargeRepos {
				fmt.Println(rName)
			}
		},
	}
)

func init() {
	dumpCmd.Flags().StringVarP(&Username, "username", "u", "", "username for git hosting account")
	dumpCmd.Flags().StringVarP(&Token, "token", "t", "", "token which is given by git provider")
	dumpCmd.Flags().StringVarP(&DestinationFolder, "destFolder", "d", "", "destination folder where repositories will be cloned to")

	dumpCmd.MarkFlagRequired("username")
	dumpCmd.MarkFlagRequired("token")
	dumpCmd.MarkFlagRequired("destFolder")
	dumpCmd.MarkFlagsRequiredTogether("username", "token", "destFolder")

	rootCmd.AddCommand(dumpCmd)
}
