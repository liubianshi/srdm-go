package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test database connection",
	Run: func(cmd *cobra.Command, args []string) {
		if Store == nil {
			fmt.Println("Database connection is nil")
			return
		}
		if err := Store.Ping(); err != nil {
			fmt.Printf("Database ping failed: %v\n", err)
		} else {
			fmt.Println("Database connection successful!")
			fmt.Printf("Database path: %s\n", Store.GetPath())
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
