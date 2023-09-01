package charts

import (
	"bytes"
	"text/template"
	"time"

	"github.com/redgoose/daikin-one/internal/db"
)

type Chart struct {
	Title      string
	Data       []db.PeriodData
	XAxisLabel string
	TempUnit   string
}

var chartTmpl = template.Must(template.ParseFiles("templates/chart.tmpl"))

func GetChartForDay(dbPath string, deviceId string, date time.Time, tempUnit string) string {
	output := ""
	data := db.GetDataForDay(dbPath, deviceId, date)

	if len(data) > 0 {
		chart := Chart{
			Title:      date.Format("January 2 2006"),
			Data:       data,
			XAxisLabel: "Hour",
			TempUnit:   tempUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForMonth(dbPath string, deviceId string, date time.Time, tempUnit string) string {
	output := ""
	data := db.GetDataForMonth(dbPath, deviceId, date)

	if len(data) > 0 {
		chart := Chart{
			Title:      date.Format("January 2006"),
			Data:       data,
			XAxisLabel: "Day",
			TempUnit:   tempUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForYear(dbPath string, deviceId string, date time.Time, tempUnit string) string {
	output := ""
	data := db.GetDataForYear(dbPath, deviceId, date)

	if len(data) > 0 {
		chart := Chart{
			Title:      date.Format("2006"),
			Data:       data,
			XAxisLabel: "Month",
			TempUnit:   tempUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}
