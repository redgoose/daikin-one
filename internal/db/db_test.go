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

	now := time.Now()
	startDate := now.AddDate(0, 0, -7)

	result := GetDataRaw(dbPath, deviceId, "outdoor_heat", startDate, now)

	if len(result) == 0 {
		t.Fatalf("No results")
	}
}

func TestGetDataForDay(t *testing.T) {
	dbPath := getDbPath()
	result := GetDataForDay(dbPath, deviceId, time.Now().AddDate(0, 0, -1))

	if len(result) == 0 {
		t.Fatalf("No results")
	}
}
