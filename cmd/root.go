package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var deviceId string
var dbPath string

var rootCmd = &cobra.Command{
	Use:   "daikin-one",
	Short: "daikin-one is a cli to interact with Daikin One devices",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", home+"/.daikin.yaml", "config file")
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

	if err := viper.ReadInConfig(); err == nil {
		configValues := []string{"integratorToken", "apiKey", "email"}
		for _, configValue := range configValues {
			if viper.GetString(configValue) == "" {
				cobra.CheckErr(fmt.Errorf("%s not defined in config", configValue))
			}
		}
	} else {
		cobra.CheckErr(err)
	}
}
