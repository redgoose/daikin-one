package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redgoose/daikin-one/internal/db"
	"github.com/redgoose/daikin-skyport"
)

// The data from the device is 2x actual. So percentages values range from 0-200.
const DAIKIN_PERCENT_MULTIPLIER = 2

var logCmd = &cobra.Command{
	Use:   "log",
	Args:  cobra.NoArgs,
	Short: "Logs device data to local SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		d := daikin.New(viper.GetString("email"), viper.GetString("password"))
		deviceInfo, err := d.GetDeviceInfo(deviceId)

		if err != nil {
			panic(err)
		}

		var data = db.DeviceData{
			DeviceId:        deviceId,
			TempIndoor:      deviceInfo.TempIndoor,
			TempOutdoor:     deviceInfo.TempOutdoor,
			HumidityIndoor:  deviceInfo.HumIndoor,
			HumidityOutdoor: deviceInfo.HumOutdoor,
			CoolSetpoint:    deviceInfo.CspHome,
			HeatSetpoint:    deviceInfo.HspHome,
			EquipmentStatus: int(deviceInfo.EquipmentStatus),
			OutdoorHeat:     float32(deviceInfo.CtOutdoorHeatRequestedDemand) / DAIKIN_PERCENT_MULTIPLIER,
			OutdoorCool:     float32(deviceInfo.CtOutdoorCoolRequestedDemand) / DAIKIN_PERCENT_MULTIPLIER,
			IndoorFan:       float32(deviceInfo.CtIFCCurrentFanActualStatus) / DAIKIN_PERCENT_MULTIPLIER,
			IndoorHeat:      float32(deviceInfo.CtIFCCurrentHeatActualStatus) / DAIKIN_PERCENT_MULTIPLIER,
		}

		db.LogData(dbPath, data)
		fmt.Println(time.Now().Format(time.RFC3339) + " - Logged data")
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
