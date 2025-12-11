package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateName         string
	updateKeys         string
	updateEngine       string
	updatePath         string
	updateSource       string
	updateDesc         string
	updateScriptFile   string
	updateScriptTag    string
	updateDescFile     string
	updateDescTag      string
	updateLogFile      string
	updateType         string
	updateLabel        string
	updateNumber       int
	updateMissNumber   int
	updateUniqueNumber int
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update data record or table",
	Long:  `Update existing data record or table.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateName == "" {
			return fmt.Errorf("--name is required")
		}

		parts := strings.Split(updateName, ":")
		if len(parts) == 2 {
			return updateTable(parts[0], parts[1])
		} else if len(parts) == 3 {
			return updateRecord(parts[0], parts[1], parts[2])
		} else {
			return fmt.Errorf("invalid name format")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&updateName, "name", "", "Record Name")
	updateCmd.Flags().StringVar(&updateKeys, "keys", "", "Primary keys of the table")
	updateCmd.Flags().StringVar(&updateEngine, "engine", "", "Data management engine")
	updateCmd.Flags().StringVar(&updatePath, "data-path", "", "Data storage path")
	updateCmd.Flags().StringVar(&updateSource, "source", "", "Data source")
	updateCmd.Flags().StringVar(&updateDesc, "description", "", "Data description")
	updateCmd.Flags().StringVar(&updateScriptFile, "script_file", "", "Data processing script file")
	updateCmd.Flags().StringVar(&updateScriptTag, "script_tag", "", "Script file version tag")
	updateCmd.Flags().StringVar(&updateDescFile, "desc_file", "", "Data analysis file")
	updateCmd.Flags().StringVar(&updateDescTag, "desc_tag", "", "Analysis file version tag")
	updateCmd.Flags().StringVar(&updateLogFile, "log_file", "", "Data usage log file")

	updateCmd.Flags().StringVar(&updateType, "type", "", "Record type")
	updateCmd.Flags().StringVar(&updateLabel, "label", "", "Data label")
	updateCmd.Flags().IntVar(&updateNumber, "number", 0, "Number of records")
	updateCmd.Flags().IntVar(&updateMissNumber, "missNumber", 0, "Number of missing values")
	updateCmd.Flags().IntVar(&updateUniqueNumber, "uniqueNumber", 0, "Number of unique values")
}

func updateTable(database, name string) error {
	fullName := database + ":" + name
	t, err := Store.GetTable(fullName)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("table not found: %s", fullName)
	}

	// Update fields if provided (not empty or 0)
	if updateKeys != "" {
		t.Keys = updateKeys
	}
	if updateEngine != "" {
		t.Engine = updateEngine
	}
	if updatePath != "" {
		t.Path = updatePath
	}
	if updateSource != "" {
		t.Source = updateSource
	}
	if updateDesc != "" {
		t.Description = updateDesc
	}
	if updateScriptFile != "" {
		t.ScriptFile = updateScriptFile
	}
	if updateScriptTag != "" {
		t.ScriptTag = updateScriptTag
	}
	if updateDescFile != "" {
		t.DescFile = updateDescFile
	}
	if updateDescTag != "" {
		t.DescTag = updateDescTag
	}
	if updateLogFile != "" {
		t.LogFile = updateLogFile
	}

	if err := Store.UpdateTable(t); err != nil {
		return err
	}
	fmt.Printf("Updated table: %s\n", t.FullName())
	return nil
}

func updateRecord(database, table, name string) error {
	fullName := database + ":" + table + ":" + name
	r, err := Store.GetRecord(fullName)
	if err != nil {
		return err
	}
	if r == nil {
		return fmt.Errorf("record not found: %s", fullName)
	}

	if updateType != "" {
		r.Type = updateType
	}
	if updateSource != "" {
		r.Source = updateSource
	}
	if updateLabel != "" {
		r.Label = updateLabel
	}
	if updateDesc != "" {
		r.Description = updateDesc
	}
	if updateNumber != 0 {
		r.Number = updateNumber
	}
	if updateMissNumber != 0 {
		r.MissNumber = updateMissNumber
	}
	if updateUniqueNumber != 0 {
		r.UniqueNumber = updateUniqueNumber
	}
	if updateScriptFile != "" {
		r.ScriptFile = updateScriptFile
	}
	if updateScriptTag != "" {
		r.ScriptTag = updateScriptTag
	}
	if updateDescFile != "" {
		r.DescFile = updateDescFile
	}
	if updateDescTag != "" {
		r.DescTag = updateDescTag
	}
	if updateLogFile != "" {
		r.LogFile = updateLogFile
	}

	if err := Store.UpdateRecord(r); err != nil {
		return err
	}
	fmt.Printf("Updated record: %s\n", r.FullName())
	return nil
}
