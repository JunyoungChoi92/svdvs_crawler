package db

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_Connect(t *testing.T) {
	// Create mock SQL DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := &Database{conn: db}
	mock.ExpectPing()

	config := &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "postgres",
		Sslmode:  "disable",
	}

	err = database.Connect(config)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLoadYaml(t *testing.T) {
	expectedConfig := &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "postgres",
		Sslmode:  "disable",
	}

	actualConfig, err := loadConfig("../../db_config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config.yaml file: %s", err)
	}

	// Compare the two Config structs
	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf("Loaded config does not match the expected config.\nExpected: %+v\nActual: %+v", expectedConfig, actualConfig)
	}

	// Or using testify assert
	assert.Equal(t, expectedConfig, actualConfig, "Loaded config does not match the expected config.")
}
