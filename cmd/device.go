package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/redgoose/daikin-skyport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deviceMode int
var deviceCoolSetpoint float32
var deviceHeatSetpoint float32

var deviceCmd = &cobra.Command{
	Use:   "device",
	Args:  cobra.NoArgs,
	Short: "Manage devices",
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Args:  cobra.NoArgs,
	Short: "Retrieves device configuration and state values",
	Run: func(cmd *cobra.Command, args []string) {
		d := daikin.New(viper.GetString("email"), viper.GetString("password"))
		info, err := d.GetDeviceInfo(deviceId)
		if err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(info, "", "\t")
		fmt.Println(string(s))
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Args:  cobra.NoArgs,
	Short: "Lists devices associated with your account",
	Run: func(cmd *cobra.Command, args []string) {
		d := daikin.New(viper.GetString("email"), viper.GetString("password"))
		devices, err := d.GetDevices()
		if err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(devices, "", "\t")
		fmt.Println(string(s))
	},
}

var modeCmd = &cobra.Command{
	Use:   "mode",
	Args:  cobra.NoArgs,
	Short: "Update device operating mode",
	Run: func(cmd *cobra.Command, args []string) {
		d := daikin.New(viper.GetString("email"), viper.GetString("password"))
		err := d.SetMode(deviceId, daikin.Mode(deviceMode))
		if err != nil {
			panic(err)
		}
	},
}

var tempCmd = &cobra.Command{
	Use:   "temp",
	Args:  cobra.NoArgs,
	Short: "Update cooling/heating setpoint(s)",
	Run: func(cmd *cobra.Command, args []string) {
		d := daikin.New(viper.GetString("email"), viper.GetString("password"))
		var params = daikin.SetTempParams{
			HeatSetpoint: deviceHeatSetpoint,
			CoolSetpoint: deviceCoolSetpoint,
		}
		err := d.SetTemp(deviceId, params)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)

	deviceCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	infoCmd.MarkFlagRequired("device-id")

	deviceCmd.AddCommand(lsCmd)

	deviceCmd.AddCommand(modeCmd)
	modeCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	modeCmd.Flags().IntVarP(&deviceMode, "mode", "", 0, "Device mode")
	modeCmd.MarkFlagRequired("device-id")
	modeCmd.MarkFlagRequired("mode")

	deviceCmd.AddCommand(tempCmd)
	tempCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	tempCmd.Flags().Float32VarP(&deviceHeatSetpoint, "heat", "", 0, "Heat setpoint")
	tempCmd.Flags().Float32VarP(&deviceCoolSetpoint, "cool", "", 0, "Cool setpoint")
	tempCmd.MarkFlagRequired("device-id")
}
