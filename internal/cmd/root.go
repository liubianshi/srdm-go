package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"srdm/internal/store"

	"github.com/spf13/cobra"
)

var (
	// Store is the global repository instance
	Store store.Repository
	// DataRepoPath data repository path
	DataRepoPath string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "srdm",
	Short: "Simple Research Data Manager",
	Long: `SRDM (Simple Research Data Manager) is a simple research data management tool.
It supports data insertion, update, deletion, viewing, extraction, query, and export.`,
	// PersistentPreRun runs before any subcommand
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// If Store is mock or already initialized, skip
		if Store != nil {
			return nil
		}

		// Initialize Logger
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		slog.SetDefault(logger)

		var err error
		// If flag is not set, try getting from environment variable
		if DataRepoPath == "" {
			DataRepoPath = os.Getenv("SRDM_DATA_REPO_PATH")
		}
		// If env is also empty, use default path
		if DataRepoPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			DataRepoPath = filepath.Join(home, "Data", "SRDM", "srdm_dataRepo.sqlite")
		}

		// Initialize database connection
		db, err := store.NewDB(DataRepoPath)
		if err != nil {
			return fmt.Errorf("could not initialize database at %s: %w", DataRepoPath, err)
		}
		Store = db
		return nil
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Define global flags
	rootCmd.PersistentFlags().StringVar(&DataRepoPath, "path", "", "Data storage location (default: $HOME/Documents/SRDM/srdm_dataRepo.sqlite)")
}
