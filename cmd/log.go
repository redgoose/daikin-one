package cmd

import (
	"fmt"
	"os"

	"github.com/redgoose/daikin-one/daikin"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Args:  cobra.NoArgs,
	Short: "Logs device metrics to local SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(deviceId, dbPath)
		fmt.Println(daikin.GetToken())
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	logCmd.MarkFlagRequired("device-id")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	logCmd.Flags().StringVarP(&dbPath, "db", "", home+"/daikin.db", "Local path to SQLite database")
}
