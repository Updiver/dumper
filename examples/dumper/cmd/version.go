package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version    string
	GitCommit  string
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "dumper library version",
		Long:  "shows dumper library version used in this cmd example",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dumper version: %s\n", Version)
			fmt.Printf("build commit: %s\n", GitCommit)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
