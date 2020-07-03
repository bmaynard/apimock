package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	c "github.com/bmaynard/apimock/pkg/cmd"
	"github.com/bmaynard/apimock/pkg/config"
	l "github.com/bmaynard/apimock/pkg/utils/logger"
)

var (
	cfgFile string
)

func NewApiMockApp() *cobra.Command {
	cobra.OnInitialize(initConfig)

	var rootCmd = &cobra.Command{
		Use:   "apimock",
		Short: "API Mock server",
		Long:  `Run an API Mock server as well as the ability to record real requests to mock later`,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.apimock.yaml)")
	rootCmd.AddCommand(c.NewCmdProxyServer())
	rootCmd.AddCommand(c.NewCmdMockServer())

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".apimock" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".apimock")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config.Configuration)

	if err != nil {
		l.Log.Errorf("unable to decode into struct, %v", err)
	}
}
