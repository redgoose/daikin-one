package charts

import (
	"bytes"
	"html/template"
	"time"

	"github.com/redgoose/daikin-one/internal/db"
	"github.com/redgoose/daikin-one/templates"
)

type Chart struct {
	Title           string
	Data            []db.PeriodData
	XAxisLabel      string
	TemperatureUnit string
}

type ChartV2 struct {
	Title           string
	Data            []db.AnyData
	XAxisLabel      string
	TemperatureUnit string
}

var chartTmpl *template.Template

// Convert PeriodData temperature fields from C to F
func convertTempsCtoF(periods []db.PeriodData) []db.PeriodData {

	cToF := func(c float32) float32 {
		return c*9/5 + 32
	}

	for i := range periods {
		periods[i].TempIndoor = cToF(periods[i].TempIndoor)
		periods[i].TempOutdoor = cToF(periods[i].TempOutdoor)
		periods[i].CoolSetpoint = cToF(periods[i].CoolSetpoint)
		periods[i].HeatSetpoint = cToF(periods[i].HeatSetpoint)
	}
	return periods
}

// Convert PeriodData temperature fields from C to F
func convertTempsCtoFv2(periods []db.AnyData) []db.AnyData {
	cToF := func(c float32) float32 {
		return c*9/5 + 32
	}
	for i := range periods {
		periods[i].Data = cToF(periods[i].Data)
	}
	return periods
}

func GetChartForField(dbPath string, deviceId string, field string, date time.Time, temperatureUnit string) string {
	output := ""

	data := db.GetDataRaw(dbPath, deviceId, field)

	// TODO: Not all data is temperature, need to maintain a map of data tpe to display type...
	// if temperatureUnit == "F" {
	// data = convertTempsCtoFv2(data)
	// }

	if len(data) > 0 {
		chart := ChartV2{
			Title:           date.Format("January 2 2006"),
			Data:            data,
			XAxisLabel:      "Hour",
			TemperatureUnit: temperatureUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForDay(dbPath string, deviceId string, date time.Time, temperatureUnit string) string {
	output := ""
	data := db.GetDataForDay(dbPath, deviceId, date)

	if temperatureUnit == "F" {
		data = convertTempsCtoF(data)
	}

	if len(data) > 0 {
		chart := Chart{
			Title:           date.Format("January 2 2006"),
			Data:            data,
			XAxisLabel:      "Hour",
			TemperatureUnit: temperatureUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForMonth(dbPath string, deviceId string, date time.Time, temperatureUnit string) string {
	output := ""
	data := db.GetDataForMonth(dbPath, deviceId, date)

	if temperatureUnit == "F" {
		data = convertTempsCtoF(data)
	}

	if len(data) > 0 {
		chart := Chart{
			Title:           date.Format("January 2006"),
			Data:            data,
			XAxisLabel:      "Day",
			TemperatureUnit: temperatureUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func GetChartForYear(dbPath string, deviceId string, date time.Time, temperatureUnit string) string {
	output := ""
	data := db.GetDataForYear(dbPath, deviceId, date)

	if temperatureUnit == "F" {
		data = convertTempsCtoF(data)
	}

	if len(data) > 0 {
		chart := Chart{
			Title:           date.Format("2006"),
			Data:            data,
			XAxisLabel:      "Month",
			TemperatureUnit: temperatureUnit,
		}

		buf := new(bytes.Buffer)
		chartTmpl.Execute(buf, chart)
		output = buf.String()
	}

	return output
}

func init() {
	var err error
	chartTmpl, err = template.ParseFS(templates.TemplatesFS, "tmpl/chart.tmpl")
	if err != nil {
		panic(err)
	}
}
