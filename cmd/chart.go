package cmd

import (
	"fmt"

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
	chartCmd.Flags().StringVar(&dbPath, "db", "daikin.db", "Path to SQLite db")

	rootCmd.AddCommand(chartCmd)
}
