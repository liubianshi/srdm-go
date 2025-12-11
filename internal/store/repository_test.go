package store

import (
	"path/filepath"
	"srdm/internal/model"
	"testing"
	"time"
)

// setupTestDB creates a test database
// Use t.TempDir() to ensure each test runs in a clean environment
func setupTestDB(t *testing.T) *DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestInsertAndGetTable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	table := &model.Table{
		Database:    "testdb",
		Name:        "testtable",
		Keys:        "id",
		Path:        "/path/to/data",
		Engine:      "SQLite3",
		Source:      "test",
		Description: "A test table",
		CreateAt:    time.Now(),
		ModifyAt:    time.Now(),
	}

	// Test insert
	if err := db.InsertTable(table); err != nil {
		t.Fatalf("InsertTable failed: %v", err)
	}

	// Test get
	retrieved, err := db.GetTable("testdb:testtable")
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetTable returned nil")
	}
	if retrieved.Name != table.Name {
		t.Errorf("Expected table name %s, got %s", table.Name, retrieved.Name)
	}
}

func TestInsertAndGetRecord(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	record := &model.Record{
		Database:    "db1",
		Table:       "tbl1",
		Name:        "rec1",
		Type:        "test_type",
		Source:      "generated",
		Label:       "label1",
		Description: "desc1",
		CreateAt:    time.Now(),
		ModifyAt:    time.Now(),
	}

	if err := db.InsertRecord(record); err != nil {
		t.Fatalf("InsertRecord failed: %v", err)
	}

	retrieved, err := db.GetRecord("db1:tbl1:rec1")
	if err != nil {
		t.Fatalf("GetRecord failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetRecord returned nil")
	}
	if retrieved.Name != record.Name {
		t.Errorf("Expected record name %s, got %s", record.Name, retrieved.Name)
	}
}

func TestSearchRecords(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	records := []model.Record{
		{Database: "db1", Table: "tbl1", Name: "rec_alpha", Type: "t1"},
		{Database: "db1", Table: "tbl1", Name: "rec_beta", Type: "t1"},
		{Database: "db1", Table: "tbl1", Name: "other_gamma", Type: "t2"},
	}

	for _, r := range records {
		if err := db.InsertRecord(&r); err != nil {
			t.Fatalf("InsertRecord failed: %v", err)
		}
	}

	// Fuzzy search "rec_%"
	results, err := db.SearchRecords("db1:tbl1:rec_%")
	if err != nil {
		t.Fatalf("SearchRecords failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 records, got %d", len(results))
	}

	// Fuzzy search "other%"
	results2, err := db.SearchRecords("db1:tbl1:other%")
	if err != nil {
		t.Fatalf("SearchRecords failed: %v", err)
	}
	if len(results2) != 1 {
		t.Errorf("Expected 1 record, got %d", len(results2))
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	table := &model.Table{
		Database: "db1",
		Name:     "tbl1",
		Keys:     "id",
	}
	record := &model.Record{
		Database: "db1",
		Table:    "tbl1",
		Name:     "rec1",
	}

	// Insert data
	db.InsertTable(table)
	db.InsertRecord(record)

	// Test delete record
	if err := db.Delete("db1:tbl1:rec1", false); err != nil {
		t.Fatalf("Delete record failed: %v", err)
	}
	r, _ := db.GetRecord("db1:tbl1:rec1")
	if r != nil {
		t.Error("Record should be deleted")
	}

	// Insert another record for cascade delete test
	db.InsertRecord(record)

	// Test cascade delete table (requires force=true)
	// Should fail if force=false
	if err := db.Delete("db1:tbl1", false); err == nil {
		t.Error("Delete table without force should fail")
	}

	if err := db.Delete("db1:tbl1", true); err != nil {
		t.Fatalf("Delete table with force failed: %v", err)
	}

	// Verify table and records are deleted
	tbl, _ := db.GetTable("db1:tbl1")
	if tbl != nil {
		t.Error("Table should be deleted")
	}
	r, _ = db.GetRecord("db1:tbl1:rec1")
	if r != nil {
		t.Error("Child record should be deleted")
	}
}
