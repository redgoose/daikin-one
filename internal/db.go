package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Metrics struct {
	DeviceId        string
	TempIndoor      float32
	TempOutdoor     float32
	CoolSetpoint    float32
	HeatSetpoint    float32
	EquipmentStatus int
}

func LogMetrics(dbPath string, metrics Metrics) {
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO daikin(timestamp, device_id, temp_indoor, temp_outdoor, cool_setpoint, heat_setpoint, equipment_status) values(?,?,?,?,?,?,?)")
	checkErr(err)

	var timestamp string = time.Now().Format(time.RFC3339)

	_, err = stmt.Exec(
		timestamp,
		metrics.DeviceId,
		metrics.TempIndoor,
		metrics.TempOutdoor,
		metrics.CoolSetpoint,
		metrics.HeatSetpoint,
		metrics.EquipmentStatus,
	)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
