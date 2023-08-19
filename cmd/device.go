package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/redgoose/daikin-one/daikin"
	"github.com/spf13/cobra"
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
		var info = daikin.GetDeviceInfo(deviceId)
		s, _ := json.MarshalIndent(info, "", "\t")
		fmt.Println(string(s))
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Args:  cobra.NoArgs,
	Short: "Lists devices associated with your account",
	Run: func(cmd *cobra.Command, args []string) {
		var locations = daikin.ListDevices()
		s, _ := json.MarshalIndent(locations, "", "\t")
		fmt.Println(string(s))
	},
}

var modeSetpointCmd = &cobra.Command{
	Use:   "mode-setpoint",
	Args:  cobra.NoArgs,
	Short: "Update device operating mode and heat/cool setpoints",
	Run: func(cmd *cobra.Command, args []string) {
		var deviceOptions = daikin.DeviceOptions{
			Mode:         deviceMode,
			HeatSetpoint: deviceHeatSetpoint,
			CoolSetpoint: deviceCoolSetpoint,
		}
		daikin.UpdateDevice(deviceId, deviceOptions)
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)

	deviceCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	infoCmd.MarkFlagRequired("device-id")

	deviceCmd.AddCommand(lsCmd)

	deviceCmd.AddCommand(modeSetpointCmd)
	modeSetpointCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	modeSetpointCmd.Flags().IntVarP(&deviceMode, "mode", "", 0, "Device mode")
	modeSetpointCmd.Flags().Float32VarP(&deviceHeatSetpoint, "heat", "", 0, "Heat setpoint")
	modeSetpointCmd.Flags().Float32VarP(&deviceCoolSetpoint, "cool", "", 0, "Cool setpoint")
	modeSetpointCmd.MarkFlagRequired("device-id")
	modeSetpointCmd.MarkFlagRequired("mode")
	modeSetpointCmd.MarkFlagRequired("heat")
	modeSetpointCmd.MarkFlagRequired("cool")
}
