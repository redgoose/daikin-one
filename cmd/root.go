package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.daikin.yaml)")
}

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgFile = ".daikin.yaml"

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error: config file `%s` not found\n", cfgFile)
		os.Exit(1)
	}
}
