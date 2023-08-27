package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	db "github.com/redgoose/daikin-one/internal"
	"github.com/spf13/cobra"
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
		var data = db.GetDataForDay(dbPath, deviceId, time.Now())

		type Foo struct {
			Title           string
			Hour            []string
			TempIndoor      []float32
			TempOutdoor     []float32
			HumidityIndoor  []float32
			HumidityOutdoor []float32
			CoolSetpoint    []float32
			HeatSetpoint    []float32
			RunTime         []int
		}

		data2 := Foo{
			Title:           time.Now().Format("2006-01-02"),
			Hour:            data.Hour,
			TempIndoor:      data.TempIndoor,
			TempOutdoor:     data.TempOutdoor,
			HumidityIndoor:  data.HumidityIndoor,
			HumidityOutdoor: data.HumidityOutdoor,
			CoolSetpoint:    data.CoolSetpoint,
			HeatSetpoint:    data.HeatSetpoint,
			RunTime:         data.RunTime,
		}

		tmpl := template.Must(template.ParseFiles("templates/chart.tmpl"))

		var doc bytes.Buffer
		tmpl.Execute(&doc, data2)
		s := doc.String()
		fmt.Println(s)

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
