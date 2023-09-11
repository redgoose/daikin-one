package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redgoose/daikin-one/daikin"
	"github.com/redgoose/daikin-one/internal/db"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Args:  cobra.NoArgs,
	Short: "Logs device data to local SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			d := daikin.New(viper.GetString("apiKey"), viper.GetString("integratorToken"), viper.GetString("email"))
			var deviceInfo = d.GetDeviceInfo(deviceId)

			var data = db.DeviceData{
				DeviceId:        deviceId,
				TempIndoor:      deviceInfo.TempIndoor,
				TempOutdoor:     deviceInfo.TempOutdoor,
				HumidityIndoor:  deviceInfo.HumIndoor,
				HumidityOutdoor: deviceInfo.HumOutdoor,
				CoolSetpoint:    deviceInfo.CoolSetpoint,
				HeatSetpoint:    deviceInfo.HeatSetpoint,
				EquipmentStatus: deviceInfo.EquipmentStatus,
			}

			db.LogData(dbPath, data)
			fmt.Println(time.Now().Format(time.RFC3339) + " - Logged data")

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
