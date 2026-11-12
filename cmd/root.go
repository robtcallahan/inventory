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

const ConfigFile = "config/config.json"

var config = &cfg.Config{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	var err error
	config, err = cfg.ReadConfig(ConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}
