package cmd

import (
	"fmt"
	"srdm/internal/model"
)

type MockRepository struct {
	Tables  map[string]*model.Table
	Records map[string]*model.Record
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Tables:  make(map[string]*model.Table),
		Records: make(map[string]*model.Record),
	}
}

func (m *MockRepository) InsertTable(t *model.Table) error {
	if _, exists := m.Tables[t.FullName()]; exists {
		return fmt.Errorf("table already exists")
	}
	m.Tables[t.FullName()] = t
	return nil
}

func (m *MockRepository) InsertRecord(r *model.Record) error {
	if _, exists := m.Records[r.FullName()]; exists {
		return fmt.Errorf("record already exists")
	}
	m.Records[r.FullName()] = r
	return nil
}

func (m *MockRepository) GetTable(name string) (*model.Table, error) {
	if t, exists := m.Tables[name]; exists {
		return t, nil
	}
	return nil, nil // Not found but no error
}

func (m *MockRepository) GetRecord(name string) (*model.Record, error) {
	if r, exists := m.Records[name]; exists {
		return r, nil
	}
	return nil, nil
}

func (m *MockRepository) UpdateTable(t *model.Table) error {
	if _, exists := m.Tables[t.FullName()]; !exists {
		return fmt.Errorf("table not found")
	}
	m.Tables[t.FullName()] = t
	return nil
}

func (m *MockRepository) UpdateRecord(r *model.Record) error {
	if _, exists := m.Records[r.FullName()]; !exists {
		return fmt.Errorf("record not found")
	}
	m.Records[r.FullName()] = r
	return nil
}

func (m *MockRepository) GetStatistics() (*model.Stats, error) {
	return &model.Stats{
		Path:        "mock.sqlite",
		TableCount:  len(m.Tables),
		RecordCount: len(m.Records),
		DbSize:      1024,
	}, nil
}



func (m *MockRepository) SearchRecords(pattern string) ([]model.Record, error) {
	var records []model.Record
	// Simple prefix match simulation for testing
	// In real DB it's LIKE 'pattern%'
	// Here assume pattern ends with %
	prefix := pattern
	if len(pattern) > 0 && pattern[len(pattern)-1] == '%' {
		prefix = pattern[:len(pattern)-1]
	}

	for _, r := range m.Records {
		if len(r.FullName()) >= len(prefix) && r.FullName()[:len(prefix)] == prefix {
			records = append(records, *r)
		}
	}
	return records, nil
}

func (m *MockRepository) Delete(name string, force bool) error {
	// Try as table
	if _, exists := m.Tables[name]; exists {
		if !force {
			return fmt.Errorf("cannot delete table %s without force flag", name)
		}
		delete(m.Tables, name)
		// Delete children (very simple impl)
		for k := range m.Records {
			if len(k) > len(name) && k[:len(name)+1] == name+":" {
				delete(m.Records, k)
			}
		}
		return nil
	}

	// Try as record
	if _, exists := m.Records[name]; exists {
		delete(m.Records, name)
		return nil
	}
	return nil
}

func (m *MockRepository) Close() error {
	return nil
}

func (m *MockRepository) Ping() error {
	return nil
}

func (m *MockRepository) GetPath() string {
	return "mock.sqlite"
}

