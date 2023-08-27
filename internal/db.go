package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	EquipmentStatusCool     = 1
	EquipmentStatusOvercool = 2
	EquipmentStatusHeat     = 3
	EquipmentStatusFan      = 4
	EquipmentStatusIdle     = 5
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
}

type DayData struct {
	Hour            []string
	TempIndoor      []float32
	TempOutdoor     []float32
	HumidityIndoor  []float32
	HumidityOutdoor []float32
	CoolSetpoint    []float32
	HeatSetpoint    []float32
	RunTime         []int
}

func LogData(dbPath string, data DeviceData) {
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO daikin(timestamp, device_id, temp_indoor, temp_outdoor, humidity_indoor, humidity_outdoor, cool_setpoint, heat_setpoint, equipment_status) values(?,?,?,?,?,?,?,?,?)")
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
	)
	checkErr(err)
}

func GetDataForDay(dbPath string, deviceId string, day time.Time) DayData {
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
	`, EquipmentStatusCool, EquipmentStatusOvercool, EquipmentStatusHeat, day.Format("2006-01-02"), deviceId)
	checkErr(err)

	defer rows.Close()

	var dayData DayData

	for rows.Next() {
		var (
			hour            string
			tempIndoor      float32
			tempOutdoor     float32
			humidityIndoor  float32
			humidityOutdoor float32
			coolSetpoint    float32
			heatSetpoint    float32
			runTime         int
		)
		err := rows.Scan(
			&hour,
			&tempIndoor,
			&tempOutdoor,
			&humidityIndoor,
			&humidityOutdoor,
			&coolSetpoint,
			&heatSetpoint,
			&runTime,
		)
		checkErr(err)

		dayData.Hour = append(dayData.Hour, hour)
		dayData.TempIndoor = append(dayData.TempIndoor, tempIndoor)
		dayData.TempOutdoor = append(dayData.TempOutdoor, tempOutdoor)
		dayData.HumidityIndoor = append(dayData.HumidityIndoor, humidityIndoor)
		dayData.HumidityOutdoor = append(dayData.HumidityOutdoor, humidityOutdoor)
		dayData.CoolSetpoint = append(dayData.CoolSetpoint, coolSetpoint)
		dayData.HeatSetpoint = append(dayData.HeatSetpoint, heatSetpoint)
		dayData.RunTime = append(dayData.RunTime, runTime)
	}

	err = rows.Err()
	checkErr(err)

	return dayData
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
