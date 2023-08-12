package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chartCmd = &cobra.Command{
	Use:   "chart",
	Args:  cobra.NoArgs,
	Short: "Generates charts based on collected log metrics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("dbPath"))
	},
}

func init() {
	rootCmd.AddCommand(chartCmd)
}
