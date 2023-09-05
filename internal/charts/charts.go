package charts

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/redgoose/daikin-one/internal/db"
)

type Chart struct {
	Title           string
	Data            []db.PeriodData
	XAxisLabel      string
	TemperatureUnit string
}

func GetChartForDay(dbPath string, deviceId string, date time.Time, temperatureUnit string) string {
	output := ""
	data := db.GetDataForDay(dbPath, deviceId, date)

	if len(data) > 0 {
		chart := Chart{
			Title:           date.Format("January 2 2006"),
			Data:            data,
			XAxisLabel:      "Hour",
			TemperatureUnit: temperatureUnit,
		}

		folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		var chartTmpl = template.Must(template.ParseFiles(filepath.Join(folder, "templates", "chart.tmpl")))

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForMonth(dbPath string, deviceId string, date time.Time, temperatureUnit string) string {
	output := ""
	data := db.GetDataForMonth(dbPath, deviceId, date)

	if len(data) > 0 {
		chart := Chart{
			Title:           date.Format("January 2006"),
			Data:            data,
			XAxisLabel:      "Day",
			TemperatureUnit: temperatureUnit,
		}

		folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		var chartTmpl = template.Must(template.ParseFiles(filepath.Join(folder, "templates", "chart.tmpl")))

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForYear(dbPath string, deviceId string, date time.Time, temperatureUnit string) string {
	output := ""
	data := db.GetDataForYear(dbPath, deviceId, date)

	if len(data) > 0 {
		chart := Chart{
			Title:           date.Format("2006"),
			Data:            data,
			XAxisLabel:      "Month",
			TemperatureUnit: temperatureUnit,
		}

		folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		var chartTmpl = template.Must(template.ParseFiles(filepath.Join(folder, "templates", "chart.tmpl")))

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}
