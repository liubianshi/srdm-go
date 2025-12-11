package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"srdm/internal/model"
	"testing"
)

func TestUpdate(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()

	// Pre-populate
	rec := &model.Record{Database: "db", Table: "t", Name: "rec1", Label: "old"}
	mockStore.InsertRecord(rec)

	rootCmd.SetArgs([]string{
		"update",
		"--name", "db:test_table:rec1", // This matches Mock Insert logical name if we constructed it right, wait.
		// In previous test we inserted directly to slice.
		// NOTE: In mock_test.go InsertRecord uses r.FullName() as key
		// db:t:rec1
	})
	
	// Let's use correct name
	rootCmd.SetArgs([]string{
		"update",
		"--name", "db:t:rec1",
		"--label", "new_label",
	})
	
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := mockStore.GetRecord("db:t:rec1")
	if updated.Label != "new_label" {
		t.Errorf("Expected new_label, got %s", updated.Label)
	}
}

func TestView(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()
	
	mockStore.InsertRecord(&model.Record{Database: "db", Table: "t", Name: "rec1", Label: "view_me"})

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"view", "db:t:rec1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("View failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	if !bytes.Contains(buf.Bytes(), []byte("view_me")) {
		t.Errorf("Output did not contain label view_me")
	}
}

func TestGet(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()
	
	// Create dummy file to get
	tmpFile, _ := os.CreateTemp("", "srdm_test")
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("content")
	tmpFile.Close()

	mockStore.InsertTable(&model.Table{
		Database: "db", Name: "t", Path: tmpFile.Name(),
	})

	outputFile := filepath.Join(os.TempDir(), "extracted_file")
	defer os.Remove(outputFile)

	rootCmd.SetArgs([]string{"get", "db:t", "--output", outputFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	content, _ := os.ReadFile(outputFile)
	if string(content) != "content" {
		t.Errorf("Content mismatch")
	}
}

func TestExport(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()
	
	mockStore.InsertRecord(&model.Record{Database: "db", Table: "t", Name: "rec1"})
	
	outputFile := filepath.Join(os.TempDir(), "export.json")
	defer os.Remove(outputFile)

	rootCmd.SetArgs([]string{"export", "db:t:%", outputFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	content, _ := os.ReadFile(outputFile)
	if !bytes.Contains(content, []byte("rec1")) {
		t.Errorf("Export file missing content")
	}
}

func TestInfo(t *testing.T) {
	mockStore := NewMockRepository()
	Store = mockStore
	defer func() { Store = nil }()

	mockStore.InsertTable(&model.Table{Database: "db", Name: "t"})
	
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"info"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	
	if !bytes.Contains(buf.Bytes(), []byte("Tables:")) {
		t.Errorf("Output missing 'Tables:'")
	}
	// Check mocked values
	// TableCount = 1 (inserted above)
	// DbSize = 1024 (hardcoded in mock)
	if !bytes.Contains(buf.Bytes(), []byte("1.0 KB")) {
		t.Errorf("Output missing size '1.0 KB', got: %s", output)
	}
}

