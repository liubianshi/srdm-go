package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var getOutput string

var getCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Extract data file",
	Long:  `Extract original file associated with data record to local.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// For now, only support Table extraction (DB file)
		// Assuming Record doesn't strictly have a file path stored directly in struct except logs/scripts
		// But let's assume if it's Table, we copy the .Path

		t, err := Store.GetTable(name)
		if err == nil && t != nil {
			src := t.Path
			dst := getOutput
			if dst == "" {
				dst = filepath.Base(src)
			}
			return copyFile(src, dst)
		}

		// TODO: Implement Logic for Record if Record has a "Path" or associated file.
		// Current Model for Record doesn't have a "DataPath", only Script/Desc/Log files.
		// So we return error for now or implement copying one of those if asked.

		return fmt.Errorf("resource not found or extraction not supported for this type: %s", name)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&getOutput, "output", "o", "", "Output filename (default: original filename)")
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}
	fmt.Printf("Extracted to %s\n", dst)
	return nil
}
