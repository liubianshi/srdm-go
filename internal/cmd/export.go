package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var exportOutput string

var exportCmd = &cobra.Command{
	Use:   "export [pattern]",
	Short: "Export metadata of data records",
	Long:  `Export metadata of matching data records to a JSON file.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := "%"
		if len(args) > 0 {
			pattern = args[0]
		}

		records, err := Store.SearchRecords(pattern)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			fmt.Fprintln(os.Stderr, "No records found to export.")
			return nil
		}

		var writer = os.Stdout
		if exportOutput != "" {
			file, err := os.Create(exportOutput)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer file.Close()
			writer = file
		}

		enc := json.NewEncoder(writer)
		enc.SetIndent("", "  ")
		if err := enc.Encode(records); err != nil {
			return fmt.Errorf("failed to encode records: %w", err)
		}

		if exportOutput != "" {
			fmt.Printf("Exported %d records to %s\n", len(records), exportOutput)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: stdout)")
}
