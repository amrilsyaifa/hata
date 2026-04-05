package cmd

import (
	"fmt"
	"os"

	"github.com/amrilsyaifa/hata/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "hata",
	Short: "A CLI tool for syncing i18n data between local files and Google Sheets",
	Long: `Hata bridges your codebase translation keys and Google Sheets,
enabling smooth collaboration between developers and non-technical translators.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", config.DefaultConfigFile, "config file path")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(diffCmd)
}

func loadConfig() (*config.Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	var err error
	cfg, err = config.Load(cfgFile)
	return cfg, err
}
