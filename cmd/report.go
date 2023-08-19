package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Args:  cobra.NoArgs,
	Short: "Generate reports",
}

var reportSummaryCmd = &cobra.Command{
	Use:   "summary",
	Args:  cobra.NoArgs,
	Short: "Generates summary report",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dbPath, deviceId)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.AddCommand(reportSummaryCmd)
	reportSummaryCmd.Flags().StringVarP(&deviceId, "device-id", "d", "", "Daikin device ID")
	reportSummaryCmd.MarkFlagRequired("device-id")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	reportSummaryCmd.Flags().StringVarP(&dbPath, "db", "", filepath.Join(home, ".daikin", "daikin.db"), "Local path to SQLite database")
}
