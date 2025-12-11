package model

import "time"

// Record represents a standard data record
// Corresponds to Record class in Perl6
type Record struct {
	Database     string    `json:"database"`     // Database name
	Table        string    `json:"table"`        // Table name
	Name         string    `json:"name"`         // Record name
	Type         string    `json:"type"`         // Record type
	Source       string    `json:"source"`       // Data source
	Label        string    `json:"label"`        // Record label
	Description  string    `json:"description"`  // Record description
	Number       int       `json:"number"`       // Number of records
	MissNumber   int       `json:"missNumber"`   // Number of missing values
	UniqueNumber int       `json:"uniqueNumber"` // Number of unique values
	ScriptFile   string    `json:"script_file"`  // Creation script
	ScriptTag    string    `json:"script_tag"`   // Tag of the creation script
	DescFile     string    `json:"desc_file"`    // Description file
	DescTag      string    `json:"desc_tag"`     // Tag of the description file
	LogFile      string    `json:"log_file"`     // Usage log file
	CreateAt     time.Time `json:"create_at"`    // Creation time
	ModifyAt     time.Time `json:"modify_at"`    // Modification time
}

// FullName returns the full name of the record
// Format: database:table:name
func (r *Record) FullName() string {
	return r.Database + ":" + r.Table + ":" + r.Name
}
