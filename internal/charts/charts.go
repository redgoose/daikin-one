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
	Title      string
	PeriodData []db.AnyData
	XAxisLabel string
	YAxisUnit  string
}

var chartTmpl *template.Template
var chartV2Tmpl *template.Template

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

func convertDbTimeToDisplayTime(periods []db.AnyData) []db.AnyData {
	for i := range periods {
		t, err := time.Parse(time.RFC3339, periods[i].Period)
		if err != nil {
			panic(err)
		}

		periods[i].Period = t.Format("Jan 02 15:04")
	}
	return periods
}

func GetChartForField(dbPath string, deviceId string, field string, temperatureUnit string) string {
	output := ""

	data := db.GetDataRaw(dbPath, deviceId, field)

	if len(data) > 0 {
		// All data in the array uses the same displayUnit, so grab from the first element.
		// If the displayUnit is C and if the config is set to F then convert the data and switch the displayUnit to F for chart.
		displayUnit := data[0].Unit
		if displayUnit == "°C" && temperatureUnit == "F" {
			data = convertTempsCtoFv2(data)
			displayUnit = "°F"
		}

		data = convertDbTimeToDisplayTime(data)

		chart := ChartV2{
			Title:      field,
			PeriodData: data,
			XAxisLabel: "Time",
			YAxisUnit:  displayUnit,
		}

		buf := new(bytes.Buffer)
		err := chartV2Tmpl.Execute(buf, chart)
		if err != nil {
			panic(err)
		}
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
	chartV2Tmpl, err = template.ParseFS(templates.TemplatesFS, "tmpl/chartV2.tmpl")
	if err != nil {
		panic(err)
	}
}
