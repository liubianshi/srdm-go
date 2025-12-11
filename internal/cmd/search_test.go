package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"srdm/internal/model"
	"testing"
)

func TestSearchRecord(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()

	// Pre-populate store
	rec := &model.Record{
		Database: "db",
		Table:    "t",
		Name:     "rec1",
		Label:    "found_me",
	}
	// FullName => db:t:rec1
	if err := mockStore.InsertRecord(rec); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Search for exact match
	rootCmd.SetArgs([]string{
		"search",
		"db:t:rec1",
		"--format", "json",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	// Parse JSON output
	var results []model.Record
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		// Output might be list of interfaces, let's try generic
		var generic []interface{}
		if err2 := json.Unmarshal(buf.Bytes(), &generic); err2 != nil {
			t.Logf("Raw output: %s", buf.String())
			t.Fatalf("Failed to parse JSON: %v", err)
		}
		// Cast manually check
		if len(generic) != 1 {
			t.Errorf("Expected 1 result, got %d", len(generic))
		}
	} else {
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		} else {
			if results[0].Label != "found_me" {
				t.Errorf("Expected label found_me, got %s", results[0].Label)
			}
		}
	}
}
