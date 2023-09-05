package cmd

import (
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/redgoose/daikin-one/internal/charts"
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
		temperatureUnit := viper.GetString("temperatureUnit")
		allCharts := ""

		// last 7 days
		for i := 0; i <= 6; i++ {
			date := time.Now().Add(time.Duration(-i*24) * time.Hour)
			allCharts += charts.GetChartForDay(dbPath, deviceId, date, temperatureUnit)
		}

		allCharts += charts.GetChartForMonth(dbPath, deviceId, time.Now(), temperatureUnit)
		allCharts += charts.GetChartForYear(dbPath, deviceId, time.Now(), temperatureUnit)

		baseTmpl := template.Must(template.ParseFiles("templates/base.tmpl"))
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

		folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		baseTmpl := template.Must(template.ParseFiles(filepath.Join(folder, "templates", "base.tmpl")))

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

		folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		baseTmpl := template.Must(template.ParseFiles(filepath.Join(folder, "templates", "base.tmpl")))

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

		folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		baseTmpl := template.Must(template.ParseFiles(filepath.Join(folder, "templates", "base.tmpl")))

		baseTmpl.Execute(os.Stdout, chart)
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
	reportCmd.AddCommand(reportDayCmd)
	reportCmd.AddCommand(reportMonthCmd)
	reportCmd.AddCommand(reportYearCmd)
}
