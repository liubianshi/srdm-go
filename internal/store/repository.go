package store

import (
	"database/sql"
	"fmt"
	"srdm/internal/model"
	"strings"
)

// Repository defines the data storage interface
type Repository interface {
	InsertTable(t *model.Table) error
	InsertRecord(r *model.Record) error
	GetTable(name string) (*model.Table, error)
	GetRecord(name string) (*model.Record, error)
	UpdateTable(t *model.Table) error
	UpdateRecord(r *model.Record) error
	GetStatistics() (*model.Stats, error)
	SearchRecords(pattern string) ([]model.Record, error)
	Delete(name string, force bool) error
	Close() error
	Ping() error
	GetPath() string
}

// InsertTable inserts a table record
func (db *DB) InsertTable(t *model.Table) error {
	query := `
	INSERT INTO data_table (
		name, keys, path, engine, source, description,
		script_file, script_tag, desc_file, desc_tag, log_file,
		create_at, modify_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := db.Exec(query,
		t.FullName(), t.Keys, t.Path, t.Engine, t.Source, t.Description,
		t.ScriptFile, t.ScriptTag, t.DescFile, t.DescTag, t.LogFile,
		t.CreateAt, t.ModifyAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert table: %w", err)
	}

	// Insert associated records
	for _, r := range t.Records {
		if err := db.InsertRecord(&r); err != nil {
			return err
		}
	}
	return nil
}

// InsertRecord inserts a regular record
func (db *DB) InsertRecord(r *model.Record) error {
	// First check if the table exists (Perl6 logic)
	// Here we assume business logic handles this, or check by Name simply
	// FullName format: database:table:name

	query := `
	INSERT INTO data_record (
		name, type, source, label, description,
		number, missNumber, uniqueNumber,
		script_file, script_tag, desc_file, desc_tag, log_file,
		create_at, modify_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := db.Exec(query,
		r.FullName(), r.Type, r.Source, r.Label, r.Description,
		r.Number, r.MissNumber, r.UniqueNumber,
		r.ScriptFile, r.ScriptTag, r.DescFile, r.DescTag, r.LogFile,
		r.CreateAt, r.ModifyAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	return nil
}

// GetTable retrieves a table by name
func (db *DB) GetTable(name string) (*model.Table, error) {
	query := `SELECT * FROM data_table WHERE name = ?`
	row := db.QueryRow(query, name)

	var t model.Table
	var fullName string
	err := row.Scan(
		&fullName, &t.Keys, &t.Path, &t.Engine, &t.Source, &t.Description,
		&t.ScriptFile, &t.ScriptTag, &t.DescFile, &t.DescTag, &t.LogFile,
		&t.CreateAt, &t.ModifyAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	// Parse FullName into Database and Name
	parts := strings.Split(fullName, ":")
	if len(parts) >= 2 {
		t.Database = parts[0]
		t.Name = parts[1]
	}

	// Get associated records
	// Associated records Name wildcard match: full_table_name:%
	records, err := db.SearchRecords(fullName + ":%")
	if err != nil {
		return nil, err
	}
	t.Records = records

	return &t, nil
}

// GetRecord retrieves a record by name
func (db *DB) GetRecord(name string) (*model.Record, error) {
	query := `SELECT * FROM data_record WHERE name = ?`
	row := db.QueryRow(query, name)

	var r model.Record
	var fullName string
	err := row.Scan(
		&fullName, &r.Type, &r.Source, &r.Label, &r.Description,
		&r.Number, &r.MissNumber, &r.UniqueNumber,
		&r.ScriptFile, &r.ScriptTag, &r.DescFile, &r.DescTag, &r.LogFile,
		&r.CreateAt, &r.ModifyAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan record: %w", err)
	}

	parts := strings.Split(fullName, ":")
	if len(parts) >= 3 {
		r.Database = parts[0]
		r.Table = parts[1]
		r.Name = parts[2]
	}

	return &r, nil
}

// SearchRecords searches records (simple LIKE implementation)
func (db *DB) SearchRecords(pattern string) ([]model.Record, error) {
	query := `SELECT * FROM data_record WHERE name LIKE ?`
	rows, err := db.Query(query, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.Record
	for rows.Next() {
		var r model.Record
		var fullName string
		err := rows.Scan(
			&fullName, &r.Type, &r.Source, &r.Label, &r.Description,
			&r.Number, &r.MissNumber, &r.UniqueNumber,
			&r.ScriptFile, &r.ScriptTag, &r.DescFile, &r.DescTag, &r.LogFile,
			&r.CreateAt, &r.ModifyAt,
		)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(fullName, ":")
		if len(parts) >= 3 {
			r.Database = parts[0]
			r.Table = parts[1]
			r.Name = parts[2]
		}
		records = append(records, r)
	}
	return records, nil
}

// Delete removes a record or table
// force: if it is a table, force remove all its records
func (db *DB) Delete(name string, force bool) error {
	// Try finding as table first
	t, err := db.GetTable(name)
	if err != nil {
		return err
	}
	if t != nil {
		if !force {
			return fmt.Errorf("cannot delete table %s without force flag", name)
		}
		// Delete all sub-records
		// Roughly using LIKE here
		if _, err := db.Exec("DELETE FROM data_record WHERE name LIKE ?", name+":%"); err != nil {
			return err
		}
		// Delete table
		if _, err := db.Exec("DELETE FROM data_table WHERE name = ?", name); err != nil {
			return err
		}
		return nil
	}

	// Try finding as record and remove
	if _, err := db.Exec("DELETE FROM data_record WHERE name = ?", name); err != nil {
		return err
	}
	return nil
}
