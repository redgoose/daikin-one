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

var startDate string
var endDate string

// The number of days to default date range to when a full range is not provided
// Could move to config file.
const defaultDateRange = -1 * 3

var reportCmd = &cobra.Command{
	Use:   "report",
	Args:  cobra.NoArgs,
	Short: "Generate reports",
}

// getTimeRange calculates the start and end dates based on provided arguments.
// If a parameter is missing it defaults the range based on provided arg or today and the defaultDateRange constant.
// Returns start and end time.Time objects.
func getTimeRange(startDateStr, endDateStr string) (time.Time, time.Time) {
	layout := "2006-01-02 15:04" // Define the date layout (Date is user provided, time is appended in this method)
	startTime := " 00:00"        // Start dates should start in morning
	endTime := " 23:59"          // End dates should end at night
	var err error
	now := time.Now()
	var startDate, endDate time.Time

	switch {
	case startDateStr == "" && endDateStr == "":
		// Neither start nor end date are set, set start date from 7 days ago and end date to now
		startDate = now.AddDate(0, 0, defaultDateRange)
		endDate = now
	case startDateStr != "" && endDateStr == "":
		// Only start date is set, default end date to now
		startDate, err = time.Parse(layout, startDateStr+startTime)
		endDate = now
	case startDateStr == "" && endDateStr != "":
		// Only end date is set, default start date to 7 days before end date
		endDate, err = time.Parse(layout, endDateStr+endTime)
		startDate = endDate.AddDate(0, 0, defaultDateRange)
	default:
		// Both start and end dates are set
		endDate, err = time.Parse(layout, endDateStr+endTime)
		startDate, err = time.Parse(layout, startDateStr+startTime)
	}
	if err != nil {
		panic(err)
	}

	return startDate, endDate
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

		startTime, endTime := getTimeRange(startDate, endDate)

		for _, field := range fields {
			allCharts += charts.GetChartForField(dbPath, deviceId, field, startTime, endTime, temperatureUnit)
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

	reportAllCmd.Flags().StringVarP(&startDate, "start", "s", "", "Start date in format YYYY-MM-DD")
	reportAllCmd.Flags().StringVarP(&endDate, "end", "e", "", "End date in format YYYY-MM-DD")

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
