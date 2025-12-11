package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestInsertTable(t *testing.T) {
	// Setup Mock
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }() // Reset after test

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command via rootCmd to ensure full chain
	// But we need to make sure args are passed correctly
	rootCmd.SetArgs([]string{
		"insert",
		"--name", "db:test_table",
		"--keys", "id",
		"--engine", "SQLite3",
	})
	
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Close writer and read output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify Output
	if buf.Len() == 0 {
		t.Log("Warning: No output captured")
	} else {
		t.Logf("Output: %s", output)
	}

	// Verify Store State
	foundTable, err := mockStore.GetTable("db:test_table")
	if err != nil {
		t.Fatalf("GetTable returned error: %v", err)
	}
	if foundTable == nil {
		t.Fatalf("Table not found in store")
	}
	if foundTable.Name != "test_table" {
		t.Errorf("Expected table name test_table, got %s", foundTable.Name)
	}
}

func TestInsertRecord(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()

	rootCmd.SetArgs([]string{
		"insert",
		"--name", "db:test_table:rec1",
		"--type", "file",
		"--label", "test_label",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	foundRecord, err := mockStore.GetRecord("db:test_table:rec1")
	if err != nil {
		t.Fatalf("GetRecord returned error: %v", err)
	}
	if foundRecord == nil {
		t.Fatalf("Record not found in store")
	}
	if foundRecord.Label != "test_label" {
		t.Errorf("Expected label test_label, got %s", foundRecord.Label)
	}
}
