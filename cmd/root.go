package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "giveth",
	Short: "A command line tool to study Giveth's data",
	Long: `This is a simple command line tool that allows us to download data from the
Giveth APIs. While we need other (i.e., on-chain) data to complete the project,
this is a good start. We [document the data structures here](./data/QUESTIONS.md).
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.giveth.yaml)")
	rootCmd.PersistentFlags().MarkHidden("config")

	rootCmd.PersistentFlags().Uint64P("round", "r", 0, "Limits the list of rounds to a single round")
	rootCmd.PersistentFlags().BoolP("update", "u", false, "If present, data commands pull data from Giveth's APIs")
	rootCmd.PersistentFlags().BoolP("script", "c", false, "If present, data commands generate bash script to query Giveth's APIs")
	rootCmd.PersistentFlags().Uint64P("sleep", "s", 0, "Instructs the tool how long to sleep between invocations")
	rootCmd.PersistentFlags().StringP("fmt", "x", "", "One of [json|csv|txt]")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "If present, certain commands will display extra data")
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".giveth" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".giveth")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
