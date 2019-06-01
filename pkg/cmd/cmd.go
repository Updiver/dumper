package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
)

const APP_VERSION = "0.1.0"

var Username string
var Password string

// NewCommand creates cobra command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "bitbucket-backup is a cli tool to dump repos from bitbucket",
		Long:    `bitbucket-backup is a cli tool to dump repos from bitbucket`,
		Version: APP_VERSION,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Flags here: %+v %+v", Username, Password)
		},
	}

	cmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "Bitbucket username")
	cmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "Bitbucket password")

	return cmd
}
