package db

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)


const deviceId = "FILL ME IN"

func getDbPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dbPath := filepath.Join(homeDir, ".daikin", "daikin.db")
	return dbPath
}

func TestGetDataRaw(t *testing.T) {
	dbPath := getDbPath()

	result := GetDataRaw(dbPath, deviceId, "outdoor_heat")

	if len(result) == 0 {
		t.Fatalf("No results")
	}
}

func TestGetDataForDay(t *testing.T) {
	dbPath := getDbPath()
	result := GetDataForDay(dbPath, deviceId, time.Now())

	if len(result) == 0 {
		t.Fatalf("No results")
	}
}
