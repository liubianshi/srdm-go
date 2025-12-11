package model

import "time"

// Table represents a data table
// Corresponds to Table class in Perl6
type Table struct {
	Database    string    `json:"database"`    // Database name
	Name        string    `json:"name"`        // Table name
	Keys        string    `json:"keys"`        // Primary keys of the table
	Path        string    `json:"path"`        // Data location
	Engine      string    `json:"engine"`      // Database file management engine (default SQLite3)
	Source      string    `json:"source"`      // Data source of the record
	Description string    `json:"description"` // Record description
	ScriptFile  string    `json:"script_file"` // Record creation script
	ScriptTag   string    `json:"script_tag"`  // Tag of the creation script
	DescFile    string    `json:"desc_file"`   // Description file
	DescTag     string    `json:"desc_tag"`    // Tag of the description file
	LogFile     string    `json:"log_file"`    // Usage log file
	CreateAt    time.Time `json:"create_at"`   // Creation time
	ModifyAt    time.Time `json:"modify_at"`   // Modification time
	Records     []Record  `json:"records"`     // List of included records
}

// FullName returns the full name of the table
// Format: database:name
func (t *Table) FullName() string {
	return t.Database + ":" + t.Name
}
