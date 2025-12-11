package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"srdm/internal/model"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	// Insert options
	insertName         string
	insertKeys         string
	insertEngine       string
	insertPath         string
	insertSource       string
	insertDesc         string
	insertScriptFile   string
	insertScriptTag    string
	insertDescFile     string
	insertDescTag      string
	insertLogFile      string
	insertType         string
	insertLabel        string
	insertNumber       int
	insertMissNumber   int
	insertUniqueNumber int
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert data record or table",
	Long:  `Insert a new data record or data table into the repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if insertName == "" {
			return fmt.Errorf("--name is required")
		}

		// Parse name to decide if it is a Table or Record
		// database:table => Table
		// database:table:record => Record
		parts := strings.Split(insertName, ":")
		if len(parts) == 2 {
			return insertTable(parts[0], parts[1])
		} else if len(parts) == 3 {
			return insertRecord(parts[0], parts[1], parts[2])
		} else {
			return fmt.Errorf("invalid name format. Use 'db:table' for table or 'db:table:record' for record")
		}
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVar(&insertName, "name", "", "Name of the record to insert (format: db:table or db:table:record)")
	insertCmd.Flags().StringVar(&insertKeys, "keys", "", "Primary keys of the table (Table only)")
	insertCmd.Flags().StringVar(&insertEngine, "engine", "SQLite3", "Data management engine (Table only)")
	insertCmd.Flags().StringVar(&insertPath, "data-path", "", "Data storage path (Table only)")
	insertCmd.Flags().StringVar(&insertSource, "source", "", "Data source")
	insertCmd.Flags().StringVar(&insertDesc, "description", "", "Data description")
	insertCmd.Flags().StringVar(&insertScriptFile, "script_file", "", "Data processing script file")
	insertCmd.Flags().StringVar(&insertScriptTag, "script_tag", "", "Script file version tag")
	insertCmd.Flags().StringVar(&insertDescFile, "desc_file", "", "Data analysis file")
	insertCmd.Flags().StringVar(&insertDescTag, "desc_tag", "", "Analysis file version tag")
	insertCmd.Flags().StringVar(&insertLogFile, "log_file", "", "Data usage log file")

	// Record specific options
	insertCmd.Flags().StringVar(&insertType, "type", "", "Record type")
	insertCmd.Flags().StringVar(&insertLabel, "label", "", "Data label")
	insertCmd.Flags().IntVar(&insertNumber, "number", 0, "Number of records")
	insertCmd.Flags().IntVar(&insertMissNumber, "missNumber", 0, "Number of missing values")
	insertCmd.Flags().IntVar(&insertUniqueNumber, "uniqueNumber", 0, "Number of unique values")
}

func insertTable(database, name string) error {
	if insertKeys == "" {
		return fmt.Errorf("--keys is required for table")
	}

	// Default path logic
	dataPath := insertPath
	if dataPath == "" {
		// Default to $HOME/DATA/DBMS/database.sqlite
		home, _ := os.UserHomeDir()
		dataPath = filepath.Join(home, "DATA", "DBMS", database+".sqlite")
	}

	table := &model.Table{
		Database:    database,
		Name:        name,
		Keys:        insertKeys,
		Path:        dataPath,
		Engine:      insertEngine,
		Source:      insertSource,
		Description: insertDesc,
		ScriptFile:  insertScriptFile,
		ScriptTag:   insertScriptTag,
		DescFile:    insertDescFile,
		DescTag:     insertDescTag,
		LogFile:     insertLogFile,
		CreateAt:    time.Now(),
		ModifyAt:    time.Now(),
	}

	if err := Store.InsertTable(table); err != nil {
		return err
	}
	fmt.Printf("Inserted table: %s\n", table.FullName())
	return nil
}

func insertRecord(database, table, name string) error {
	record := &model.Record{
		Database:     database,
		Table:        table,
		Name:         name,
		Type:         insertType,
		Source:       insertSource,
		Label:        insertLabel,
		Description:  insertDesc,
		Number:       insertNumber,
		MissNumber:   insertMissNumber,
		UniqueNumber: insertUniqueNumber,
		ScriptFile:   insertScriptFile,
		ScriptTag:    insertScriptTag,
		DescFile:     insertDescFile,
		DescTag:      insertDescTag,
		LogFile:      insertLogFile,
		CreateAt:     time.Now(),
		ModifyAt:     time.Now(),
	}

	if err := Store.InsertRecord(record); err != nil {
		return err
	}
	fmt.Printf("Inserted record: %s\n", record.FullName())
	return nil
}
