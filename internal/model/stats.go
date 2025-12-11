package model

import "time"

// Stats contains database usage statistics
type Stats struct {
	Path        string `json:"path"`
	TableCount  int    `json:"table_count"`
	RecordCount int    `json:"record_count"`
	DbSize        int64     `json:"db_size"` // in bytes
	LastUpdated   time.Time `json:"last_updated"`
	SqliteVersion string    `json:"sqlite_version"`
	TablesList    []string  `json:"tables_list"`
}
