package backup

import (
	"fmt"

	"github.com/spf13/cobra"
	"bitbucket-backup/pkg/backup"
)

// NewCommand creates new command and returns new cmd pointer
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Backup command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is a backup command")
			flag := cmd.Flag("username")
			// fmt.Printf("%#+v", cmd)
			fmt.Println(flag.Value)

			backup.Some()
		},
	}
}
