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
var saveFolder string

// NewCommand creates cobra command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "dumper is a cli tool to dump repos from bitbucket",
		Long:    `dumper is a cli tool to dump repos from bitbucket`,
		Version: AppVersion,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Flags here: %+v %+v", username, password)
		},
	}

	cmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Bitbucket username")
	cmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Bitbucket password")
	cmd.PersistentFlags().StringVarP(&interval, "interval", "i", "", "Interval. How frequently make requests (in seconds)")
	cmd.PersistentFlags().StringVarP(&saveFolder, "saveFolder", "f", "", "Location where to save cloned repos")

	return cmd
}
