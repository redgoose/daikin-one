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
		tempUnit := viper.GetString("temperatureUnit")
		chartsString := ""

		// last 7 days
		for i := 0; i <= 6; i++ {
			date := time.Now().Add(time.Duration(-i*24) * time.Hour)
			chartsString += charts.GetChartForDay(dbPath, deviceId, date, tempUnit)
		}

		chartsString += charts.GetChartForMonth(dbPath, deviceId, time.Now(), tempUnit)
		chartsString += charts.GetChartForYear(dbPath, deviceId, time.Now(), tempUnit)

		baseTmpl := template.Must(template.ParseFiles("templates/base.tmpl"))
		baseTmpl.Execute(os.Stdout, chartsString)
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
