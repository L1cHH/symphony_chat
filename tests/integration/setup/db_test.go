package setup

import "testing"

func TestNewTestDB(t *testing.T) {
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("failed to create test database connection: %v", err)
	}
	
	defer db.Close()

	var result int
	err = db.DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if result != 1 {
		t.Errorf("expected result to be 1, got %d", result)
	}
}