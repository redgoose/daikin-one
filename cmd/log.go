package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "log",
	Args:  cobra.NoArgs,
	Short: "Logs thermostat metrics to local SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(deviceId)
	},
}

func init() {
	listCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	listCmd.MarkFlagRequired("device-id")

	listCmd.Flags().StringVar(&dbPath, "db", "daikin.db", "Path to SQLite db")

	rootCmd.AddCommand(listCmd)
}
