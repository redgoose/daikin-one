package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
		fmt.Println(deviceId)
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Args:  cobra.NoArgs,
	Short: "Lists devices associated with your account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List devices")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Args:  cobra.NoArgs,
	Short: "Update device operating mode and heat/cool setpoints",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(deviceId)
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	infoCmd.MarkFlagRequired("device-id")

	deviceCmd.AddCommand(lsCmd)

	deviceCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	updateCmd.MarkFlagRequired("device-id")
}
