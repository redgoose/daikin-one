package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var chartCmd = &cobra.Command{
	Use:   "chart",
	Args:  cobra.NoArgs,
	Short: "Generates charts based on collected log metrics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dbPath)
	},
}

func init() {
	rootCmd.AddCommand(chartCmd)
	chartCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	chartCmd.MarkFlagRequired("device-id")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	chartCmd.Flags().StringVarP(&dbPath, "db", "", home+"/daikin.db", "Local path to SQLite database")
}
