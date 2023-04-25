package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "dumper library version",
		Long:  "shows dumper library version used in this cmd example",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: need build var for version here
			const version = "0.0.1"
			fmt.Printf("dumper version: %s\n", version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
