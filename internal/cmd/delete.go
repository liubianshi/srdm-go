package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [names]",
	Short: "Delete data record or table",
	Long:  `Delete data record or table by name. To delete a table, use --force option.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range args {
			if err := Store.Delete(name, deleteForce); err != nil {
				return fmt.Errorf("failed to delete %s: %w", name, err)
			}
			fmt.Printf("Deleted: %s\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force delete table and all included records")
}
