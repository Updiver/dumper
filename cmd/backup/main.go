package main

import (
	"bitbucket-backup/pkg/cli/backup"
	"bitbucket-backup/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewCommand()

	// adding commands for cli
	rootCmd.AddCommand(
		backup.NewCommand(),
	)

	rootCmd.Execute()
}
