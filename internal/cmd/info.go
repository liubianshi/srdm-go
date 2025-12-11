package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show database statistics",
	Long:  `Display summary statistics of the currently connected database, including table count, record count, and file size.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := Store.GetStatistics()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Helper closure for printing rows
		printRow := func(label string, value any) {
			fmt.Fprintf(w, "%s\t%v\n", Colorize(Cyan, label+":"), value)
		}

		printRow("Database Path", stats.Path)
		printRow("SQLite Version", stats.SqliteVersion)
		printRow("Tables", stats.TableCount)
		printRow("Records", stats.RecordCount)
		printRow("Size", formatBytes(stats.DbSize))

		lastUpdatedStr := "Never"
		if !stats.LastUpdated.IsZero() {
			lastUpdatedStr = stats.LastUpdated.Format(time.RFC3339)
		}
		printRow("Last Updated", lastUpdatedStr)

		printRow("Check Time", time.Now().Format(time.RFC3339))

		if len(stats.TablesList) > 0 {
			fmt.Fprintln(w, Colorize(Cyan, "Tables List:"))
			for _, t := range stats.TablesList {
				fmt.Fprintf(w, "  - %s\n", t)
			}
		}

		w.Flush()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
