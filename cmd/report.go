package cmd

import (
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/redgoose/daikin-one/internal/charts"
	"github.com/redgoose/daikin-one/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var baseTmpl *template.Template

var reportCmd = &cobra.Command{
	Use:   "report",
	Args:  cobra.NoArgs,
	Short: "Generate reports",
}

var reportAllCmd = &cobra.Command{
	Use:   "all",
	Args:  cobra.NoArgs,
	Short: "Generates report of all data, separating them one per graph",
	Run: func(cmd *cobra.Command, args []string) {
		temperatureUnit := viper.GetString("temperatureUnit")
		allCharts := ""

		fields := []string{
			"temp_outdoor",
			"temp_indoor",
			"humidity_outdoor",
			"humidity_indoor",
			"cool_setpoint",
			"heat_setpoint",
			"outdoor_heat",
			"outdoor_cool",
			"indoor_fan",
			"indoor_heat",
			// "equipment_status", -- I think this is fully covered by other charts now?
		}

		for _, field := range fields {
			allCharts += charts.GetChartForField(dbPath, deviceId, field, temperatureUnit)
		}

		baseTmpl.Execute(os.Stdout, allCharts)
	},
}

var reportSummaryCmd = &cobra.Command{
	Use:   "summary",
	Args:  cobra.NoArgs,
	Short: "Generates summary report",
	Run: func(cmd *cobra.Command, args []string) {
		temperatureUnit := viper.GetString("temperatureUnit")
		allCharts := ""

		// last 7 days
		for i := 0; i <= 6; i++ {
			date := time.Now().Add(time.Duration(-i*24) * time.Hour)
			allCharts += charts.GetChartForDay(dbPath, deviceId, date, temperatureUnit)
		}

		allCharts += charts.GetChartForMonth(dbPath, deviceId, time.Now(), temperatureUnit)
		allCharts += charts.GetChartForYear(dbPath, deviceId, time.Now(), temperatureUnit)

		baseTmpl.Execute(os.Stdout, allCharts)
	},
}

var reportDayCmd = &cobra.Command{
	Use:   "day",
	Args:  cobra.ExactArgs(1),
	Short: "Generates report for given day",
	Run: func(cmd *cobra.Command, args []string) {
		date, err := time.Parse("2006-01-02", args[0])
		cobra.CheckErr(err)

		temperatureUnit := viper.GetString("temperatureUnit")

		chart := charts.GetChartForDay(dbPath, deviceId, date, temperatureUnit)
		baseTmpl.Execute(os.Stdout, chart)
	},
}

var reportMonthCmd = &cobra.Command{
	Use:   "month",
	Args:  cobra.ExactArgs(1),
	Short: "Generates report for given month",
	Run: func(cmd *cobra.Command, args []string) {
		date, err := time.Parse("2006-01", args[0])
		cobra.CheckErr(err)

		temperatureUnit := viper.GetString("temperatureUnit")

		chart := charts.GetChartForMonth(dbPath, deviceId, date, temperatureUnit)
		baseTmpl.Execute(os.Stdout, chart)
	},
}

var reportYearCmd = &cobra.Command{
	Use:   "year",
	Args:  cobra.ExactArgs(1),
	Short: "Generates report for given year",
	Run: func(cmd *cobra.Command, args []string) {
		date, err := time.Parse("2006", args[0])
		cobra.CheckErr(err)

		temperatureUnit := viper.GetString("temperatureUnit")

		chart := charts.GetChartForYear(dbPath, deviceId, date, temperatureUnit)
		baseTmpl.Execute(os.Stdout, chart)
	},
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	reportCmd.PersistentFlags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	reportCmd.PersistentFlags().StringVarP(&dbPath, "db", "", filepath.Join(home, ".daikin", "daikin.db"), "Local path to SQLite database")
	reportCmd.MarkPersistentFlagRequired("device-id")

	reportCmd.AddCommand(reportAllCmd)
	reportCmd.AddCommand(reportSummaryCmd)
	reportCmd.AddCommand(reportDayCmd)
	reportCmd.AddCommand(reportMonthCmd)
	reportCmd.AddCommand(reportYearCmd)

	rootCmd.AddCommand(reportCmd)

	baseTmpl, err = template.ParseFS(templates.TemplatesFS, "tmpl/base.tmpl")
	if err != nil {
		panic(err)
	}
}
