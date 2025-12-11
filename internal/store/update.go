package store

import (
	"fmt"
	"srdm/internal/model"
)

// UpdateTable updates table information
func (db *DB) UpdateTable(t *model.Table) error {
	query := `
	UPDATE data_table SET 
		keys = ?, path = ?, engine = ?, source = ?, description = ?,
		script_file = ?, script_tag = ?, desc_file = ?, desc_tag = ?, log_file = ?,
		modify_at = DATETIME('NOW', 'LOCALTIME')
	WHERE name = ?;
	`
	res, err := db.Exec(query,
		t.Keys, t.Path, t.Engine, t.Source, t.Description,
		t.ScriptFile, t.ScriptTag, t.DescFile, t.DescTag, t.LogFile,
		t.FullName(), // Name in DB includes db prefix if we stored it that way, but wait.
	)
	if err != nil {
		return fmt.Errorf("failed to update table: %w", err)
	}

	// Check if any row updated
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("table not found: %s", t.FullName())
	}
	return nil
}

// UpdateRecord updates record information
func (db *DB) UpdateRecord(r *model.Record) error {
	query := `
	UPDATE data_record SET 
		type = ?, source = ?, label = ?, description = ?,
		number = ?, missNumber = ?, uniqueNumber = ?,
		script_file = ?, script_tag = ?, desc_file = ?, desc_tag = ?, log_file = ?,
		modify_at = DATETIME('NOW', 'LOCALTIME')
	WHERE name = ?;
	`
	res, err := db.Exec(query,
		r.Type, r.Source, r.Label, r.Description,
		r.Number, r.MissNumber, r.UniqueNumber,
		r.ScriptFile, r.ScriptTag, r.DescFile, r.DescTag, r.LogFile,
		r.FullName(),
	)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("record not found: %s", r.FullName())
	}
	return nil
}
