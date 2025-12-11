package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [name]",
	Short: "View data record details",
	Long:  `View details of a specific data record or table.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		t, err := Store.GetTable(name)
		if err == nil && t != nil {
			fmt.Printf("Table: %s\n", t.FullName())
			fmt.Printf("  Database:    %s\n", t.Database)
			fmt.Printf("  Name:        %s\n", t.Name)
			fmt.Printf("  Keys:        %s\n", t.Keys)
			fmt.Printf("  Path:        %s\n", t.Path)
			fmt.Printf("  Engine:      %s\n", t.Engine)
			fmt.Printf("  Description: %s\n", t.Description)
			fmt.Printf("  Source:      %s\n", t.Source)
			fmt.Printf("  CreateAt:    %s\n", t.CreateAt)
			fmt.Printf("  ModifyAt:    %s\n", t.ModifyAt)
			fmt.Printf("  Records:     %d\n", len(t.Records))
			return nil
		}

		r, err := Store.GetRecord(name)
		if err == nil && r != nil {
			fmt.Printf("Record: %s\n", r.FullName())
			fmt.Printf("  Database:    %s\n", r.Database)
			fmt.Printf("  Table:       %s\n", r.Table)
			fmt.Printf("  Name:        %s\n", r.Name)
			fmt.Printf("  Type:        %s\n", r.Type)
			fmt.Printf("  Label:       %s\n", r.Label)
			fmt.Printf("  Source:      %s\n", r.Source)
			fmt.Printf("  Description: %s\n", r.Description)
			fmt.Printf("  Stats:       N=%d, Miss=%d, Unique=%d\n", r.Number, r.MissNumber, r.UniqueNumber)
			fmt.Printf("  CreateAt:    %s\n", r.CreateAt)
			fmt.Printf("  ModifyAt:    %s\n", r.ModifyAt)
			return nil
		}

		return fmt.Errorf("not found: %s", name)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
