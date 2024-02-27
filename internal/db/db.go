package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redgoose/daikin-skyport"
)

type DeviceData struct {
	DeviceId        string
	TempIndoor      float32
	TempOutdoor     float32
	HumidityIndoor  int
	HumidityOutdoor int
	CoolSetpoint    float32
	HeatSetpoint    float32
	EquipmentStatus int
	OutdoorHeat     float32 // demand for heat from heat pump in %
	OutdoorCool     float32 // demand for cool from heat pump in %
	IndoorFan       float32 // hvac fan actual usage in %
	IndoorHeat      float32 // furnace heat actual usage %
}

type PeriodData struct {
	Period          string
	TempIndoor      float32
	TempOutdoor     float32
	HumidityIndoor  float32
	HumidityOutdoor float32
	CoolSetpoint    float32
	HeatSetpoint    float32
	RunTime         int
}

func LogData(dbPath string, data DeviceData) {
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	sqlStatement := `INSERT INTO daikin (
		timestamp,
		device_id,
		temp_indoor,
		temp_outdoor,
		humidity_indoor,
		humidity_outdoor,
		cool_setpoint,
		heat_setpoint,
		equipment_status,
		outdoor_heat,
		outdoor_cool,
		indoor_fan,
		indoor_heat
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	stmt, err := db.Prepare(sqlStatement)
	checkErr(err)

	var timestamp string = time.Now().Format(time.RFC3339)

	_, err = stmt.Exec(
		timestamp,
		data.DeviceId,
		data.TempIndoor,
		data.TempOutdoor,
		data.HumidityIndoor,
		data.HumidityOutdoor,
		data.CoolSetpoint,
		data.HeatSetpoint,
		data.EquipmentStatus,
		data.OutdoorHeat,
		data.OutdoorCool,
		data.IndoorFan,
		data.IndoorHeat,
	)
	checkErr(err)
}

func GetDataForDay(dbPath string, deviceId string, day time.Time) []PeriodData {
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	rows, err := db.Query(`
		select
			substr(timestamp, 12, 2) || ":00" as hour,
			round(avg(temp_indoor), 2) as temp_indoor,
			round(avg(temp_outdoor), 2) as temp_outdoor,
			round(avg(humidity_indoor), 2) as humidity_indoor,
			round(avg(humidity_outdoor), 2) as humidity_outdoor,
			round(avg(cool_setpoint), 2) as cool_setpoint,
			round(avg(heat_setpoint), 2) as heat_setpoint,
			sum(
				case when equipment_status = ? or equipment_status = ? or equipment_status = ?
				then 5
				else 0
				end
			) as run_time
		from daikin
		where substr(timestamp, 0, 11) = ?
		and device_id = ?
		group by substr(timestamp, 0, 14);
	`, daikin.EquipmentStatusCool, daikin.EquipmentStatusOvercool, daikin.EquipmentStatusHeat, day.Format("2006-01-02"), deviceId)
	checkErr(err)

	defer rows.Close()

	var dayData []PeriodData

	for rows.Next() {
		var data PeriodData

		err := rows.Scan(
			&data.Period,
			&data.TempIndoor,
			&data.TempOutdoor,
			&data.HumidityIndoor,
			&data.HumidityOutdoor,
			&data.CoolSetpoint,
			&data.HeatSetpoint,
			&data.RunTime,
		)
		checkErr(err)
		dayData = append(dayData, data)
	}

	err = rows.Err()
	checkErr(err)

	return dayData
}

func GetDataForMonth(dbPath string, deviceId string, day time.Time) []PeriodData {
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	rows, err := db.Query(`
		select
			substr(timestamp, 9, 2) as day,
			round(avg(temp_indoor), 2) as temp_indoor,
			round(avg(temp_outdoor), 2) as temp_outdoor,
			round(avg(humidity_indoor), 2) as humidity_indoor,
			round(avg(humidity_outdoor), 2) as humidity_outdoor,
			round(avg(cool_setpoint), 2) as cool_setpoint,
			round(avg(heat_setpoint), 2) as heat_setpoint,
			sum(
				case when equipment_status = ? or equipment_status = ? or equipment_status = ?
				then 5
				else 0
				end
			) as run_time
		from daikin
		where  substr(timestamp, 0, 8) = ?
		and device_id = ?
		group by substr(timestamp, 0, 11);
	`, daikin.EquipmentStatusCool, daikin.EquipmentStatusOvercool, daikin.EquipmentStatusHeat, day.Format("2006-01"), deviceId)
	checkErr(err)

	defer rows.Close()

	var monthData []PeriodData

	for rows.Next() {
		var data PeriodData

		err := rows.Scan(
			&data.Period,
			&data.TempIndoor,
			&data.TempOutdoor,
			&data.HumidityIndoor,
			&data.HumidityOutdoor,
			&data.CoolSetpoint,
			&data.HeatSetpoint,
			&data.RunTime,
		)
		checkErr(err)
		monthData = append(monthData, data)
	}

	err = rows.Err()
	checkErr(err)

	return monthData
}

func GetDataForYear(dbPath string, deviceId string, day time.Time) []PeriodData {
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	rows, err := db.Query(`
		select
			substr(timestamp, 6, 2) as month,
			round(avg(temp_indoor), 2) as temp_indoor,
			round(avg(temp_outdoor), 2) as temp_outdoor,
			round(avg(humidity_indoor), 2) as humidity_indoor,
			round(avg(humidity_outdoor), 2) as humidity_outdoor,
			round(avg(cool_setpoint), 2) as cool_setpoint,
			round(avg(heat_setpoint), 2) as heat_setpoint,
			sum(
				case when equipment_status = ? or equipment_status = ? or equipment_status = ?
				then 5
				else 0
				end
				) as run_time
		from daikin
		where  substr(timestamp, 0, 5) = ?
		and device_id = ?
		group by substr(timestamp, 0, 8);
	`, daikin.EquipmentStatusCool, daikin.EquipmentStatusOvercool, daikin.EquipmentStatusHeat, day.Format("2006"), deviceId)
	checkErr(err)

	defer rows.Close()

	var yearData []PeriodData

	for rows.Next() {
		var data PeriodData

		err := rows.Scan(
			&data.Period,
			&data.TempIndoor,
			&data.TempOutdoor,
			&data.HumidityIndoor,
			&data.HumidityOutdoor,
			&data.CoolSetpoint,
			&data.HeatSetpoint,
			&data.RunTime,
		)
		checkErr(err)
		yearData = append(yearData, data)
	}

	err = rows.Err()
	checkErr(err)

	return yearData
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
