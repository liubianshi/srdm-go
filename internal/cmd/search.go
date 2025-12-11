package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	searchMode       string
	searchFormat     string
	searchOutputFile string
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [names]",
	Short: "Query data records",
	Long:  `Query data records. Supports exact match by name or fuzzy search.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []interface{}

		if len(args) > 0 {
			for _, name := range args {
				// Try fetching as Table
				t, err := Store.GetTable(name)
				if err == nil && t != nil {
					results = append(results, t)
					continue
				}

				// Try fetching as Record
				r, err := Store.GetRecord(name)
				if err == nil && r != nil {
					results = append(results, r)
					continue
				}

				// Try fuzzy search
				records, err := Store.SearchRecords(name + "%")
				if err == nil && len(records) > 0 {
					for _, rec := range records {
						results = append(results, rec)
					}
					continue
				}

				fmt.Fprintf(os.Stderr, "not found: %s\n", name)
			}
		} else {
			// If no args provided, show help
			return cmd.Help()
		}

		if len(results) == 0 {
			return nil
		}

		// Output handling
		if searchFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(results); err != nil {
				return err
			}
		} else {
			// Default text output (simplified)
			for _, res := range results {
				fmt.Printf("%+v\n", res)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchMode, "mode", "detail", "Display mode (detail, name-only, oneline)")
	searchCmd.Flags().StringVar(&searchFormat, "format", "json", "Output format (json, text)")
	searchCmd.Flags().StringVar(&searchOutputFile, "output-file", "", "Output file")
}
