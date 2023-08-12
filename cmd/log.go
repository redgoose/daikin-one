package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Args:  cobra.NoArgs,
	Short: "Logs thermostat metrics to local SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("deviceId"))
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
