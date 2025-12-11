package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"srdm/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps sql.DB object and provides database operations
// It maintains a reference to the underlying sql.DB connection
// and the file path to the SQLite database
type DB struct {
	*sql.DB
	Path string
}

// GetPath returns the file system path to the database
func (db *DB) GetPath() string {
	return db.Path
}

// GetStatistics retrieves comprehensive database statistics including
// file size, table/record counts, SQLite version, and last update time
// Returns a Stats object containing all collected metrics
func (db *DB) GetStatistics() (*model.Stats, error) {
	stats := &model.Stats{
		Path: db.Path,
	}

	// Retrieve database file size from filesystem
	if fi, err := os.Stat(db.Path); err == nil {
		stats.DbSize = fi.Size()
	}

	// Query counts and version information in a single transaction for consistency
	queries := []struct {
		sql    string
		dest   interface{}
		errMsg string
	}{
		{"SELECT COUNT(*) FROM data_table", &stats.TableCount, "failed to count tables"},
		{"SELECT COUNT(*) FROM data_record", &stats.RecordCount, "failed to count records"},
		{"SELECT sqlite_version()", &stats.SqliteVersion, "failed to get sqlite version"},
	}

	for _, q := range queries {
		if err := db.QueryRow(q.sql).Scan(q.dest); err != nil {
			return nil, fmt.Errorf("%s: %w", q.errMsg, err)
		}
	}

	// Determine the most recent modification time across both tables
	var lastTableModify, lastRecordModify sql.NullTime

	// Query last modification times (ignore errors for empty tables)
	_ = db.QueryRow("SELECT MAX(modify_at) FROM data_table").Scan(&lastTableModify)
	_ = db.QueryRow("SELECT MAX(modify_at) FROM data_record").Scan(&lastRecordModify)

	// Set LastUpdated to the most recent modification time
	if lastTableModify.Valid {
		stats.LastUpdated = lastTableModify.Time
	}
	if lastRecordModify.Valid && lastRecordModify.Time.After(stats.LastUpdated) {
		stats.LastUpdated = lastRecordModify.Time
	}

	// Retrieve list of all table names sorted alphabetically
	rows, err := db.Query("SELECT name FROM data_table ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, name)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table names: %w", err)
	}

	stats.TablesList = tables

	return stats, nil
}

// NewDB creates and initializes a new DB instance
// Creates the database directory if it doesn't exist
// Initializes the database schema with required tables and indexes
// dbPath: The file system path where the SQLite database will be created/opened
func NewDB(dbPath string) (*DB, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Open SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify database connectivity
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize DB wrapper
	sdb := &DB{DB: db, Path: dbPath}

	// Create schema tables and indexes
	if err := sdb.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return sdb, nil
}

// initSchema creates the required database tables and indexes
// Creates data_table for storing table metadata
// Creates data_record for storing record/column metadata
// Both tables include automatic timestamp tracking
func (db *DB) initSchema() error {
	// Define schema for data_table (stores table-level metadata)
	tableSchema := `
	CREATE TABLE IF NOT EXISTS data_table (
		name            VARCHAR PRIMARY KEY,
		keys            VARCHAR NOT NULL,
		path            VARCHAR NOT NULL,
		engine          VARCHAR NOT NULL DEFAULT 'SQLite3',
		source          VARCHAR,
		description     VARCHAR,
		script_file     VARCHAR,
		script_tag      VARCHAR,
		desc_file       VARCHAR,
		desc_tag        VARCHAR,
		log_file        VARCHAR,
		create_at       TIMESTAMP NOT NULL DEFAULT (DATETIME('NOW', 'LOCALTIME')),
		modify_at       TIMESTAMP NOT NULL DEFAULT (DATETIME('NOW', 'LOCALTIME'))
	);
	CREATE INDEX IF NOT EXISTS data_table_name ON data_table (name);
	`
	if _, err := db.Exec(tableSchema); err != nil {
		return fmt.Errorf("failed to create data_table: %w", err)
	}

	// Define schema for data_record (stores record/column-level metadata)
	recordSchema := `
	CREATE TABLE IF NOT EXISTS data_record (
		name         VARCHAR PRIMARY KEY,
		type         VARCHAR NOT NULL,
		source       VARCHAR NOT NULL DEFAULT 'unknown',
		label        VARCHAR NOT NULL,
		description  VARCHAR,
		number       INTEGER,
		missNumber   INTEGER,
		uniqueNumber INTEGER,
		script_file  VARCHAR,
		script_tag   VARCHAR,
		desc_file    VARCHAR,
		desc_tag     VARCHAR,
		log_file     VARCHAR,
		create_at    TIMESTAMP NOT NULL DEFAULT (DATETIME('NOW', 'LOCALTIME')),
		modify_at    TIMESTAMP NOT NULL DEFAULT (DATETIME('NOW', 'LOCALTIME'))
	);
	CREATE INDEX IF NOT EXISTS data_record_name ON data_record (name);
	`
	if _, err := db.Exec(recordSchema); err != nil {
		return fmt.Errorf("failed to create data_record: %w", err)
	}

	return nil
}
