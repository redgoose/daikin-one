package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/redgoose/daikin-one/daikin"
	db "github.com/redgoose/daikin-one/internal"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Args:  cobra.NoArgs,
	Short: "Logs device metrics to local SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			var deviceInfo = daikin.GetDeviceInfo(deviceId)

			var metric = db.Metrics{
				DeviceId:        deviceId,
				TempIndoor:      deviceInfo.TempIndoor,
				TempOutdoor:     deviceInfo.TempOutdoor,
				HumidityIndoor:  deviceInfo.HumIndoor,
				HumidityOutdoor: deviceInfo.HumOutdoor,
				CoolSetpoint:    deviceInfo.CoolSetpoint,
				HeatSetpoint:    deviceInfo.HeatSetpoint,
				EquipmentStatus: deviceInfo.EquipmentStatus,
			}

			db.LogMetrics(dbPath, metric)
			fmt.Println(time.Now().Format(time.RFC3339) + " - Logged metrics")

			time.Sleep(5 * time.Minute)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	logCmd.MarkFlagRequired("device-id")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	logCmd.Flags().StringVarP(&dbPath, "db", "", filepath.Join(home, ".daikin", "daikin.db"), "Local path to SQLite database")
}
