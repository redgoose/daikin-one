package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
	"time"

	db "github.com/redgoose/daikin-one/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Args:  cobra.NoArgs,
	Short: "Generate reports",
}

var reportSummaryCmd = &cobra.Command{
	Use:   "summary",
	Args:  cobra.NoArgs,
	Short: "Generates summary report",
	Run: func(cmd *cobra.Command, args []string) {
		tempUnit := viper.GetString("temperatureUnit")

		type Chart struct {
			Title      string
			Data       []db.PeriodData
			XAxisLabel string
			TempUnit   string
		}

		chartTmpl := template.Must(template.ParseFiles("templates/chart.tmpl"))
		charts := ""

		// last 7 days
		for i := 0; i <= 6; i++ {
			date := time.Now().Add(time.Duration(-i*24) * time.Hour)
			data := db.GetDataForDay(dbPath, deviceId, date)

			if len(data) == 0 {
				continue
			}

			chart := Chart{
				Title:      date.Format("January 2 2006"),
				Data:       data,
				XAxisLabel: "Hour",
				TempUnit:   tempUnit,
			}

			buf := new(bytes.Buffer)
			chartTmpl.Execute(buf, chart)
			charts += buf.String()
		}

		// current month
		monthData := db.GetDataForMonth(dbPath, deviceId, time.Now())
		if len(monthData) > 0 {
			chart := Chart{
				Title:      time.Now().Format("January 2006"),
				Data:       monthData,
				XAxisLabel: "Day",
				TempUnit:   tempUnit,
			}

			buf := new(bytes.Buffer)
			chartTmpl.Execute(buf, chart)
			charts += buf.String()
		}

		// current year
		yearData := db.GetDataForYear(dbPath, deviceId, time.Now())
		if len(yearData) > 0 {
			chart := Chart{
				Title:      time.Now().Format("2006"),
				Data:       yearData,
				XAxisLabel: "Month",
				TempUnit:   tempUnit,
			}

			buf := new(bytes.Buffer)
			chartTmpl.Execute(buf, chart)
			charts += buf.String()
		}

		baseTmpl := template.Must(template.ParseFiles("templates/base.tmpl"))
		baseTmpl.Execute(os.Stdout, charts)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	rootCmd.PersistentFlags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "db", "", filepath.Join(home, ".daikin", "daikin.db"), "Local path to SQLite database")
	rootCmd.MarkPersistentFlagRequired("device-id")

	reportCmd.AddCommand(reportSummaryCmd)
}
