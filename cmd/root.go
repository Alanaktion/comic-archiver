package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "comic-archiver",
		Short: "A tool for archiving web comics",
		Long:  `Comic Archiver is a tool for downloading and serving web comics.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/comic-archiver.yaml)")
}

func initConfig() {
	viper.SetDefault("ComicsDir", "comics")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Read config in ~/.config/comic-archiver.yaml
		viper.AddConfigPath("$HOME/.config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("comic-archiver")
	}

	viper.SafeWriteConfig()
	viper.AutomaticEnv()

	viper.ReadInConfig()
	dir := viper.GetString("ComicsDir")
	fmt.Println("Using comic dir:", dir)
	os.MkdirAll(dir, os.ModePerm)
	if err := os.Chdir(viper.GetString("ComicsDir")); err != nil {
		log.Fatal(err)
	}
}
