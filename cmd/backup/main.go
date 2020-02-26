package main

import (
	"dumper/pkg/cli/backup"
	"dumper/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewCommand()

	// adding commands for cli
	rootCmd.AddCommand(
		backup.NewCommand(),
	)

	rootCmd.Execute()
}
