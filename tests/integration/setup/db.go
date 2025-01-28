package setup

import (
	"database/sql"
	"fmt"
	"os"
	database "symphony_chat/internal/infrastructure/database"
)

type TestDB struct {
	DB *sql.DB
}

var TestDBInstance *TestDB

func NewTestDB() (*TestDB, error) {

	if TestDBInstance != nil {
		return TestDBInstance, nil
	}

	config := database.PostgresConfig {
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("TEST_DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		DBName:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  os.Getenv("TEST_DB_SSLMODE"),
	}

	db, err := database.NewPostgresConnection(config)
	if err != nil {
		return nil, err
	}

	TestDBInstance = &TestDB{
		DB: db,
	}

	return TestDBInstance, nil
}

func (tbd *TestDB) TruncateAllTables() error {
	rows, err := tbd.DB.Query(`
		SELECT table_name 
		FROM pg_tables
		WHERE table_schema = 'public'
	`)

	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate over rows: %w", err)
	}

	if _, err := tbd.DB.Exec("SET CONSTRAINTS ALL DEFERRED"); err != nil {
		return fmt.Errorf("failed to set constraints: %w", err)
	}

	for _, table := range tables {
		_, err := tbd.DB.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			return fmt.Errorf("failed to truncate table: %w", err)
		}
	}

	return nil
}

func (tdb *TestDB) Close() error {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
	return nil
}