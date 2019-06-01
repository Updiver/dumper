package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
)

// AppVersion indicated application version
const AppVersion = "0.1.0"

var username string
var password string
var interval string

// NewCommand creates cobra command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "bitbucket-backup is a cli tool to dump repos from bitbucket",
		Long:    `bitbucket-backup is a cli tool to dump repos from bitbucket`,
		Version: AppVersion,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Flags here: %+v %+v", username, password)
		},
	}

	cmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Bitbucket username")
	cmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Bitbucket password")
	cmd.PersistentFlags().StringVarP(&interval, "interval", "i", "", "Interval. How frequently make requests (in seconds)")

	return cmd
}
