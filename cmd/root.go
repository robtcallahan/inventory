package cmd

import (
	"fmt"
	cfg "inventory/pkg/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "inventory",
	Short:         "Get a hold of your personal things",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long:          ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
}

var configFile string

var config = &cfg.Config{}
var options = &cfg.Options{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config/config.json", "Configuration file path")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		config, err = cfg.ReadConfig(configFile)
		return err
	}
}
