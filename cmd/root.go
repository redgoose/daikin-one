package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbPath string
var deviceId string

var rootCmd = &cobra.Command{
	Use:   "daikin-one",
	Short: "daikin-one is a cli to interact with Daikin One thermostats",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
